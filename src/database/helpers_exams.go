package database

import (
	"ExamSphere/src/core/utils/logging"
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/jackc/pgx/v5"
)

// CreateNewExam creates a new exam in the database,
func CreateNewExam(data *NewExamData) (*ExamInfo, error) {
	if data.Price == "" {
		data.Price = DefaultExamPrice
	}

	data.ExamTitle = strings.TrimSpace(data.ExamTitle)
	data.ExamDescription = strings.TrimSpace(data.ExamDescription)

	info := &ExamInfo{
		CourseId:        data.CourseId,
		ExamTitle:       data.ExamTitle,
		ExamDescription: data.ExamDescription,
		Price:           data.Price,
		CreatedBy:       data.CreatedBy,
		IsPublic:        data.IsPublic,
		Duration:        data.Duration,
		ExamDate:        data.ExamDate,
		CreatedAt:       time.Now(),
		mut:             &sync.RWMutex{},
	}

	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT create_exam_info(
			p_course_id := $1,
			p_exam_title := $2,
			p_exam_description := $3,
			p_price := $4,
			p_created_by := $5,
			p_is_public := $6,
			p_duration := $7,
			p_exam_date := $8
		)`,
		info.CourseId,
		info.ExamTitle,
		info.ExamDescription,
		info.Price,
		info.CreatedBy,
		info.IsPublic,
		info.Duration,
		data.ExamDate.Format(ExamDateLayout),
	).Scan(&info.ExamId)
	if err != nil {
		return nil, err
	}

	examsInfoMap.Add(info.ExamId, info)
	return info, nil
}

// GetExamInfo gets an exam from the database.
func GetExamInfo(examId int) (*ExamInfo, error) {
	info := examsInfoMap.Get(examId)
	if info != nil && info != valueExamNotFound && info.ExamId == examId {
		return info, nil
	}

	info = &ExamInfo{
		ExamId: examId,
		mut:    &sync.RWMutex{},
	}
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT exam_id,
			course_id, 
			exam_title,
			exam_description,
			price, 
			created_at, 
			exam_date, 
			duration, 
			created_by, 
			is_public
		FROM exam_info WHERE exam_id = $1`,
		examId,
	).Scan(
		&info.ExamId,
		&info.CourseId,
		&info.ExamTitle,
		&info.ExamDescription,
		&info.Price,
		&info.CreatedAt,
		&info.ExamDate,
		&info.Duration,
		&info.CreatedBy,
		&info.IsPublic,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			examsInfoMap.Add(examId, valueExamNotFound)
			return nil, ErrExamNotFound
		}

		return nil, err
	}

	examsInfoMap.Add(info.ExamId, info)
	return info, nil
}

// SearchExam searches for exams in the database.
func SearchExam(data *SearchExamsData) (*SearchExamResult, error) {
	publicWhere := ""
	if data.PublicOnly {
		publicWhere = " AND is_public = true "
	}
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT exam_id, 
			course_id, 
			exam_title, 
			exam_description, 
			price, 
			created_at, 
			exam_date, 
			duration, 
			created_by, 
			is_public
		FROM exam_info
		WHERE exam_title ILIKE '%' || $1 || '%'`+publicWhere+`
		ORDER BY exam_date DESC
		LIMIT $2 OFFSET $3`,
		"%"+data.SearchQuery+"%",
		data.Limit,
		data.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []*SearchedExamInfo
	for rows.Next() {
		info := &SearchedExamInfo{}
		err = rows.Scan(
			&info.ExamId,
			&info.CourseId,
			&info.ExamTitle,
			&info.ExamDescription,
			&info.Price,
			&info.CreatedAt,
			&info.ExamDate,
			&info.Duration,
			&info.CreatedBy,
			&info.IsPublic,
		)
		if err != nil {
			return nil, err
		}

		exams = append(exams, info)
	}

	return &SearchExamResult{
		Exams: exams,
	}, nil
}

// EditExamInfo edits the information of an exam.
func EditExamInfo(data *EditExamInfoData) (*ExamInfo, error) {
	info, err := GetExamInfo(data.ExamId)
	if err != nil {
		return nil, err
	} else if info == nil {
		return nil, ErrExamNotFound
	}

	data.ExamTitle = strings.TrimSpace(data.ExamTitle)
	data.ExamDescription = strings.TrimSpace(data.ExamDescription)

	info.ExamTitle = data.ExamTitle
	info.ExamDescription = data.ExamDescription
	info.Price = data.Price
	info.IsPublic = data.IsPublic
	info.Duration = data.Duration
	info.ExamDate = data.ExamDate

	_, err = DefaultContainer.db.Exec(context.Background(),
		`UPDATE exam_info SET
			exam_title = $1,
			exam_description = $2,
			price = $3,
			is_public = $4,
			duration = $5,
			exam_date = $6
		WHERE exam_id = $7`,
		info.ExamTitle,
		info.ExamDescription,
		info.Price,
		info.IsPublic,
		info.Duration,
		info.ExamDate,
		info.ExamId,
	)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// GetExamInfoOrNil gets the exam info or nil if not found.
func GetExamInfoOrNil(examId int) *ExamInfo {
	info, err := GetExamInfo(examId)
	if err != nil && err != ErrExamNotFound {
		logging.UnexpectedError("GetExamInfoOrNil: failed to get exam info:", err)
		return nil
	}

	return info
}

// HasExamStarted returns true if the exam has started.
func HasExamStarted(examId int) (bool, error) {
	var hasStarted bool
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT has_exam_started($1)`,
		examId,
	).Scan(&hasStarted)
	if err != nil {
		return false, err
	}

	return hasStarted, nil
}

// HasExamFinished returns true if the exam has finished.
func HasExamFinished(examId int) (bool, error) {
	var hasFinished bool
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT has_exam_finished($1)`,
		examId,
	).Scan(&hasFinished)
	if err != nil {
		return false, err
	}

	return hasFinished, nil
}

// GetExamStartsIn returns the time in minutes until the exam starts.
func GetExamStartsIn(examId int) (int, error) {
	var startsIn int
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT get_exam_starts_in($1)`,
		examId,
	).Scan(&startsIn)
	if err != nil {
		return 0, err
	}

	return startsIn, nil
}

// GetExamFinishesIn returns the time in minutes until the exam finishes.
func GetExamFinishesIn(examId int) (int, error) {
	var finishesIn int
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT get_exam_finishes_in($1)`,
		examId,
	).Scan(&finishesIn)
	if err != nil {
		return 0, err
	}

	return finishesIn, nil
}

// GetExamQuestionsCount returns the count of questions in the exam.
func GetExamQuestionsCount(examId int) int {
	var count int
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM exam_question WHERE exam_id = $1`,
		examId,
	).Scan(&count)
	if err != nil && err != pgx.ErrNoRows {
		logging.UnexpectedError("GetExamQuestionsCount: failed to query database:", err)
		return 0
	}

	return count
}

// CreateNewExamQuestion creates a new exam question in the database,
// using the plpgsql function create_exam_question.
func CreateNewExamQuestion(data *NewExamQuestionData) (*ExamQuestion, error) {
	examInfo, err := GetExamInfo(data.ExamId)
	if err != nil {
		return nil, err
	} else if examInfo == nil {
		return nil, ErrExamNotFound
	}

	info := &ExamQuestion{
		ExamId:        data.ExamId,
		QuestionTitle: data.QuestionTitle,
		Description:   data.Description,
		Option1:       data.Option1,
		Option2:       data.Option2,
		Option3:       data.Option3,
		Option4:       data.Option4,
	}

	err = DefaultContainer.db.QueryRow(context.Background(),
		`SELECT create_exam_question(
			p_exam_id := $1,
			p_question_title := $2,
			p_description := $3,
			p_option1 := $4,
			p_option2 := $5,
			p_option3 := $6,
			p_option4 := $7
		)`,
		info.ExamId,
		info.QuestionTitle,
		info.Description,
		info.Option1,
		info.Option2,
		info.Option3,
		info.Option4,
	).Scan(&info.QuestionId)
	if err != nil {
		return nil, err
	}

	examInfo.AddQuestion(info)
	examQuestionsMap.Add(info.QuestionId, info)

	return info, nil
}

// EditExamQuestion edits an exam question in the database.
func EditExamQuestion(data *EditExamQuestionData) (*ExamQuestion, error) {
	examInfo := GetExamInfoOrNil(data.ExamId)
	if examInfo == nil {
		return nil, ErrExamNotFound
	}

	info, err := GetExamQuestion(data.ExamId, data.QuestionId)
	if err != nil {
		return nil, err
	} else if info == nil {
		return nil, ErrExamQuestionNotFound
	}

	info.QuestionTitle = data.QuestionTitle
	info.Description = data.Description
	info.Option1 = data.Option1
	info.Option2 = data.Option2
	info.Option3 = data.Option3
	info.Option4 = data.Option4

	_, err = DefaultContainer.db.Exec(context.Background(),
		`UPDATE exam_question SET
			question_title = $1,
			description = $2,
			option1 = $3,
			option2 = $4,
			option3 = $5,
			option4 = $6
		WHERE question_id = $7`,
		info.QuestionTitle,
		info.Description,
		info.Option1,
		info.Option2,
		info.Option3,
		info.Option4,
		info.QuestionId,
	)
	if err != nil {
		return nil, err
	}

	// if the question is cached, update the cached question
	if len(examInfo.GetQuestions()) > 0 {
		examInfo.AddQuestion(info)
	}

	return info, nil
}

// GetExamQuestion gets an exam question from the database.
func GetExamQuestion(examId, questionId int) (*ExamQuestion, error) {
	examInfo, err := GetExamInfo(examId)
	if err != nil {
		return nil, err
	} else if examInfo == nil {
		return nil, ErrExamNotFound
	}

	info := examQuestionsMap.Get(questionId)
	if info != nil && info != valueExamQuestionNotFound && info.QuestionId == questionId {
		return info, nil
	}

	info = &ExamQuestion{
		QuestionId: questionId,
		ExamId:     examId,
	}
	err = DefaultContainer.db.QueryRow(context.Background(),
		`SELECT question_id, 
			exam_id, 
			question_title, 
			description, 
			option1, 
			option2, 
			option3, 
			option4, 
			created_at
		FROM exam_question WHERE question_id = $1`,
		questionId,
	).Scan(
		&info.QuestionId,
		&info.ExamId,
		&info.QuestionTitle,
		&info.Description,
		&info.Option1,
		&info.Option2,
		&info.Option3,
		&info.Option4,
		&info.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			examQuestionsMap.Add(questionId, valueExamQuestionNotFound)
			return nil, ErrExamQuestionNotFound
		}

		return nil, err
	}

	// AddQuestion is safe to be used multiple times even if the question is
	// already added, because it has a check to prevent duplicates.
	examInfo.AddQuestion(info)
	examQuestionsMap.Add(info.QuestionId, info)
	return info, nil
}

// GetExamQuestions gets all questions of an exam from the database.
func GetExamQuestions(data *GetExamQuestionsData) ([]*ExamQuestion, error) {
	examInfo, err := GetExamInfo(data.ExamId)
	if err != nil {
		return nil, err
	} else if examInfo == nil {
		return nil, ErrExamNotFound
	}

	addedQuestions := examInfo.GetQuestions()
	if len(addedQuestions) > 0 {
		// just use the cached questions
		return addedQuestions, nil
	}

	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT question_id, 
			exam_id, 
			question_title, 
			description, 
			option1, 
			option2, 
			option3, 
			option4, 
			created_at
		FROM exam_question WHERE exam_id = $1
		ORDER BY question_id
		LIMIT $2 OFFSET $3`,
		data.ExamId,
		data.Limit,
		data.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*ExamQuestion
	for rows.Next() {
		info := &ExamQuestion{}
		err = rows.Scan(
			&info.QuestionId,
			&info.ExamId,
			&info.QuestionTitle,
			&info.Description,
			&info.Option1,
			&info.Option2,
			&info.Option3,
			&info.Option4,
			&info.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		examInfo.AddQuestion(info)
		examQuestionsMap.Add(info.QuestionId, info)
		questions = append(questions, info)
	}

	return questions, nil
}

// HasParticipatedInExam returns true if the user has participated in the exam.
// It uses the plpgsql function has_participated_in_exam.
func HasParticipatedInExam(userId string, examId int) bool {
	var hasParticipated bool
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT has_participated_in_exam($1, $2)`,
		examId,
		userId,
	).Scan(&hasParticipated)
	if err != nil {
		logging.UnexpectedError("HasParticipatedInExam: failed to query database:", err)
		return false
	}

	return hasParticipated
}

// CanParticipateInExam returns true if the user can participate in the exam.
func CanParticipateInExam(userId string, examId int) (bool, error) {
	var canParticipate bool
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT can_participate_in_exam($1, $2)`,
		examId,
		userId,
	).Scan(&canParticipate)
	if err != nil {
		return false, err
	}

	return canParticipate, nil
}

// CanParticipateInExamOrFalse returns true if the user can participate in the exam.
// It will also returns false if there is an error.
func CanParticipateInExamOrFalse(userId string, examId int) bool {
	canParticipate, err := CanParticipateInExam(userId, examId)
	if err != nil {
		logging.UnexpectedError("CanParticipateInExamOrFalse: failed to check participation:", err)
		return false
	}

	return canParticipate
}

// GetGivenExam gets the information of a given exam.
func GetGivenExam(userId string, examId int) (*GivenExam, error) {
	uniqueId := userId + KeySepChar + ssg.ToBase10(examId)
	info := givenExamsMap.Get(uniqueId)
	if info != nil && info != valueGivenExamNotFound &&
		info.ExamId == examId && info.UserId == userId {
		return info, nil
	}

	info = &GivenExam{
		UserId: userId,
		ExamId: examId,
	}
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT user_id,
			exam_id, 
			price, 
			added_by, 
			scored_by, 
			created_at, 
			final_score
		FROM given_exam WHERE user_id = $1 AND exam_id = $2`,
		userId,
		examId,
	).Scan(
		&info.UserId,
		&info.ExamId,
		&info.Price,
		&info.AddedBy,
		&info.ScoredBy,
		&info.CreatedAt,
		&info.FinalScore,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			givenExamsMap.Add(uniqueId, valueGivenExamNotFound)
			return nil, ErrGivenExamNotFound
		}

		return nil, err
	}

	givenExamsMap.Add(uniqueId, info)
	return info, nil
}

// AddUserInExam adds a user to an exam.
func AddUserInExam(data *NewGivenExamData) (*GivenExam, error) {
	if data.Price == "" {
		data.Price = DefaultExamPrice
	}

	uniqueId := data.UserId + KeySepChar + ssg.ToBase10(data.ExamId)
	info := givenExamsMap.Get(uniqueId)
	if info != nil && info != valueGivenExamNotFound &&
		info.ExamId == data.ExamId && info.UserId == data.UserId {
		return info, nil
	}

	info = &GivenExam{
		UserId:    data.UserId,
		ExamId:    data.ExamId,
		Price:     data.Price,
		AddedBy:   data.AddedBy,
		CreatedAt: time.Now(),
	}

	// 	-- Example usage:
	// --    CALL add_user_in_exam(
	// --        p_user_id := 'user123',
	// --        p_exam_id := 1001,
	// --        p_price := '0T',
	// --        p_added_by := 'admin'
	// --    );
	_, err := DefaultContainer.db.Exec(context.Background(),
		`CALL add_user_in_exam(
			p_user_id := $1,
			p_exam_id := $2,
			p_price := $3,
			p_added_by := $4
		)`,
		info.UserId,
		info.ExamId,
		info.Price,
		info.AddedBy,
	)
	if err != nil {
		return nil, err
	}

	givenExamsMap.Add(uniqueId, info)
	return info, nil
}

// SetScoreForUserInExam sets the final score for a user in an exam.
// It uses the sp set_score_for_user_in_exam.
func SetScoreForUserInExam(data *NewScoreData) (*GivenExam, error) {
	info, err := GetGivenExam(data.UserId, data.ExamId)
	if err != nil {
		return nil, err
	} else if info == nil {
		return nil, ErrGivenExamNotFound
	}

	info.FinalScore = ssg.Clone(&data.FinalScore)
	info.ScoredBy = ssg.Clone(&data.ScoredBy)

	_, err = DefaultContainer.db.Exec(context.Background(),
		`CALL set_score_for_user_in_exam(
			p_exam_id := $1,
			p_user_id := $2,
			p_final_score := $3,
			p_scored_by := $4
		)`,
		info.ExamId,
		info.UserId,
		info.FinalScore,
		info.ScoredBy,
	)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// GetMostRecentExams returns the most recent exams.
// It uses this sql command (just an example):
// --   SELECT * FROM most_recent_exams_view LIMIT 10 OFFSET 0;
func GetMostRecentExams(data *GetMostRecentExamsData) ([]*MostRecentExamInfo, error) {
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT exam_id,
			course_id,
			exam_title,
			exam_description,
			price,
			created_at,
			exam_date,
			duration,
			created_by,
			is_public
		FROM most_recent_exams_view LIMIT $1 OFFSET $2`,
		data.Limit,
		data.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []*MostRecentExamInfo
	for rows.Next() {
		info := &MostRecentExamInfo{}
		err = rows.Scan(
			&info.ExamId,
			&info.CourseId,
			&info.ExamTitle,
			&info.ExamDescription,
			&info.Price,
			&info.CreatedAt,
			&info.ExamDate,
			&info.Duration,
			&info.CreatedBy,
			&info.IsPublic,
		)
		if err != nil {
			return nil, err
		}

		exams = append(exams, info)
	}

	return exams, nil
}

// GetGivenAnswer gets the given answer of a user for a question in an exam.
func GetGivenAnswer(data *GetGivenAnswerData) (*GivenAnswerInfo, error) {
	uniqueId := ssg.ToBase10(data.ExamId) + KeySepChar +
		ssg.ToBase10(data.QuestionId) + KeySepChar +
		data.UserId
	info := givenAnswersMap.Get(uniqueId)
	if info != nil && info != valueGivenAnswerNotFound &&
		info.ExamId == data.ExamId &&
		info.QuestionId == data.QuestionId &&
		info.AnsweredBy == data.UserId {
		return info, nil
	} else if info == valueGivenAnswerNotFound {
		return nil, ErrGivenAnswerNotFound
	}

	info = &GivenAnswerInfo{}
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT exam_id, 
			question_id, 
			answered_by, 
			chosen_option,
			seconds_taken,
			answer_text,
			answered_at
		FROM given_answer WHERE exam_id = $1 AND question_id = $2 AND answered_by = $3`,
		data.ExamId,
		data.QuestionId,
		data.UserId,
	).Scan(
		&info.ExamId,
		&info.QuestionId,
		&info.AnsweredBy,
		&info.ChosenOption,
		&info.SecondsTaken,
		&info.AnswerText,
		&info.AnsweredAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			givenAnswersMap.Add(uniqueId, valueGivenAnswerNotFound)
			return nil, ErrGivenAnswerNotFound
		}

		return nil, err
	}

	givenAnswersMap.Add(uniqueId, info)
	return info, nil
}

// GetGivenAnswerOrNil gets the given answer or nil if not found.
// It will also log the error if the error is something unexpected.
func GetGivenAnswerOrNil(data *GetGivenAnswerData) *GivenAnswerInfo {
	info, err := GetGivenAnswer(data)
	if err != nil && err != ErrGivenAnswerNotFound {
		logging.UnexpectedError("GetGivenAnswerOrNil: failed to get given answer:", err)
		return nil
	}

	return info
}

// AnswerQuestion answers a question in an exam.
// It uses the plpgsql function give_answer_to_exam_question.
func AnswerQuestion(data *AnswerQuestionData) (*GivenAnswerInfo, error) {
	if data.ChosenOption == nil && data.AnswerText == nil {
		return nil, ErrInvalidAnswer
	}

	uniqueId := ssg.ToBase10(data.ExamId) + KeySepChar +
		ssg.ToBase10(data.QuestionId) + KeySepChar +
		data.AnsweredBy
	info := givenAnswersMap.Get(uniqueId)
	if info == nil {
		info = &GivenAnswerInfo{
			ExamId:     data.ExamId,
			QuestionId: data.QuestionId,
			AnsweredBy: data.AnsweredBy,
		}
	}

	info.ChosenOption = ssg.Clone(data.ChosenOption)
	info.SecondsTaken = data.SecondsTaken
	info.AnswerText = ssg.Clone(data.AnswerText)
	info.AnsweredAt = time.Now()

	_, err := DefaultContainer.db.Exec(context.Background(),
		`SELECT give_answer_to_exam_question(
			p_exam_id := $1,
			p_question_id := $2,
			p_answered_by := $3,
			p_chosen_option := $4,
			p_seconds_taken := $5,
			p_answer_text := $6
		)`,
		info.ExamId,
		info.QuestionId,
		info.AnsweredBy,
		info.ChosenOption,
		info.SecondsTaken,
		info.AnswerText,
	)
	if err != nil {
		logging.UnexpectedError("AnswerQuestion: failed to answer question:", err)
		return nil, err
	}

	givenAnswersMap.Add(uniqueId, info)
	return info, nil
}

// GetUserOngoingExams gets the ongoing exams of a user.
func GetUserOngoingExams(userId string) ([]*UserOngoingExamInfo, error) {
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT exam_id, exam_title, exam_date
		FROM user_ongoing_exams WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrExamNotFound
		}

		return nil, err
	}
	defer rows.Close()

	var exams []*UserOngoingExamInfo
	for rows.Next() {
		info := &UserOngoingExamInfo{}
		err = rows.Scan(
			&info.ExamId,
			&info.ExamTitle,
			&info.StartTime,
		)
		if err != nil {
			return nil, err
		}

		exams = append(exams, info)
	}

	return exams, nil
}

// GetUserOngoingExamsOrNil gets the ongoing exams of a user or nil if not found.
func GetUserOngoingExamsOrNil(userId string) []*UserOngoingExamInfo {
	exams, err := GetUserOngoingExams(userId)
	if err != nil && err != pgx.ErrNoRows && err != ErrExamNotFound {
		logging.UnexpectedError("GetUserOngoingExamsOrNil: failed to get ongoing exams:", err)
		return nil
	}

	return exams
}

// GetUserExamsHistory gets the past exams of a user.
func GetUserExamsHistory(opts *GetUserExamsHistoryOptions) ([]*UserPastExamInfo, error) {
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT exam_id, exam_title, exam_date
		FROM user_exams_history WHERE user_id = $1 LIMIT $2 OFFSET $3`,
		opts.UserId, opts.Limit, opts.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []*UserPastExamInfo
	for rows.Next() {
		info := &UserPastExamInfo{}
		err = rows.Scan(
			&info.ExamId,
			&info.ExamTitle,
			&info.StartedAt,
		)
		if err != nil {
			return nil, err
		}

		exams = append(exams, info)
	}

	return exams, nil
}

// GetUserExamsHistoryOrNil gets the past exams of a user or nil if not found.
func GetUserExamsHistoryOrNil(opts *GetUserExamsHistoryOptions) []*UserPastExamInfo {
	exams, err := GetUserExamsHistory(opts)
	if err != nil && err != pgx.ErrNoRows {
		logging.UnexpectedError("GetUserOngoingExamsOrNil: failed to get ongoing exams:", err)
		return nil
	}

	return exams
}

// GetExamParticipants gets all the participants of an exam.
func GetExamParticipants(opts *GetExamParticipantsOptions) ([]*GivenExam, error) {
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT user_id, 
			exam_id, 
			price, 
			added_by, 
			scored_by, 
			created_at, 
			final_score
		FROM given_exam WHERE exam_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		opts.ExamId,
		opts.Limit,
		opts.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []*GivenExam
	for rows.Next() {
		info := &GivenExam{}
		err = rows.Scan(
			&info.UserId,
			&info.ExamId,
			&info.Price,
			&info.AddedBy,
			&info.ScoredBy,
			&info.CreatedAt,
			&info.FinalScore,
		)
		if err != nil {
			return nil, err
		}

		exams = append(exams, info)
	}

	return exams, nil
}

// GetExamParticipantsOrNil gets the participants of an exam or nil if not found.
func GetExamParticipantsOrNil(opts *GetExamParticipantsOptions) []*GivenExam {
	exams, err := GetExamParticipants(opts)
	if err != nil && err != pgx.ErrNoRows {
		logging.UnexpectedError("GetExamParticipantsOrNil: failed to get exam participants:", err)
		return nil
	}

	return exams
}

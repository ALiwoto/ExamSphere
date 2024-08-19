package database

import (
	"context"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/jackc/pgx/v5"
)

// CreateNewExam creates a new exam in the database,
func CreateNewExam(data *NewExamData) (*ExamInfo, error) {
	if data.Price == "" {
		data.Price = DefaultExamPrice
	}
	info := &ExamInfo{
		CourseId:  data.CourseId,
		Price:     data.Price,
		CreatedBy: data.CreatedBy,
		IsPublic:  data.IsPublic,
		Duration:  data.Duration,
	}

	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT create_exam_info(
			p_course_id := $1,
			p_price := $2,
			p_created_by := $3,
			p_is_public := $4,
			p_duration := $5,
			p_exam_date := $6
		)`,
		info.CourseId,
		info.Price,
		info.CreatedBy,
		info.IsPublic,
		info.Duration,
		data.ExamDate,
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

	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT exam_id, course_id, price, created_at, exam_date, duration, created_by, is_public
		FROM exam_info WHERE exam_id = $1`,
		examId,
	).Scan(
		&info.ExamId,
		&info.CourseId,
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

	err = DefaultContainer.db.QueryRow(context.Background(),
		`SELECT question_id, exam_id, question_title, description, option1, option2, option3, option4, created_at
		FROM exam_questions WHERE question_id = $1`,
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
func GetExamQuestions(examId int) ([]*ExamQuestion, error) {
	examInfo, err := GetExamInfo(examId)
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
		`SELECT question_id, exam_id, question_title, description, option1, option2, option3, option4, created_at
		FROM exam_questions WHERE exam_id = $1`,
		examId,
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
func HasParticipatedInExam(examId int, userId string) (bool, error) {
	var hasParticipated bool
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT has_participated_in_exam($1, $2)`,
		examId,
		userId,
	).Scan(&hasParticipated)
	if err != nil {
		return false, err
	}

	return hasParticipated, nil
}

// CanParticipateInExam returns true if the user can participate in the exam.
func CanParticipateInExam(examId int, userId string) (bool, error) {
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

// GetGivenExam gets the information of a given exam.
func GetGivenExam(examId int, userId string) (*GivenExam, error) {
	uniqueId := userId + KeySepChar + ssg.ToBase10(examId)
	info := givenExamsMap.Get(uniqueId)
	if info != nil && info != valueGivenExamNotFound &&
		info.ExamId == examId && info.UserId == userId {
		return info, nil
	}

	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT user_id, exam_id, price, added_by, scored_by, created_at, final_score
		FROM given_exams WHERE user_id = $1 AND exam_id = $2`,
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
	info := &GivenExam{
		UserId:  data.UserId,
		ExamId:  data.ExamId,
		Price:   data.Price,
		AddedBy: data.AddedBy,
	}

	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT add_user_in_exam(
			p_user_id := $1,
			p_exam_id := $2,
			p_price := $3,
			p_added_by := $4
		)`,
		info.UserId,
		info.ExamId,
		info.Price,
		info.AddedBy,
	).Scan(&info.CreatedAt)
	if err != nil {
		return nil, err
	}

	givenExamsMap.Add(uniqueId, info)
	return info, nil
}

// SetScoreForUserInExam sets the final score for a user in an exam.
// It uses the plpgsql function set_score_for_user_in_exam.
func SetScoreForUserInExam(data *NewScoreData) (*GivenExam, error) {
	info, err := GetGivenExam(data.ExamId, data.UserId)
	if err != nil {
		return nil, err
	} else if info == nil {
		return nil, ErrGivenExamNotFound
	}

	info.FinalScore = ssg.Clone(&data.FinalScore)
	info.ScoredBy = ssg.Clone(&data.ScoredBy)

	_, err = DefaultContainer.db.Exec(context.Background(),
		`SELECT set_score_for_user_in_exam(
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

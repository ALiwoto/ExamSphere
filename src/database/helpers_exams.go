package database

import (
	"ExamSphere/src/core/utils/logging"
	"context"
	"strings"
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
		CourseId:            data.CourseId,
		ExamTitle:           data.ExamTitle,
		ExamDescription:     data.ExamDescription,
		Price:               data.Price,
		CreatedBy:           data.CreatedBy,
		IsPublic:            data.IsPublic,
		Duration:            data.Duration,
		ExamDate:            data.ExamDate,
		CreatedAt:           time.Now(),
		IsStrict:            data.IsStrict,
		IsSampleExam:        data.IsSampleExam,
		MaxQuestionsSeconds: data.MaxQuestionsSeconds,
		NeedsVideoCall:      data.NeedsVideoCall,
		NeedsVoiceCall:      data.NeedsVoiceCall,
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
			p_exam_date := $8,
			p_is_strict := $9,
			p_is_sample_exam := $10,
			p_max_questions_seconds := $11,
			p_needs_video_call := $12,
			p_needs_voice_call := $13
		)`,
		info.CourseId,                        // 1
		info.ExamTitle,                       // 2
		info.ExamDescription,                 // 3
		info.Price,                           // 4
		info.CreatedBy,                       // 5
		info.IsPublic,                        // 6
		info.Duration,                        // 7
		info.ExamDate.Format(ExamDateLayout), // 8
		info.IsStrict,                        // 9
		info.IsSampleExam,                    // 10
		info.MaxQuestionsSeconds,             // 11
		info.NeedsVideoCall,                  // 12
		info.NeedsVoiceCall,                  // 13
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
			is_public,
			is_strict,
			is_sample_exam,
			max_questions_seconds,
			needs_video_call,
			needs_voice_call
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
		&info.IsStrict,
		&info.IsSampleExam,
		&info.MaxQuestionsSeconds,
		&info.NeedsVideoCall,
		&info.NeedsVoiceCall,
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
	extraWhereClause := ""
	if data.PublicOnly {
		extraWhereClause = " AND is_public = TRUE "
	}
	if data.SampleExams {
		extraWhereClause = " AND is_sample_exam = TRUE "
	} else {
		extraWhereClause = " AND is_sample_exam = FALSE "
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
			is_public,
			is_strict,
			is_sample_exam,
			max_questions_seconds,
			needs_video_call,
			needs_voice_call
		FROM exam_info
		WHERE exam_title ILIKE '%' || $1 || '%'`+extraWhereClause+`
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
			&info.IsStrict,
			&info.IsSampleExam,
			&info.MaxQuestionsSeconds,
			&info.NeedsVideoCall,
			&info.NeedsVoiceCall,
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
	info.IsStrict = data.IsStrict
	info.NeedsVideoCall = data.NeedsVideoCall
	info.NeedsVoiceCall = data.NeedsVoiceCall

	_, err = DefaultContainer.db.Exec(context.Background(),
		`UPDATE exam_info SET
			exam_title = $1,
			exam_description = $2,
			price = $3,
			is_public = $4,
			duration = $5,
			exam_date = $6
			is_strict = $7,
			needs_video_call = $8,
			needs_voice_call = $9
		WHERE exam_id = $10`,
		info.ExamTitle,       // 1
		info.ExamDescription, // 2
		info.Price,           // 3
		info.IsPublic,        // 4
		info.Duration,        // 5
		info.ExamDate,        // 6
		info.IsStrict,        // 7
		info.NeedsVideoCall,  // 8
		info.NeedsVoiceCall,  // 9
		info.ExamId,          // 10
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
		`SELECT COALESCE(SUM(
            CASE 
                WHEN is_pointer = false THEN 1 
                WHEN is_pointer = true THEN pointer_count 
                ELSE 0 
            END
        ), 0) AS total_count
        FROM exam_question 
        WHERE exam_id = $1`,
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

	// if the current exam itself is a sample exam, it cannot have
	// pointer to other exams.
	if examInfo.IsSampleExam && data.IsPointer {
		return nil, ErrSampleExamCannotHavePointer
	}

	if data.IsPointer {
		if data.PointerToExamId == nil || *data.PointerToExamId == 0 {
			return nil, ErrInvalidPointerToExamId
		}

		pointedToExamInfo, err := GetExamInfo(*data.PointerToExamId)
		if err != nil && err != ErrExamNotFound {
			logging.UnexpectedError("CreateNewExamQuestion: failed to get pointed to exam info:", err)
			return nil, err
		} else if pointedToExamInfo == nil {
			return nil, ErrInvalidPointerToExamId
		}

		// the exam that we are pointing to, must be a sample exam
		// TODO: in the future, add accessibility check here
		if !pointedToExamInfo.IsSampleExam {
			return nil, ErrPointedExamMustBeSample
		}
	}

	info := &ExamQuestion{
		ExamId:          data.ExamId,
		QuestionTitle:   data.QuestionTitle,
		Description:     data.Description,
		Option1:         data.Option1,
		Option2:         data.Option2,
		Option3:         data.Option3,
		Option4:         data.Option4,
		IsPointer:       data.IsPointer,
		PointerCount:    data.PointerCount,
		PointerToExamId: ssg.Clone(data.PointerToExamId),
	}

	err = DefaultContainer.db.QueryRow(context.Background(),
		`SELECT create_exam_question(
			p_exam_id := $1,
			p_question_title := $2,
			p_description := $3,
			p_option1 := $4,
			p_option2 := $5,
			p_option3 := $6,
			p_option4 := $7,
			p_is_pointer := $8,
			p_pointer_count := $9,
			p_pointer_to_exam_id := $10
		)`,
		info.ExamId,          // 1
		info.QuestionTitle,   // 2
		info.Description,     // 3
		info.Option1,         // 4
		info.Option2,         // 5
		info.Option3,         // 6
		info.Option4,         // 7
		info.IsPointer,       // 8
		info.PointerCount,    // 9
		info.PointerToExamId, // 10
	).Scan(&info.QuestionId)
	if err != nil {
		return nil, err
	}

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

	if !info.IsPointer {
		if data.PointerCount > 0 {
			return nil, ErrNonPointerQuestionCannotHavePointerCount
		} else if data.PointerToExamId != nil {
			if *data.PointerToExamId == 0 {
				data.PointerToExamId = nil
			} else {
				return nil, ErrNonPointerQuestionCannotHavePointedToExam
			}
		}
	} else {
		// the question is a pointer question, it must have a pointer to exam
		if data.PointerToExamId == nil || *data.PointerToExamId == 0 {
			return nil, ErrInvalidPointerToExamId
		}

		if data.PointerCount <= 0 {
			return nil, ErrInvalidPointerCount
		}

		pointedToExamInfo, err := GetExamInfo(*data.PointerToExamId)
		if err != nil && err != ErrExamNotFound {
			logging.UnexpectedError("CreateNewExamQuestion: failed to get pointed to exam info:", err)
			return nil, err
		} else if pointedToExamInfo == nil {
			return nil, ErrInvalidPointerToExamId
		}

		// the exam that we are pointing to, must be a sample exam
		// TODO: in the future, add accessibility check here
		if !pointedToExamInfo.IsSampleExam {
			return nil, ErrPointedExamMustBeSample
		}
	}

	info.QuestionTitle = data.QuestionTitle
	info.Description = data.Description
	info.Option1 = data.Option1
	info.Option2 = data.Option2
	info.Option3 = data.Option3
	info.Option4 = data.Option4
	info.PointerCount = data.PointerCount
	info.PointerToExamId = ssg.Clone(data.PointerToExamId)

	_, err = DefaultContainer.db.Exec(context.Background(),
		`UPDATE exam_question SET
			question_title = $1,
			description = $2,
			option1 = $3,
			option2 = $4,
			option3 = $5,
			option4 = $6
			pointer_count = $7,
			pointer_to_exam_id = $8
		WHERE question_id = $9`,
		info.QuestionTitle,   // 1
		info.Description,     // 2
		info.Option1,         // 3
		info.Option2,         // 4
		info.Option3,         // 5
		info.Option4,         // 6
		info.PointerCount,    // 7
		info.PointerToExamId, // 8
		info.QuestionId,      // 9
	)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// GetExamQuestion gets an exam question from the database.
// NOTE: The exam id you should be passing here should be the current
// exam's id; if the question is actually a referenced question from a
// sample exam, you should NOT pass exam id of that sample exam.
func GetExamQuestion(examId, questionId int) (*ExamQuestion, error) {
	// this exam info shall not belong to the current question's exam if and
	// only if the question's origin exam is a sample exam; otherwise we will
	// be having a logical error here (and we should return an error).
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
			created_at,
			is_pointer,
			pointer_count,
			pointer_to_exam_id
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
		&info.IsPointer,
		&info.PointerCount,
		&info.PointerToExamId,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			examQuestionsMap.Add(questionId, valueExamQuestionNotFound)
			return nil, ErrExamQuestionNotFound
		}

		return nil, err
	}

	if info.IsPointer {
		// then the exam id should match the current exam id
		if info.ExamId != examId {
			logging.UnexpectedError("GetExamQuestion: question is a pointer, but exam id does not match")
			return nil, ErrExamQuestionNotFound
		}
	} else {
		// then the origin exam, should be a sample exam if the exam id does
		// not match
		if info.ExamId != examId {
			originExamInfo, err := GetExamInfo(info.ExamId)
			if err != nil && err != ErrExamNotFound {
				return nil, err
			} else if originExamInfo == nil {
				logging.UnexpectedError("GetExamQuestion: failed to get origin exam info")
				return nil, ErrExamQuestionNotFound
			}

			if !originExamInfo.IsSampleExam {
				logging.UnexpectedError("GetExamQuestion: question is not a pointer, but origin exam is not a sample exam")
				return nil, ErrExamQuestionNotFound
			}
		}
	}

	examQuestionsMap.Add(info.QuestionId, info)
	return info, nil
}

// MarkGivenAnswersAsSeen marks the given answers as seen.
// the way to run the pg sql function is:
// -- SELECT mark_given_answers_as_seen(
// --     p_exam_id := 1,
// --     p_question_ids := ARRAY[1, 2, 3, 4, 5],
// --     p_answered_by := 'user123'
// -- );
func MarkGivenAnswersAsSeen(data *MarkGivenAnswersAsSeenData) error {
	questionIds := "{"
	for i, id := range data.QuestionIds {
		questionIds += ssg.ToBase10(id)
		if i < len(data.QuestionIds)-1 {
			questionIds += ","
		}
	}
	questionIds += "}"

	_, err := DefaultContainer.db.Exec(context.Background(),
		`SELECT mark_given_answers_as_seen(
			p_exam_id := $1,
			p_question_ids := $2::integer[],
			p_answered_by := $3
		)`,
		data.ExamId,
		questionIds,
		data.AnsweredBy,
	)
	return err
}

// GetExamQuestions gets all questions of an exam from the database.
func GetExamQuestions(data *GetExamQuestionsData) ([]*ExamQuestion, error) {
	if data.UserId == "" {
		return nil, ErrInvalidParticipantUserId
	}

	examInfo, err := GetExamInfo(data.ExamId)
	if err != nil {
		return nil, err
	} else if examInfo == nil {
		return nil, ErrExamNotFound
	}

	// addedQuestions := examInfo.GetQuestions()
	// if len(addedQuestions) > 0 {
	// 	// just use the cached questions
	// 	return addedQuestions, nil
	// }

	var rows pgx.Rows
	if data.ResolvePointers {
		rows, err = DefaultContainer.db.Query(context.Background(),
			`SELECT question_id, 
				exam_id, 
				question_title, 
				description, 
				option1, 
				option2, 
				option3, 
				option4, 
				created_at,
				is_pointer,
				pointer_count,
				pointer_to_exam_id
			FROM get_exam_questions_for_participant(
				p_exam_id := $1,
				p_participant_id := $2
			)
			ORDER BY question_id
			LIMIT $3 OFFSET $4`,
			data.ExamId, // 1
			data.UserId, // 2
			data.Limit,  // 3
			data.Offset, // 4
		)
	} else {
		// the raw results (pointers are NOT resolved)
		rows, err = DefaultContainer.db.Query(context.Background(),
			`SELECT question_id, 
				exam_id, 
				question_title, 
				description, 
				option1, 
				option2, 
				option3, 
				option4, 
				created_at,
				is_pointer,
				pointer_count,
				pointer_to_exam_id
			FROM exam_question WHERE exam_id = $1
			ORDER BY question_id
			LIMIT $2 OFFSET $3`,
			data.ExamId, // 1
			data.Limit,  // 2
			data.Offset, // 3
		)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*ExamQuestion
	var questionIds []int
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
			&info.IsPointer,
			&info.PointerCount,
			&info.PointerToExamId,
		)
		if err != nil {
			return nil, err
		}

		if !examQuestionsMap.Exists(info.QuestionId) {
			examQuestionsMap.Add(info.QuestionId, info)
		}
		questions = append(questions, info)
		questionIds = append(questionIds, info.QuestionId)
	}

	if data.MarkAsSeen {
		// resolving the pointers mean actually getting the answers
		// for someone who is participating in the exam.
		err = MarkGivenAnswersAsSeen(&MarkGivenAnswersAsSeenData{
			ExamId:      data.ExamId,
			AnsweredBy:  data.UserId,
			QuestionIds: questionIds,
		})
		if err != nil {
			logging.UnexpectedError("GetExamQuestions: failed to mark answers as seen:", err)
			return nil, err
		}
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
			answered_at,
			seen_at
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
		&info.SeenAt,
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
	if info == nil || info == valueGivenAnswerNotFound || info.ExamId != data.ExamId {
		info = &GivenAnswerInfo{
			ExamId:     data.ExamId,
			QuestionId: data.QuestionId,
			AnsweredBy: data.AnsweredBy,
		}
	}

	examInfo := GetExamInfoOrNil(data.ExamId)
	if examInfo == nil {
		return nil, ErrExamNotFound
	}

	if examInfo.IsStrict {
		if examInfo.MaxQuestionsSeconds > 0 &&
			info.SeenAt != nil && !(*info.SeenAt).IsZero() {
			if time.Since(*info.SeenAt).Seconds() > float64(examInfo.MaxQuestionsSeconds) {
				return nil, ErrExamQuestionTimeExceeded
			}
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

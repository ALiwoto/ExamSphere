package database

import (
	"time"
)

// ExamInfo is a struct that represents the information of an exam.
type ExamInfo struct {
	ExamId              int       `json:"exam_id"`
	CourseId            int       `json:"course_id"`
	ExamTitle           string    `json:"exam_title"`
	ExamDescription     string    `json:"exam_description"`
	Price               string    `json:"price"`
	CreatedAt           time.Time `json:"created_at"`
	ExamDate            time.Time `json:"exam_date"`
	Duration            int       `json:"duration"`
	CreatedBy           string    `json:"created_by"`
	IsPublic            bool      `json:"is_public"`
	IsStrict            bool      `json:"is_strict"`
	IsSampleExam        bool      `json:"is_sample_exam"`
	MaxQuestionsSeconds int       `json:"max_questions_seconds"`
	NeedsVideoCall      bool      `json:"needs_video_call"`
	NeedsVoiceCall      bool      `json:"needs_voice_call"`
}

// SearchExamsData is a struct that represents the data needed to search for exams.
type SearchExamsData struct {
	SearchQuery string `json:"search_query"`
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
	PublicOnly  bool   `json:"public_only"`
	SampleExams bool   `json:"sample_exams"`
}

type SearchExamResult struct {
	Exams []*SearchedExamInfo `json:"exams"`
}

type SearchedExamInfo struct {
	ExamId              int       `json:"exam_id"`
	CourseId            int       `json:"course_id"`
	ExamTitle           string    `json:"exam_title"`
	ExamDescription     string    `json:"exam_description"`
	Price               string    `json:"price"`
	CreatedAt           time.Time `json:"created_at"`
	ExamDate            time.Time `json:"exam_date"`
	Duration            int       `json:"duration"`
	CreatedBy           string    `json:"created_by"`
	IsPublic            bool      `json:"is_public"`
	IsStrict            bool      `json:"is_strict"`
	IsSampleExam        bool      `json:"is_sample_exam"`
	MaxQuestionsSeconds int       `json:"max_questions_seconds"`
	NeedsVideoCall      bool      `json:"needs_video_call"`
	NeedsVoiceCall      bool      `json:"needs_voice_call"`
}

type MarkGivenAnswersAsSeenData struct {
	ExamId      int    `json:"exam_id"`
	AnsweredBy  string `json:"user_id"`
	QuestionIds []int  `json:"question_ids"`
}

type GetExamQuestionsData struct {
	ExamId int    `json:"exam_id"`
	UserId string `json:"user_id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`

	// ResolvePointers is a boolean value that indicates whether the pointers
	// in the exam questions should be resolved or not.
	// If this is set to false, the results will basically be the raw results which
	// are inside of the exam questions table;
	// otherwise (if the user is participating inside of the exam), the pointers will be
	// actually resolved and the returned exam questions might be different for each
	// user.
	ResolvePointers bool `json:"resolve_pointers"`
	MarkAsSeen      bool `json:"mark_as_seen"`
}

// NewExamData is a struct that represents the data needed to create a new exam.
type NewExamData struct {
	CourseId            int       `json:"course_id"`
	ExamTitle           string    `json:"exam_title"`
	ExamDescription     string    `json:"exam_description"`
	Price               string    `json:"price"`
	CreatedBy           string    `json:"created_by"`
	IsPublic            bool      `json:"is_public"`
	Duration            int       `json:"duration"`
	ExamDate            time.Time `json:"exam_date"`
	IsStrict            bool      `json:"is_strict"`
	IsSampleExam        bool      `json:"is_sample_exam"`
	MaxQuestionsSeconds int       `json:"max_questions_seconds"`
	NeedsVideoCall      bool      `json:"needs_video_call"`
	NeedsVoiceCall      bool      `json:"needs_voice_call"`
}

type EditExamInfoData struct {
	ExamId          int       `json:"exam_id"`
	CourseId        int       `json:"course_id"`
	ExamTitle       string    `json:"exam_title"`
	ExamDescription string    `json:"exam_description"`
	Price           string    `json:"price"`
	IsPublic        bool      `json:"is_public"`
	Duration        int       `json:"duration"`
	ExamDate        time.Time `json:"exam_date"`
	IsStrict        bool      `json:"is_strict"`
	NeedsVideoCall  bool      `json:"needs_video_call"`
	NeedsVoiceCall  bool      `json:"needs_voice_call"`
}

// ExamQuestion is a struct that represents the information of an exam question.
type ExamQuestion struct {
	QuestionId      int       `json:"question_id"`
	ExamId          int       `json:"exam_id"`
	QuestionTitle   string    `json:"question_title"`
	Description     *string   `json:"description"`
	Option1         *string   `json:"option1"`
	Option2         *string   `json:"option2"`
	Option3         *string   `json:"option3"`
	Option4         *string   `json:"option4"`
	CreatedAt       time.Time `json:"created_at"`
	IsPointer       bool      `json:"is_pointer"`
	PointerCount    int       `json:"pointer_count"`
	PointerToExamId *int      `json:"pointer_to_exam_id"`
}

// NewExamQuestionData is a struct that represents the data needed to create a new exam question.
type NewExamQuestionData struct {
	ExamId          int     `json:"exam_id"`
	QuestionTitle   string  `json:"question_title"`
	Description     *string `json:"description"`
	Option1         *string `json:"option1"`
	Option2         *string `json:"option2"`
	Option3         *string `json:"option3"`
	Option4         *string `json:"option4"`
	IsPointer       bool    `json:"is_pointer"`
	PointerCount    int     `json:"pointer_count"`
	PointerToExamId *int    `json:"pointer_to_exam_id"`
}

// EditExamQuestionData is a struct that represents the data needed to edit an exam question.
type EditExamQuestionData struct {
	QuestionId      int     `json:"question_id"`
	ExamId          int     `json:"exam_id"`
	QuestionTitle   string  `json:"question_title"`
	Description     *string `json:"description"`
	Option1         *string `json:"option1"`
	Option2         *string `json:"option2"`
	Option3         *string `json:"option3"`
	Option4         *string `json:"option4"`
	PointerCount    int     `json:"pointer_count"`
	PointerToExamId *int    `json:"pointer_to_exam_id"`
}

// NewScoreData is a struct that represents the data needed to create
// a new score for a user in an exam.
type NewScoreData struct {
	ExamId     int    `json:"exam_id"`
	UserId     string `json:"user_id"`
	FinalScore string `json:"final_score"`
	ScoredBy   string `json:"scored_by"`
}

// GivenExam is a struct that represents the information of an exam
// that a certain user has participated in.
// Please note that when an admin or a teacher forcefully adds a user
// to an exam, a record will be created for the user in this table.
// And if in that case, the user does not participate in the exam, their
// final score can be set to 0 by the admin or teacher.
type GivenExam struct {
	UserId     string    `json:"user_id"`
	ExamId     int       `json:"exam_id"`
	Price      string    `json:"price"`
	AddedBy    *string   `json:"added_by"`
	ScoredBy   *string   `json:"scored_by"`
	CreatedAt  time.Time `json:"created_at"`
	FinalScore *string   `json:"final_score"`
}

// NewGivenExamData is a struct that represents the data needed to
// create a new given exam.
type NewGivenExamData struct {
	UserId  string  `json:"user_id"`
	ExamId  int     `json:"exam_id"`
	Price   string  `json:"price"`
	AddedBy *string `json:"added_by"`
}

type GetMostRecentExamsData struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type MostRecentExamInfo struct {
	ExamId          int       `json:"exam_id"`
	CourseId        int       `json:"course_id"`
	ExamTitle       string    `json:"exam_title"`
	ExamDescription string    `json:"exam_description"`
	Price           string    `json:"price"`
	CreatedAt       time.Time `json:"created_at"`
	ExamDate        time.Time `json:"exam_date"`
	Duration        int       `json:"duration"`
	CreatedBy       string    `json:"created_by"`
	IsPublic        bool      `json:"is_public"`
}

type GetGivenAnswerData struct {
	ExamId     int    `json:"exam_id"`
	QuestionId int    `json:"question_id"`
	UserId     string `json:"user_id"`
}

type GivenAnswerInfo struct {
	ExamId       int        `json:"exam_id"`
	QuestionId   int        `json:"question_id"`
	AnsweredBy   string     `json:"answered_by"`
	ChosenOption *string    `json:"chosen_option"`
	SecondsTaken int        `json:"seconds_taken"`
	AnswerText   *string    `json:"answer_text"`
	AnsweredAt   time.Time  `json:"answered_at"`
	SeenAt       *time.Time `json:"seen_at"`
}

type AnswerQuestionData struct {
	ExamId       int     `json:"exam_id"`
	QuestionId   int     `json:"question_id"`
	AnsweredBy   string  `json:"answered_by"`
	ChosenOption *string `json:"chosen_option"`
	SecondsTaken int     `json:"seconds_taken"`
	AnswerText   *string `json:"answer_text"`
}

type GetUserExamsHistoryOptions struct {
	UserId string `json:"user_id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type UserOngoingExamInfo struct {
	ExamId    int       `json:"exam_id"`
	ExamTitle string    `json:"exam_title"`
	StartTime time.Time `json:"start_time"`
}

type UserPastExamInfo struct {
	ExamId    int       `json:"exam_id"`
	ExamTitle string    `json:"exam_title"`
	StartedAt time.Time `json:"started_at"`
}

type UserFutureExamInfo struct {
	ExamId    int       `json:"exam_id"`
	ExamTitle string    `json:"exam_title"`
	StartTime time.Time `json:"start_time"`
}

type GetExamParticipantsOptions struct {
	ExamId int `json:"exam_id"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

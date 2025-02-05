package examHandlers

import "time"

type CreateExamData struct {
	CourseId            int    `json:"course_id" validate:"required"`
	ExamTitle           string `json:"exam_title" validate:"required"`
	ExamDescription     string `json:"exam_description" validate:"required"`
	Price               string `json:"price" default:"0T"`
	IsPublic            bool   `json:"is_public" default:"false"`
	Duration            int    `json:"duration" default:"60"`
	ExamDate            int64  `json:"exam_date"`
	IsStrict            bool   `json:"is_strict" default:"false"`
	IsSampleExam        bool   `json:"is_sample_exam" default:"false"`
	MaxQuestionsSeconds int    `json:"max_questions_seconds" default:"0"`
	NeedsVideoCall      bool   `json:"needs_video_call" default:"false"`
	NeedsVoiceCall      bool   `json:"needs_voice_call" default:"false"`
} // @name CreateExamData

type CreateExamResult struct {
	ExamId    int       `json:"exam_id"`
	CourseId  int       `json:"course_id"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	ExamDate  time.Time `json:"exam_date"`
	Duration  int       `json:"duration"`
	CreatedBy string    `json:"created_by"`
	IsPublic  bool      `json:"is_public"`
} // @name CreateExamResult

type SearchExamData struct {
	SearchQuery string `json:"search_query" validate:"required"`
	Offset      int    `json:"offset" validate:"required"`
	Limit       int    `json:"limit" validate:"required"`
	SampleExams bool   `json:"sample_exams" default:"false"`
} // @name SearchExamData

type SearchExamResult struct {
	Exams []*SearchedExamInfo `json:"exams"`
} // @name SearchExamResult

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
} // @name SearchedExamInfo

type EditExamData struct {
	ExamId              int    `json:"exam_id"`
	CourseId            int    `json:"course_id"`
	ExamTitle           string `json:"exam_title"`
	ExamDescription     string `json:"exam_description"`
	Price               string `json:"price" default:"0T"`
	IsPublic            bool   `json:"is_public" default:"false"`
	Duration            int    `json:"duration" default:"60"`
	ExamDate            int64  `json:"exam_date"`
	IsStrict            bool   `json:"is_strict"`
	MaxQuestionsSeconds int    `json:"max_questions_seconds"`
	NeedsVideoCall      bool   `json:"needs_video_call"`
	NeedsVoiceCall      bool   `json:"needs_voice_call"`
} // @name EditExamData

type EditExamResult struct {
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
} // @name EditExamResult

type GetExamInfoResult struct {
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
	HasStarted          bool      `json:"has_started"`
	HasParticipated     bool      `json:"has_participated" default:"false"`
	CanParticipate      bool      `json:"can_participate" default:"false"`
	CanEditQuestion     bool      `json:"can_edit_question" default:"false"`
	CanAddOthersToExam  bool      `json:"can_add_others_to_exam" default:"false"`
	HasFinished         bool      `json:"has_finished" default:"false"`
	StartsIn            int       `json:"starts_in" default:"0"`
	FinishesIn          int       `json:"finishes_in" default:"0"`
	QuestionCount       int       `json:"question_count" default:"0"`
	IsStrict            bool      `json:"is_strict"`
	IsSampleExam        bool      `json:"is_sample_exam"`
	MaxQuestionsSeconds int       `json:"max_questions_seconds"`
	NeedsVideoCall      bool      `json:"needs_video_call"`
	NeedsVoiceCall      bool      `json:"needs_voice_call"`
} // @name GetExamInfoResult

type GetExamQuestionsData struct {
	Pov    string `json:"pov"` // Point of view
	ExamId int    `json:"exam_id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
} // @name GetExamQuestionsData

type GetExamQuestionsResult struct {
	Pov       string              `json:"pov"` // Point of view
	ExamId    int                 `json:"exam_id"`
	Questions []*ExamQuestionInfo `json:"questions"`
} // @name GetExamQuestionsResult

type ExamQuestionInfo struct {
	QuestionId      int                   `json:"question_id"`
	QuestionTitle   string                `json:"question_title"`
	Description     *string               `json:"description"`
	Option1         *string               `json:"option1"`
	Option2         *string               `json:"option2"`
	Option3         *string               `json:"option3"`
	Option4         *string               `json:"option4"`
	CreatedAt       time.Time             `json:"created_at"`
	IsPointer       bool                  `json:"is_pointer" default:"false"`
	PointerCount    int                   `json:"pointer_count" default:"0"`
	PointerToExamId *int                  `json:"pointer_to_exam_id"`
	UserAnswer      *AnsweredQuestionInfo `json:"user_answer"`
} // @name ExamQuestionInfo

type AnsweredQuestionInfo struct {
	UserId       string     `json:"user_id"`
	QuestionId   int        `json:"question_id"`
	ChosenOption *string    `json:"chosen_option"`
	SecondsTaken int        `json:"seconds_taken"`
	AnswerText   *string    `json:"answer"`
	SeenAt       *time.Time `json:"seen_at"`
} // @name AnsweredQuestionInfo

type ParticipateExamData struct {
	// UserId is the user who is trying to participate in the exam.
	// If the user is trying to participate in the exam themselves,
	// this field should be set to their own user id.
	UserId string `json:"user_id"`

	// ExamId is the exam the user is trying to participate in.
	ExamId int `json:"exam_id"`

	// Price is the price of the exam that user has already paid.
	Price string `json:"price"`
} // @name ParticipateExamData

type ParticipateExamResult struct {
	ExamId        int       `json:"exam_id"`
	UserId        string    `json:"user_id"`
	Price         string    `json:"price"`
	AddedBy       *string   `json:"added_by"`
	CreatedAt     time.Time `json:"created_at"`
	StartsIn      int       `json:"starts_in" default:"0"`
	FinishesIn    int       `json:"finishes_in" default:"0"`
	QuestionCount int       `json:"question_count" default:"0"`
} // @name ParticipateExamResult

type AnswerQuestionData struct {
	ExamId       int     `json:"exam_id"`
	QuestionId   int     `json:"question_id"`
	ChosenOption *string `json:"chosen_option"`
	SecondsTaken int     `json:"seconds_taken"`
	AnswerText   *string `json:"answer_text"`
} // @name AnswerQuestionData

type AnswerQuestionResult struct {
	ExamId       int        `json:"exam_id"`
	QuestionId   int        `json:"question_id"`
	AnsweredBy   string     `json:"answered_by"`
	AnsweredAt   time.Time  `json:"answered_at"`
	SeenAt       *time.Time `json:"seen_at"`
	ChosenOption *string    `json:"chosen_option"`
	SecondsTaken int        `json:"seconds_taken"`
	AnswerText   *string    `json:"answer_text"`
} // @name AnswerQuestionResult

type SetExamScoreData struct {
	// ExamId is the exam we are trying to give this score to.
	ExamId int `json:"exam_id"`

	// UserId is the person we are trying to give this score to.
	UserId string `json:"user_id"`

	// Score is the score we are trying to give to the user.
	Score string `json:"score"`
} // @name SetExamScoreData

type SetExamScoreResult struct {
	ExamId   int    `json:"exam_id"`
	UserId   string `json:"user_id"`
	Score    string `json:"score"`
	ScoredBy string `json:"scored_by"`
} // @name SetExamScoreResult

type GetGivenExamData struct {
	UserId string `json:"user_id"`
	ExamId int    `json:"exam_id"`
} // @name GetGivenExamData

type GetGivenExamResult struct {
	UserId     string    `json:"user_id"`
	ExamId     int       `json:"exam_id"`
	Price      string    `json:"price"`
	AddedBy    *string   `json:"added_by"`
	ScoredBy   *string   `json:"scored_by"`
	CreatedAt  time.Time `json:"created_at"`
	FinalScore *string   `json:"final_score"`
} // @name GetGivenExamData

type GetUserOngoingExamsResult struct {
	Exams []*UserOngoingExamInfo `json:"exams"`
} // @name GetUserOngoingExamsResult

type UserOngoingExamInfo struct {
	ExamId    int       `json:"exam_id"`
	ExamTitle int       `json:"exam_title"`
	StartTime time.Time `json:"start_time"`
} // @name UserOngoingExamInfo

type GetUsersExamHistoryData struct {
	UserId string `json:"user_id"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
} // @name GetUsersExamHistoryData

type GetUsersExamHistoryResult struct {
	Exams []*UserExamHistoryInfo `json:"exams"`
} // @name GetUsersExamHistoryResult

type UserExamHistoryInfo struct {
	ExamId    int       `json:"exam_id"`
	ExamTitle string    `json:"exam_title"`
	StartedAt time.Time `json:"started_at"`
} // @name UserExamHistoryInfo

type UserFutureExamInfo struct {
	ExamId    int       `json:"exam_id"`
	ExamTitle string    `json:"exam_title"`
	ExamDate  time.Time `json:"exam_date"`
} // @name UserFutureExamInfo

type GetUserFutureExamsResult struct {
	Exams []*UserFutureExamInfo `json:"exams"`
} // @name GetUserFutureExamsResult

type CreateExamQuestionData struct {
	ExamId          int     `json:"exam_id"`
	QuestionTitle   string  `json:"question_title"`
	Description     *string `json:"description"`
	Option1         *string `json:"option1"`
	Option2         *string `json:"option2"`
	Option3         *string `json:"option3"`
	Option4         *string `json:"option4"`
	IsPointer       bool    `json:"is_pointer" default:"false"`
	PointerCount    int     `json:"pointer_count" default:"0"`
	PointerToExamId *int    `json:"pointer_to_exam_id"`
} // @name CreateExamQuestionData

type CreateExamQuestionResult struct {
	ExamId          int       `json:"exam_id"`
	QuestionId      int       `json:"question_id"`
	QuestionTitle   string    `json:"question_title"`
	Description     *string   `json:"description"`
	Option1         *string   `json:"option1"`
	Option2         *string   `json:"option2"`
	Option3         *string   `json:"option3"`
	Option4         *string   `json:"option4"`
	CreatedAt       time.Time `json:"created_at"`
	IsPointer       bool      `json:"is_pointer" default:"false"`
	PointerCount    int       `json:"pointer_count" default:"0"`
	PointerToExamId *int      `json:"pointer_to_exam_id"`
} // @name CreateExamQuestionResult

type EditExamQuestionData struct {
	QuestionId      int     `json:"question_id"`
	ExamId          int     `json:"exam_id"`
	QuestionTitle   string  `json:"question_title"`
	Description     *string `json:"description"`
	Option1         *string `json:"option1"`
	Option2         *string `json:"option2"`
	Option3         *string `json:"option3"`
	Option4         *string `json:"option4"`
	PointerCount    int     `json:"pointer_count" default:"0"`
	PointerToExamId *int    `json:"pointer_to_exam_id"`
} // @name EditExamQuestionData

type EditExamQuestionResult struct {
	QuestionId      int       `json:"question_id"`
	ExamId          int       `json:"exam_id"`
	QuestionTitle   string    `json:"question_title"`
	Description     *string   `json:"description"`
	Option1         *string   `json:"option1"`
	Option2         *string   `json:"option2"`
	Option3         *string   `json:"option3"`
	Option4         *string   `json:"option4"`
	CreatedAt       time.Time `json:"created_at"`
	IsPointer       bool      `json:"is_pointer" default:"false"`
	PointerCount    int       `json:"pointer_count"`
	PointerToExamId *int      `json:"pointer_to_exam_id"`
} // @name EditExamQuestionResult

type GetExamParticipantsData struct {
	ExamId int `json:"exam_id"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
} // @name GetExamParticipantsData

type GetExamParticipantsResult struct {
	ExamId       int                    `json:"exam_id"`
	Participants []*ExamParticipantInfo `json:"participants"`
	CanSetScore  bool                   `json:"can_set_score" default:"false"`
} // @name GetExamParticipantsResult

type ExamParticipantInfo struct {
	UserId     string    `json:"user_id"`
	FullName   string    `json:"full_name"`
	ExamId     int       `json:"exam_id"`
	Price      string    `json:"price"`
	FinalScore *string   `json:"final_score"`
	AddedBy    *string   `json:"added_by"`
	ScoredBy   *string   `json:"scored_by"`
	CreatedAt  time.Time `json:"created_at"`
} // @name ExamParticipantInfo

package examHandlers

import "time"

type CreateExamData struct {
	CourseId        int    `json:"course_id"`
	ExamTitle       string `json:"exam_title"`
	ExamDescription string `json:"exam_description"`
	Price           string `json:"price" default:"0T"`
	IsPublic        bool   `json:"is_public" default:"false"`
	Duration        int    `json:"duration" default:"60"`
	ExamDate        int64  `json:"exam_date"`
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

type GetExamInfoResult struct {
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
	HasStarted      bool      `json:"has_started"`
	HasFinished     bool      `json:"has_finished"`
	StartsIn        int       `json:"starts_in"`
	FinishesIn      int       `json:"finishes_in"`
	QuestionCount   int       `json:"question_count"`
} // @name GetExamInfoResult

type GetExamQuestionsData struct {
	ExamId int `json:"exam_id"`
} // @name GetExamQuestionsData

type GetExamQuestionsResult struct {
	ExamId    int                 `json:"exam_id"`
	Questions []*ExamQuestionInfo `json:"questions"`
} // @name GetExamQuestionsResult

type ExamQuestionInfo struct {
	QuestionId    int                   `json:"question_id"`
	QuestionTitle string                `json:"question_title"`
	Description   *string               `json:"description"`
	Option1       *string               `json:"option1"`
	Option2       *string               `json:"option2"`
	Option3       *string               `json:"option3"`
	Option4       *string               `json:"option4"`
	CreatedAt     string                `json:"created_at"`
	UserAnswer    *AnsweredQuestionInfo `json:"user_answer"`
} // @name ExamQuestionInfo

type AnsweredQuestionInfo struct {
	QuestionId   int     `json:"question_id"`
	ChosenOption *string `json:"chosen_option"`
	SecondsTaken int     `json:"seconds_taken"`
	AnswerText   *string `json:"answer"`
} // @name AnsweredQuestionInfo

type AnswerQuestionData struct {
	ExamId       int     `json:"exam_id"`
	QuestionId   int     `json:"question_id"`
	ChosenOption *string `json:"chosen_option"`
	SecondsTaken int     `json:"seconds_taken"`
	AnswerText   *string `json:"answer_text"`
} // @name AnswerQuestionData

type AnswerQuestionResult struct {
	ExamId     int       `json:"exam_id"`
	QuestionId int       `json:"question_id"`
	AnsweredBy string    `json:"answered_by"`
	AnsweredAt time.Time `json:"answered_at"`
} // @name AnswerQuestionResult

type SetScoreData struct {
	// ExamId is the exam we are trying to give this score to.
	ExamId int `json:"exam_id"`

	// UserId is the person we are trying to give this score to.
	UserId string `json:"user_id"`

	// Score is the score we are trying to give to the user.
	Score string `json:"score"`
} // @name SetScoreData

type SetScoreResult struct {
	ExamId   int    `json:"exam_id"`
	UserId   string `json:"user_id"`
	Score    string `json:"score"`
	ScoredBy string `json:"scored_by"`
} // @name SetScoreResult

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
	ExamTitle int       `json:"course_id"`
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

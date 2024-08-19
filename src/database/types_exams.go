package database

import (
	"sync"
	"time"
)

// ExamInfo is a struct that represents the information of an exam.
type ExamInfo struct {
	ExamId    int       `json:"exam_id"`
	CourseId  int       `json:"course_id"`
	Price     string    `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	ExamDate  time.Time `json:"exam_date"`
	Duration  int       `json:"duration"`
	CreatedBy string    `json:"created_by"`
	IsPublic  bool      `json:"is_public"`

	mut *sync.RWMutex `json:"-"`
	// Questions is a slice of exam questions.
	Questions []*ExamQuestion `json:"-"`
}

// NewExamData is a struct that represents the data needed to create a new exam.
type NewExamData struct {
	CourseId  int    `json:"course_id"`
	Price     string `json:"price"`
	CreatedBy string `json:"created_by"`
	IsPublic  bool   `json:"is_public"`
	Duration  int    `json:"duration"`
	ExamDate  string `json:"exam_date"`
}

// ExamQuestion is a struct that represents the information of an exam question.
type ExamQuestion struct {
	QuestionId    int     `json:"question_id"`
	ExamId        int     `json:"exam_id"`
	QuestionTitle string  `json:"question_title"`
	Description   *string `json:"description"`
	Option1       *string `json:"option1"`
	Option2       *string `json:"option2"`
	Option3       *string `json:"option3"`
	Option4       *string `json:"option4"`
	CreatedAt     string  `json:"created_at"`
}

// NewExamQuestionData is a struct that represents the data needed to create a new exam question.
type NewExamQuestionData struct {
	ExamId        int     `json:"exam_id"`
	QuestionTitle string  `json:"question_title"`
	Description   *string `json:"description"`
	Option1       *string `json:"option1"`
	Option2       *string `json:"option2"`
	Option3       *string `json:"option3"`
	Option4       *string `json:"option4"`
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

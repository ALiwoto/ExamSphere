package database

import (
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

func (e *ExamInfo) HasExamStarted() bool {
	return time.Now().After(e.ExamDate)
}

func (e *ExamInfo) HasExamFinished() bool {
	return time.Now().After(e.ExamDate.Add(time.Minute * time.Duration(e.Duration)))
}

func (e *ExamInfo) ExamStartsIn() int {
	until := time.Until(e.ExamDate)
	if until < 0 {
		return 0
	}
	return int(until.Minutes())
}

func (e *ExamInfo) ExamFinishesIn() int {
	return int(time.Until(e.ExamDate.Add(time.Minute * time.Duration(e.Duration))).Minutes())
}

//-------------------------------------------------------------

func (e *ExamQuestion) GetUniqueId() string {
	return ssg.ToBase10(e.ExamId) + KeySepChar + ssg.ToBase10(e.QuestionId)
}

// HasOption checks if the given option is one of the options of the question.
func (e *ExamQuestion) HasOption(option string) bool {
	return (e.Option1 != nil && *e.Option1 == option) ||
		(e.Option2 != nil && *e.Option2 == option) ||
		(e.Option3 != nil && *e.Option3 == option) ||
		(e.Option4 != nil && *e.Option4 == option)
}

//-------------------------------------------------------------

func (g *GivenExam) GetUniqueId() string {
	return g.UserId + KeySepChar + ssg.ToBase10(g.ExamId)
}

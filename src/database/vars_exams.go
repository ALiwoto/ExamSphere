package database

import (
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

var (
	examsInfoMap = func() *ssg.SafeEMap[int, ExamInfo] {
		m := ssg.NewSafeEMap[int, ExamInfo]()
		m.SetExpiration(time.Hour * 3)
		m.SetInterval(time.Hour * 12)
		m.EnableChecking()

		return m
	}()

	examQuestionsMap = func() *ssg.SafeEMap[int, ExamQuestion] {
		m := ssg.NewSafeEMap[int, ExamQuestion]()
		m.SetExpiration(time.Hour * 3)
		m.SetInterval(time.Hour * 12)
		m.EnableChecking()

		return m
	}()

	givenExamsMap = func() *ssg.SafeEMap[string, GivenExam] {
		m := ssg.NewSafeEMap[string, GivenExam]()
		m.SetExpiration(time.Hour * 3)
		m.SetInterval(time.Hour * 12)
		m.EnableChecking()

		return m
	}()

	givenAnswersMap = func() *ssg.SafeEMap[string, GivenAnswerInfo] {
		m := ssg.NewSafeEMap[string, GivenAnswerInfo]()
		m.SetExpiration(time.Hour * 3)
		m.SetInterval(time.Hour * 12)
		m.EnableChecking()

		return m
	}()
)

var (
	valueExamNotFound         = &ExamInfo{}
	valueExamQuestionNotFound = &ExamQuestion{}
	valueGivenExamNotFound    = &GivenExam{}
	valueGivenAnswerNotFound  = &GivenAnswerInfo{}
)

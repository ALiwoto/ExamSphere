package database

import "github.com/ALiwoto/ssg/ssg"

func (e *ExamInfo) lock() {
	e.mut.Lock()
}

func (e *ExamInfo) unlock() {
	e.mut.Unlock()
}

func (e *ExamInfo) RLock() {
	e.mut.RLock()
}

func (e *ExamInfo) RUnlock() {
	e.mut.RUnlock()
}

func (e *ExamInfo) AddQuestion(q *ExamQuestion) {
	e.lock()
	defer e.unlock()

	for i := 0; i < len(e.Questions); i++ {
		if e.Questions[i].QuestionId == q.QuestionId {
			// since question ids are unique, we can replace the question
			e.Questions[i] = q
			return
		}
	}

	e.Questions = append(e.Questions, q)
}

func (e *ExamInfo) GetQuestion(questionId int) *ExamQuestion {
	e.RLock()
	defer e.RUnlock()

	for _, q := range e.Questions {
		if q.QuestionId == questionId {
			return q
		}
	}

	return nil
}

func (e *ExamInfo) GetQuestions() []*ExamQuestion {
	e.RLock()
	defer e.RUnlock()

	return e.Questions
}

//-------------------------------------------------------------

func (g *GivenExam) GetUniqueId() string {
	return g.UserId + KeySepChar + ssg.ToBase10(g.ExamId)
}

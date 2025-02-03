package examHandlers

import (
	"time"
)

func (d *CreateExamData) IsValid() bool {
	if d.IsSampleExam {
		d.ExamDate = 0
		d.NeedsVideoCall = false
		d.NeedsVoiceCall = false
		d.Duration = 0
		return true
	}

	return d.CourseId != 0 &&
		d.Duration > 0 &&
		d.ExamDate >= time.Now().UTC().Unix()
}

//-------------------------------------------------------------

func (d *EditExamData) IsValid() bool {
	return d.ExamId != 0 &&
		d.CourseId != 0 &&
		d.Price != "" &&
		d.Duration > 0
}

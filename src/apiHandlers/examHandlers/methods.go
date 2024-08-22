package examHandlers

import "time"

func (d *CreateExamData) IsValid() bool {
	return d.CourseId != 0 &&
		d.Price != "" &&
		d.Duration > 0 &&
		d.ExamDate >= time.Now().UTC().Unix()
}

package userHandlers

import (
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

var (
	requestRateLimitMap = func() *ssg.SafeEMap[string, userRequestEntry] {
		m := ssg.NewSafeEMap[string, userRequestEntry]()
		m.SetExpiration(time.Minute * 40)
		m.SetInterval(time.Hour)
		m.EnableChecking()

		return m
	}()
)

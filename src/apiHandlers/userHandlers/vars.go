package userHandlers

import (
	"regexp"
	"sync"
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

	changePasswordRequestMap = func() *ssg.SafeEMap[string, changePasswordRequestEntry] {
		m := ssg.NewSafeEMap[string, changePasswordRequestEntry]()

		// NOTE: the expiration parameter here is not the time where token parameters
		// are valid, the token parameters should be valid for *few minutes* only.
		// (which we can also consider set it in config file, for now it's hardcoded).
		// So the time set here, is actually the time-frame where we limit the user for
		// changing the password too frequently.
		// E.g. if this is set to 2 hours, and the user change their password for more than
		// 10 times in 2 hours, they will be rate-limited.
		// Their rate-limit will be removed after 2 hours (when this entry gets removed from
		// the cache).
		m.SetExpiration(time.Hour)
		m.SetInterval(time.Hour)
		m.EnableChecking()

		return m
	}()

	passwordChangeRqGenerator = ssg.NewNumIdGenerator[int32]()
)

var (
	createUserMutex = &sync.Mutex{}
)

var (
	emailRegex = regexp.MustCompile(`^[A-Za-z0-9._+%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$`)
)

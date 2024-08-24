package tests

import (
	"ExamSphere/src/core/utils/emailUtils"
	"fmt"
	"testing"
)

func TestChangePasswordTemplateEn(t *testing.T) {
	ok := fmt.Sprintf(emailUtils.PasswordChangeEmailTemplate_en, "John Doe", "http://example.com")

	println(ok)
}

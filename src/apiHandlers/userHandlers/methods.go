package userHandlers

import (
	"ExamSphere/src/core/utils/hashing"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/ALiwoto/ssg/ssg"
)

func (e *changePasswordRequestEntry) GetRedirectAddress(baseAddr string) string {
	// URL encode the query parameters
	encodedRq := url.QueryEscape(e.RqId)
	encodedRt := url.QueryEscape(e.RTParam)
	encodedLt := url.QueryEscape(hex.EncodeToString([]byte(ssg.ToBase16(e.LTNum))))

	// Construct the full URL with query parameters
	fullURL := fmt.Sprintf("%s?rq=%s&rt=%s&lt=%s", baseAddr, encodedRq, encodedRt, encodedLt)

	return fullURL
}

func (e *changePasswordRequestEntry) Verify(data *ConfirmChangePasswordData) bool {
	// check if the request id matches
	if e.RqId != data.RqId {
		return false
	}

	// check if the RTParam matches
	if e.RTParam != data.RTParam {
		return false
	}

	// check if the RTVerifier is correct
	if !hashing.CompareSHA256(data.RTHash, e.RTParam) {
		return false
	}

	if !hashing.CompareSHA256(data.RTVerifier, "LT:"+ssg.ToBase10(e.LTNum)) {
		return false
	}

	return true
}

//---------------------------------------------------------

func (e *confirmAccountRequestEntry) GetRedirectAddress(baseAddr string) string {
	// URL encode the query parameters
	encodedConfirmToken := url.QueryEscape(e.ConfirmToken)
	encodedRLToken := url.QueryEscape(e.RLToken)

	// Construct the full URL with query parameters
	fullURL := fmt.Sprintf(
		"%s?confirmToken=%s&rlToken=%s&lt=%s",
		baseAddr,
		encodedConfirmToken,
		encodedRLToken,

		// later on we should change this; if the user is not confirming
		// their account through the email.
		"emailConfirmation",
	)

	return fullURL
}

//-------------------------------------------------------------

func (d *ConfirmAccountData) IsValid() bool {
	return d != nil &&
		d.ConfirmToken != "" &&
		d.RLToken != "" &&
		d.LTToken != ""
}

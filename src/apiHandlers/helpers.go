package apiHandlers

import (
	"OnlineExams/src/core/appValues"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetJWTClaimsStr(token string, key []byte) jwt.MapClaims {
	user, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil || user == nil {
		return nil
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		return nil
	}

	return claims
}

func GetJWTClaimsInfoStr(token string, key []byte) *appValues.JWTClaimsInfo {
	claims := GetJWTClaimsStr(token, key)
	if claims == nil {
		return nil
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return nil
	}

	isRefresh, ok := claims["refresh"].(bool)
	if !ok {
		return nil
	}

	authHash, ok := claims["auth_hash"].(string)
	if !ok || authHash == "" {
		return nil
	}

	expF, ok := claims["exp"].(float64)
	exp := int64(expF)
	if !ok || exp == 0 || time.Unix(exp, 0).Before(time.Now()) {
		return nil
	}

	return &appValues.JWTClaimsInfo{
		UserId:   userId,
		Refresh:  isRefresh,
		AuthHash: authHash,
		Exp:      exp,
	}
}

func GetJWTClaims(c *fiber.Ctx) jwt.MapClaims {
	user, ok := c.Locals("user").(*jwt.Token)
	if !ok || user == nil {
		return nil
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		return nil
	}

	return claims
}

func GetJWTClaimsInfo(c *fiber.Ctx) *appValues.JWTClaimsInfo {
	claims := GetJWTClaims(c)
	if claims == nil {
		return nil
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return nil
	}

	isRefresh, ok := claims["refresh"].(bool)
	if !ok {
		return nil
	}

	authHash, ok := claims["auth_hash"].(string)
	if !ok || authHash == "" {
		return nil
	}

	expF, ok := claims["exp"].(float64)
	exp := int64(expF)
	if !ok || exp == 0 || time.Unix(exp, 0).Before(time.Now()) {
		return nil
	}

	return &appValues.JWTClaimsInfo{
		UserId:   userId,
		Refresh:  isRefresh,
		AuthHash: authHash,
		Exp:      exp,
	}
}

func SendResult(c *fiber.Ctx, data any) error {
	return c.JSON(&EndpointResponse{
		Success: true,
		Result:  data,
	})
}

func SendError(status int, c *fiber.Ctx, err *EndpointError) error {
	err.Date = time.Now().Format(time.RFC3339)
	return c.Status(status).JSON(&EndpointResponse{
		Success: false,
		Error:   err,
	})
}

func SendErrMalformedJWT(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeMalformedJWT,
		Message:   ErrMalformedJWT,
		Origin:    c.Path(),
	})
}

func SendErrInvalidJWT(c *fiber.Ctx) error {
	return SendError(fiber.StatusUnauthorized, c, &EndpointError{
		ErrorCode: ErrCodeInvalidJWT,
		Message:   ErrInvalidJWT,
		Origin:    c.Path(),
	})
}

func SendErrInvalidBodyData(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidBodyData,
		Message:   ErrInvalidBodyData,
		Origin:    c.Path(),
	})
}

func SendErrInvalidUsernamePass(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidUsernamePass,
		Message:   ErrInvalidUsernamePass,
		Origin:    c.Path(),
	})
}

func SendErrInvalidAuth(c *fiber.Ctx) error {
	return SendError(fiber.StatusUnauthorized, c, &EndpointError{
		ErrorCode: ErrCodeInvalidAuth,
		Message:   ErrInvalidAuth,
		Origin:    c.Path(),
	})
}

func SendErrPermissionDenied(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodePermissionDenied,
		Message:   ErrPermissionDenied,
		Origin:    c.Path(),
	})
}

func SendErrInvalidInputPass(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidInputPass,
		Message:   ErrInvalidInputPass,
		Origin:    c.Path(),
	})
}

func SendErrUsernameExists(c *fiber.Ctx) error {
	return SendError(fiber.StatusConflict, c, &EndpointError{
		ErrorCode: ErrCodeUsernameExists,
		Message:   ErrUsernameExists,
		Origin:    c.Path(),
	})
}

func SendErrInternalServerError(c *fiber.Ctx) error {
	return SendError(fiber.StatusInternalServerError, c, &EndpointError{
		ErrorCode: ErrCodeInternalServerError,
		Message:   ErrInternalServerError,
		Origin:    c.Path(),
	})
}

func SendErrInvalidFileData(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidFileData,
		Message:   ErrInvalidFileData,
		Origin:    c.Path(),
	})
}

func SendErrInvalidPhoneNumber(c *fiber.Ctx, phone string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidPhoneNumber,
		Message:   fmt.Sprintf(ErrInvalidPhoneNumber, phone),
		Origin:    c.Path(),
	})
}

func SendErrPhoneNumberAlreadyImported(c *fiber.Ctx, phone string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodePhoneNumberAlreadyImported,
		Message:   fmt.Sprintf(ErrPhoneNumberAlreadyImported, phone),
		Origin:    c.Path(),
	})
}

func SendErrInvalidUsername(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidUsername,
		Message:   ErrInvalidUsername,
		Origin:    c.Path(),
	})
}

func SendErrNoPhonesDonated(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeNoPhonesDonated,
		Message:   ErrNoPhonesDonated,
		Origin:    c.Path(),
	})
}

func SendErrAgentNotConnected(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeAgentNotConnected,
		Message:   ErrAgentNotConnected,
		Origin:    c.Path(),
	})
}

func SendErrInvalidPagination(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidPagination,
		Message:   ErrInvalidPagination,
		Origin:    c.Path(),
	})
}

func SendErrMaxContactImportLimit(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeMaxContactImportLimit,
		Message:   ErrMaxContactImportLimit,
		Origin:    c.Path(),
	})
}

func SendErrPhoneNumberNotFound(c *fiber.Ctx, phone string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodePhoneNumberNotFound,
		Message:   fmt.Sprintf(ErrPhoneNumberNotFound, phone),
		Origin:    c.Path(),
	})
}

func SendErrParameterRequired(c *fiber.Ctx, paramName string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeParameterRequired,
		Message:   fmt.Sprintf(ErrParameterRequired, paramName),
		Origin:    c.Path(),
	})
}

func SendErrUserBanned(c *fiber.Ctx) error {
	return SendError(fiber.StatusForbidden, c, &EndpointError{
		ErrorCode: ErrCodeUserBanned,
		Message:   ErrUserBanned,
		Origin:    c.Path(),
	})
}

func SendErrLabelInfoNotFound(c *fiber.Ctx, labelId int) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeLabelInfoNotFound,
		Message:   fmt.Sprintf(ErrLabelInfoNotFound, labelId),
		Origin:    c.Path(),
	})
}

func SendErrLabelAlreadyApplied(c *fiber.Ctx, labelId int, phone string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeLabelAlreadyApplied,
		Message:   fmt.Sprintf(ErrLabelAlreadyApplied, labelId, phone),
		Origin:    c.Path(),
	})
}

func SendErrLabelAlreadyExistsByName(c *fiber.Ctx, labelName string) error {
	return SendError(fiber.StatusConflict, c, &EndpointError{
		ErrorCode: ErrCodeLabelAlreadyExistsByName,
		Message:   fmt.Sprintf(ErrLabelAlreadyExistsByName, labelName),
		Origin:    c.Path(),
	})
}

func SendErrTooManyChatLabelInfo(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeTooManyChatLabelInfo,
		Message:   ErrTooManyChatLabelInfo,
		Origin:    c.Path(),
	})
}

func SendErrLabelNameTooLong(c *fiber.Ctx, maxLen int) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeLabelNameTooLong,
		Message:   fmt.Sprintf(ErrLabelNameTooLong, maxLen),
		Origin:    c.Path(),
	})
}

func SendErrLabelDescriptionTooLong(c *fiber.Ctx, maxLen int) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeLabelDescriptionTooLong,
		Message:   fmt.Sprintf(ErrLabelDescriptionTooLong, maxLen),
		Origin:    c.Path(),
	})
}

func SendErrInvalidColor(c *fiber.Ctx, color string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidColor,
		Message:   fmt.Sprintf(ErrInvalidColor, color),
		Origin:    c.Path(),
	})
}

func SendErrLabelNotApplied(c *fiber.Ctx, labelId int, phone string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeLabelNotApplied,
		Message:   fmt.Sprintf(ErrLabelNotApplied, labelId, phone),
		Origin:    c.Path(),
	})
}

func SendErrCannotDeleteBuiltInLabel(c *fiber.Ctx, labelId int) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeCannotDeleteBuiltInLabel,
		Message:   fmt.Sprintf(ErrCannotDeleteBuiltInLabel, labelId),
		Origin:    c.Path(),
	})
}

func SendErrDuplicatePhoneNumber(c *fiber.Ctx, phone string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeDuplicatePhoneNumber,
		Message:   fmt.Sprintf(ErrDuplicatePhoneNumber, phone),
		Origin:    c.Path(),
	})
}

func SendErrPhoneNotWorking(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodePhoneNotWorking,
		Message:   ErrPhoneNotWorking,
		Origin:    c.Path(),
	})
}

func SendErrInvalidPmsPass(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidPmsPass,
		Message:   ErrInvalidPmsPass,
		Origin:    c.Path(),
	})
}

func SendErrInvalidAgentId(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidAgentId,
		Message:   ErrInvalidAgentId,
		Origin:    c.Path(),
	})
}

func SendErrInvalidAppSettingName(c *fiber.Ctx, settingName string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidAppSettingName,
		Message:   fmt.Sprintf(ErrInvalidAppSettingName, settingName),
		Origin:    c.Path(),
	})
}

func SendErrAppSettingNotFound(c *fiber.Ctx, settingName string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeAppSettingNotFound,
		Message:   fmt.Sprintf(ErrAppSettingNotFound, settingName),
		Origin:    c.Path(),
	})
}

func SendErrTextEmpty(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeTextEmpty,
		Message:   ErrTextEmpty,
		Origin:    c.Path(),
	})
}

func SendErrTextTooLong(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeTextTooLong,
		Message:   ErrTextTooLong,
		Origin:    c.Path(),
	})
}

func SendErrInvalidClientRId(c *fiber.Ctx, rId string) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidClientRId,
		Message:   fmt.Sprintf(ErrInvalidClientRId, rId),
		Origin:    c.Path(),
	})
}

func SendErrInvalidCaptcha(c *fiber.Ctx) error {
	return SendError(fiber.StatusBadRequest, c, &EndpointError{
		ErrorCode: ErrCodeInvalidCaptcha,
		Message:   ErrInvalidCaptcha,
		Origin:    c.Path(),
	})
}

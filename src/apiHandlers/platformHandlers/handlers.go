package platformHandlers

import (
	"ExamSphere/src/apiHandlers"
	"ExamSphere/src/apiUtils"
	"ExamSphere/src/core/utils/logging"

	"github.com/gofiber/fiber/v2"
	fUtils "github.com/gofiber/fiber/v2/utils"
)

// GetPlatformLogsV1 godoc
// @Summary Get platform logs
// @Description Allows a client to get platform logs
// @ID GetPlatformLogsV1
// @Tags Platform
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param storage query string false "Storage name"
// @Success 200 {object} apiHandlers.EndpointResponse{result=GetPlatformLogsResult}
// @Router /api/v1/platform/logs [get]
func GetPlatformLogsV1(c *fiber.Ctx) error {
	userInfo, err := apiUtils.GetUserInfo(c)
	if err != nil {
		return err
	} else if userInfo == nil {
		// sending error json is already handled in GetUserInfo
		return nil
	}

	if !userInfo.CanGetPlatformLogs() {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	// get storage name from query
	storageName := fUtils.CopyString(c.Query("storage"))

	logs, err := logging.GetAllLogEntries(storageName)
	if err != nil {
		logging.UnexpectedError("GetPlatformLogsV1: GetAllLogEntries:", err)
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &GetPlatformLogsResult{
		Logs: logs,
	})
}

package topicHandlers

import (
	"ExamSphere/src/apiHandlers"
	"ExamSphere/src/database"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CreateTopicV1 godoc
// @Summary Create a new topic
// @Description Create a new topic
// @ID createTopicV1
// @Tags Topic
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param data body CreateNewTopicData true "Data needed to create a new topic"
// @Success 200 {object} apiHandlers.EndpointResponse{result=CreateNewTopicResult}
// @Router /api/v1/topic/create [post]
func CreateTopicV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	} else if !userInfo.CanCreateNewTopic() {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	data := &CreateNewTopicData{}
	if err := c.BodyParser(data); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if !data.IsValid() {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	data.TopicName = strings.TrimSpace(data.TopicName)
	topics, _ := database.GetTopicInfoByName(data.TopicName)
	if len(topics) > 0 {
		return apiHandlers.SendErrTopicNameExists(c)
	}

	topicInfo, err := database.CreateNewTopic(&database.NewTopicData{
		TopicName: data.TopicName,
	})
	if err != nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &CreateNewTopicResult{
		TopicId:   topicInfo.TopicId,
		TopicName: topicInfo.TopicName,
	})
}

// SearchTopicV1 godoc
// @Summary Search for topics
// @Description Search for topics
// @ID searchTopicV1
// @Tags Topic
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param data body SearchTopicData true "Data needed to search for topics"
// @Success 200 {object} apiHandlers.EndpointResponse{result=SearchTopicResult}
// @Router /api/v1/topic/search [post]
func SearchTopicV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	} else if !userInfo.CanSearchTopic() {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	data := &SearchTopicData{}
	if err := c.BodyParser(data); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	topicsInfo, err := database.SearchTopics(data.TopicName)
	if err != nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	sTopics := make([]*SearchedTopicInfo, 0, len(topicsInfo))
	for _, topicInfo := range topicsInfo {
		sTopics = append(sTopics, &SearchedTopicInfo{
			TopicId:   topicInfo.TopicId,
			TopicName: topicInfo.TopicName,
		})
	}

	return apiHandlers.SendResult(c, &SearchTopicResult{
		Topics: sTopics,
	})
}

// GetTopicInfoV1 godoc
// @Summary Get topic info
// @Description Get topic info
// @ID getTopicInfoV1
// @Tags Topic
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id query int true "Topic ID"
// @Success 200 {object} apiHandlers.EndpointResponse{result=GetTopicInfoResult}
// @Router /api/v1/topic/info [get]
func GetTopicInfoV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	} else if !userInfo.CanGetTopicInfo() {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	topicId := c.QueryInt("id")
	if topicId == 0 {
		return apiHandlers.SendErrParameterRequired(c, "id")
	}

	topicInfo, err := database.GetTopicInfo(topicId)
	if err != nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &GetTopicInfoResult{
		TopicId:   topicInfo.TopicId,
		TopicName: topicInfo.TopicName,
	})
}

// GetUserTopicStatV1 godoc
// @Summary Get user topic stat
// @Description Get user topic stat
// @ID getUserTopicStatV1
// @Tags Topic
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param data body GetUserTopicStatData true "Data needed to get user topic stat"
// @Success 200 {object} apiHandlers.EndpointResponse{result=GetUserTopicStatResult}
// @Router /api/v1/topic/userTopicStat [post]
func GetUserTopicStatV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	}

	data := &GetUserTopicStatData{}
	if err := c.BodyParser(data); err != nil {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	if data.UserId == "" || data.TopicId == 0 {
		return apiHandlers.SendErrInvalidBodyData(c)
	}

	stat, err := database.GetUserTopicStat(data.UserId, data.TopicId)
	if err != nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, &GetUserTopicStatResult{
		Stat: &UserTopicStatInfo{
			TopicId:      stat.TopicId,
			UserId:       stat.UserId,
			CurrentExp:   stat.CurrentExp,
			TotalExp:     stat.TotalExp,
			CurrentLevel: stat.CurrentLevel,
			LastVisited:  stat.LastVisited,
		},
	})
}

// GetAllUserTopicStatsV1 godoc
// @Summary Get all user topic stats
// @Description Get all user topic stats
// @ID getAllUserTopicStatsV1
// @Tags Topic
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} apiHandlers.EndpointResponse{result=GetAllUserTopicStatsResult}
// @Router /api/v1/topic/allUserTopicStats [get]
func GetAllUserTopicStatsV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	}

	stats, err := database.GetAllUserTopicStats(claimInfo.UserId)
	if err != nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	sStats := make([]*UserTopicStatInfo, 0, len(stats))
	for _, stat := range stats {
		sStats = append(sStats, &UserTopicStatInfo{
			TopicId:      stat.TopicId,
			UserId:       stat.UserId,
			CurrentExp:   stat.CurrentExp,
			TotalExp:     stat.TotalExp,
			CurrentLevel: stat.CurrentLevel,
			LastVisited:  stat.LastVisited,
		})
	}

	return apiHandlers.SendResult(c, &GetAllUserTopicStatsResult{
		UserId: claimInfo.UserId,
		Stats:  sStats,
	})
}

// DeleteTopicV1 godoc
// @Summary Delete a topic
// @Description Allows moderators to delete a topic. All courses and exams related to the topic will be deleted as well.
// @ID deleteTopicV1
// @Tags Topic
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization token"
// @Param id query int true "Topic ID"
// @Success 200 {object} apiHandlers.EndpointResponse{result=bool}
// @Router /api/v1/topic/delete [delete]
func DeleteTopicV1(c *fiber.Ctx) error {
	claimInfo := apiHandlers.GetJWTClaimsInfo(c)
	if claimInfo == nil {
		return apiHandlers.SendErrInvalidJWT(c)
	}

	userInfo := database.GetUserInfoByAuthHash(
		claimInfo.UserId, claimInfo.AuthHash,
	)
	if userInfo == nil {
		return apiHandlers.SendErrInvalidAuth(c)
	} else if !userInfo.CanDeleteTopic() {
		return apiHandlers.SendErrPermissionDenied(c)
	}

	topicId := c.QueryInt("id")
	if topicId == 0 {
		return apiHandlers.SendErrParameterRequired(c, "id")
	}

	err := database.DeleteTopicById(topicId)
	if err != nil {
		return apiHandlers.SendErrInternalServerError(c)
	}

	return apiHandlers.SendResult(c, true)
}

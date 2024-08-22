package topicHandlers

import "time"

type CreateNewTopicData struct {
	TopicName string `json:"topic_name"`
} // @name CreateNewTopicData

type CreateNewTopicResult struct {
	TopicId   int    `json:"topic_id"`
	TopicName string `json:"topic_name"`
} // @name CreateNewTopicResult

type UserTopicStatInfo struct {
	TopicId      int       `json:"topic_id"`
	UserId       string    `json:"user_id"`
	CurrentExp   int       `json:"current_exp"`
	TotalExp     int       `json:"total_exp"`
	CurrentLevel int       `json:"current_level"`
	LastVisited  time.Time `json:"last_visited"`
} // @name UserTopicStatInfo

type GetUserTopicStatData struct {
	UserId  string `json:"user_id"`
	TopicId int    `json:"topic_id"`
} // @name GetUserTopicStatData

type GetUserTopicStatResult struct {
	Stat *UserTopicStatInfo `json:"stat"`
} // @name GetUserTopicStatResult

type GetAllUserTopicStatsResult struct {
	UserId string               `json:"user_id"`
	Stats  []*UserTopicStatInfo `json:"stats"`
} // @name GetAllUserTopicStatsResult

type SearchTopicData struct {
	TopicName string `json:"topic_name"`
} // @name SearchTopicData

type SearchTopicResult struct {
	Topics []*SearchedTopicInfo `json:"topics"`
} // @name SearchTopicResult

type SearchedTopicInfo struct {
	TopicId   int    `json:"topic_id"`
	TopicName string `json:"topic_name"`
} // @name SearchedTopicInfo

type GetTopicInfoResult struct {
	TopicId   int    `json:"topic_id"`
	TopicName string `json:"topic_name"`
} // @name GetTopicInfoResult

package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// CreateNewTopic creates a new topic in the database,
// using the plpgsql function create_topic_info.
func CreateNewTopic(topicData *NewTopicData) (*TopicInfo, error) {
	info := &TopicInfo{
		TopicName: topicData.TopicName,
	}

	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT create_topic_info($1)`,
		info.TopicName,
	).Scan(&info.TopicId)
	if err != nil {
		return nil, err
	}

	topicsInfoMap.Add(info.TopicId, info)

	return info, nil
}

// GetTopicInfo gets a topic from the database.
func GetTopicInfo(topicId int) (*TopicInfo, error) {
	info := topicsInfoMap.Get(topicId)
	if info != nil && info != valueTopicNotFound && info.TopicId == topicId {
		return info, nil
	}

	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT topic_id, topic_name
		FROM topic_info WHERE topic_id = $1`,
		topicId,
	).Scan(
		&info.TopicId,
		&info.TopicName,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrTopicNotFound
		}

		return nil, err
	}

	topicsInfoMap.Add(info.TopicId, info)
	return info, nil
}

// SearchTopics searches for topics in the database.
func SearchTopics(topicName string) ([]*TopicInfo, error) {
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT topic_id, topic_name
		FROM topic_info WHERE topic_name ILIKE '%' || $1 || '%'`,
		topicName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topics := make([]*TopicInfo, 0)
	for rows.Next() {
		info := &TopicInfo{}
		err = rows.Scan(
			&info.TopicId,
			&info.TopicName,
		)
		if err != nil {
			return nil, err
		}

		topics = append(topics, info)
	}

	return topics, nil
}

// GetAllUserTopicStats gets all the topic stats for a user.
func GetAllUserTopicStats(userId string) ([]*UserTopicStat, error) {
	rows, err := DefaultContainer.db.Query(context.Background(),
		`SELECT user_id, topic_id, current_exp, total_exp, current_level, last_visited
		FROM user_topic_stat WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*UserTopicStat, 0)
	for rows.Next() {
		stat := &UserTopicStat{}
		err = rows.Scan(
			&stat.UserId,
			&stat.TopicId,
			&stat.CurrentExp,
			&stat.TotalExp,
			&stat.CurrentLevel,
			&stat.LastVisited,
		)
		if err != nil {
			return nil, err
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

// GetUserTopicStat gets the topic stat for a user and a topic.
func GetUserTopicStat(userId string, topicId int) (*UserTopicStat, error) {
	stat := &UserTopicStat{}
	err := DefaultContainer.db.QueryRow(context.Background(),
		`SELECT user_id, topic_id, current_exp, total_exp, current_level, last_visited
		FROM user_topic_stat WHERE user_id = $1 AND topic_id = $2`,
		userId, topicId,
	).Scan(
		&stat.UserId,
		&stat.TopicId,
		&stat.CurrentExp,
		&stat.TotalExp,
		&stat.CurrentLevel,
		&stat.LastVisited,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserTopicStatNotFound
		}

		return nil, err
	}

	return stat, nil
}

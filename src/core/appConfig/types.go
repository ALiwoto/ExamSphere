package appConfig

import "time"

type Minute = time.Duration

type PlatformConfig struct {
	OwnerUsername string `key:"owner_username"`
	OwnerPassword string `key:"owner_password"`
	SudoToken     string `key:"sudo_token"`
	BindAddress   string `key:"bind_address" default:":8080"`
	Debug         bool   `key:"debug"`

	CertFile    string `key:"cert_file"`
	CertKeyFile string `key:"cert_key_file"`

	AccessTokenSigningKey  string `key:"access_token_signing_key"`
	RefreshTokenSigningKey string `key:"refresh_token_signing_key"`

	PostgresHost     string `key:"postgres_host"`
	PostgresUser     string `key:"postgres_user"`
	PostgresPassword string `key:"postgres_password"`
	PostgresDb       string `key:"postgres_db"`

	MaxOutgoingMessagesCount      int    `key:"max_outgoing_messages_count" default:"15"`
	MaxDailyNewConversationsCount int    `key:"max_daily_new_conversations_count" default:"45"`
	NewConversationMinDelay       Minute `key:"new_conversation_min_delay" default:"2"`
	MaxRequestTillRateLimit       int    `key:"max_request_till_rate_limit" default:"15"`
	MaxRateLimitDuration          Minute `key:"max_rate_limit_duration" default:"3"`
	RateLimitPunishmentDuration   Minute `key:"rate_limit_punishment_duration" default:"5"`
	AdminStatsCacheDuration       Minute `key:"admin_stats_cache_duration" default:"3"`
}

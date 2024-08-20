package appConfig

import "time"

type Minute = time.Duration

type PlatformConfig struct {
	OwnerUsername string `key:"owner_username"`
	OwnerPassword string `key:"owner_password"`
	SudoToken     string `key:"sudo_token"`
	BindAddress   string `key:"bind_address" default:":8080"`

	// If you are behind a load-balancer such as cloudflare or
	// nginx, you should enable this option to get the real client IP,
	// for nginx, it's usually "X-Real-IP" and for cloudflare, it's
	// "CF-Connecting-IP".
	// Otherwise, just set it to an empty string.
	IPProxyHeader       string `key:"ip_proxy_header"`
	Debug               bool   `key:"debug"`
	SwaggerInstanceName string `key:"swagger_instance_name"`
	SwaggerTitle        string `key:"swagger_title"`
	SwaggerBaseURL      string `key:"swagger_base_url"`

	CertFile    string `key:"cert_file"`
	CertKeyFile string `key:"cert_key_file"`

	AccessTokenSigningKey  string `key:"access_token_signing_key"`
	RefreshTokenSigningKey string `key:"refresh_token_signing_key"`

	PostgresHost     string `key:"postgres_host"`
	PostgresUser     string `key:"postgres_user"`
	PostgresPassword string `key:"postgres_password"`
	PostgresDb       string `key:"postgres_db"`

	EmailFrom             string `key:"email_from"`
	EmailHost             string `key:"email_host"`
	EmailUser             string `key:"email_user"`
	EmailPort             int    `key:"email_port" default:"587"`
	EmailPass             string `key:"email_pass"`
	ChangePassBaseUrl     string `key:"change_pass_base_url"`
	ConfirmAccountBaseUrl string `key:"confirm_account_base_url"`

	MaxOutgoingMessagesCount      int    `key:"max_outgoing_messages_count" default:"15"`
	MaxDailyNewConversationsCount int    `key:"max_daily_new_conversations_count" default:"45"`
	NewConversationMinDelay       Minute `key:"new_conversation_min_delay" default:"2"`
	MaxRequestTillRateLimit       int    `key:"max_request_till_rate_limit" default:"15"`
	MaxRateLimitDuration          Minute `key:"max_rate_limit_duration" default:"3"`
	RateLimitPunishmentDuration   Minute `key:"rate_limit_punishment_duration" default:"5"`
	AdminStatsCacheDuration       Minute `key:"admin_stats_cache_duration" default:"3"`
}

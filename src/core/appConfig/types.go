package appConfig

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
}

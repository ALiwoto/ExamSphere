package appConfig

import (
	"OnlineExams/src/core/utils/hashing"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

func LoadConfig() error {
	return LoadConfigFromFile("config.ini:virtual")
}

func LoadConfigFromFile(fileName string) error {
	if TheConfig != nil {
		return nil
	}
	var config = &PlatformConfig{}

	err := ssg.ParseConfig(config, fileName)
	if err != nil {
		return err
	}

	if config.AccessTokenSigningKey != "" {
		AccessTokenSigningKey, err = base64.StdEncoding.DecodeString(config.AccessTokenSigningKey)
		if err != nil {
			return errors.New("failed to decode access token signing key:" + err.Error())
		}
	}

	if config.RefreshTokenSigningKey != "" {
		RefreshTokenSigningKey, err = base64.StdEncoding.DecodeString(config.RefreshTokenSigningKey)
		if err != nil {
			return errors.New("failed to decode refresh token signing key:" + err.Error())
		}
	}

	TheConfig = config

	return nil
}

func IsOwner(username, password string) bool {
	if TheConfig == nil {
		return false
	}

	// compare them as SHA256
	return strings.EqualFold(TheConfig.OwnerUsername, username) &&
		hashing.CompareSHA256(TheConfig.OwnerPassword, password)
}

func IsOwnerUsername(username string) bool {
	if TheConfig == nil {
		return false
	}

	return strings.EqualFold(TheConfig.OwnerUsername, username)
}

func GetCertFile() string {
	if TheConfig == nil {
		return ""
	}

	return TheConfig.CertFile
}

func GetCertKeyFile() string {
	if TheConfig == nil {
		return ""
	}

	return TheConfig.CertKeyFile
}

func GetDBUrl() string {
	if TheConfig == nil {
		return ""
	}

	return "host=" + TheConfig.PostgresHost +
		" user=" + TheConfig.PostgresUser +
		" password=" + TheConfig.PostgresPassword +
		" dbname=" + TheConfig.PostgresDb +
		" port=5432 sslmode=disable TimeZone=UTC"
}

func IsDebug() bool {
	if TheConfig == nil {
		return false
	}

	return TheConfig.Debug
}

func GetIPProxyHeader() string {
	if TheConfig == nil {
		return ""
	}

	return TheConfig.IPProxyHeader
}

func GetSwaggerInstanceName() string {
	if TheConfig == nil || TheConfig.SwaggerInstanceName == "" {
		return "ExamSphere Swagger Documentation"
	}

	return TheConfig.SwaggerInstanceName
}

func GetSwaggerTitle() string {
	if TheConfig == nil || TheConfig.SwaggerTitle == "" {
		return "ExamSphere Swagger Documentation"
	}

	return TheConfig.SwaggerTitle
}

func GetSwaggerBaseURL() string {
	if TheConfig == nil {
		return ""
	}

	return TheConfig.SwaggerBaseURL
}

func IsSudoToken(value string) bool {
	if TheConfig == nil || TheConfig.SudoToken == "" {
		return false
	}

	return TheConfig.SudoToken == value
}

func GetMaxRequestTillRateLimit() int {
	if TheConfig == nil {
		return 5
	}

	return TheConfig.MaxRequestTillRateLimit
}

func GetMaxRateLimitDuration() time.Duration {
	if TheConfig == nil {
		return 3 * time.Minute
	}

	return TheConfig.MaxRateLimitDuration * time.Minute
}

func GetRateLimitPunishmentDuration() time.Duration {
	if TheConfig == nil {
		return 5 * time.Minute
	}

	return TheConfig.RateLimitPunishmentDuration * time.Minute
}

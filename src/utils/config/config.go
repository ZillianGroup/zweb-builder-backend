package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/caarlos0/env"
)

const DEPLOY_MODE_SELF_HOST = "self-host"
const DEPLOY_MODE_CLOUD = "cloud"
const DEPLOY_MODE_CLOUD_TEST = "cloud-test"
const DEPLOY_MODE_CLOUD_BETA = "cloud-beta"
const DEPLOY_MODE_CLOUD_PRODUCTION = "cloud-production"
const DRIVE_TYPE_AWS = "aws"
const DRIVE_TYPE_DO = "do"
const DRIVE_TYPE_MINIO = "minio"
const PROTOCOL_WEBSOCKET = "ws"
const PROTOCOL_WEBSOCKET_OVER_TLS = "wss"

var instance *Config
var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		var err error
		if instance == nil {
			instance, err = getConfig() // not thread safe
			if err != nil {
				panic(err)
			}
		}
	})
	return instance
}

type Config struct {
	// server config
	ServerHost                    string `env:"ZWEB_SERVER_HOST"                    envDefault:"0.0.0.0"`
	ServerPort                    string `env:"ZWEB_SERVER_PORT"                    envDefault:"8001"`
	InternalServerPort            string `env:"ZWEB_SERVER_INTERNAL_PORT"           envDefault:"9005"`
	ServerMode                    string `env:"ZWEB_SERVER_MODE"                    envDefault:"debug"`
	DeployMode                    string `env:"ZWEB_DEPLOY_MODE"                    envDefault:"self-host"`
	SecretKey                     string `env:"ZWEB_SECRET_KEY"                     envDefault:"8xEMrWkBARcDDYQ"`
	WebsocketServerHost           string `env:"ZWEB_WEBSOCKET_SERVER_HOST"          envDefault:"0.0.0.0"`
	WebsocketServerPort           string `env:"ZWEB_WEBSOCKET_SERVER_PORT"          envDefault:"8002"`
	WebsocketServerConnectionHost string `env:"ZWEB_WEBSOCKET_CONNECTION_HOST"      envDefault:"0.0.0.0"`
	WebsocketServerConnectionPort string `env:"ZWEB_WEBSOCKET_CONNECTION_PORT"      envDefault:"80"`
	WSSEnabled                    string `env:"ZWEB_WSS_ENABLED"                    envDefault:"false"`

	// key for idconvertor
	RandomKey string `env:"ZWEB_RANDOM_KEY"  envDefault:"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"`
	// storage config
	PostgresAddr     string `env:"ZWEB_PG_ADDR" envDefault:"localhost"`
	PostgresPort     string `env:"ZWEB_PG_PORT" envDefault:"5432"`
	PostgresUser     string `env:"ZWEB_PG_USER" envDefault:"zweb_builder"`
	PostgresPassword string `env:"ZWEB_PG_PASSWORD" envDefault:"71De5JllWSetLYU"`
	PostgresDatabase string `env:"ZWEB_PG_DATABASE" envDefault:"zweb_builder"`
	// cache config
	RedisAddr     string `env:"ZWEB_REDIS_ADDR" envDefault:"localhost"`
	RedisPort     string `env:"ZWEB_REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"ZWEB_REDIS_PASSWORD" envDefault:"zweb2022"`
	RedisDatabase int    `env:"ZWEB_REDIS_DATABASE" envDefault:"0"`
	// drive config
	DriveType             string `env:"ZWEB_DRIVE_TYPE"               envDefault:""`
	DriveAccessKeyID      string `env:"ZWEB_DRIVE_ACCESS_KEY_ID"      envDefault:""`
	DriveAccessKeySecret  string `env:"ZWEB_DRIVE_ACCESS_KEY_SECRET"  envDefault:""`
	DriveRegion           string `env:"ZWEB_DRIVE_REGION"             envDefault:""`
	DriveEndpoint         string `env:"ZWEB_DRIVE_ENDPOINT"           envDefault:""`
	DriveSystemBucketName string `env:"ZWEB_DRIVE_SYSTEM_BUCKET_NAME" envDefault:"zweb-cloud"`
	DriveTeamBucketName   string `env:"ZWEB_DRIVE_TEAM_BUCKET_NAME"   envDefault:"zweb-cloud-team"`
	DriveUploadTimeoutRaw string `env:"ZWEB_DRIVE_UPLOAD_TIMEOUT"     envDefault:"30s"`
	DriveUploadTimeout    time.Duration
	// supervisor API
	ZwebSupervisorInternalRestAPI string `env:"ZWEB_SUPERVISOR_INTERNAL_API"     envDefault:"http://127.0.0.1:9001/api/v1"`

	// peripheral API
	ZwebPeripheralAPI string `env:"ZWEB_PERIPHERAL_API" envDefault:"https://peripheral-api.zwebsoft.com/v1/"`
	// resource manager API
	ZwebResourceManagerRestAPI         string `env:"ZWEB_RESOURCE_MANAGER_API"     envDefault:"http://zweb-resource-manager-backend:8006"`
	ZwebResourceManagerInternalRestAPI string `env:"ZWEB_RESOURCE_MANAGER_INTERNAL_API"     envDefault:"http://zweb-resource-manager-backend-internal:9004"`
	// zweb marketplace config
	ZwebMarketplaceInternalRestAPI string `env:"ZWEB_MARKETPLACE_INTERNAL_API"     envDefault:"http://zweb-marketplace-backend-internal:9003/api/v1"`
	// token for internal api
	ControlToken string `env:"ZWEB_CONTROL_TOKEN"     envDefault:""`
	// google config
	ZwebGoogleSheetsClientID     string `env:"ZWEB_GS_CLIENT_ID"           envDefault:""`
	ZwebGoogleSheetsClientSecret string `env:"ZWEB_GS_CLIENT_SECRET"       envDefault:""`
	ZwebGoogleSheetsRedirectURI  string `env:"ZWEB_GS_REDIRECT_URI"        envDefault:""`
}

func getConfig() (*Config, error) {
	// fetch
	cfg := &Config{}
	err := env.Parse(cfg)
	// process data
	var errInParseDuration error
	cfg.DriveUploadTimeout, errInParseDuration = time.ParseDuration(cfg.DriveUploadTimeoutRaw)
	if errInParseDuration != nil {
		return nil, errInParseDuration
	}
	// ok
	fmt.Printf("----------------\n")
	fmt.Printf("%+v\n", cfg)
	fmt.Printf("%+v\n", err)

	return cfg, err
}

func (c *Config) IsSelfHostMode() bool {
	return c.DeployMode == DEPLOY_MODE_SELF_HOST
}

func (c *Config) IsCloudMode() bool {
	if c.DeployMode == DEPLOY_MODE_CLOUD || c.DeployMode == DEPLOY_MODE_CLOUD_TEST || c.DeployMode == DEPLOY_MODE_CLOUD_BETA || c.DeployMode == DEPLOY_MODE_CLOUD_PRODUCTION {
		return true
	}
	return false
}

func (c *Config) IsCloudTestMode() bool {
	return c.DeployMode == DEPLOY_MODE_CLOUD_TEST
}

func (c *Config) IsCloudBetaMode() bool {
	return c.DeployMode == DEPLOY_MODE_CLOUD_BETA
}

func (c *Config) IsCloudProductionMode() bool {
	return c.DeployMode == DEPLOY_MODE_CLOUD_PRODUCTION
}

func (c *Config) GetWebScoketServerListenAddress() string {
	return c.WebsocketServerHost + ":" + c.WebsocketServerPort
}

func (c *Config) GetWebScoketServerConnectionAddress() string {
	return c.WebsocketServerConnectionHost + ":" + c.WebsocketServerConnectionPort
}

func (c *Config) GetWebsocketProtocol() string {
	if c.WSSEnabled == "true" {
		return PROTOCOL_WEBSOCKET_OVER_TLS
	}
	return PROTOCOL_WEBSOCKET
}

func (c *Config) GetRuntimeEnv() string {
	if c.IsCloudBetaMode() {
		return DEPLOY_MODE_CLOUD_BETA
	} else if c.IsCloudProductionMode() {
		return DEPLOY_MODE_CLOUD_PRODUCTION
	} else {
		return DEPLOY_MODE_CLOUD_TEST
	}
}

func (c *Config) GetSecretKey() string {
	return c.SecretKey
}

func (c *Config) GetRandomKey() string {
	return c.RandomKey
}

func (c *Config) GetPostgresAddr() string {
	return c.PostgresAddr
}

func (c *Config) GetPostgresPort() string {
	return c.PostgresPort
}

func (c *Config) GetPostgresUser() string {
	return c.PostgresUser
}

func (c *Config) GetPostgresPassword() string {
	return c.PostgresPassword
}

func (c *Config) GetPostgresDatabase() string {
	return c.PostgresDatabase
}

func (c *Config) GetRedisAddr() string {
	return c.RedisAddr
}

func (c *Config) GetRedisPort() string {
	return c.RedisPort
}

func (c *Config) GetRedisPassword() string {
	return c.RedisPassword
}

func (c *Config) GetRedisDatabase() int {
	return c.RedisDatabase
}

func (c *Config) GetDriveType() string {
	return c.DriveType
}

func (c *Config) IsAWSTypeDrive() bool {
	if c.DriveType == DRIVE_TYPE_AWS || c.DriveType == DRIVE_TYPE_DO {
		return true
	}
	return false
}

func (c *Config) IsMINIODrive() bool {
	return c.DriveType == DRIVE_TYPE_MINIO
}

func (c *Config) GetAWSS3Endpoint() string {
	return c.DriveEndpoint
}

func (c *Config) GetAWSS3AccessKeyID() string {
	return c.DriveAccessKeyID
}

func (c *Config) GetAWSS3AccessKeySecret() string {
	return c.DriveAccessKeySecret
}

func (c *Config) GetAWSS3Region() string {
	return c.DriveRegion
}

func (c *Config) GetAWSS3SystemBucketName() string {
	return c.DriveSystemBucketName
}

func (c *Config) GetAWSS3TeamBucketName() string {
	return c.DriveTeamBucketName
}

func (c *Config) GetAWSS3Timeout() time.Duration {
	return c.DriveUploadTimeout
}

func (c *Config) GetMINIOAccessKeyID() string {
	return c.DriveAccessKeyID
}

func (c *Config) GetMINIOAccessKeySecret() string {
	return c.DriveAccessKeySecret
}

func (c *Config) GetMINIOEndpoint() string {
	return c.DriveEndpoint
}

func (c *Config) GetMINIOSystemBucketName() string {
	return c.DriveSystemBucketName
}

func (c *Config) GetMINIOTeamBucketName() string {
	return c.DriveTeamBucketName
}

func (c *Config) GetMINIOTimeout() time.Duration {
	return c.DriveUploadTimeout
}

func (c *Config) GetControlToken() string {
	return c.ControlToken
}

func (c *Config) GetZwebSupervisorInternalRestAPI() string {
	return c.ZwebSupervisorInternalRestAPI
}

func (c *Config) GetZwebPeripheralAPI() string {
	return c.ZwebPeripheralAPI
}

func (c *Config) GetZwebResourceManagerRestAPI() string {
	return c.ZwebResourceManagerRestAPI
}

func (c *Config) GetZwebResourceManagerInternalRestAPI() string {
	return c.ZwebResourceManagerInternalRestAPI
}

func (c *Config) GetZwebMarketplaceInternalRestAPI() string {
	return c.ZwebMarketplaceInternalRestAPI
}

func (c *Config) GetZwebGoogleSheetsClientID() string {
	return c.ZwebGoogleSheetsClientID
}

func (c *Config) GetZwebGoogleSheetsClientSecret() string {
	return c.ZwebGoogleSheetsClientSecret
}

func (c *Config) GetZwebGoogleSheetsRedirectURI() string {
	return c.ZwebGoogleSheetsRedirectURI
}

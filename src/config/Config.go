package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

var GRPCTimeOut int = 10
var GlobalConfig Config
var VoiceMatchThreshold float32 = 0.6

const MaxVoiceLength = 6 // secconds
const UserTokenMaxUsageCount = 1000
const UserTokenMaxVoiceLoginCount = 3
const UserTokenMaxFaceLoginCount = 3

const UserAccountLockTimeout = 1

const MaxVectorRecordsPerUser = 10

const FaceMultipleSamplesCount = 5 // must be odd, n-1 login samples AND 1 mfa sample

const DateLayout = "02-Jan-06 15:04:05"
const QueryDateLayout = "02-01-2006"

func GetConfigStruct(cfg *Config) {
	readConfigFile(cfg)
	readConfigEnv(cfg)
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readConfigFile(cfg *Config) {
	fmt.Print("Parsing Config from config.yml ...\n")
	f, err := os.Open("config.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
	fmt.Printf("Parsed, current config %+v\n", cfg)
}

func readConfigEnv(cfg *Config) {
	fmt.Print("Parsing Config from environment variables (if exists)...")
	err := envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
	fmt.Printf("Parsed, current config %+v\n", cfg)
}

type Config struct {
	Version string `yaml:"version"`
	HTTP    struct {
		Port string `yaml:"port" envconfig:"HTTP_PORT"`
		Cors bool   `yaml:"cors" envconfig:"HTTP_CORS"`
	} `yaml:"http"`
	DB struct {
		Host     string `yaml:"host" envconfig:"DB_HOST"`
		Port     string `yaml:"port" envconfig:"DB_PORT"`
		Username string `yaml:"username" envconfig:"DB_USER"`
		Password string `yaml:"password" envconfig:"DB_PASSWORD"`
		Database string `yaml:"database" envconfig:"DB_DATABASE"`
	} `yaml:"db"`
	Minio struct {
		Address    string `yaml:"address" envconfig:"MINIO_ADDRESS"`
		AccessKey  string `yaml:"accesskey" envconfig:"MINIO_ACCESS_KEY"`
		SecretKey  string `yaml:"secretkey" envconfig:"MINIO_SECRET_KEY"`
		BucketName string `yaml:"bucketname" envconfig:"MINIO_BUCKET_NAME"`
		SSL        bool   `yaml:"ssl" envconfig:"MINIO_SSL"`
	} `yaml:"minio"`
}

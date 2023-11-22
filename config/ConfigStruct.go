package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/tkanos/gonfig"
	"os"
)

var ApplicationConfiguration Configuration

type Configuration interface {
	GetServer() Server
	GetPostgresql() Postgresql
	GetLogFile() []string
}

type Server struct {
	Protocol     string
	Host         string
	Port         string
	Version      string
	PrefixPath   string
	Application  string
	SignatureKey string
}

type Postgresql struct {
	Driver            string
	Address           string
	DefaultSchema     string
	MaxOpenConnection int
	MaxIdleConnection int
}

func GenerateConfiguration(argument string) {
	var (
		err      error
		enviName string
	)

	if argument == "local" {
		err = godotenv.Load()
		if err != nil {
			fmt.Println(`Error load file env -> `, err.Error())
			os.Exit(2)
		}

		enviName = os.Getenv("ENVDEV")
		temp := LocalConfig{}
		err = gonfig.GetConf(enviName+"/config_local.json", &temp)
		if err != nil {
			fmt.Println(`Error get config local -> `, err)
			os.Exit(2)
		}

		err = envconfig.Process(enviName+"/config_local.json", &temp)
		if err != nil {
			fmt.Println(`Error get param config local -> `, err)
			os.Exit(2)
		}
		ApplicationConfiguration = &temp
	}
}

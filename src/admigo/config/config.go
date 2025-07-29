package config

import (
	"encoding/json"
	"os"

	"admigo/applog"
	"admigo/common"
)

type dbConfig struct {
	Host     string `json:"host"`
	Driver   string `json:"driver"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
	Sslmode  string `json:"sslmode"`
}

type mailConfig struct {
	From     string `json:"from"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	GotoURL  string `json:"gotourl"`
}

// Configuration type for config of admigo
type Configuration struct {
	Appname  string      `json:"appname"`
	Address  string      `json:"address"`
	Port     int         `json:"port"`
	Static   string      `json:"static"`
	Acme     bool        `json:"acme"`
	Acmehost []string    `json:"acmehost"`
	DirCache string      `json:"dirCache"`
	Crt      string      `json:"crt,omitempty"`
	Key      string      `json:"key,omitempty"`
	Lang     string      `json:"lang"`
	Db       *dbConfig   `json:"db"`
	Mail     *mailConfig `json:"mail"`
}

var (
	config   *Configuration
	csrf_key string
)

// LoadConfig loads config
func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		applog.Danger("Cannot open config file", err)
	}

	decoder := json.NewDecoder(file)
	config = &Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		applog.Danger("Cannot get configuration from file", err)
	}
}

// Env returns config
func Env(reload bool) *Configuration {
	if reload {
		LoadConfig()
	}
	return config
}

func SetCsrf() {
	csrf_key = common.CreateUID()
}

func GetKeyCSRF() string {
	return csrf_key
}

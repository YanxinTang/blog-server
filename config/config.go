package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ConfigStruct struct {
	Site struct {
		Name string `json:"name"`
	}
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbName"`
	} `json:"database"`
}

var Config ConfigStruct

var SigninKey = []byte("blog")

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	configFilePath := filepath.Join(filepath.Dir(ex), "./config/config.json")
	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Load config error %v:", err)
	}
	if err := json.Unmarshal(data, &Config); err != nil {
		log.Fatalf("Read config error %v:", err)
	}
}

package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ConfigStruct struct {
	Site struct {
		Name string `json:"name"`
	}
	Auth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"auth"`
	Mysql struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
	} `json:"mysql"`
}

// Claims creates a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserID uint64 `json:"userID"`
	jwt.StandardClaims
}

const (
	TokenPrefix            = "Bearer "
	TokenPrefixLength      = len(TokenPrefix)
	RefreshTokenExpiration = 7 * 24 * time.Hour
)

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

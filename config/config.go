package config

import (
	"github.com/BurntSushi/toml"
	"log"
)




type GRPC struct {
	BaseURL    string `toml:"BaseURL"`
	PrivateKey string `toml:"PrivateKey"`
	Mysql mysql
}

type mysql struct {
	URI string `toml:"uri"`
}

var GRPC1 = new(GRPC)

func GetConfig() *GRPC {
	return GRPC1
}

func InitConfig(p string) {
	// p := "./config.toml"
	if _, err := toml.DecodeFile(p, GRPC1); err != nil {
		log.Fatalf("获取配置文件失败， %s", err.Error())
	}
}

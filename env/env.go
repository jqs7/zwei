package env

import (
	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	Debug    bool
	Token    string
	Address  string `default:"127.0.0.1"`
	Port     string `default:"5432"`
	User     string `default:"jqs7"`
	Database string `default:"zwei"`
}

var Spec Specification

func Init() error {
	err := envconfig.Process("zwei", &Spec)
	return err
}

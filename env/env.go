package env

import (
	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	Debug    bool `default:"false"`
	Token    string
	Address  string `default:"127.0.0.1"`
	Port     string `default:"5432"`
	User     string `default:"jqs7"`
	Password string
	Database string `default:"zwei"`
}

var Spec Specification

func Init() {
	envconfig.MustProcess("zwei", &Spec)
}

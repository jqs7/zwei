module github.com/jqs7/zwei

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-pg/migrations v6.2.0+incompatible
	github.com/go-pg/pg v6.15.0+incompatible
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/hanguofeng/gocaptcha v1.0.7
	github.com/json-iterator/go v1.1.5
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.8.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.2.2 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20181012144002-a92615f3c490 // indirect
)

replace github.com/go-pg/pg v6.15.0+incompatible => github.com/jqs7/pg v0.0.0-20181014041559-1b1319d49317

replace github.com/hanguofeng/gocaptcha v1.0.7 => github.com/jqs7/gocaptcha v1.0.8-0.20181014100812-c7bcbe23fde4

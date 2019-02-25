module github.com/jqs7/zwei

require (
	github.com/go-pg/migrations v6.7.2+incompatible
	github.com/go-pg/pg v6.15.1+incompatible
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/google/wire v0.2.0
	github.com/hanguofeng/gocaptcha v1.0.7
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/onsi/gomega v1.4.2 // indirect
	github.com/skip2/go-qrcode v0.0.0-20171229120447-cf5f9fa2f0d8
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/hanguofeng/gocaptcha v1.0.7 => github.com/jqs7/gocaptcha v1.0.8-0.20181014100812-c7bcbe23fde4

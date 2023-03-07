module github.com/jqs7/zwei

require (
	github.com/go-pg/migrations v6.3.0+incompatible
	github.com/go-pg/pg v6.15.1+incompatible
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/google/wire v0.2.1
	github.com/hanguofeng/gocaptcha v1.0.7
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/skip2/go-qrcode v0.0.0-20171229120447-cf5f9fa2f0d8
)

require (
	github.com/bradfitz/gomemcache v0.0.0-20180710155616-bc664df96737 // indirect
	github.com/hanguofeng/config v1.0.0 // indirect
	github.com/hanguofeng/freetype-go-mirror v0.0.0-20140928112427-cfb10e2cb6de // indirect
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/onsi/ginkgo v1.6.0 // indirect
	github.com/onsi/gomega v1.4.2 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20190320223903-b7391e95e576 // indirect
	gopkg.in/bufio.v1 v1.0.0-20140618132640-567b2bfa514e // indirect
	gopkg.in/redis.v2 v2.3.2 // indirect
	mellium.im/sasl v0.3.1 // indirect
)

replace github.com/hanguofeng/gocaptcha v1.0.7 => github.com/jqs7/gocaptcha v1.0.8-0.20181014100812-c7bcbe23fde4

module github.com/kalaGN/gincron

go 1.16

require (
	github.com/garyburd/redigo v1.6.2
	github.com/go-ini/ini v1.62.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/json-iterator/go v1.1.9
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/kalaGN/gincron/src/common => ./src/common

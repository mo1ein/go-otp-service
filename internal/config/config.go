package config

type Config struct {
	HTTP     HTTP
	Database Database
}

type HTTP struct {
	APIHost string
	APIPort int
}

type Database struct {
	Postgres Postgres
	Redis    Redis
}

type Postgres struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	Database int
}

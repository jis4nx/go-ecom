package config

type Database struct {
	DBNAME string
	DBUSER string
	DBHOST string
	DBPORT string
	DBPASS string
}

type Server struct {
	HOST string
	PORT string
}

type RabbitMQ struct {
	USER     string
	PASSWORD string
	HOST     string
	VHOST    string
}

type Services struct {
	ProductServer Server
}

type Config struct {
	DB       Database
	RQ       RabbitMQ
	Services Services
}

func LoadServer(env map[string]string, defaultPort string) Server {
	port := env["PRODUCT_PORT"]

	if port == "" {
		port = "4001"
	}
	productServer := Server{PORT: port}
	return productServer
}

func GetProductServer(env map[string]string) Server {
	productServer := LoadServer(env, "4001")
	return productServer
}

func LoadConfig(env map[string]string) Config {
	cfg := Config{
		DB: Database{
			DBNAME: env["DBNAME"],
			DBUSER: env["DBUSER"],
			DBHOST: env["DBHOST"],
			DBPORT: env["DBPORT"],
			DBPASS: env["DBPASS"],
		},
		Services: Services{
			ProductServer: GetProductServer(env),
		},
		RQ: RabbitMQ{
			USER:     env["RABBIT_USER"],
			PASSWORD: env["RABBIT_PASS"],
			VHOST:    env["RABBIT_VHOST"],
			HOST:     env["RABBIT_HOST"],
		},
	}
	return cfg
}

// Setting Config to access Globally
var globalConfig *Config

func SetGlobalConfig(c *Config) {
	globalConfig = c
}

func GetConfig() *Config {
	return globalConfig
}

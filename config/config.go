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

func LoadServer(env map[string]string, defautHost, defaultPort string) Server {
	host := env["PRODUCT_HOST"]
	port := env["PRODUCT_PORT"]

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "4001"
	}
	productServer := Server{HOST: host, PORT: port}
	return productServer
}

func GetProductServer(env map[string]string) Server {
	productServer := LoadServer(env, "localhost", "4001")
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

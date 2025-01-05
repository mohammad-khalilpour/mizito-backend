package env

type Config struct {
	MongoCollection  string
	MongoDatabase    string
	AppPort          string `envDefault:":8080"`
	RedisProjectsDB  string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPass     string
	PostgresDatabase string
}

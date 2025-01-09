package env

type Config struct {
	MongoCollection  string `envDefault:"messages"`
	MongoDatabase    string `envDefault:"mizito"`
	AppPort          string `envDefault:":8080"`
	RedisHost        string `envDefault:"localhost"`
	RedisPort        string `envDefault:"6379"`
	RedisUsername    string
	RedisPassword    string
	RedisProjectsDB  string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPass     string
	PostgresDatabase string
	MongoDBHost      string `envDefault:"mongodb://localhost:27017"`
}

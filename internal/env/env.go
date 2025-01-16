package env

type Config struct {
	MongoCollection     string `envDefault:"messages"`
	MongoDatabase       string `envDefault:"mizito"`
	AppPort             string `envDefault:":8080"`
	AuthorizationSecret string `evnDefault:"testing12345678910"`
	RedisHost           string `envDefault:"localhost"`
	RedisPort           string `envDefault:"6379"`
	RedisUsername       string
	RedisPassword       string
	RedisProjectsDB     string
	PostgresHost        string `envDefault:"localhost"`
	PostgresPort        string `envDefault:"5432"`
	PostgresUser        string `envDefault:"postgres"`
	PostgresPass        string `envDefault:"postgres"`
	PostgresDatabase    string `envDefault:"postgres"`
	MongoDBHost         string `envDefault:"mongodb://localhost:27017"`
}

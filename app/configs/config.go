package configs

type Configuration struct {
	App     ApplicationConfiguration
	MongoDb DatabaseConfiguration
	Redis   RedisConfiguration
}

type ApplicationConfiguration struct {
	Env  string
	Port string
}

type DatabaseConfiguration struct {
	Connection   string
	SessionName  string
	DatabaseName string
}

type RedisConfiguration struct {
	Connection  string
	SessionName SessionConfiguration
}

type SessionConfiguration struct {
	UserSession string
}

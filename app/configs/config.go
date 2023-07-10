package configs

type Configuration struct {
	App     ApplicationConfiguration
	MongoDb DatabaseConfiguration
}

type ApplicationConfiguration struct {
	Env  string
	Port string
}

type DatabaseConfiguration struct {
	Connection   string
	Sessionname  string
	Databasename string
}

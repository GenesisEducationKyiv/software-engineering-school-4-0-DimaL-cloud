package configs

type RabbitMQ struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Queue    struct {
		Customer string `yaml:"customer"`
	} `yaml:"queue"`
}

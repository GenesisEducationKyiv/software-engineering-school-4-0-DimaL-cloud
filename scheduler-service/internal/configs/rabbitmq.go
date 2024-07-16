package configs

type RabbitMQ struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Queue    struct {
		RateNotificationCron string `yaml:"rate-notification-cron"`
	} `yaml:"queue"`
}

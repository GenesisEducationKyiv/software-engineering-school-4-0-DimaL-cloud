package configs

type DB struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	DBName     string `yaml:"name"`
	SSLMode    string `yaml:"ssl_mode"`
	DriverName string `yaml:"driver_name"`
	SearchPath string `yaml:"search_path"`
}

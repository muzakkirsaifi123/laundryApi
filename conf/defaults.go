package conf

import (
	config "github.com/spf13/viper"
	"log"
)

func SetDefaults() {

	config.SetConfigFile("./application.yaml")
	err := config.ReadInConfig()
	config.SetDefault("pidFile", "./pid")
	config.SetDefault("version", "0.1")
	config.SetDefault("server_hostname", "localhost")
	config.SetDefault("server_port", "9090")
	config.SetDefault("database_name", "postgres")
	config.SetDefault("database_hostname", "localhost")
	config.SetDefault("database_port", "5432")
	config.SetDefault("database_username", "")
	config.SetDefault("database_password", "Jeanalevante9423")

	if err != nil {
		err = config.WriteConfigAs("application.yaml")

		if err != nil {
			log.Fatal(err)
		}
	}

}

package cmd

import (
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/jchenriquez/laundromat/store"
	"github.com/jchenriquez/laundromat/store/queries"
	cli "github.com/spf13/cobra"
	config "github.com/spf13/viper"

	"github.com/jchenriquez/laundromat/controllers"
	"github.com/jchenriquez/laundromat/server"
)

func ClientDB() (*store.Client, error) {
	dbName := config.Get("database_name").(string)
	dbHostname := config.GetString("database_hostname")
	dbPort := config.GetString("database_port")
	dbUsername := config.GetString("database_username")
	dbPassword := config.GetString("database_password")

	port, err := strconv.Atoi(dbPort)

	if err != nil {
		return nil, err
	}

	return store.New(dbUsername, dbHostname, dbPassword, dbName, port), nil
}

func addStartCommand() {
	rootCmd.AddCommand(
		&cli.Command{
			Use:   "start",
			Short: "Starts Server",
			Long:  `Starts the server by discovering config [.yaml] file from project directory, if it can't find one it will create config file with defaults with name application.yml'`,
			Args:  cli.NoArgs,
			Run: func(cmd *cli.Command, args []string) {
				errorChan := make(chan os.Signal, 1)
				hostName := config.GetString("server_hostname")
				port := config.GetString("server_port")
				mServer := server.New(hostName, port)
				db, err := ClientDB()

				if err != nil {
					log.Fatal(err)
					return
				}

				controllers.AddControllers(mServer.Router, db)

				if err != nil {
					log.Fatal(err)
				}
				go func() {
					log.Fatal(mServer.ListenAndServe())
				}()

				signal.Notify(errorChan, os.Interrupt, os.Kill)

				<-errorChan

				err = queries.DeleteSessions(db)
				if err != nil {
					log.Fatal(err)
				}
				log.Fatal(mServer.Stop())
			},
		},
	)

}

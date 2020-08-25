package cmd

import (
	"github.com/jchenriquez/laundromat/apis"
	"github.com/jchenriquez/laundromat/server"
	cli "github.com/spf13/cobra"
	config "github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
)

func addStartCommand() {
	rootCmd.AddCommand(&cli.Command{
		Use:   "start",
		Short: "Starts Server",
		Long:  `Starts the server by discovering config [.yaml] file from project directory, if it can't find one it will create config file with defaults with name application.yaml'`,
		Args:  cli.NoArgs,
		Run: func(cmd *cli.Command, args []string) {
			errorChan := make(chan os.Signal, 1)
			hostName := config.GetString("server_hostname")
			port := config.GetString("server_port")
			mServer := server.New(hostName, port)
			apis.AddApis(mServer.Router)

			go func() {
				log.Fatal(mServer.ListenAndServe())
			}()

			signal.Notify(errorChan, os.Interrupt, os.Kill)

			<-errorChan

			log.Fatal(mServer.Stop())
		},
	})

}

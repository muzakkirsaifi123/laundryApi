package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	config "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	PreRunE: func(cmd *cobra.Command, args []string) error {
		pidFile := config.GetString("pidFile")

		if pidFile != "" {
			file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_EXCL|os.O_TRUNC|os.O_WRONLY, 0666)

			if err != nil {
				return err
			}

			_, err = fmt.Fprint(file, os.Getegid())

			return err

		}

		return nil
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		pidFile := config.GetString("pidFile")

		if pidFile != "" {
			return os.Remove(pidFile)
		}

		return nil
	},
}

func Execute() {
	addStartCommand()
	addVersionCommand()
	log.Fatal(rootCmd.Execute())
}

package cmd

import (
	"fmt"
	cli "github.com/spf13/cobra"
	config "github.com/spf13/viper"
)

func addVersionCommand() {
	rootCmd.AddCommand(&cli.Command{
		Use:   "version",
		Short: "Prints app version",
		Args:  cli.NoArgs,
		Run: func(cmd *cli.Command, args []string) {
			fmt.Printf("version %s", config.GetString("version"))
		},
	})
}

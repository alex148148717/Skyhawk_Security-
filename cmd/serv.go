/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"skyhawk/internal/player_logs"

	"github.com/spf13/cobra"
)

// servCmdCmd represents the playerstats command
var servCmd = &cobra.Command{
	Use:   "serv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("servCmd called")

		player_logs.NewApp().Run()
	},
}

func init() {
	rootCmd.AddCommand(servCmd)

}

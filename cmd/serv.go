/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"skyhawk/internal/player_logs"

	"github.com/spf13/cobra"
)

type ServConfig struct {
	Port           int
	DSN            string
	Region         string
	Endpoint       string
	CacheTableName string
	IsDAX          bool
	DaxHostPorts   string
}

var servConfig ServConfig

// servCmdCmd represents the playerstats command
var servCmd = &cobra.Command{
	Use:   "serv",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("servCmd called")

		ctx := cmd.Context()

		return player_logs.New(ctx, servConfig.Port, servConfig.DSN, servConfig.Region, servConfig.Endpoint, servConfig.CacheTableName, servConfig.IsDAX, []string{servConfig.DaxHostPorts})

	},
}

func init() {
	rootCmd.AddCommand(servCmd)
	servCmd.PersistentFlags().IntVar(&servConfig.Port, "port", 8081, "port ...")
	servCmd.PersistentFlags().StringVar(&servConfig.DSN, "dsn", "host=localhost user=user password=pass dbname=player_stats port=5444 sslmode=disable TimeZone=UTC", "for local : host=localhost user=user password=pass dbname=player_stats port=5444 sslmode=disable TimeZone=UTC")
	servCmd.PersistentFlags().StringVar(&servConfig.Region, "region", "us-west-2", "amazon Region")
	servCmd.PersistentFlags().StringVar(&servConfig.Endpoint, "endpoint", "http://localhost:8000", "for production pass empty")
	servCmd.PersistentFlags().StringVar(&servConfig.CacheTableName, "cacheTableName", "cache", "aws dynamodb table ")
	servCmd.PersistentFlags().BoolVar(&servConfig.IsDAX, "dax", false, "aws DAX over dynamodb ")
	servCmd.PersistentFlags().StringVar(&servConfig.DaxHostPorts, "daxhp", "", ",aws DAX over dynamodb ")

}

package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/spf13/cobra"
	"log"
)

type MigrationsCmdConfig struct {
	Region   string
	Endpoint string
}

var migrationsCmdConfig MigrationsCmdConfig

// migrationsCmd represents the playerstats command
var migrationsCmd = &cobra.Command{
	Use:   "migrations",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return startMigrations(cmd.Context(), migrationsCmdConfig)
	},
}

func init() {
	rootCmd.AddCommand(migrationsCmd)

	migrationsCmd.PersistentFlags().StringVar(&migrationsCmdConfig.Region, "region", "us-east-1", "example us-east-1")
	migrationsCmd.PersistentFlags().StringVar(&migrationsCmdConfig.Endpoint, "endpoint", "", "example us-east-1")

}

func ensurePagesTableExists(svc *dynamodb.DynamoDB, tableName string) error {
	_, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			fmt.Println("🔧 Table doesn't exist, creating:", tableName)
			input := &dynamodb.CreateTableInput{
				TableName: aws.String(tableName),
				AttributeDefinitions: []*dynamodb.AttributeDefinition{
					{
						AttributeName: aws.String("id"),
						AttributeType: aws.String("S"),
					},
				},
				KeySchema: []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("id"),
						KeyType:       aws.String("HASH"),
					},
				},
				BillingMode: aws.String("PAY_PER_REQUEST"),
			}

			_, err := svc.CreateTable(input)
			if err != nil {
				return fmt.Errorf("❌ failed to create table: %w", err)
			}

			fmt.Println("✅ Table created:", tableName)
			return nil
		}
		return fmt.Errorf("❌ failed to describe table: %w", err)
	}
	fmt.Println("✅ Table already exists:", tableName)
	return nil
}
func startMigrations(ctx context.Context, migrationsCmdConfig MigrationsCmdConfig) error {

	cfg := &aws.Config{
		Region: aws.String(migrationsCmdConfig.Region),
	}

	if migrationsCmdConfig.Endpoint != "" {
		cfg.Endpoint = aws.String(migrationsCmdConfig.Endpoint)
		cfg.Credentials = credentials.NewStaticCredentials("fake", "fake", "")
	}

	sess := session.Must(session.NewSession(cfg))

	svc := dynamodb.New(sess)

	err := ensurePagesTableExists(svc, "cache")
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	return err
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"testing"
)

func TestMigrations(t *testing.T) {
	ctx := context.Background()
	migrationsCmdConfig := MigrationsCmdConfig{
		Endpoint: "http://localhost:8000",
		Region:   "us-west-2",
	}
	err := startMigrations(ctx, migrationsCmdConfig)
	if err != nil {
		t.Error(err)
	}
}

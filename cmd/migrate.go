package cmd

import (
	"github.com/oakeshq/go-starter/config"
	"github.com/oakeshq/go-starter/pkg/storage/dialer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var migrateCmd = &cobra.Command{
	Use:  "migrate",
	Long: "Migrate database strucutures. This will create new tables and add missing columns and indexes.",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadConfig(cmd)
		if err != nil {
			zap.L().Sugar().Fatalf("Unable to read config %v", err)
		}

		db, err := dialer.Connect(c)
		if err != nil {
			zap.L().Sugar().Fatalf("Error opening database: %+v", err)
		}

		if err = dialer.Migrate(db); err != nil {
			zap.L().Sugar().Fatalf("Error migrating database:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}

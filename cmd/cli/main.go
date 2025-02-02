package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfg *Config

var rootCmd = &cobra.Command{
	Use:   "hamster",
	Short: "Hamster is a media manager",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cliCfg, err := LoadConfig("hamster-cli.toml")
		if err != nil {
			return fmt.Errorf("failed to load hamster-cli.toml: %w", err)
		}

		cfg = cliCfg
		return nil
	},
}

func init() {
	rootCmd.AddCommand(libraryCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

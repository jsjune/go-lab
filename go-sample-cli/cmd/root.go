package cmd

import (
	"fmt"
	"os"

	"go-sample-cli/store"

	"github.com/spf13/cobra"
)

var globalStore *store.Store

var rootCmd = &cobra.Command{
	Use:   "item",
	Short: "아이템 관리 CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		s, err := store.New()
		if err != nil {
			return err
		}
		globalStore = s
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(deleteCmd)
}

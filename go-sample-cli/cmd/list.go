package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "아이템 전체 목록 출력",
	RunE: func(cmd *cobra.Command, args []string) error {
		items := globalStore.List()
		if len(items) == 0 {
			fmt.Println("아이템이 없습니다.")
			return nil
		}
		fmt.Fprintf(os.Stdout, "%-5s %-20s %s\n", "ID", "NAME", "DESCRIPTION")
		fmt.Println("----------------------------------------------")
		for _, item := range items {
			fmt.Fprintf(os.Stdout, "%-5d %-20s %s\n", item.ID, item.Name, item.Description)
		}
		return nil
	},
}

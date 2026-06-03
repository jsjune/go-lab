package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	createName        string
	createDescription string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "아이템 생성",
	Example: `  item create --name "apple" --description "red fruit"
  item create -n "banana"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if createName == "" {
			return fmt.Errorf("--name 은 필수입니다")
		}
		item, err := globalStore.Create(createName, createDescription)
		if err != nil {
			return err
		}
		fmt.Printf("생성 완료 — ID: %d, Name: %s\n", item.ID, item.Name)
		return nil
	},
}

func init() {
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "아이템 이름 (필수)")
	createCmd.Flags().StringVarP(&createDescription, "description", "d", "", "아이템 설명")
}

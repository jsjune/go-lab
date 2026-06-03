package cmd

import (
	"errors"
	"fmt"
	"go-sample-cli/store"
	"strconv"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "get <id>",
	Short:   "아이템 단건 조회",
	Example: "  item get 1",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("id는 숫자여야 합니다")
		}
		item, err := globalStore.Get(id)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("id %d 아이템을 찾을 수 없습니다", id)
		}
		if err != nil {
			return err
		}
		fmt.Printf("ID:          %d\n", item.ID)
		fmt.Printf("Name:        %s\n", item.Name)
		fmt.Printf("Description: %s\n", item.Description)
		fmt.Printf("CreatedAt:   %s\n", item.CreatedAt.Format("2006-01-02 15:04:05"))
		return nil
	},
}

package cmd

import (
	"errors"
	"fmt"
	"go-sample-cli/store"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <id>",
	Short:   "아이템 삭제",
	Example: "  item delete 1",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("id는 숫자여야 합니다")
		}
		if err := globalStore.Delete(id); errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("id %d 아이템을 찾을 수 없습니다", id)
		} else if err != nil {
			return err
		}
		fmt.Printf("id %d 삭제 완료\n", id)
		return nil
	},
}

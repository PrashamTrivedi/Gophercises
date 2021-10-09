package taskcmd

import (
	"fmt"
	"os"
	"strings"
	"taskdb"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")

		_, err := taskdb.StoreTask(task)
		if err != nil {
			fmt.Println("Error in adding tasks", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Task '%s' is added to your list\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

}

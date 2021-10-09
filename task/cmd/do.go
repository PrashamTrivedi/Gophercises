package taskcmd

import (
	"fmt"
	"os"
	"strconv"
	"taskdb"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task as Done",

	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse: ", arg)
			} else {

				ids = append(ids, id)
			}

		}
		tasks, err := taskdb.ListTasks()
		if err != nil {
			fmt.Println("Something went wrong", err.Error())
			os.Exit(1)
		}
		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid task number")
				continue
			}
			task := tasks[id-1]
			err := taskdb.DeleteTask(task.Key)
			if err != nil {
				fmt.Printf("Faild to mark \"%d\" completed, Error: %s\n", id, err)
			} else {
				fmt.Printf("Task \"%s\" completed\n", task.Value)
			}

		}

	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}

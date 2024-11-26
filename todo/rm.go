package todo

import (
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var cmdRemove = &cobra.Command{
	Use:     "rm [TODO]",
	Short:   "Remove one todo from the list of todos",
	Long:    "Removes selected incomplete todo from the list of todos. The selected todo to remove must correspond to its list number displayed by the 'list' command. The argument must be a number and cannot be less than 1, or more than the amount of incompleted todos.",
	Example: "todo rm 1",
	Args:    cobra.ExactArgs(1),
	Run:     remove,
}

func remove(cmd *cobra.Command, args []string) {
	db := getDB(cmd)
	defer db.Close()

	todo, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		cmd.PrintErrf("Argument \"%s\" is not valid. See todo rm --help for more info\n", args[0])
		os.Exit(1)
	}

	if err = db.Update(removeTodo(int(todo))); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}

func removeTodo(numTodo int) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		if _, _, err := removeTodoByIndex(numTodo, tx); err != nil {
			return err
		}
		return nil
	}
}

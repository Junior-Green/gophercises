package todo

import (
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var cmdDo = &cobra.Command{
	Use:     "do [TODO]",
	Short:   "Mark selected todo as completed",
	Long:    "Marks the selected todo as completed. The selected todo to mark as completed must correspond to its list number displayed by the 'list' command. The argument must be a number and cannot be less than 1, or more than the amount of incompleted todos.",
	Example: "todo do 1",
	Args:    cobra.ExactArgs(1),
	Run:     do,
}

func do(cmd *cobra.Command, args []string) {
	db := getDB(cmd)
	defer db.Close()

	todo, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		cmd.PrintErrf("Argument \"%s\" is not valid. See todo do --help for more info\n", args[0])
		os.Exit(1)
	}

	if err = db.Update(completeTodo(int(todo))); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}

func completeTodo(numTodo int) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		var (
			todoKey []byte
			todoVal []byte
			err     error
		)

		todoKey, todoVal, err = removeTodoByIndex(numTodo, tx)
		if err != nil {
			return err
		}

		completed, err := tx.CreateBucketIfNotExists([]byte(completedBucket))
		if err != nil {
			return err
		}

		if err = completed.Put(todoKey, todoVal); err != nil {
			return err
		}

		return nil
	}
}

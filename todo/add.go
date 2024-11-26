package todo

import (
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var cmdAdd = &cobra.Command{
	Use:     "add [TODO]",
	Short:   "Add a task to your list of todos",
	Long:    "Add one todo to your list. If it already exists it this command will be ignored. If the todo has any spaces wrap it in qoutes (\"\").",
	Example: "todo add \"clean dishes\"",
	Args:    cobra.ExactArgs(1),
	Run:     add,
}

func add(cmd *cobra.Command, args []string) {
	db := getDB(cmd)
	defer db.Close()

	if err := db.Update(addTodo(args[0])); err != nil {
		cmd.PrintErrln("Error occured while saving todo")
		os.Exit(1)
	}
}

func addTodo(todo string) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(todoBucket))
		if err != nil {
			return err
		}

		b.Put([]byte(time.Now().Format(time.RFC3339)), []byte(todo))

		return nil
	}
}

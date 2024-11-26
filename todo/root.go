package todo

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var dbPath string

const fileMode fs.FileMode = 0600
const completedBucket string = "complete"
const todoBucket string = "incomplete"

var cmdRoot = &cobra.Command{
	Use:   "todo [COMMAND]",
	Short: "Command line todo list",
	Long:  "CLI program that lets you manage and organize tasks that persists on disk",
}

type invalidIndexError struct {
	message string
}

func (err invalidIndexError) Error() string {
	return err.message
}

func Execute() error {
	return cmdRoot.Execute()
}

func init() {
	dir, err := homedir.Dir()
	
	if err != nil {
		fmt.Fprintln(os.Stderr, "No home directory detected. Aborting")
		os.Exit(1)
	}

	dbPath = fmt.Sprintf("%s/todo.db", dir)

	cmdRoot.AddCommand(cmdAdd)
	cmdRoot.AddCommand(cmdDo)
	cmdRoot.AddCommand(cmdList)
	cmdRoot.AddCommand(cmdCompleted)
	cmdRoot.AddCommand(cmdRemove)
}

func getDB(cmd *cobra.Command) *bolt.DB {
	db, err := bolt.Open(dbPath, fileMode, nil)
	if err != nil {
		cmd.PrintErr("Cannot open database")
		os.Exit(1)
	}

	return db
}

func removeTodoByIndex(numTodo int, tx *bolt.Tx) ([]byte, []byte, error) {
	var (
		key []byte
		val []byte
	)
	todos, err := tx.CreateBucketIfNotExists([]byte(todoBucket))

	if err != nil {
		return nil, nil, err
	}

	cursor, i := todos.Cursor(), 1

	for k, v := cursor.First(); k != nil; k, _ = cursor.Next() {
		if i == numTodo {
			key, val = k, v
			break
		}
	}

	if key == nil {
		return nil, nil, invalidIndexError{"Invalid index. See todo --help for more info"}
	}

	if err = todos.Delete(key); err != nil {
		return nil, nil, err
	}

	return key, val, nil
}

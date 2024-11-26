package todo

import (
	"fmt"
	"io"
	"os"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var cmdList = &cobra.Command{
	Use:     "list",
	Short:   "Show all incomplete todos",
	Long:    "Show a list of all incomplete todos added by using the 'add' command. List will be numbered starting from 1 which can be used to reference a todo in 'do' or 'rm' commands",
	Example: "todo list",
	Args:    cobra.NoArgs,
	Run:     list,
}

func list(cmd *cobra.Command, args []string) {
	db := getDB(cmd)
	defer db.Close()

	if err := db.Update(writeTodos(os.Stdout)); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}

func writeTodos(w io.Writer) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {

		b, err := tx.CreateBucketIfNotExists([]byte(todoBucket))
		if err != nil {
			fmt.Println(err)
			return err
		}

		c, i := b.Cursor(), 1

		if k, _ := c.First(); k == nil {
			w.Write([]byte("You have no tasks. See todo add --help\n"))
			return nil
		}

		output := "You have the following tasks:\n"

		for k, v := c.First(); k != nil; k, v = c.Next() {
			output += fmt.Sprintf("%d. %s\n", i, string(v))
			i++
		}

		if _, err := w.Write([]byte(output)); err != nil {
			return err
		}

		return nil
	}
}

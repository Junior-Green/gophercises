package todo

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
)

var cmdCompleted = &cobra.Command{
	Use:     "completed",
	Short:   "Show all completed todos",
	Long:    "Shows a list of all todos marked using that 'do' command",
	Example: "todo completed",
	Args:    cobra.NoArgs,
	Run:     completed,
}

func completed(cmd *cobra.Command, args []string) {
	db := getDB(cmd)
	defer db.Close()

	if err := db.Update(writeCompleted(os.Stdout)); err != nil {
		cmd.PrintErr(err)
		os.Exit(1)
	}
}

func writeCompleted(w io.Writer) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(completedBucket))
		if err != nil {
			return err
		}

		var (
			c      = b.Cursor()
			now    = time.Now()
			output = "You have finished the following tasks today:\n"
		)

		if k, _ := c.First(); k == nil {
			w.Write([]byte("You have no completed tasks. See todo do --help\n"))
			return nil
		}

		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 0, 1)

		startStr, endStr := start.Format(time.RFC3339), end.Format(time.RFC3339)

		for k, v := c.Seek([]byte(startStr)); k != nil && bytes.Compare(k, []byte(endStr)) <= 0; k, v = c.Next() {
			output += fmt.Sprintf("- %s\n", string(v))
		}

		if _, err := w.Write([]byte(output)); err != nil {
			return err
		}

		return nil
	}
}

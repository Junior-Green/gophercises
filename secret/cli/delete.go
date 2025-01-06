package secret

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdDelete = &cobra.Command{
	Use:  "delete [KEY]",
	Args: cobra.ExactArgs(1),
	Run:  deletePair,
}

func deletePair(cmd *cobra.Command, args []string) {
	store, err := getSecretStore()
	if err != nil {
		cmd.PrintErrf("%v\n", err)
		os.Exit(1)
	}

	if len(args) != 1 {
		panic("wrong number of arguments")
	}

	if err := store.Delete(args[0]); err != nil {
		cmd.PrintErrf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("key %q deleted\n", args[0])
}

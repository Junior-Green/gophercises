package secret

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdSet = &cobra.Command{
	Use:  "set [KEY] [VALUE]",
	Args: cobra.ExactArgs(2),
	Run:  setPair,
}

func setPair(cmd *cobra.Command, args []string) {
	store, err := getSecretStore()
	if err != nil {
		cmd.PrintErrf("%v\n", err)
		os.Exit(1)
	}

	if len(args) != 2 {
		panic("wrong number of arguments")
	}

	if err := store.Set(args[0], args[1]); err != nil {
		cmd.PrintErrf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("key %q set\n", args[0])
}

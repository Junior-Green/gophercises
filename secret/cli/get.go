package secret

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cmdGet = &cobra.Command{
	Use:  "get [KEY]",
	Args: cobra.ExactArgs(1),
	Run:  getPair,
}

func getPair(cmd *cobra.Command, args []string) {
	store, err := getSecretStore()
	if err != nil {
		cmd.PrintErrf("%v\n", err)
		os.Exit(1)
	}

	if len(args) != 1 {
		panic("wrong number of arguments")
	}

	secret, err := store.Get(args[0])
	if err != nil {
		cmd.PrintErrf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("secret: %s\n", secret)
}

package secret

import (
	"fmt"

	"github.com/Junior-Green/gophercises/secret"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const filename = ".secrets"

var cmdRoot = &cobra.Command{
	Use:   "secret [COMMAND]",
	Short: "Store secrets which are encrypted and persisted to local storage",
	Long:  "CLI program that lets you store key-value records which are encrypted using AES-GCM encryption. Encryption key used must be 16,24, or 32 bytes long",
}

func Execute() error {
	return cmdRoot.Execute()
}

func init() {
	cmdRoot.AddCommand(cmdSet)
	cmdRoot.AddCommand(cmdGet)
	cmdRoot.AddCommand(cmdDelete)
}

func getSecretStore() (*secret.Store, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	filepath := fmt.Sprintf("%s/%s", dir, filename)
	return &secret.Store{Filepath: filepath}, nil
}

package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/chanzuckerberg/chamber/store"
	"github.com/spf13/cobra"
)

var (
	recoverCmd = &cobra.Command{
		Use:   "recover <service> <key>",
		Short: "recover a previously deleted parameter",
		Args:  cobra.ExactArgs(2),
		RunE:  recover,
	}
)

func init() {
	RootCmd.AddCommand(recoverCmd)
}

func recover(cmd *cobra.Command, args []string) error {
	service := strings.ToLower(args[0])
	if err := validateService(service); err != nil {
		return errors.Wrap(err, "Failed to validate service")
	}

	key := strings.ToLower(args[1])
	if err := validateKey(key); err != nil {
		return errors.Wrap(err, "Failed to validate key")
	}

	secretStore := store.NewSSMStore(numRetries)
	secretId := store.SecretId{
		Service: service,
		Key:     key,
	}

	return secretStore.UntagDeleted(secretId)
}

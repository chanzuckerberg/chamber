package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/chanzuckerberg/chamber/store"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list <service>",
	Short: "List the secrets set for a service",
	Args:  cobra.ExactArgs(1),
	RunE:  list,
}

var (
	includeDeleted bool
	withValues bool
)

func init() {
	listCmd.Flags().BoolVarP(&withValues, "expand", "e", false, "Expand parameter list with values")
	listCmd.Flags().BoolVarP(&includeDeleted, "deleted", "d", false, "Include parameters marked as deleted")
	RootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	service := strings.ToLower(args[0])
	if err := validateService(service); err != nil {
		return errors.Wrap(err, "Failed to validate service")
	}

	secretStore := store.NewSSMStore(numRetries)
	secrets, err := secretStore.List(service, withValues, includeDeleted)
	if err != nil {
		return errors.Wrap(err, "Failed to list store contents")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprint(w, "Key\tVersion\tLastModified\tUser")
	if includeDeleted {
		fmt.Fprint(w, "\tDeleted")
	}
	if withValues {
		fmt.Fprint(w, "\tValue")
	}
	fmt.Fprintln(w, "")

	for _, secret := range secrets {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s",
			key(secret.Meta.Key),
			secret.Meta.Version,
			secret.Meta.Created.Local().Format(ShortTimeFormat),
			secret.Meta.CreatedBy)
		if includeDeleted {
			fmt.Fprintf(w, "\t%t", secret.Deleted)
		}
		if withValues {
			fmt.Fprintf(w, "\t%s", *secret.Value)
		}
		fmt.Fprintln(w, "")
	}

	w.Flush()
	return nil
}

func key(s string) string {
	_, noPaths := os.LookupEnv("CHAMBER_NO_PATHS")
	if !noPaths {
		tokens := strings.Split(s, "/")
		secretKey := tokens[2]
		return secretKey
	}

	tokens := strings.Split(s, ".")
	secretKey := tokens[1]
	return secretKey
}

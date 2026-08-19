package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLoginCmd() *cobra.Command {
	var uid, token string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate and get access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if uid == "" || token == "" {
				return fmt.Errorf("both --uid and --token are required")
			}

			c := newClient()
			pair, err := c.Login(uid, token)
			if err != nil {
				return err
			}

			fmt.Printf("Success! You are now authenticated.\n")
			fmt.Printf("Access Token: %s\n", pair.AccessToken)
			fmt.Printf("\nExport for subsequent commands:\n")
			fmt.Printf("  export LOADSG_TOKEN=%s\n", pair.AccessToken)
			return nil
		},
	}

	cmd.Flags().StringVarP(&uid, "uid", "u", "", "User UID")
	cmd.Flags().StringVarP(&token, "token", "t", "", "User secret")
	_ = cmd.MarkFlagRequired("uid")
	_ = cmd.MarkFlagRequired("token")

	return cmd
}
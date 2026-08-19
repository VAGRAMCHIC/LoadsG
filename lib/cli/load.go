package cli

import "github.com/spf13/cobra"

func newLoadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load",
		Short: "Manage load jobs",
	}
	cmd.AddCommand(newLoadCreateCmd())
	return cmd
}

//func listLoadCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "load",
//		Short: "Manage load jobs",
//	}
//	cmd.AddCommand(newLoadListCmd())
//	return cmd
//}
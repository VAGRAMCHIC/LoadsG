package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show LoadsG version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("LoadsG CLI\n")
		fmt.Printf("Version:   %s\n", Version)
		fmt.Printf("Commit:    %s\n", Commit)
		fmt.Printf("BuildDate: %s\n", BuildDate)
		fmt.Printf("Go:        %s\n", runtime.Version())
		fmt.Printf("OS/Arch:   %s/%s\n", runtime.GOOS, runtime.GOARCH)
		return nil
	},
}
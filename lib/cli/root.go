package cli

import (
	"os"

	"github.com/spf13/cobra"
	"loadsg/lib/client"
)

var (
	flagAddress string
	flagToken   string
	flagFormat  string
)

var rootCmd = &cobra.Command{
	Use:           "loadsg",
	Short:         "LoadsG command line interface",
	Long:          `LoadsG CLI is a command line client for the LoadsG load testing platform.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagAddress, "address", "a",
		envOr("LOADSG_ADDR", "http://127.0.0.1:8080"),
		"LoadsG API address")
	rootCmd.PersistentFlags().StringVarP(&flagToken, "token", "t",
		os.Getenv("LOADSG_TOKEN"),
		"Access token (or set LOADSG_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&flagFormat, "format", "table",
		"Output format: table|json")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(newLoginCmd())
	rootCmd.AddCommand(newLoadCmd())
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newClient() *client.Client {
	c := client.New(flagAddress)
	if flagToken != "" {
		c.SetToken(flagToken)
	}
	return c
}
package cmd

import (
	"os"

	"github.com/jaydeluca/benchmark-metrics/internal"
	"github.com/spf13/cobra"
)

var (
	repo  string
	owner string

	rootCmd = &cobra.Command{
		Use:   "benchmark-metrics",
		Short: "Pulls, parses, and visualizes benchmark metrics from github",
		Long:  `Pulls, parses, and visualizes benchmark metrics from github.`,
		Run: func(cmd *cobra.Command, args []string) {
			internal.Run()
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&repo, "repo", "opentelemetry-java-instrumentation", "Repository")
	rootCmd.PersistentFlags().StringVar(&owner, "owner", "opentelemetry", "Github owner")
}

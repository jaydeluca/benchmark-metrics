package cmd

/*
Copyright © 2023 Jay DeLuca

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	internal.Run()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&repo, "repo", "opentelemetry-java-instrumentation", "Repository")
	rootCmd.PersistentFlags().StringVar(&owner, "owner", "opentelemetry", "Github owner")
}

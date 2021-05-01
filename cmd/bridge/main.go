package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "bridge",
	}

	monitorCmd := &cobra.Command{
		Use: "monitor",
	}

	fdbCmd := &cobra.Command{
		Use: "fdb",
		Run: func(cmd *cobra.Command, args []string) {
			monitorFdb()
		},
	}

	monitorCmd.AddCommand(fdbCmd)
	rootCmd.AddCommand(monitorCmd)
	rootCmd.Execute()
}

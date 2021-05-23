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
	monitorCmd.AddCommand(&cobra.Command{
		Use: "fdb",
		Run: func(cmd *cobra.Command, args []string) {
			monitorFdb()
		},
	})

	fdbCmd := &cobra.Command{
		Use: "fdb",
	}
	fdbCmd.AddCommand(&cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			listFdb()
		},
	})

	rootCmd.AddCommand(fdbCmd)
	rootCmd.AddCommand(monitorCmd)
	rootCmd.Execute()
}

package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "tent"}
	rootCmd.AddCommand(cmdDiscover(), cmdAuth())
	rootCmd.Execute()
}
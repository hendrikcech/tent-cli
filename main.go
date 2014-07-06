package main

import (
	"github.com/spf13/cobra"
	"github.com/hendrikcech/tent/config"
)

func main() {
	_, err := config.Read()
	if err != nil {
		panic(err)
	}
	var rootCmd = &cobra.Command{Use: "tent"}
	rootCmd.AddCommand(cmdDiscover(), cmdAuth())
	rootCmd.Execute()
}
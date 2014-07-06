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
	cmdProfiles := CmdProfiles()
	cmdProfiles.AddCommand(CmdProfilesAdd(), CmdProfilesList(), CmdProfilesRemove())
	rootCmd.AddCommand(CmdDiscover(), CmdAuth(), cmdProfiles)

	rootCmd.Execute()
}
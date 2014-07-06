package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "tent"}
	cmdProfiles := CmdProfiles()
	cmdProfiles.AddCommand(CmdProfilesAdd(), CmdProfilesList(), CmdProfilesRemove())
	rootCmd.AddCommand(CmdDiscover(), CmdAuth(), cmdProfiles)

	rootCmd.Execute()
}

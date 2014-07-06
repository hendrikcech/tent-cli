package main

import (
	"github.com/spf13/cobra"
	"github.com/hendrikcech/tent/config"
	"fmt"
)

func main() {
	c := config.Config{}
	if err := c.Read(); err != nil {
		fmt.Println(err)
		return
	}

	rootCmd := &cobra.Command{Use: "tent"}
	cmdProfiles := CmdProfiles()
	cmdProfiles.AddCommand(CmdProfilesAdd(&c), CmdProfilesList(&c), CmdProfilesRemove(&c))
	rootCmd.AddCommand(CmdDiscover(), CmdAuth(&c), cmdProfiles)
	rootCmd.AddCommand(CmdQuery(&c))

	rootCmd.Execute()
}

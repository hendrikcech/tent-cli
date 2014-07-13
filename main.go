package main

import (
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
)

func main() {
	c := config.Config{}
	if err := c.Read(); err != nil {
		fmt.Println(err)
		return
	}

	rootCmd := &cobra.Command{Use: "tent"}

	cmdProfiles := CmdProfiles(&c)
	cmdProfiles.AddCommand(CmdProfilesAdd(&c), CmdProfilesRemove(&c), CmdProfilesDefault(&c))

	cmdSchemas := CmdSchemas(&c)
	cmdSchemas.AddCommand(CmdSchemasAdd(&c), CmdSchemasRemove(&c))

	rootCmd.AddCommand(CmdDiscover(), CmdAuth(&c), cmdProfiles, cmdSchemas, CmdQuery(&c), CmdGet(&c), CmdCreate(&c), CmdDelete(&c))

	rootCmd.Execute()
}

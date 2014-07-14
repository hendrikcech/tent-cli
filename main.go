package main

import (
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
	"net/url"
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

func isURL(s string) bool {
	if _, err := url.ParseRequestURI(s); err != nil {
		return false
	}
	return true
}
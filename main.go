package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

func main() {
	c := Config{}
	if err := c.Read(); err != nil {
		fmt.Println(err)
		return
	}

	rootCmd := &cobra.Command{Use: "tent"}

	cmdProfiles := CmdProfiles(&c)
	cmdProfiles.AddCommand(
		CmdProfilesAdd(&c),
		CmdProfilesRemove(&c),
		CmdProfilesDefault(&c),
	)

	cmdSchemas := CmdSchemas(&c)
	cmdSchemas.AddCommand(
		CmdSchemasAdd(&c),
		CmdSchemasRemove(&c),
	)

	rootCmd.AddCommand(
		CmdDiscover(),
		CmdAuth(&c),
		CmdCreate(&c),
		CmdGet(&c),
		CmdQuery(&c),
		CmdDelete(&c),
		cmdProfiles,
		cmdSchemas,
	)

	rootCmd.Execute()
}

func isURL(s string) bool {
	if _, err := url.ParseRequestURI(s); err != nil {
		return false
	}
	return true
}

func showHelpAndExit(cmd *cobra.Command) {
	cmd.Help()
	os.Exit(1)
}

func maybeExit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

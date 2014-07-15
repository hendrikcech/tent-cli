package main

import (
	"github.com/spf13/cobra"
)

func CmdDelete(c *Config) *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:   "delete <post_id> [<version_id>]",
		Short: "Delete a post.",
		Long: `
Delete a post with the given id created by the current entity.
The entire post will be deleted if <version_id> is not specified.`,
		Run: func(cmd *cobra.Command, args []string) {
			id := ""
			version := ""

			switch len(args) {
			case 1:
				id = args[0]
			case 2:
				id = args[0]
				version = args[1]
			default:
				showHelpAndExit(cmd)
			}

			p, err := c.Profile(profile)
			maybeExit(err)
			c := p.Client()

			_, err = c.DeletePost(id, version, true)
			maybeExit(err)
		},
	}

	setUseFlag(&profile, cmd)

	return cmd
}

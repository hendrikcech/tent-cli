package main

import (
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
	// "github.com/tent/tent-client-go"
)

func CmdDelete(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <post_id> [<version_id>]",
		Short: "Delete a post.",
		Long:  `
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
				cmd.Help()
				return
			}

			p, err := c.DefaultProfile()
			if err != nil {
				fmt.Println(err)
				return
			}
			c := p.Client()

			_, err = c.DeletePost(id, version, true)
			if err != nil {
				fmt.Println(err)
				return
			}
		},
	}

	return cmd
}

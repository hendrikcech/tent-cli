package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

func CmdGet(c *Config) *cobra.Command {
	var profile string

	cmd := &cobra.Command{
		Use:   "get [<entity>] <post_id> [<version>]",
		Short: "Get a single post.",
		Long: `
Get a single post by its' post id.
<entity> defaults to the current profiles entity; <version> to the latest version known by the server.`,
		Run: func(cmd *cobra.Command, args []string) {
			p, err := c.Profile(profile)
			maybeExit(err)

			entity := p.Entity
			id := ""
			version := ""

			switch len(args) {
			case 1:
				id = args[0]
			case 2:
				if i, profile := c.ProfileByName(args[0]); i > -1 {
					entity = profile.Entity
					id = args[1]
				} else if !isURL(args[0]) {
					// not an url -> id
					id = args[0]
					version = args[1]
				} else {
					entity = args[0]
					id = args[1]
				}
			case 3:
				entity = args[0]
				id = args[1]
				version = args[2]
			default:
				showHelpAndExit(cmd)
			}

			res, err := p.Client().GetPost(entity, id, version, nil)
			maybeExit(err)

			o, err := json.MarshalIndent(res.Post, "", "  ")
			maybeExit(err)

			fmt.Println(string(o))
		},
	}

	setUseFlag(&profile, cmd)

	return cmd
}

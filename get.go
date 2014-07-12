package main

import (
	"encoding/json"
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
	"net/url"
)

func CmdGet(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [entity] id [version]",
		Short: "Get a single post",
		Long:  "Get a single post. Only the post id is required. Entity defaults to the current profiles entity.",
		Run: func(cmd *cobra.Command, args []string) {
			p, err := c.DefaultProfile()
			if err != nil {
				fmt.Println(err)
				return
			}

			entity := p.Entity
			id := ""
			version := ""

			switch len(args) {
			case 1:
				id = args[0]
			case 2:
				if _, err := url.ParseRequestURI(args[0]); err != nil {
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
				cmd.Help()
				return
			}

			c := p.Client()
			res, err := c.GetPost(entity, id, version, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			o, err := json.MarshalIndent(res.Post, "", "  ")
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(string(o))
		},
	}

	return cmd
}
package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tent/tent-client-go"
	"strings"
)

func CmdUpdate(c *Config) *cobra.Command {
	var publishedAt string
	var public string
	var fragment string
	var useProfile string

	cmd := &cobra.Command{
		Use:   "update [<entity>] <post_id> <content>",
		Short: "Update an existing post.",
		Long: `
Create a new version of an existing post.
To make this as convenient as possible, the latest version of the specified post will be fetched first. Changes will then be applied on top and send back to the server.

Usage examples:
1. Simplest example. This will create a new version with the same properties as the previous version, except for the "version" object. 
update CiURoW

2. Update content. The format is equal to the create command content syntax. Notice that this doesn't just add properties to the content object but overwrites it completely. 
update CiURoW new="text" or json:='["object", "content"]'

3. It's also possible to create a version based on a post from a different entity. Either spell out the entity url or use a profile name.
update https://entity.cupcake.is CiURoW --public=false
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				showHelpAndExit(cmd)
			}

			profile, err := c.Profile(useProfile)
			maybeExit(err)

			entity := profile.Entity
			var id string
			var contentStart int

			if i, match := c.ProfileByName(args[0]); i > -1 {
				entity = match.Entity
			} else if isURL(args[0]) {
				entity = args[0]
			} else {
				id = args[0]
				contentStart = 1
			}

			if id == "" {
				if len(args) < 2 {
					showHelpAndExit(cmd)
				}
				id = args[1]
				contentStart = 2
			}

			client := profile.Client()

			res, err := client.GetPost(entity, id, "", nil)
			maybeExit(err)

			p := res.Post

			p.Version = &tent.PostVersion{
				Parents: []tent.PostVersionParent{
					{
						Entity:  p.Entity,
						Post:    p.ID,
						Version: p.Version.ID,
					},
				},
			}

			p.Entity = entity
			p.OriginalEntity = ""

			if fragment != "" {
				p.Type = strings.Split(p.Type, "#")[0] + "#" + strings.TrimLeft(fragment, "#")
			}

			p.Content = buildContent(args[contentStart:])

			if public != "" {
				if p.Permissions == nil {
					p.Permissions = &tent.PostPermissions{}
				}
				pBool := public == "true"
				p.Permissions.PublicFlag = &pBool
			}

			p.App = nil

			p.ReceivedAt = nil
			if publishedAt != "" {
				err = res.Post.PublishedAt.UnmarshalJSON([]byte(publishedAt))
				maybeExit(err)
			}

			err = client.CreatePost(p)
			maybeExit(err)

			o, err := json.MarshalIndent(res.Post, "", "  ")
			maybeExit(err)

			fmt.Println(string(o))
		},
	}

	cmd.Flags().StringVarP(&publishedAt, "publishedAt", "a", "", "Define published_at metadata. Pass a unix timestamp in milliseconds.")
	cmd.Flags().StringVarP(&public, "public", "p", "", "Set basic visibility of post.")
	cmd.Flags().StringVarP(&fragment, "fragment", "f", "", "Update post fragment. It's not possible to change the full type.")
	setUseFlag(&useProfile, cmd)

	return cmd
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
	"github.com/tent/tent-client-go"
	"net/url"
	"strings"
)

func CmdCreate(c *config.Config) *cobra.Command {
	var publishedAt string
	var public bool

	cmd := &cobra.Command{
		Use:   "create [<type> <content> | <json>]",
		Short: "Create a new post",
		Long: `Create a new post. Usage example:
create https://example.com/types/song/v0# name="Also Sprach Zarathustra" composor="Richard Strauss"

Use the ":=" operator to create a nested structure:
create https://example.com/types/place/v0# name=":)" location='{"lat": "-41.290975", "lon": "174.792864"}'

You can also directly pass the full post:
create '{"type": "https://example.com/types/person/v0#", "licenses": [{"url": "https://some.license"}], "content": { "name": "Joy" }}' 
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}

			var post *tent.Post

			if _, err := url.ParseRequestURI(args[0]); err != nil {
				err = json.Unmarshal([]byte(args[0]), &post)
				if err != nil {
					fmt.Println("Invalid post type or post json.")
					return
				}
			} else {
				if !strings.Contains(args[0], "#") {
					fmt.Println(`Post type must have a fragment. Place a "#" at the end.`)
					return
				}

				post = &tent.Post{
					Type:    args[0],
					Content: buildContent(args[1:]),
					Permissions: &tent.PostPermissions{
						PublicFlag: &public,
					},
				}

				if publishedAt != "" {
					post.PublishedAt = &tent.UnixTime{}
					if err := post.PublishedAt.UnmarshalJSON([]byte(publishedAt)); err != nil {
						fmt.Println(err)
						return
					}
				}
			}

			p, err := c.DefaultProfile()
			if err != nil {
				fmt.Println(err)
				return
			}

			if err = p.Client().CreatePost(post); err != nil {
				fmt.Println(err)
				return
			}
		},
	}

	cmd.Flags().StringVarP(&publishedAt, "publishedAt", "", "", "Define published_at metadata. Pass unix timestamp in milliseconds.")
	cmd.Flags().BoolVarP(&public, "public", "p", true, "Set basic visibility of post.")

	return cmd
}

func buildContent(args []string) []byte {
	buf := bytes.NewBufferString("{")
	for n, arg := range args {
		runes := []rune(arg)
		for i, r := range runes {
			if r == ':' && runes[i+1] == '=' {
				buf.WriteString(fmt.Sprintf(`"%s": %s`, string(runes[:i]), string(runes[i+2:])))
				break
			}
			if r == '=' {
				buf.WriteString(fmt.Sprintf(`"%s": "%s"`, string(runes[:i]), string(runes[i+1:])))
				break
			}
		}
		if n < len(args)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("}")
	return buf.Bytes()
}

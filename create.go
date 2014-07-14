package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tent/tent-client-go"
	"strings"
)

func CmdCreate(c *Config) *cobra.Command {
	var publishedAt string
	var public bool

	cmd := &cobra.Command{
		Use:   "create [<type> <content> | <json>]",
		Short: "Create a new post.",
		Long: `Create a new post. Usage example:
create https://example.com/types/song/v0# name="Also Sprach Zarathustra" composor="Richard Strauss"

Use the ":=" operator to create a nested structure:
create https://example.com/types/place/v0# name=":)" location='{"lat": "-41.290975", "lon": "174.792864"}'

You can also directly pass the full post:
create '{"type": "https://example.com/types/person/v0#", "licenses": [{"url": "https://some.license"}], "content": { "name": "Joy" }}'
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				showHelpAndExit(cmd)
			}

			var err error
			var post *tent.Post
			postType := args[0]

			i, s := c.SchemaByName(args[0])
			if i > -1 {
				postType = s.MergeFragment(args[0])
			}

			if i > -1 || isURL(args[0]) {
				if !strings.Contains(postType, "#") {
					exitWithError(`Post types are required to have a fragment. Place a "#" at the end.`)
				}

				post = &tent.Post{
					Type:    postType,
					Content: buildContent(args[1:]),
					Permissions: &tent.PostPermissions{
						PublicFlag: &public,
					},
				}

				if publishedAt != "" {
					post.PublishedAt = &tent.UnixTime{}
					err = post.PublishedAt.UnmarshalJSON([]byte(publishedAt))
					maybeExit(err)
				}
			} else {
				err = json.Unmarshal([]byte(args[0]), &post)
				if err != nil {
					exitWithError("Invalid post type or post json.")
				}
			}

			p, err := c.DefaultProfile()
			maybeExit(err)

			err = p.Client().CreatePost(post)
			maybeExit(err)

			o, err := json.MarshalIndent(post, "", "  ")
			maybeExit(err)

			fmt.Println(string(o))
		},
	}

	cmd.Flags().StringVarP(&publishedAt, "publishedAt", "", "", "Define published_at metadata. Pass a unix timestamp in milliseconds.")
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

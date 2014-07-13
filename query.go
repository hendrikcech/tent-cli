package main

import (
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
	"github.com/stevedomin/termtable"
	"github.com/tent/tent-client-go"
	"strconv"
	"strings"
	"time"
	"errors"
	"net/url"
)

func CmdQuery(c *config.Config) *cobra.Command {
	var limit int
	var since string    // 1234567890,version
	var before string   // 1234567890,version
	var until string    // 1234567890,version
	var entities string // entityone,entitytwo
	var types string    // typeone,typetwo
	var maxRefs int
	// var mentions string // mentionone,mentiontwo

	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query the posts feed",
		Long:  "Query the posts feed.",
		Run: func(cmd *cobra.Command, args []string) {
			p, err := c.DefaultProfile()
			if err != nil {
				fmt.Println(err)
				return
			}
			client := p.Client()

			q := tent.NewPostsFeedQuery()

			q.Limit(limit)
			q.MaxRefs(maxRefs)

			timeSetter := map[string]func(time.Time, string) *tent.PostsFeedQuery{
				since:  q.Since,
				before: q.Before,
				until:  q.Until,
			}
			for arg, setter := range timeSetter {
				err = splitAndSetTimeValue(arg, setter)
				if err != nil {
					fmt.Println("Invalid value", err)
					return
				}
			}

			if entities != "" {
				entitiesStr, err := splitAndMaybeReplace(entities, func(name string) (string, error) {
					if i, p := c.ProfileByName(name); i > -1 {
						return p.Entity, nil
					}
					if _, err := url.ParseRequestURI(name); err != nil {
						return "", errors.New(fmt.Sprintf(`Profile "%v" not found.`, name))
					}
					return name, nil
				})
				if err != nil {
					fmt.Println(err)
					return
				}
				q.Set("entities", entitiesStr)
			}
			if types != "" {
				typesStr, err := splitAndMaybeReplace(types, func(name string) (string, error) {
					n := strings.Split(name, "#")
					if i, p := c.SchemaByName(n[0]); i > -1 {
						if len(n) == 2 && !strings.Contains(p.PostType, "#") {
							return p.PostType + "#" + n[1], nil
						}
						return p.PostType, nil
					}
					if _, err := url.ParseRequestURI(name); err != nil {
						return "", errors.New(fmt.Sprintf(`Schema "%v" not found.`, name))
					}
					return name, nil
				})
				if err != nil {
					fmt.Println(err)
					return
				}
				q.Set("types", typesStr)
			}

			res, err := client.GetFeed(q, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			t := termtable.NewTable(nil, nil)
			t.SetHeader([]string{"ID", "ENTITY", "TYPE", "PUBLISHED_AT"})

			for _, p := range res.Posts {
				// layout := "2006-01-02 15:04" // p.PublishedAt.Format(layout)
				pubAt := strconv.FormatInt(p.PublishedAt.UnixMillis(), 10)
				t.AddRow([]string{p.ID, p.Entity, p.Type, pubAt})
			}

			fmt.Println(t.Render())
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 25, "Cap number of posts returned")
	cmd.Flags().StringVarP(&since, "since", "s", "", "Only posts since timestamp and/or version. Format: `123456789,versionid`")
	cmd.Flags().StringVarP(&before, "before", "b", "", "Only posts before timestamp and/or version. Format: `123456789,versionid`")
	cmd.Flags().StringVarP(&until, "until", "u", "", "Only posts until timestamp and/or version. Format: `123456789,versionid`")
	cmd.Flags().StringVarP(&entities, "entities", "e", "", "Only posts from specific entities.")
	cmd.Flags().StringVarP(&types, "types", "t", "", "Only posts with these types.")
	cmd.Flags().IntVarP(&maxRefs, "maxrefs", "r", 5, "Cap number of inlined refs per post.")
	// cmd.Flags().StringVarP(&mentions, "mentions", "m", "", "Only posts which mention specific entities.")

	return cmd
}

func splitAndSetTimeValue(arg string, setter func(time.Time, string) *tent.PostsFeedQuery) error {
	if arg != "" {
		s := strings.Split(arg, ",")

		t := tent.UnixTime{}
		err := t.UnmarshalJSON([]byte(s[0]))
		if err != nil {
			return err
		}
		if len(s) == 1 {
			setter(t.Time, "")
		} else {
			setter(t.Time, s[1])
		}
	}
	return nil
}

func splitAndMaybeReplace(s string, get func(string) (string, error)) (string, error) {
	res := []string{}
	for _, part := range strings.Split(s, ",") {
		p, err := get(part)
		if err != nil {
			return "", err
		}
		res = append(res, p)
	}
	return strings.Join(res, ","), nil
}
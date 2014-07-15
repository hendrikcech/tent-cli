package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stevedomin/termtable"
	"github.com/tent/tent-client-go"
	"strconv"
	"strings"
	"time"
)

func CmdQuery(c *Config) *cobra.Command {
	var limit int
	var since string    // 1234567890,version
	var before string   // 1234567890,version
	var until string    // 1234567890,version
	var entities string // entityone,entitytwo
	var types string    // typeone,typetwo
	var maxRefs int
	// var mentions string // mentionone,mentiontwo

	var profile string

	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query the posts feed.",
		Long: `
Query the posts feed of the default profile.
Find more information about the available parameters here: https://tent.io/docs/api#postsfeed
Join multiple values with commata, i.e. when using --entities or --types.

A note about --types and fragments:
- "--types=https://tent.io/types/status/v0" matches all fragments
- "--types=https://tent.io/types/status/v0#" just matches "https://tent.io/types/status/v0#"
- "--types=https://tent.io/types/status/v0#reply" just matches "https://tent.io/types/status/v0#reply"
`,
		Run: func(cmd *cobra.Command, args []string) {
			p, err := c.Profile(profile)
			maybeExit(err)
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
				maybeExit(err)
			}

			if entities != "" {
				entitiesStr, err := splitAndMaybeReplace(entities, func(name string) (string, error) {
					if i, p := c.ProfileByName(name); i > -1 {
						return p.Entity, nil
					}
					if !isURL(name) {
						return "", fmt.Errorf(`Profile "%v" not found.`, name)
					}
					return name, nil
				})
				maybeExit(err)
				q.Set("entities", entitiesStr)
			}
			if types != "" {
				typesStr, err := splitAndMaybeReplace(types, func(name string) (string, error) {
					if i, s := c.SchemaByName(name); i > -1 {
						return s.MergeFragment(name), nil
					}
					if !isURL(name) {
						return "", fmt.Errorf(`Schema "%v" not found.`, name)
					}
					return name, nil
				})
				maybeExit(err)
				q.Set("types", typesStr)
			}

			res, err := client.GetFeed(q, nil)
			maybeExit(err)

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

	cmd.Flags().IntVarP(&limit, "limit", "l", 25, "Cap number of returned posts.")
	cmd.Flags().StringVarP(&since, "since", "s", "", "Only posts since timestamp and/or version. Format: `123456789,versionid`")
	cmd.Flags().StringVarP(&before, "before", "b", "", "Only posts before timestamp and/or version. Format: `123456789,versionid`")
	cmd.Flags().StringVarP(&until, "until", "u", "", "Only posts until timestamp and/or version. Format: `123456789,versionid`")
	cmd.Flags().StringVarP(&entities, "entities", "e", "", "Only posts from specific entities.")
	cmd.Flags().StringVarP(&types, "types", "t", "", "Only posts with specific types.")
	cmd.Flags().IntVarP(&maxRefs, "maxrefs", "r", 5, "Cap number of inlined refs per post.")
	// cmd.Flags().StringVarP(&mentions, "mentions", "m", "", "Only posts which mention specific entities.")

	setUseFlag(&profile, cmd)

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

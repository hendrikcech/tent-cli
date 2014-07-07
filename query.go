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
)

var CmdQuery = func(c *config.Config) *cobra.Command {
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
				q.Set("entities", entities)
			}
			if types != "" {
				q.Set("types", types)
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

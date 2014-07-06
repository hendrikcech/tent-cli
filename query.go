package main

import (
	"github.com/hendrikcech/tent/config"
	"github.com/spf13/cobra"
	"github.com/tent/tent-client-go"
	"fmt"
	"github.com/stevedomin/termtable"
)

var CmdQuery = func(c *config.Config) *cobra.Command {
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
			q := tent.NewPostsFeedQuery().Limit(25)
			client := p.Client()
			res, err := client.GetFeed(q, nil)
			if err != nil {
				fmt.Println(err)
				return
			}
			// fmt.Printf("%+v\n", res.Posts[0])
			
			layout := "2006-01-02 15:04"
			
			t := termtable.NewTable(nil, nil)
			t.SetHeader([]string{"ID", "ENTITY", "TYPE", "PUBLISHED_AT"})

			for _, p := range res.Posts {
				t.AddRow([]string{p.ID, p.Entity, p.Type, p.PublishedAt.Format(layout)})
			}

			fmt.Println(t.Render())
		},
	}

	// cmd.Flags().StringVarP(&name, "name", "n", "Tent CLI", "Name of app")

	return cmd
}
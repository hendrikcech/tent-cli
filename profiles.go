package main

import (
	"fmt"
	"github.com/hendrikcech/tent/config"
	"github.com/spf13/cobra"
	"github.com/stevedomin/termtable"
	"github.com/tent/tent-client-go"
)

var CmdProfiles = func() *cobra.Command {
	return &cobra.Command{
		Use:   "profiles [add|list|remove]",
		Short: "Manage your tent profiles",
		Long:  "Add, list or remove tent profiles.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}

var CmdProfilesAdd = func() *cobra.Command {
	id := ""
	key := ""
	app := ""

	cmd := &cobra.Command{
		Use:   "add [name] [entity]",
		Short: "Add a profile",
		Long:  "Add a profile.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}

			name := args[0]
			entity := args[1]

			c := config.Config{}
			if err := c.Read(); err != nil {
				fmt.Println(err)
				return
			}

			if i, _ := c.ByName(name); i > -1 {
				fmt.Printf("Profile \"%v\" already exists.\n", name)
				return
			}

			meta, err := tent.Discover(entity)
			if err != nil {
				fmt.Println(err)
				return
			}

			c.Profiles = append(c.Profiles, config.ProfileConfig{
				Name:    name,
				Entity:  entity,
				Servers: meta.Servers,
				ID:      id,
				Key:     key,
				App:     app,
			})

			if err = c.Write(); err != nil {
				fmt.Println(err)
			}
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "Access token")
	cmd.Flags().StringVarP(&key, "key", "k", "", "Hawk key")
	cmd.Flags().StringVarP(&app, "app", "a", "", "App id")

	return cmd
}
var CmdProfilesList = func() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List your tent profiles",
		Long:  "List your tent profiles.",
		Run: func(cmd *cobra.Command, args []string) {
			c := config.Config{}
			if err := c.Read(); err != nil {
				fmt.Println(err)
				return
			}

			t := termtable.NewTable(nil, nil)
			t.SetHeader([]string{"NAME", "ENTITY", "ID", "KEY", "APP"})

			for _, p := range c.Profiles {
				t.AddRow([]string{p.Name, p.Entity, p.ID, p.Key, p.App})
			}

			fmt.Println(t.Render())
		},
	}
}
var CmdProfilesRemove = func() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [name]",
		Short: "Remove a profile",
		Long:  "Remove a profile by its name.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}
			name := args[0]

			c := config.Config{}
			if err := c.Read(); err != nil {
				fmt.Println(err)
				return
			}

			i, _ := c.ByName(name)
			if i == -1 {
				return
			}
			c.Profiles = append(c.Profiles[:i], c.Profiles[i+1:]...)

			if err := c.Write(); err != nil {
				fmt.Println(err)
			}
		},
	}
}

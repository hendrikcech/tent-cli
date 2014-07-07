package main

import (
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
	"github.com/stevedomin/termtable"
	"github.com/tent/tent-client-go"
)

var CmdProfiles = func(c *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "profiles [add|remove]",
		Short: "Manage your profiles",
		Long:  "Profiles are identified by a unique name and save the entity uri plus credentials.\nAdd, remove or set a default profile with this command.",
		Run: func(cmd *cobra.Command, args []string) {
			t := termtable.NewTable(nil, nil)
			t.SetHeader([]string{"NAME", "ENTITY", "ID", "KEY", "APP"})

			for _, p := range c.Profiles {
				t.AddRow([]string{p.Name, p.Entity, p.ID, p.Key, p.App})
			}

			fmt.Println(t.Render())
		},
	}
}
var CmdProfilesAdd = func(c *config.Config) *cobra.Command {
	var id string
	var key string
	var app string

	cmd := &cobra.Command{
		Use:   "add profile_name entity",
		Short: "Create a new profile",
		Long:  "Create a new profile named `profile_name` associated with `entity`.\nCredentials can either be specified with flags or by running `tent auth profile_name`.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}

			name := args[0]
			entity := args[1]

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
var CmdProfilesRemove = func(c *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "remove profile_name",
		Short: "Remove a profile",
		Long:  "Remove a profile by its name.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}
			name := args[0]

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
var CmdProfilesDefault = func(c *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "default [name]",
		Short: "Echo or set the default profile",
		Long:  "Echo or set the default profile.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 { // echo default profile
				if c.Default == "" {
					fmt.Println("No default profile set.")
					return
				}
				fmt.Printf("Default profile is \"%v\"\n", c.Default)
			} else { // set default profile
				i, _ := c.ByName(args[0])
				if i == -1 {
					fmt.Println("No profile named \"%v\" existent.", args[0])
					return
				}
				c.Default = args[0]
				c.Write()
			}
		},
	}
}

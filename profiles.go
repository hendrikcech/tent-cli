package main

import (
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
	"github.com/stevedomin/termtable"
	"github.com/tent/tent-client-go"
)

func CmdProfiles(c *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "profiles",
		Short: "Manage entity profiles.",
		Long:  `
List, add or remove profiles or change the default.
Profiles save entity uris' and credentials. They are identified by a unique name.
The default profile is used by other commands like create, query and get.`,
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

func CmdProfilesAdd(c *config.Config) *cobra.Command {
	var id string
	var key string
	var app string

	cmd := &cobra.Command{
		Use:   "add <profile_name> <entity>",
		Short: "Create a new profile.",
		Long:  `
Create a new profile named <profile_name> that's associated with <entity>.
Credentials can either be specified with flags or by running ` + "`auth <profile_name>`.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}

			name := args[0]
			entity := args[1]

			if i, _ := c.ProfileByName(name); i > -1 {
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

			if c.Default == "" && len(c.Profiles) == 1 {
				c.Default = name
			}

			if err = c.Write(); err != nil {
				fmt.Println(err)
			}
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "Save hawk_id (or `access_token` in tent post).")
	cmd.Flags().StringVarP(&key, "key", "k", "", "Save hawk_key.")
	cmd.Flags().StringVarP(&app, "app", "a", "", "Save app id.")

	return cmd
}

func CmdProfilesRemove(c *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <profile_name>",
		Short: "Remove a profile.",
		Long:  "Remove a profile by its name.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}
			name := args[0]

			i, _ := c.ProfileByName(name)
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

func CmdProfilesDefault(c *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "default [<profile_name>]",
		Short: "Output or set the default profile.",
		Long:  `
Output or set the default profile.
This profile will be used by other commands like create, get or delete.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 { // echo default profile
				if c.Default == "" {
					fmt.Println("No default profile set.")
					return
				}
				fmt.Printf("Default profile is \"%v\"\n", c.Default)
			} else { // set default profile
				i, _ := c.ProfileByName(args[0])
				if i == -1 {
					fmt.Printf("No profile named \"%v\" existent.\n", args[0])
					return
				}
				c.Default = args[0]
				c.Write()
			}
		},
	}
}

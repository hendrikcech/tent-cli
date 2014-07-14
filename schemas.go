package main

import (
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
	"github.com/stevedomin/termtable"
	"strings"
)

func CmdSchemas(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schemas",
		Short: "Manage post schemas.",
		Long:  `
List, add or remove post schemas.
Schemas keep you from typing post types over and over by mapping them to short names.
This simple feature might get expanded when strict post schemas get introduced with Tent 0.4.

Post types can optionally be saved with fragments (e.g. "schemas add status https://tent.io/types/status/v0#reply"). These extensions can be overwritten when used: "create status#different_fragment".
`,
		Run: func(cmd *cobra.Command, args []string) {
			t := termtable.NewTable(nil, nil)
			t.SetHeader([]string{"NAME", "POSTTYPE"})

			for _, s := range c.Schemas {
				t.AddRow([]string{s.Name, s.PostType})
			}

			fmt.Println(t.Render())
		},
	}
	return cmd
}

func CmdSchemasAdd(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <schema_name> <post_type>",
		Short: "Add a new post schema.",
		Long:  "Add a new post schema.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}

			name := args[0]
			postType := args[1]

			if i, _ := c.SchemaByName(name); i > -1 {
				fmt.Printf("Schema \"%v\" already exists.\n", name)
				return
			}

			if strings.Contains(name, "#") {
				fmt.Println(`Schema names can't contain "#".`)
				return
			}

			c.Schemas = append(c.Schemas, config.SchemaConfig{
				Name:    name,
				PostType: postType,
			})

			if err := c.Write(); err != nil {
				fmt.Println(err)
			}
		},
	}
	return cmd
}

func CmdSchemasRemove(c *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <schema_name>",
		Short: "Remove a post schema.",
		Long:  "Remove a schema by its name.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}

			i, _ := c.SchemaByName(args[0])
			if i == -1 {
				return
			}
			c.Schemas = append(c.Schemas[:i], c.Schemas[i+1:]...)

			if err := c.Write(); err != nil {
				fmt.Println(err)
			}
		},
	}
}
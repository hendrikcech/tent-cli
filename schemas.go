package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stevedomin/termtable"
	"strings"
)

func CmdSchemas(c *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "schemas",
		Short: "Manage post schemas.",
		Long: `
List, set or remove post schemas.
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
}

func CmdSchemasSet(c *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "set <schema_name> <post_type>",
		Short: "Define a post schema.",
		Long:  "Add a new or overwrite an existing post schema.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				showHelpAndExit(cmd)
			}

			name := args[0]
			postType := args[1]

			if strings.Contains(name, "#") {
				exitWithError(`Schema names can't contain "#".`)
			}

			if i, s := c.SchemaByName(name); i > -1 {
				// schema already exists
				s.Name = name
				s.PostType = postType
			} else {
				c.Schemas = append(c.Schemas, SchemaConfig{
					Name:     name,
					PostType: postType,
				})
			}

			err := c.Write()
			maybeExit(err)
		},
	}
}

func CmdSchemasRemove(c *Config) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <schema_name>",
		Short: "Remove a post schema.",
		Long:  "Remove a schema by its name.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				showHelpAndExit(cmd)
			}

			i, _ := c.SchemaByName(args[0])
			if i == -1 {
				return
			}
			c.Schemas = append(c.Schemas[:i], c.Schemas[i+1:]...)

			err := c.Write()
			maybeExit(err)
		},
	}
}

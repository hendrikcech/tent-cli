package main

import (
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
	"github.com/stevedomin/termtable"
	"strings"
)

func CmdSchemas(c *config.Config) *cobra.Command {
	return &cobra.Command{
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
}

func CmdSchemasAdd(c *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "add <schema_name> <post_type>",
		Short: "Add a new post schema.",
		Long:  "Add a new post schema.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				showHelpAndExit(cmd)
			}

			name := args[0]
			postType := args[1]

			if i, _ := c.SchemaByName(name); i > -1 {
				exitWithError(fmt.Sprintf("Schema \"%v\" already exists.", name))
			}

			if strings.Contains(name, "#") {
				exitWithError(`Schema names can't contain "#".`)
			}

			c.Schemas = append(c.Schemas, config.SchemaConfig{
				Name:    name,
				PostType: postType,
			})

			err := c.Write()
			maybeExit(err)
		},
	}
}

func CmdSchemasRemove(c *config.Config) *cobra.Command {
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
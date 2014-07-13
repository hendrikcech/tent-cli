package main

import (
	"fmt"
	"github.com/hendrikcech/tent-cli/config"
	"github.com/spf13/cobra"
	"github.com/stevedomin/termtable"
)

func CmdSchemas(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schemas",
		Short: "Mange post schemas",
		Long:  "Manage post schemas.",
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
		Short: "Add a new post schema",
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
		Short: "Remove a schema",
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
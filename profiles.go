package main

import (
	"github.com/spf13/cobra"
	"github.com/hendrikcech/tent/config"
	"github.com/stevedomin/termtable"
	"fmt"
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
	return &cobra.Command{
		Use:   "add [entity]",
		Short: "Add an entity to your tent profiles",
		Long:  "Add an entity to your tent profiles.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("not implemented yet")
		},
	}
}
var CmdProfilesList = func() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List your tent profiles",
		Long:  "List your tent profiles.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := config.Read()
			if err != nil {
				fmt.Println(err)
				return
			}

			t := termtable.NewTable(nil, nil)
			t.SetHeader([]string{"ENTITY", "SHORT", "ID", "KEY", "APP"})

			for _, e := range c.Entities {
				t.AddRow([]string{e.Entity, e.Short, e.ID, e.Key, e.App})
			}

			fmt.Println(t.Render())
		},
	}
}
var CmdProfilesRemove = func() *cobra.Command {
	return &cobra.Command{
		Use:   "remove",
		Short: "Remove a tent profile",
		Long:  "Remove a tent profile.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("not implemented yet")
		},
	}
}
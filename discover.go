package main

import (
	"github.com/spf13/cobra"
	"github.com/tent/tent-client-go"
	"fmt"
	"encoding/json"
)

var CmdDiscover = func() *cobra.Command {
	return &cobra.Command{
		Use:   "discover [url to discover]",
		Short: "Discover an url",
		Long:  "Get the associated entity for an url.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				return
			}
			meta, err := tent.Discover(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			o, err := json.MarshalIndent(meta, "", "  ")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(o))
		},
	}
}
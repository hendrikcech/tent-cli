package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tent/tent-client-go"
)

var CmdDiscover = func() *cobra.Command {
	return &cobra.Command{
		Use:   "discover url",
		Short: "Discover an url",
		Long:  "Discover an url and output any associated meta posts.",
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

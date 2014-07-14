package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tent/tent-client-go"
)

func CmdDiscover() *cobra.Command {
	return &cobra.Command{
		Use:   "discover <url>",
		Short: "Get the meta post that is associated with an url.",
		Long:  "Discover an url and output any associated meta posts.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				showHelpAndExit(cmd)
			}
			meta, err := tent.Discover(args[0])
			maybeExit(err)
			o, err := json.MarshalIndent(meta, "", "  ")
			maybeExit(err)

			fmt.Println(string(o))
		},
	}
}

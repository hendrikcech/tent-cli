package main

import (
	"github.com/spf13/cobra"
	"github.com/tent/tent-client-go"
	"github.com/skratchdot/open-golang/open"
	"strings"
	"fmt"
)
/*auth [entity] --write=all --read=all --name=Tent CLI --url=testapp.com --scopes=permissions*/


var CmdAuth = func() *cobra.Command {
	name := "Tent CLI"
	url := "https://app.example.com"
	write := "all"
	read := "all"
	scopes := "permissions"
	
	cmd := &cobra.Command{
		Use:   "auth [entity]",
		Short: "Get new credentials for an entity",
		Long:  "Get new credentials for an entity.",
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

			client := &tent.Client{Servers: meta.Servers}

			app := &tent.App{
				Name: name,
				URL:  url,
				RedirectURI: "https://app.example.com/oauth",
			}
			
			if write != "" || read != "" {
				app.Types = tent.AppTypes{}

				if write != "" {
					app.Types.Write = strings.Split(write, ",")
				}
				if read != "" {
					app.Types.Read = strings.Split(read, ",")
				}
			}

			if scopes != "" {
				app.Scopes = strings.Split(scopes, ",")
			}

			post := tent.NewAppPost(app)

			err = client.CreatePost(post)
			if err != nil {
				fmt.Println(err)
				return
			}

			client.Credentials, _, err = post.LinkedCredentials()
			if err != nil {
				fmt.Println(err)
				return
			}

			// redirect url
			oauthURL := meta.Servers[0].URLs.OAuthURL(post.ID, "randomState")
			err = open.Run(oauthURL)
			if err != nil {
				fmt.Println(oauthURL)
			}

			// wait for code input
			var code string
			_, err = fmt.Scanf("%s", &code)
			if err != nil {
				fmt.Println(err)
				return
			}

			// request access token
			// client.Credentials, err = client.RequestAccessToken(code)
			c, err := client.RequestAccessToken(code)
			if err != nil {
				fmt.Println(err)
				return
			}
			
			fmt.Printf(`{
  "id": "%v",
  "key": "%v",
  "algorithm": "sha256",
  "token_type": ""https://tent.io/oauth/hawk-token"
}`, c.ID, c.Key)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "Tent CLI", "Name of app")
	// cmd.Flags().StringVarP(&url, "url", "u", "tentcliapp.com", "App url")
	cmd.Flags().StringVarP(&write, "write", "w", "all", "Write permissions")
	cmd.Flags().StringVarP(&read, "read", "r", "all", "Read permissions")
	cmd.Flags().StringVarP(&scopes, "scopes", "s", "permissions", "Scopes")

	return cmd
}	

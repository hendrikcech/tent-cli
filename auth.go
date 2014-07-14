package main

import (
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/tent/tent-client-go"
	"strings"
)

func CmdAuth(c *Config) *cobra.Command {
	var name string
	url := "https://app.example.com"
	var write string
	var read string
	var scopes string

	cmd := &cobra.Command{
		Use:   "auth [<entity>|<profile_name>]",
		Short: "Authorize a new app.",
		Long: `
Create a new app on the tent server of the specified entity and output the credentials.
If <profile_name> is given, the profile will be updated with the obtained tokens.
Join multiple values with commata, i.e. when using --read or --scopes.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				showHelpAndExit(cmd)
			}

			var servers []tent.MetaPostServer
			_, p := c.ProfileByName(args[0])
			if p.Name != "" { // existing profile name passed
				servers = p.Servers
			} else {
				meta, err := tent.Discover(args[0])
				maybeExit(err)
				servers = meta.Servers
			}

			client := &tent.Client{Servers: servers}

			app := &tent.App{
				Name:        name,
				URL:         url,
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

			err := client.CreatePost(post)
			maybeExit(err)

			client.Credentials, _, err = post.LinkedCredentials()
			maybeExit(err)

			// redirect url
			oauthURL := servers[0].URLs.OAuthURL(post.ID, "randomState")
			err = open.Run(oauthURL)
			fmt.Printf("Visit this url, accept, paste the code back in and press enter:\n%v\n", oauthURL)
			// if err != nil {
			// 	// fmt.Println(oauthURL)
			// }

			// wait for code input
			var code string
			_, err = fmt.Scanf("%s", &code)
			maybeExit(err)

			// request access token
			// client.Credentials, err = client.RequestAccessToken(code)
			tokens, err := client.RequestAccessToken(code)
			maybeExit(err)

			tmpl := `
{
  "access_token": "%v",
  "hawk_key": "%v",
  "hawk_algorithm": "sha256",
  "token_type": ""https://tent.io/oauth/hawk-token"
}
`
			fmt.Printf(tmpl, tokens.ID, tokens.Key)

			if p.Name != "" {
				p.ID = tokens.ID
				p.Key = tokens.Key
				p.App = tokens.App
				c.Write()
				fmt.Printf("Saved to profile \"%v\".\n", p.Name)
			}
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "Tent CLI", "The applications' name.")
	cmd.Flags().StringVarP(&read, "read", "r", "all", "Read permissions.")
	cmd.Flags().StringVarP(&scopes, "scopes", "s", "permissions", `Scopes ("permissions"!).`)
	cmd.Flags().StringVarP(&write, "write", "w", "all", "Write permissions.")
	// cmd.Flags().StringVarP(&url, "url", "u", "tentcliapp.com", "App url")

	return cmd
}

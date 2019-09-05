package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin"
)

// DNSResolverPlugin - a simple plugin to resolve bosh based hostname's IP address
type DNSResolverPlugin struct{}

// Run - the entry point for the plugin
func (c *DNSResolverPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	// command resolve
	if args[0] == "resolve" {
		if len(args) < 2 {
			fmt.Println("Insufficient arguments. Please provide the host name to be resolved")
			os.Exit(-1)
		}

		result, err := cliConnection.IsLoggedIn()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}

		if !result {
			fmt.Println("Not logged in")
			os.Exit(-1)
		}

		currentOrg, err := cliConnection.GetCurrentOrg()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}

		currentSpace, err := cliConnection.GetCurrentSpace()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}

		fmt.Printf("Resolving in Org - %s / Space - %s\n\n", currentOrg.Name, currentSpace.Name)

		apps, err := cliConnection.GetApps()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-1)
		}

		// Count apps in org, and handle error if org has no apps deployed
		appsCount := len(apps)
		if appsCount == 0 {
			fmt.Printf("Sorry, this plugin requires at least one deployed app. No apps found under org-%s space-%s\n ", currentOrg.Name, currentSpace.Name)
			os.Exit(-1)
		}

		// for deployed apps, find the first one that runs. app need to be running so we can connect to it with SSH.
		var i int
		var vtargetApp string
		for i = 0; i < appsCount; i++ {
			// debug
			// fmt.Printf("Checking app %s has State %s\n", apps[i].Name, apps[i].State)
			if apps[i].State == "started" {
				vtargetApp = apps[i].Name

				fmt.Printf("Checking if %s app is eligible to resolve host\n", vtargetApp)

				_, err := cliConnection.CliCommand("ssh", vtargetApp, "-c", "host "+args[1]+" > /dev/null;echo;echo IP lookup  result - ;echo `dig +short "+args[1]+"`")
				if err == nil {
					// ssh works!
					// Resolution successful - Not sure on the output from ssh though
					os.Exit(0)
				}
			}
		}

		// Error handling if the org has no running apps
		fmt.Printf("Sorry, this plugin requires at least one ssh enabled started app. No apps were eligible in org-%s space-%s\n ", currentOrg.Name, currentSpace.Name)
		os.Exit(-1)
	}
}

// GetMetadata - the method to return metadata about the plugin
func (c *DNSResolverPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "DNSResolverPlugin",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 1,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "resolve",
				HelpText: "takes a bosh based host name & returns its IP address",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "cf resolve <bosh based host name>",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(DNSResolverPlugin))
}

package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"

	log "github.com/Sirupsen/logrus"
)

var logLevel = 5
var client ucp.Client

var importPath, exportPath, action string

var top, exampleFile bool

func init() {
	diverCmd.AddCommand(UCPRoot)

	ucpLogin.Flags().StringVar(&client.Username, "username", os.Getenv("DIVER_USERNAME"), "Username that has permissions to authenticate to Docker EE")
	ucpLogin.Flags().StringVar(&client.Password, "password", os.Getenv("DIVER_PASSWORD"), "Password allowing a user to authenticate to Docker EE")
	ucpLogin.Flags().StringVar(&client.UCPURL, "url", os.Getenv("DIVER_URL"), "URL for Docker EE, e.g. https://10.0.0.1")
	ignoreCert := strings.ToLower(os.Getenv("DIVER_INSECURE")) == "true"

	ucpLogin.Flags().BoolVar(&client.IgnoreCert, "ignorecert", ignoreCert, "Ignore x509 certificate")

	ucpLogin.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	// Container flags
	ucpContainer.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
	ucpContainer.Flags().BoolVar(&top, "top", false, "Enable TOP for watching running containers")

	UCPRoot.AddCommand(ucpContainer)
	UCPRoot.AddCommand(ucpCliBundle)
	UCPRoot.AddCommand(ucpNetwork)
	UCPRoot.AddCommand(ucpLogin)

	// Sub commands
	ucpContainer.AddCommand(ucpContainerTop)
	ucpContainer.AddCommand(ucpContainerList)

}

// UCPRoot - This is the root of all UCP commands / flags
var UCPRoot = &cobra.Command{
	Use:   "ucp",
	Short: "Universal Control Plane ",
	Run: func(cmd *cobra.Command, args []string) {

		existingClient, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			cmd.Help()
			log.Warn("Unable to find existing session, please login")
			return
		}
		currentAccount, err := existingClient.AuthStatus()
		if err != nil {
			cmd.Help()
			log.Warn("Session has expired, please login")
			return
		}
		log.Infof("Current user [%s]", currentAccount.Name)
		return
	},
}

// UCPRoot - This is the root of all UCP commands / flags
var ucpLogin = &cobra.Command{
	Use:   "login",
	Short: "Authenticate against the Universal Control Pane",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		err := client.Connect()

		// Check if connection was succesful
		if err != nil {
			log.Fatalf("%v", err)
		} else {
			// If succesfull write the token and annouce as succesful
			err = client.WriteToken()
			if err != nil {
				log.Errorf("%v", err)
			}
			log.Infof("Succesfully logged into [%s]", client.UCPURL)
		}
	},
}

var ucpContainer = &cobra.Command{
	Use:   "containers",
	Short: "Interact with containers",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
	},
}

var ucpContainerTop = &cobra.Command{
	Use:   "top",
	Short: "A list of containers and their CPU usage like the top command on linux",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.ContainerTop()
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}

var ucpContainerList = &cobra.Command{
	Use:   "list",
	Short: "List all containers across all nodes in UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.GetContainerNames()
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}

var ucpCliBundle = &cobra.Command{
	Use:   "client-bundle",
	Short: "Download the client bundle for UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.GetClientBundle()
		if err != nil {
			log.Fatalf("%v", err)
		}

	},
}

var ucpNetwork = &cobra.Command{
	Use:   "network",
	Short: "Interact with container networks",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		err = client.GetNetworks()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

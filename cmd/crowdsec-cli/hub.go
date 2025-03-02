package main

import (
	"fmt"

	"github.com/crowdsecurity/crowdsec/pkg/cwhub"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewHubCmd() *cobra.Command {
	/* ---- HUB COMMAND */
	var cmdHub = &cobra.Command{
		Use:   "hub [action]",
		Short: "Manage Hub",
		Long: `
Hub management

List/update parsers/scenarios/postoverflows/collections from [Crowdsec Hub](https://hub.crowdsec.net).
Hub is manage by cscli, to get latest hub files from [Crowdsec Hub](https://hub.crowdsec.net), you need to update.
		`,
		Example: `
cscli hub list   # List all installed configurations
cscli hub update # Download list of available configurations from the hub
		`,
		Args:              cobra.ExactArgs(0),
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if csConfig.Cscli == nil {
				return fmt.Errorf("you must configure cli before interacting with hub")
			}

			return nil
		},
	}
	cmdHub.PersistentFlags().StringVarP(&cwhub.HubBranch, "branch", "b", "", "Use given branch from hub")

	var cmdHubList = &cobra.Command{
		Use:               "list [-a]",
		Short:             "List installed configs",
		Args:              cobra.ExactArgs(0),
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {

			if err := csConfig.LoadHub(); err != nil {
				log.Fatalf(err.Error())
			}
			if err := cwhub.GetHubIdx(csConfig.Hub); err != nil {
				log.Fatalf("Failed to get Hub index : %v", err)
				log.Infoln("Run 'sudo cscli hub update' to get the hub index")
			}
			//use LocalSync to get warnings about tainted / outdated items
			_, warn := cwhub.LocalSync(csConfig.Hub)
			for _, v := range warn {
				log.Info(v)
			}
			cwhub.DisplaySummary()
			ListItems([]string{
				cwhub.COLLECTIONS, cwhub.PARSERS, cwhub.SCENARIOS, cwhub.PARSERS_OVFLW,
			}, args, true, false, all)
		},
	}
	cmdHubList.PersistentFlags().BoolVarP(&all, "all", "a", false, "List disabled items as well")
	cmdHub.AddCommand(cmdHubList)

	var cmdHubUpdate = &cobra.Command{
		Use:   "update",
		Short: "Fetch available configs from hub",
		Long: `
Fetches the [.index.json](https://github.com/crowdsecurity/hub/blob/master/.index.json) file from hub, containing the list of available configs.
`,
		Args:              cobra.ExactArgs(0),
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if csConfig.Cscli == nil {
				return fmt.Errorf("you must configure cli before interacting with hub")
			}

			if err := setHubBranch(); err != nil {
				return fmt.Errorf("error while setting hub branch: %s", err)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := csConfig.LoadHub(); err != nil {
				log.Fatalf(err.Error())
			}
			if err := cwhub.UpdateHubIdx(csConfig.Hub); err != nil {
				log.Fatalf("Failed to get Hub index : %v", err)
			}
			//use LocalSync to get warnings about tainted / outdated items
			_, warn := cwhub.LocalSync(csConfig.Hub)
			for _, v := range warn {
				log.Info(v)
			}
		},
	}
	cmdHub.AddCommand(cmdHubUpdate)

	var cmdHubUpgrade = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade all configs installed from hub",
		Long: `
Upgrade all configs installed from Crowdsec Hub. Run 'sudo cscli hub update' if you want the latest versions available.
`,
		Args:              cobra.ExactArgs(0),
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if csConfig.Cscli == nil {
				return fmt.Errorf("you must configure cli before interacting with hub")
			}

			if err := setHubBranch(); err != nil {
				return fmt.Errorf("error while setting hub branch: %s", err)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := csConfig.LoadHub(); err != nil {
				log.Fatalf(err.Error())
			}
			if err := cwhub.GetHubIdx(csConfig.Hub); err != nil {
				log.Fatalf("Failed to get Hub index : %v", err)
				log.Infoln("Run 'sudo cscli hub update' to get the hub index")
			}

			log.Infof("Upgrading collections")
			cwhub.UpgradeConfig(csConfig, cwhub.COLLECTIONS, "", forceAction)
			log.Infof("Upgrading parsers")
			cwhub.UpgradeConfig(csConfig, cwhub.PARSERS, "", forceAction)
			log.Infof("Upgrading scenarios")
			cwhub.UpgradeConfig(csConfig, cwhub.SCENARIOS, "", forceAction)
			log.Infof("Upgrading postoverflows")
			cwhub.UpgradeConfig(csConfig, cwhub.PARSERS_OVFLW, "", forceAction)
		},
	}
	cmdHubUpgrade.PersistentFlags().BoolVar(&forceAction, "force", false, "Force upgrade : Overwrite tainted and outdated files")
	cmdHub.AddCommand(cmdHubUpgrade)
	return cmdHub
}

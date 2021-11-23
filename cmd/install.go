package cmd

import (
	"github.com/satisfactoryhost/satisfy/pkg/system"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var installLog = log.WithFields(log.Fields{
	"command": "install",
})
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your application",
	Long:  `A brief description of your application`,
	Run: func(cmd *cobra.Command, args []string) {
		installLog.Info("starting installation")
		installLog.Info("setting up system client")
		sysClient := system.NewClient(installLog)
		installLog.Infof("starting new %s install", sysClient.ID)

		if err := sysClient.PackageService.InstallPrereqs(); err != nil {
			installLog.Fatal(err)
		}

		if err := sysClient.PackageService.InstallSteam(); err != nil {
			installLog.Fatal(err)
		}

		if err := sysClient.PackageService.CheckInstall(); err != nil {
			installLog.Fatal(err)
		}

		installLog.Info("creating Satisfactory service")
		if err := sysClient.CreateUser(); err != nil {
			installLog.Fatal(err)
		}
		if err := sysClient.CreateService(); err != nil {
			installLog.Fatal(err)
		}
		installLog.Info("✅ done creating Satisfactory service")

		installLog.Info("✅ installation complete")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

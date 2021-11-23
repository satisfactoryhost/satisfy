package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "satisfy",
	Short: "A brief description of your application",
	Long:  `A brief description of your application`,
}

func init() {
	rootCmd.PersistentFlags().StringP("user", "u", "satisfactory", "the user to run the Satisfactory server as")
	rootCmd.PersistentFlags().StringP("group", "g", "satisfactory", "the group for the Satisfactory server user and assets")
}

func Execute() error {
	return rootCmd.Execute()
}

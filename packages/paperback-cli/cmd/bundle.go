/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	paperback "paperback.moe/paperback-cli/api"
)

// bundleCmd represents the bundle command
var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		baseDir, _ := cmd.Flags().GetString("cwd")
		outFolder, _ := cmd.Flags().GetString("folder")

		paperback.Bundle(baseDir, paperback.BundlerConfiguration{
			OutputFolder: outFolder,
		})
	},
}

func init() {
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(bundleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bundleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	bundleCmd.Flags().StringP("folder", "o", "", "Output folder")
	bundleCmd.Flags().String("cwd", ex, "Output folder")
}

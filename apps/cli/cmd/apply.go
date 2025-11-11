/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"cli/pkg/manifest"

	"cli/pkg/resource"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration from a file",
	Long:  `Apply a configuration from a file to the API server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, err := cmd.Flags().GetString("filename")
		if err != nil {
			return err
		}
		fmt.Printf("Applying configuration from %s\n", filename)
		manifest, err := manifest.ParseManifest(filename)
		if err != nil {
			return err
		}
		fmt.Printf("Manifest: %+v\n", manifest)
		return resource.ApplyResource(manifest)
	},
}

func init() {

	applyCmd.Flags().StringP("filename", "f", "", "The filename of the configuration file to apply")
	rootCmd.AddCommand(applyCmd)

}

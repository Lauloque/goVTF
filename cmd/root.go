/* SPDX-License-Identifier: GPL-3.0-or-later */
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	outputDir   string
	alphaFormat string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goVTF",
	Short: "An ersatz of VTFLib and VTFCmd that should work on linux",
	Long: `goVTF is a humble attempt to imit ate VTFLib and VTFCmd, VALVe's
library and its CLI to create a "Valve Texture File" (VTF) used in VALVe games.
But usable in Linux.

First implementation should just be able to take a common picture file and
output a VTF usable as a "spray" in TF2`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goVTF.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "./", "Output directory")
}

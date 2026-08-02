/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file [path]",
	Short: "Input file path",
	Long:  `Input file path to a png, jpg or tga image to process.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input_path := args[0]

		_, err := os.Stat(input_path)
		if os.IsNotExist(err) {
			fmt.Printf("Couldn't find '%s'\n", input_path)
			return
		} else if err != nil {
			fmt.Println("Error: ", err)
		}

		input_ext := filepath.Ext(input_path)
		switch input_ext {
		case "png", "jpg", "jpeg":
			fmt.Printf("Found '%s' image in '%s'\n", input_ext, input_path)
		default:
			fmt.Printf("Unsupported file foramt '%s' found in '%s'\n", input_ext, input_path)
		}

	},
}

func init() {
	rootCmd.AddCommand(fileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

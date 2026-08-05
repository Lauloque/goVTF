/* SPDX-License-Identifier: GPL-3.0-or-later */
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Lauloque/goVTF/imageutils"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file [path]",
	Short: "Input file path",
	Long:  `Input file path to a png, jpg or tga image to process.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input_path := args[0]

		_, err := os.Stat(input_path)
		if os.IsNotExist(err) {
			return fmt.Errorf("Couldn't find '%s'", input_path)
		} else if err != nil {
			return fmt.Errorf("Error: %w", err)
		}

		input_ext := filepath.Ext(input_path)
		switch input_ext {
		case ".png", ".jpg", ".jpeg":
			fmt.Printf("Found '%s' image in '%s'\n", input_ext, input_path)
		default:
			fmt.Printf("Unsupported file foramt '%s' found in '%s'\n", input_ext, input_path)
		}

		output_path := outputDir

		_, err = os.Stat(output_path)
		if os.IsNotExist(err) {
			return fmt.Errorf("Couldn't find output directory '%s'", output_path)
		} else if os.IsPermission(err) {
			return fmt.Errorf("Permission denied for output directory '%s'", output_path)
		} else if err != nil {
			return fmt.Errorf("Error: %w", err)
		}

		fmt.Printf("Output directory '%s'\n", output_path)

		img, format, err := imageutils.Load(input_path)
		if err != nil {
			return fmt.Errorf("Error: %w", err)
		}
		fmt.Printf("Loaded %s\n", format)
		fmt.Printf("Bounds: %v\n", img.Bounds())

		return nil
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

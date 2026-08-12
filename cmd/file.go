/* SPDX-License-Identifier: GPL-3.0-or-later */
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lauloque/goVTF/imageutils"
	"github.com/Lauloque/goVTF/vtf"
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

		// -----------------------------------------------------
		// End of inputs validation, will have to clean up later
		// -----------------------------------------------------

		// Image loading
		img, format, err := imageutils.Load(input_path)
		if err != nil {
			return fmt.Errorf("Error loading image: %w", err)
		}
		fmt.Printf("Loaded %s\n", format)
		fmt.Printf("Bounds: %v\n", img.Bounds())

		// Texture loading
		tex := imageutils.LoadTexture(input_path)
		if tex == nil {
			return fmt.Errorf("Failed to load texture from '%s'", input_path)
		}
		fmt.Printf("Loaded %dx%d texture\n", tex.Width, tex.Height)

		// Output file creation
		input_base := filepath.Base(input_path)
		input_stem := strings.TrimSuffix(input_base, input_ext)
		output_file := filepath.Join(output_path, input_stem+".vtf")
		f, err := os.Create(output_file)
		if err != nil {
			return fmt.Errorf("Error creating output file: %w", err)
		}
		defer f.Close()

		// Output file writing
		if err := vtf.Write(f, tex); err != nil {
			return fmt.Errorf("Error writing VTF: %w", err)
		}

		fmt.Printf("Successfully wrote '%s'\n", output_file)
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

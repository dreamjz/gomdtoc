/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	// flags
	skipDirs   []string
	recursive  bool
	titleField string
	// cmd
	rootCmd = &cobra.Command{
		Use:   "gomdtoc",
		Short: "CLI program to generate toc for markdown notes",
		Long:  `CLI program to generate toc for markdown notes directory`,
		Run: func(cmd *cobra.Command, args []string) {
			//root := "."
			root := "E:\\tmp\\go-temp"
			if len(args) > 0 {
				root = args[0]
			}
			GenerateTOCFile(root)
			log.Printf("Skip Dirs: %v", skipDirs)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&skipDirs, "skip", "s", []string{}, "--skip dir_name1,dir_name2, ...; skip specified directories")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "--recursive; generate TOC file for every sub-directory ")
	rootCmd.PersistentFlags().StringVarP(&titleField, "title", "t", "title", "--title title_field, specify the title field in frontmatter ")
}

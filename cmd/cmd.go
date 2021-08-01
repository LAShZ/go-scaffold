package cmd

import (
	"embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-scaffold",
	Short: "A auto project scaffold generator for go",
	Long:  "A auto project scaffold generator for go, will generate a golang project using gin and gorm, and with a default ping router",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run go-scaffold -h for help")
	},
}
var Template embed.FS

func init() {
	rootCmd.AddCommand(initCmd)
}

func Execute(template embed.FS) {
	Template = template
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

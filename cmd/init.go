package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var ProjectName string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate scaffold",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			ProjectName = args[len(args)-1]
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if ProjectName == "" {
			path, err := os.Getwd()
			if err != nil {
				log.Println("ERROR: get current filepath failed!")
				return
			}
			pathSlice := strings.Split(path, "/")
			ProjectName = pathSlice[len(pathSlice)-1]
		}
	},
}

func init() {
	initCmd.Flags().StringVarP(&ProjectName, "name", "n", "", "specify project name")
}

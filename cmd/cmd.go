package cmd

import (
	"embed"
	"fmt"
	"os"

	"github.com/LAShZ/go-scaffold/pkg/tempfs"
	"github.com/spf13/cobra"
)

const ASCIIART string = `
 ________  ________  ________  ________ ________ ________  ___       ________     
|\   ____\|\   ____\|\   __  \|\  _____\\  _____\\   __  \|\  \     |\   ___ \    
\ \  \___|\ \  \___|\ \  \|\  \ \  \__/\ \  \__/\ \  \|\  \ \  \    \ \  \_|\ \   
 \ \_____  \ \  \    \ \   __  \ \   __\\ \   __\\ \  \\\  \ \  \    \ \  \ \\ \  
  \|____|\  \ \  \____\ \  \ \  \ \  \_| \ \  \_| \ \  \\\  \ \  \____\ \  \_\\ \ 
    ____\_\  \ \_______\ \__\ \__\ \__\   \ \__\   \ \_______\ \_______\ \_______\
   |\_________\|_______|\|__|\|__|\|__|    \|__|    \|_______|\|_______|\|_______|
   \|_________|             
                                                                                  
		`

var rootCmd = &cobra.Command{
	Use:   "go-scaffold",
	Short: "A auto project scaffold generator for go",
	Long:  "A auto project scaffold generator for go, will generate a golang project using gin and gorm, and with a default ping router",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ASCIIART)
		fmt.Println("run go-scaffold -h for help")
	},
}
var Template tempfs.TempFS

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

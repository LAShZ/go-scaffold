package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/LAShZ/go-scaffold/config"
	"github.com/LAShZ/go-scaffold/pkg"
	"github.com/spf13/cobra"
)

var ProjectName string
var configFile string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate scaffold",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			ProjectName = args[len(args)-1]
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(ASCIIART)
		if configFile != "" {
			config.Setup(configFile)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if ProjectName == "" {
			path, err := os.Getwd()
			if err != nil {
				log.Println("ERROR: get current filepath failed!")
				return
			}
			pathSlice := strings.Split(path, string(filepath.Separator))
			ProjectName = pathSlice[len(pathSlice)-1]
		}
		fmt.Printf("Creating project: %s\n", ProjectName)
		generateProject(Template)
	},
}

func init() {
	initCmd.Flags().StringVarP(&ProjectName, "name", "n", "", "specify project name")
	initCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./config.toml", "specify personal config file path")
}

func generateProject(fs fs.FS) {
	var opts []pkg.Options
	if config.Info.Log.Use {
		opts = append(opts, pkg.WithLogger(config.Info.Log.Logger))
	}
	if config.Info.ORM.Use {
		opts = append(opts, pkg.WithORM(config.Info.ORM.Frame))
	}
	if config.Info.Web.Use {
		opts = append(opts, pkg.WithWeb(config.Info.Web.Frame))
	}
	opts = append(opts, pkg.UseRedis(config.Info.Redis.Use))
	pg := pkg.NewProjectGenerator(fs, opts...)
	pg.Generate()
}
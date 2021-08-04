package cmd

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/LAShZ/go-scaffold/config"
	"github.com/LAShZ/go-scaffold/pkg"
	"github.com/spf13/cobra"
)

var configFile string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate scaffold",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			pkg.ProjectName = args[len(args)-1]
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
		if pkg.ProjectName == "" {
			path, err := os.Getwd()
			if err != nil {
				log.Println("ERROR: get current filepath failed!")
				return
			}
			pathSlice := strings.Split(path, string(filepath.Separator))
			pkg.ProjectName = pathSlice[len(pathSlice)-1]
		}
		if config.Info.Project.Name != "" {
			pkg.ProjectName = config.Info.Project.Name
		}
		if config.Info.Project.Module != "" {
			pkg.ProjectModule = config.Info.Project.Module
		} else if pkg.ProjectModule == "" {
			pkg.ProjectModule = pkg.ProjectName
		}
		if config.Info.Project.GoVersion != "" {
			pkg.GoVersion = config.Info.Project.GoVersion
		} else {
			version := runtime.Version()
			version = strings.TrimPrefix(version, "go")
			version = strings.Join(strings.Split(version, ".")[:2], ".")
			pkg.GoVersion = version
		}
		fmt.Printf("\033[0;31mCreating project: \033[0m%s\n", pkg.ProjectName)
		generateProject(Template)
	},
}

func init() {
	initCmd.Flags().StringVarP(&pkg.ProjectName, "name", "n", "", "specify project name")
	initCmd.PersistentFlags().StringVarP(&pkg.ProjectModule, "module", "m", "", "specify project module")
	initCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "specify personal config file path")
	initCmd.PersistentFlags().BoolVarP(&pkg.Verbose, "verbose", "v", false, "show details of building")
}

func generateProject(fs embed.FS) {
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

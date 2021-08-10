package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/LAShZ/go-scaffold/config"
	"github.com/LAShZ/go-scaffold/pkg/engine"
	"github.com/LAShZ/go-scaffold/pkg/tempfs"
	"github.com/spf13/cobra"
)

var configFile string
var tempDir string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "generate scaffold",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			engine.ProjectName = args[len(args)-1]
		}
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println(ASCIIART)
		if configFile != "" {
			config.Setup(configFile)
		}
		if tempDir != "" {
			ufs, err := tempfs.NewUserTempFS(tempDir)
			if err != nil {
				fmt.Println("Read specify template dir failed, err:" + err.Error())
				os.Exit(1)
			}
			Template = ufs
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if engine.ProjectName == "" {
			path, err := os.Getwd()
			if err != nil {
				log.Println("ERROR: get current filepath failed!")
				return
			}
			pathSlice := strings.Split(path, string(filepath.Separator))
			engine.ProjectName = pathSlice[len(pathSlice)-1]
		}
		if config.Info.Project.Name != "" {
			engine.ProjectName = config.Info.Project.Name
		}
		if config.Info.Project.Module != "" {
			engine.ProjectModule = config.Info.Project.Module
		} else if engine.ProjectModule == "" {
			engine.ProjectModule = engine.ProjectName
		}
		if config.Info.Project.GoVersion != "" {
			engine.GoVersion = config.Info.Project.GoVersion
		} else {
			version := runtime.Version()
			version = strings.TrimPrefix(version, "go")
			version = strings.Join(strings.Split(version, ".")[:2], ".")
			engine.GoVersion = version
		}
		fmt.Printf("\033[0;31mCreating project: \033[0m%s\n", engine.ProjectName)
		generateProject(Template)
	},
}

func init() {
	initCmd.Flags().StringVarP(&engine.ProjectName, "name", "n", "", "specify project name")
	initCmd.PersistentFlags().StringVarP(&engine.ProjectModule, "module", "m", "", "specify project module")
	initCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "specify personal config file path")
	initCmd.PersistentFlags().BoolVarP(&engine.Verbose, "verbose", "v", false, "show details of building")
	initCmd.PersistentFlags().StringVar(&tempDir, "template", "", "specify personal template directory")
}

func generateProject(fs tempfs.TempFS) {
	var opts []engine.Options
	if config.Info.Log.Use {
		opts = append(opts, engine.WithLogger(config.Info.Log.Logger))
	}
	if config.Info.ORM.Use {
		opts = append(opts, engine.WithORM(config.Info.ORM.Frame))
	}
	if config.Info.Web.Use {
		opts = append(opts, engine.WithWeb(config.Info.Web.Frame))
	}
	opts = append(opts, engine.UseRedis(config.Info.Redis.Use))
	pg := engine.NewProjectGenerator(fs, opts...)
	pg.Generate()
}

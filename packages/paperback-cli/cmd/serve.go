/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
	paperback "paperback.moe/paperback-cli/api"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		baseDir, _ := cmd.Flags().GetString("cwd")
		outFolder, _ := cmd.Flags().GetString("folder")
		bundlesDir := filepath.Join(baseDir, outFolder, "bundles")

		bundleCmd.Run(cmd, args)

		srv := paperback.RunServer(paperback.ServerConfiguration{
			Directory: bundlesDir,
		})

		buf := bufio.NewReader(os.Stdin)

		tsRegex := regexp.MustCompile(`\.(ts|js|json)$`)
		watchServer := watcher.New()
		watchServer.SetMaxEvents(1)
		watchServer.FilterOps(watcher.Write)
		watchServer.AddFilterHook(watcher.RegexFilterHook(tsRegex, false))

		go func() {
			for {
				select {
				case <-watchServer.Event:
					log.Printf("watcher: typescript files changed")
					bundleCmd.Run(cmd, args)
				case err := <-watchServer.Error:
					log.Fatalln(err)
				case <-watchServer.Closed:
					return
				}
			}
		}()

		if err := watchServer.Add("."); err != nil {
			log.Fatalln(err)
		}

		// Watch test_folder recursively for changes.
		if err := watchServer.AddRecursive(filepath.Join(baseDir, "src")); err != nil {
			log.Fatalln(err)
		}

		go func() {
			if err := watchServer.Start(time.Millisecond * 300); err != nil {
				log.Fatalln(err)
			}
		}()

		log.Printf("serve: r = restart; s = stop;")
		running := true
		for running {
			sentence, err := buf.ReadBytes('\n')
			if err != nil {
				fmt.Println(err)
				break
			} else {
				switch strings.Trim(string(sentence), " \n") {
				case "s", "stop":
					running = false
				case "r", "restart":
					bundleCmd.Run(cmd, args)
				}
			}
		}

		paperback.ShutdownServer(srv)
	},
}

func init() {
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serveCmd.Flags().StringP("folder", "o", "", "Output folder")
	serveCmd.Flags().String("cwd", ex, "Output folder")
}

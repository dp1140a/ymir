/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"ymir/pkg/logger"
	"ymir/pkg/server"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the ymir server",
	Long:  `Describe required config params etc here`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func serve() {
	//Init Logger
	err := logger.InitLogger()
	if err != nil {
		log.Error("cant initialize logger")
	}
	log.Infof("Starting Ymir Server")

	//Init Server
	s, err := server.NewServer()
	if err != nil {
		log.Fatal("Cannot start server.  Shutting down: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Start Shutdown channel
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		sig := <-ch
		log.Info("Signal caught. Shutting down... Reason: ", sig)
		cancel()
		log.Println("Server gracefully stopped")
	}()

	// Spawn and start server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		s.ServeAPI(ctx)
	}()
	wg.Wait()
}

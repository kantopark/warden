package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/valve"
	"github.com/spf13/viper"

	"warden/router"
)

//func main() {
//	configSetup()
//
//	man, _ := docker.NewManager()
//	log.Println("Running build image")
//	err := man.BuildImage(docker.ImageBuildOptions{
//		Name:     "python-simple-proj",
//		Alias:    "",
//		GitURL:   "https://github.com/kantopark-tpl/python-simple.git",
//		Hash:     "",
//		Username: "danielbok",
//		Password: "167b0061e33c5ef5731c2f66bc4a7a387923af36",
//		Handler:  "main.entry_func",
//		RunEnv:   "python",
//	})
//	if err != nil {
//		log.Fatalln(err)
//	}
//}

func main() {
	configSetup()

	valv := valve.New()
	baseCtx := valv.Context()
	r := router.NewApp()

	srv := http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: chi.ServerBaseContext(baseCtx, r)}

	c := make(chan os.Signal, 1)
	gracePeriod := viper.GetDuration("server.graceperiod")
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down..")

			// sends a shutdown context to the context into the server
			if err := valv.Shutdown(gracePeriod); err != nil {
				log.Println(err)
			}

			// create context with timeout. Shuts down automatically after 3 seconds
			ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
			defer cancel()

			// start http shutdown
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatalf("error encountered when shutting down server: %s\n", err)
			}

			select {
			case <-time.After(gracePeriod + 2*time.Second):
				log.Println("shutting down early even though not all processes were killed")
			case <-ctx.Done():
			}
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		panic(fmt.Errorf("Error occured during server start: %s\n", err))
	}
}

func configSetup() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/kantopark/warden")
	viper.AddConfigPath("C:\\kantopark\\warden")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error reading in config: %s\n", err))
	}
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/valve"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"warden/application"
	"warden/config"
	"warden/store"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config.ReadConfig()
}

func main() {
	pflag.Bool("dev", false, "If true, runs application in development mode")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatalln(errors.Wrap(err, "error binding flags to viper"))
	}

	useDev := viper.GetBool("dev") // retrieve values from viper instead of pflag
	if useDev {
		runDev()
	} else {
		runProd()
	}
}

// Routine to setup application with some trial data. Used in local development
func localSetup() (CleanUp func()) {
	db, err := store.NewStore()
	if err != nil {
		log.Fatalln(err)
	}

	// add user
	user, err := db.UserCreate(store.UserBody{
		Email:    "daniel.bok@outlook.com",
		Username: "admin",
		Password: "admin",
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Create a python project application
	projectName := "py-sample-proj"
	proj, err := db.ProjectCreate(
		"https://github.com/kantopark-tpl/python-simple",
		projectName,
		"A sample python project",
		*user)
	if err != nil {
		log.Fatalln(err)
	}

	// Create an Instance
	hash := "95bfc3515452bfafeb2e04f948ac26d1e2a871c8"
	inst, err := db.InstanceCreate(hash, "dev", proj.Name)
	if err != nil {
		log.Fatalln(err)
	}

	return func() {
		if err := db.ProjectDelete(projectName); err != nil {
			log.Fatalln(err)
		}
		if err := db.InstanceDelete(inst.ProjectID, inst.CommitHash); err != nil {
			log.Fatalln(err)
		}
		log.Println("Clean up successfully")
	}
}

// Runs application in development mode
func runDev() {
	cb := localSetup()

	app := application.NewApp()
	defer app.Close()

	valv := valve.New()
	baseCtx := valv.Context()

	srv := http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: chi.ServerBaseContext(baseCtx, app.Router())}

	c := make(chan os.Signal, 1)
	gracePeriod := viper.GetDuration("server.graceperiod")
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("warden is shutting down..")
			cb()

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

	srv.ListenAndServe()
	log.Println("warden has shutdown")
}

// Runs application in production mode
func runProd() {
	app := application.NewApp()
	defer app.Close()

	valv := valve.New()
	baseCtx := valv.Context()

	srv := http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: chi.ServerBaseContext(baseCtx, app.Router())}

	c := make(chan os.Signal, 1)
	gracePeriod := viper.GetDuration("server.graceperiod")
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("warden is shutting down..")

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

	srv.ListenAndServe()
	log.Println("warden has shutdown")
}

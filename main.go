package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/spf13/viper"

	"warden/server"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/kantopark/warden")
	viper.AddConfigPath("C:\\kantopark\\warden")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error reading in config: %s\n", err))
	}

	srv := server.NewServer()
	images, _ := srv.Docker.ImageList(context.Background(), types.ImageListOptions{})

	log.Printf("Number of images: %d\n", len(images))
	for _, image := range images {
		log.Printf("%+v\n", image)
	}
	log.Println("End of story")
}

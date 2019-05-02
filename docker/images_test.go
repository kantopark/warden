package docker

import "log"

func ExampleClient_FindImageByName() {
	cli, _ := NewClient()

	// Matches all images such as busybox:latest, busytruck:lat01
	img, _ := cli.FindImageByName("busy*:lat*")
	log.Printf("%+v\n", img)

	// Matches only image redis:5.0.4-alpine, returns an error if image
	// can't be found
	_, err := cli.FindImageByName("redis:5.0.4-alpine")
	if err != nil {
		log.Fatalln(err)
	}
}

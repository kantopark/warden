package docker

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

func streamResponse(body io.ReadCloser) {
	content, _ := ioutil.ReadAll(body)

	var message []string
	for _, bline := range bytes.Split(content, []byte("\n")) {
		bline = bytes.TrimSpace(bline)
		if len(bline) == 0 {
			message = append(message, "")
		}
		var msg map[string]string
		json.Unmarshal(bline, &msg)
		mm := ""
		for _, m := range msg {
			mm += m
		}
		message = append(message, mm)
	}
	log.Println(strings.Join(message, ""))
}

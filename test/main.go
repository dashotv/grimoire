package main

import (
	"encoding/json"
	"fmt"
	"os"

	ptn "github.com/middelink/go-parse-torrent-name"
	"github.com/sirupsen/logrus"
)

func main() {
	info, err := ptn.Parse(os.Args[1])
	if err != nil {
		logrus.Errorf("error during parse: %s", err)
	}
	b, err := json.Marshal(info)
	if err != nil {
		logrus.Errorf("json failed to marshal: %s", err)
	}
	fmt.Printf(string(b))
}

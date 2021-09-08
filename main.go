package main

import (
	"log"
	"os"
	"strings"

	"9fans.net/go/acme"
)

var dir_tags = []byte(" als_start f gg ")
var file_tags = []byte(" als_def als_refs als_impls ")

func main() {
	// add to existing windows
	wins, err := acme.Windows()
	if err != nil {
		log.Fatal(err)
	}
	for _, win := range wins {
		err = add_tag(win.ID)
		if err != nil {
			log.Fatal(err)
		}
	}

	// watching new windows
	logF, err := acme.Log()
	if err != nil {
		log.Fatal(err)
	}
	for {
		event, err := logF.Read()
		if err != nil {
			log.Fatal(err)
		}

		if event.Op == "new" {
			err = add_tag(event.ID)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func add_tag(win_id int) error {
	win, err := acme.Open(win_id, nil)
	if err != nil {
		return err
	}
	defer win.CloseFiles()

	// select tags according to the path type
	var tags []byte
	s, err := win.ReadAll("tag")
	if err != nil {
		return err
	}
	path := strings.Fields(string(s))[0]
	if info, err := os.Stat(path); err != nil {
		return nil
	} else if info.IsDir() {
		tags = dir_tags
	} else {
		tags = file_tags
	}

	_, err = win.Write("tag", tags)
	if err != nil {
		return err
	}

	return nil
}

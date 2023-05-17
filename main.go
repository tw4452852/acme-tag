package main

import (
	"log"
	"os"
	"strings"

	"9fans.net/go/acme"
)

var dir_tags = []string{"gf"}
var file_tags = []string{"Get", "als_def", "als_refs"}

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
	var tags = []string{""}
	s, err := win.ReadAll("tag")
	if err != nil {
		return err
	}

	tagline := string(s)
	path := strings.Fields(tagline)[0]
	if info, err := os.Stat(path); err != nil && path[0] != '/' {
		return nil
	} else if info != nil && info.IsDir() {
		for _, tag := range dir_tags {
			if !strings.Contains(tagline, tag) {
				tags = append(tags, tag)
			}
		}
	} else {
		for _, tag := range file_tags {
			if !strings.Contains(tagline, tag) {
				tags = append(tags, tag)
			}
		}
	}

	_, err = win.Write("tag", []byte(strings.Join(tags, " ")))
	if err != nil {
		return err
	}

	return nil
}

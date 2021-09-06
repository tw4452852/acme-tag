package main

import (
	"log"
	"os"
	"strings"

	"9fans.net/go/acme"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("%s [tag1] [tag2] ...\n", os.Args[0])
	}

	tags := os.Args[1:]

	// add to existing windows
	wins, err := acme.Windows()
	if err != nil {
		log.Fatal(err)
	}
	for _, win := range wins {
		err = add_tag(win.ID, tags)
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
			err = add_tag(event.ID, tags)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func add_tag(win_id int, tags []string) error {
	win, err := acme.Open(win_id, nil)
	if err != nil {
		return err
	}
	defer win.CloseFiles()

	// get existed tags
	s, err := win.ReadAll("tag")
	if err != nil {
		return err
	}
	currentTags := map[string]struct{}{}
	for _, tag := range strings.Fields(string(s)) {
		currentTags[tag] = struct{}{};
	}

	add := []string{}
	for _, tag := range tags {
		if _, ok := currentTags[tag]; !ok {
			add = append(add, tag)
		}
	}

	if len(add) == 0 {
		return nil
	}

	_, err = win.Write("tag", []byte(strings.Join(add, " ")))
	if err != nil {
		return err
	}

	return nil
}
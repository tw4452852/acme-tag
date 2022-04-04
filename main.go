package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"9fans.net/go/acme"
)

var dir_tags = []string{" f "}
var file_tags = []string{"Get", "Put", "als_refs"}

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
	var tags []string
	s, err := win.ReadAll("tag")
	if err != nil {
		return err
	}

	tagline := string(s)
	path := strings.Fields(tagline)[0]
	is_file := false
	if info, err := os.Stat(path); err != nil {
		return nil
	} else if info.IsDir() {
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
		is_file = true
	}

	_, err = win.Write("tag", []byte(strings.Join(tags, " ")))
	if err != nil {
		return err
	}

	if is_file {
		go captureMiddleClick(win_id)
	}
	return nil
}

func captureMiddleClick(win_id int) error {
	win, err := acme.Open(win_id, nil)
	if err != nil {
		return err
	}
	defer win.CloseFiles()

	for event := range win.EventChan() {
		//log.Printf("%d: %c%c %d %d %#x %q\n", win_id, event.C1, event.C2, event.OrigQ0, event.Q0, event.Flag, event.Text)

		// middle click on the file's body which isn't recognized as acme's builtin cmd triggers als_def
		if event.C1 == 'M' && event.C2 == 'X' && event.Flag&0x1 == 0 {
			cmd := exec.Command("als", "def")
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, fmt.Sprintf("acme_pos0=%d", event.OrigQ0))
			cmd.Run()
			continue
		}

		win.WriteEvent(event)

		// <ctrl-s> for `Put` shortcut, this works but will block undo/redo functionality
		if false {
			if event.C1 == 'K' && event.C2 == 'I' && string(event.Text) == "\x13" {
				win.Write("addr", []byte(fmt.Sprintf("#%d,#%d", event.OrigQ0, event.OrigQ0+1)))
				win.Write("data", []byte(""))
				win.Write("ctl", []byte("put"))
			}
		}
	}

	//log.Printf("%d: exit\n", win_id)
	return nil
}

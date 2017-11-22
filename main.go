package main

import (
	"log"
	"path/filepath"
	"encoding/json"
	"os"
	"bytes"
	"io/ioutil"
	"crypto/sha1"
	"github.com/fsnotify/fsnotify"
	"strings"
)

type WatcherStruct struct {
	To string
	From string
	Watch string
}

func main() {
	var paths WatcherStruct
	f, err := ioutil.ReadFile("watcherConfig.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(f, &paths)
	if err != nil {
		log.Fatal(err)
	}
	// create object to watch files from specific directory
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close() // this runs after main() exits

	// run function in another channel
	done := make(chan bool)
	go func() {
		for { // runs like a while loop in other langs
			select {
			case event := <-watcher.Events:
				// run this stuff when a new file is created
				if (event.Op&fsnotify.Create == fsnotify.Create) && strings.HasSuffix(event.Name, ".pdf") {
					// list all files in the input file path
					dirs, erro := filepath.Glob(filepath.Join(paths.From, "*.pdf"))
					if erro != nil {
						log.Fatal(erro)
					}
					file, err := os.Open(event.Name)
					if err != nil {
						log.Fatal(err)
					}
					// read in 4kb of data to compute a checksum
					filebytes := make([]byte, 4096)
					_, err = file.Read(filebytes)
					if err != nil {
						log.Fatal(err)
					}
					file.Close()
					// compute the checksum
					chksum := sha1.Sum(filebytes)
					notequal := true
					// loop thru each file in input path until a file matches the checksum of the
					// created file in the mendeley directory
					for i := 0; i < len(dirs) && notequal; i++ {
						// open file and read 4kb
						f2, err := os.Open(dirs[i])
						if err != nil {
							log.Fatal(err)
						}
						fb2 := make([]byte, 4096)
						_, err = f2.Read(fb2)
						if err != nil {
							log.Fatal(err)
						}
						f2.Close()
						// compute checksum
						chksum2 := sha1.Sum(fb2)
						// compare the bytes themselves
						if bytes.Equal(chksum2[:], chksum[:]) {
							log.Println("These two files are equal:", dirs[i], event.Name)
							notequal = true
							// newPath := filepath.Join(paths.To, filepath.Base(dirs[i]))
							// log.Println("Moving", dirs[i], "to", newPath)
							log.Println("Deleting", dirs[i])
							os.Remove(dirs[i])
						}
					}
					log.Println("created file:", event.Name)
				}
			case err:= <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// add the mendeley dir to watch for created files
	err = watcher.Add(paths.Watch)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

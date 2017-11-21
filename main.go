package main

import (
	"log"
	"path/filepath"
	"os"
	"bytes"
	"crypto/sha1"
	"github.com/fsnotify/fsnotify"
)

const comparePath = "/Users/wgillis/Dropbox (HMS)/lab-datta/papers to read"
const watchDir = "/Users/wgillis/Dropbox (HMS)/mendeley"

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	// archive, err := filepath.Abs(filepath.Join(comparePath, ".."))
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close() // this runs after main() exits

	done := make(chan bool)
	go func() {
		for { // runs like a while loop in other langs
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					dirs, erro := filepath.Glob(filepath.Join(comparePath, "*.pdf"))
					if erro != nil {
						log.Fatal(erro)
					}
					file, err := os.Open(event.Name)
					if err != nil {
						log.Fatal(err)
					}
					filebytes := make([]byte, 4096)
					_, err = file.Read(filebytes)
					if err != nil {
						log.Fatal(err)
					}
					chksum := sha1.Sum(filebytes)
					file.Close()
					notequal := true
					for i := 0; i < len(dirs) && notequal; i++ {
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
						chksum2 := sha1.Sum(fb2)
						if bytes.Equal(chksum2[:], chksum[:]) {
							log.Println("These two files are equal:", dirs[i], event.Name)
							notequal = true
							newPath := filepath.Join(comparePath, "..", "paper-archive", filepath.Base(dirs[i]))
							log.Println("Moving", dirs[i], "to", newPath)
							os.Rename(dirs[i], newPath)
						}
					}
					log.Println("created file:", event.Name)
				}
			case err:= <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

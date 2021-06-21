package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
)

func main() {
	baseDir := "./my-files"
	items, _ := ioutil.ReadDir(baseDir)
	report := make(map[string](chan int))
	waitGroup := &sync.WaitGroup{}
	var mutex = &sync.Mutex{}

	for _, item := range items {

		if item.IsDir() {
			waitGroup.Add(1)
			go iterateOnDirectory(filepath.Join(baseDir, item.Name()), report, waitGroup, mutex)
		} else {
			fileExtension := filepath.Ext(item.Name())
			waitGroup.Add(1)
			go incrimentExtension(fileExtension, report, waitGroup, mutex)
		}

	}
	waitGroup.Wait()
	for key, val := range report {
		fmt.Println(key, ": ", <-val)
	}
}

func iterateOnDirectory(path string, rep map[string](chan int), wg *sync.WaitGroup, mutex *sync.Mutex) {
	subitems, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, subitem := range subitems {
		if subitem.IsDir() {
			wg.Add(1)
			go iterateOnDirectory(filepath.Join(path, subitem.Name()), rep, wg, mutex)
		} else {
			wg.Add(1)
			fileExtension := filepath.Ext(subitem.Name())
			go incrimentExtension(fileExtension, rep, wg, mutex)
		}
	}
	wg.Done()
}

func incrimentExtension(extensionName string, rep map[string](chan int), wg *sync.WaitGroup, mutex *sync.Mutex) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, isExist := rep[extensionName]; !isExist {
		rep[extensionName] = make(chan int, 2)
		rep[extensionName] <- 1
		wg.Done()
		return
	}
	count := <-rep[extensionName]
	rep[extensionName] <- count + 1
	wg.Done()
}

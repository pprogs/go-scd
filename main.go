package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	config  *conf
	lg      *log.Logger
	watcher *fsnotify.Watcher
	err     error
	logFile *os.File
)

func main() {

	confFlag := flag.String("C", "config.json", "Path to config json file")
	rulesFlag := flag.String("R", "rules.json", "Path to rules json file")
	flag.Parse()

	if *confFlag == "" || *rulesFlag == "" {
		flag.PrintDefaults()
		return
	}

	confFile := *confFlag
	rulesFile := *rulesFlag

	config = &conf{mu: &sync.RWMutex{}}
	if _, err = ReadJSON(confFile, config); err != nil {
		CheckIfError(err)
	}

	if config.LogFile != "" {
		if logFile, err = os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
			CheckIfError(err)
		}
		lg = log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags|log.Lshortfile)
	} else {
		lg = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	}

	lg.Printf("config: %v\n", config)

	reloadRules(rulesFile)

	if watcher, err = fsnotify.NewWatcher(); err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	runWatcher(rulesFile)

	runHTTP(config, processGitHooks)
}

package main

import (
	"runtime"
	"sync"
)

type ifHandler func(*command, *ruleIf) bool

var (
	mutex  = &sync.Mutex{}
	ifHndl = make(map[string]ifHandler)
)

func init() {
	ifHndl["file exists"] = ifHandlerFileExists
}

func execCommands(commands []command) {

	//there should be only one!!! (c)
	//... for now ...

	mutex.Lock()
	defer mutex.Unlock()

	lg.Printf("Executing commands\n")

	for idx := range commands {

		runtime.Gosched()
		com := &commands[idx]

		lg.Printf("Running... %s\n", com.Description)

		if com.Ifs != nil && !checkIf(com) {
			lg.Printf("Skipping...\n")
			continue
		}

		switch com.Type {
		case "git":
			lg.Printf("Run git command...\n")
			gitCommand(com)
		case "exec":
			lg.Printf("Run exec command...\n")
			execCommand(com)
		default:
			lg.Printf("Unknown command = %s\n", com.Type)
		}
	}

	lg.Printf("Done commands\n")
}

func checkIf(c *command) bool {

	var ret bool
	var ok bool
	var h ifHandler

	lastRet := true

	for idx := range c.Ifs {
		r := &c.Ifs[idx]

		if r.Logic == "OR" && idx > 0 && lastRet {
			return true
		}
		if h, ok = ifHndl[r.Condition]; !ok {
			continue
		}

		ret = h(c, r)

		if idx > 0 {
			if r.Logic == "" || r.Logic == "AND" {
				ret = lastRet && ret
			}
			if r.Logic == "OR" {
				ret = lastRet || ret
			}
		}

		lastRet = ret
	}

	return lastRet
}

func ifHandlerFileExists(c *command, r *ruleIf) bool {
	ret, _ := PathExists(r.Value)
	return ret
}

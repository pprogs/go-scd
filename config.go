package main

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type ruleIf struct {
	Logic     string `json:"logic"`
	Condition string `json:"cond"`
	Value     string `json:"value"`
}

type rule struct {
	Name     string `json:"rule_name"`
	Trigger  string `json:"rule_trigger"`
	Commands string `json:"command_set"`

	GitRepo   string `json:"git_repo"`
	GitBranch string `json:"git_branch"`
}

type command struct {
	Type        string `json:"type"`
	Description string `json:"descr"`

	//ifs
	Ifs []ruleIf `json:"if"`

	//git
	Token      string `json:"token"`
	LocalRepo  string `json:"localRepo"`
	RemoteRepo string `json:"remoteRepo"`
	Branch     string `json:"branch"`

	//exec
	WorkDir string   `json:"workDir"`
	CmdName string   `json:"cmdName"`
	CmdArgs []string `json:"cmdArgs"`
}

type conf struct {
	HTTPAddr string `json:"http_endpoint"`
	GitHook  string `json:"git_webhook"`
	LogFile  string `json:"log_file"`

	rules    []rule
	commands map[string][]command

	mu *sync.RWMutex
}

func runWatcher(ruleFile string) {
	if err = watcher.Add(ruleFile); err != nil {
		CheckIfError(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					lg.Printf("rules file watcher error 1\n")
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					go reloadRules(ruleFile)
				}
			case <-watcher.Errors:
				lg.Printf("rules file watcher error 2\n")
				return
			}
		}
	}()
}

func reloadRules(ruleFile string) {

	lg.Printf("reloading rules file %s\n", ruleFile)

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	lg.Printf("CWD: %s\n", dir)

	type tmp struct {
		Rules    []rule               `json:"rules"`
		Commands map[string][]command `json:"commands"`
	}

	t := &tmp{}

	if _, err := ReadJSON(ruleFile, t); err != nil {
		CheckIfError(err)
	}

	config.mu.Lock()
	config.rules = t.Rules
	config.commands = t.Commands
	config.mu.Unlock()
}

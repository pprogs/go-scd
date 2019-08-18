package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type hookPayload struct {
	Ref string `json:"ref"`

	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		URL      string `json:"git_url"`
		Branch   string `json:"default_branch"`
	} `json:"repository"`
}

//getHookHandler returns http handler than will decode github payload and
//call processHookPayload function
func getHookHandler(processHookPayload func(*hookPayload)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pl, err := handleHook(w, r)

		lg.Printf("check hook pl\n")

		if err == nil && pl != nil && processHookPayload != nil {
			lg.Printf("run processHookPayload\n")
			go processHookPayload(pl)
		} else {
			lg.Printf("something wrong with hook!! %v\n", err)
		}
	}
}

//handleHook http handler for hook endpoint
func handleHook(w http.ResponseWriter, r *http.Request) (*hookPayload, error) {

	lg.Printf("Hook handler run\n")

	defer io.Copy(ioutil.Discard, r.Body)

	var pl = &hookPayload{}
	err := json.NewDecoder(r.Body).Decode(pl)

	if err == io.EOF {
		lg.Printf("empty body in hook\n")
		return nil, nil
	}
	CheckIfError(err)

	event := r.Header.Get("X-GitHub-Event")
	delivery := r.Header.Get("X-GitHub-Delivery")
	sig := r.Header.Get("X-Hub-Signature")

	lg.Printf("Payload: (%v)\n", pl)
	lg.Printf("Headers: (%s) (%s) (%s)\n", event, delivery, sig)

	return pl, nil
}

//processGitHooks searches for rule for this payload and exec commands if found one
func processGitHooks(pl *hookPayload) {

	lg.Printf("Processing git hook (%s) (%s)\n", pl.Repository.FullName, pl.Ref)

	if pl.Ref == "" || pl.Repository.FullName == "" {
		lg.Printf("wrong payload. Exiting\n")
		return
	}

	config.mu.RLock()
	defer config.mu.RUnlock()

	for idx := range config.rules {
		r := &config.rules[idx]

		if r.Trigger != "git_hook" {
			continue
		}

		if pl.Repository.FullName == r.GitRepo && pl.Ref == r.GitBranch {

			lg.Printf("Found rule (%s) for git hook\n", r.Name)
			if r.Commands != "" {
				execCommands(config.commands[r.Commands])
				return
			}
		}
	}

	lg.Printf("No rules found\n")
}

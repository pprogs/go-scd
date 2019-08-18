package main

import (
	"net/http"
)

func runHTTP(c *conf, processHookPayload func(*hookPayload)) {

	http.HandleFunc(c.GitHook, getHookHandler(processHookPayload))

	lg.Printf("Starting server on %s\n", c.HTTPAddr)
	//daemon.SdNotify(false, "READY=1")

	err := http.ListenAndServe(c.HTTPAddr, nil)
	CheckIfError(err)
}

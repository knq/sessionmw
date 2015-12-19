// example/redis/redis.go
package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/knq/sessionmw"
	"github.com/knq/sessionmw/redisstore"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
)

func main() {
	rs, err := redisstore.New("redis://localhost:6379")
	if err != nil {
		log.Fatalln(err)
	}

	// create session middleware
	sessConfig := &sessionmw.Config{
		Secret:      "LymWKG0UvJFCiXLHdeYJTR1xaAcRvrf7",
		BlockSecret: "NxyECgzxiYdMhMbsBrUcAAbyBuqKDrpp",

		Store: rs,
	}

	// create goji mux and add sessionmw
	mux := goji.NewMux()
	mux.UseC(sessConfig.Handler)

	// add handlers
	mux.HandleFuncC(pat.Get("/set/:name"), func(ctxt context.Context, res http.ResponseWriter, req *http.Request) {
		val := pat.Param(ctxt, "name")
		sessionmw.Set(ctxt, "name", val)
		http.Error(res, fmt.Sprintf("name saved as '%s'.", html.EscapeString(val)), http.StatusOK)
	})
	mux.HandleFuncC(pat.Get("/"), func(ctxt context.Context, res http.ResponseWriter, req *http.Request) {
		var name = "[no name]"
		val, _ := sessionmw.Get(ctxt, "name")
		if n, ok := val.(string); ok {
			name = n
		}
		http.Error(res, fmt.Sprintf("hello '%s'", html.EscapeString(name)), http.StatusOK)
	})

	// serve
	http.ListenAndServe(":3000", mux)
}

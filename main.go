package main

import (
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADRESS", false, ":8080", "Bind address for the server")

func main() {

	env.Parse()

	// create a new serve mux
	sm := mux.NewRouter()

	_ = sm

}

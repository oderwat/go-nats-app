package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // needs to be better for a real app

	// check the system for int64 fitting into int (so we can convert int64 to int safely)
	if uint64(^uint(0)) < ^uint64(0) {
		panic("int does not fit int64")
	}

	// This is behind a WASM build tag so the business logic of the UI does not increase
	// the size of our server executable.
	// It also has an empty component with its own Render in "server.go"
	Front()

	// this concludes the part which goes into the front-end
	app.RunWhenOnBrowser()

	// I declare this here because of the "logic" to find it in
	// the backend code still is a mindbender for me :)
	ah := &app.Handler{
		Name:         "Go-Nats-App",
		Lang:         "de",
		Author:       "Hans Raaf - METATEXX GmbH",
		Title:        "Go Nats App",
		Description:  "NATS in a PWA",
		Image:        "/web/logo-512.png",
		LoadingLabel: "Loading...",
		Icon: app.Icon{
			Default:  "/web/logo-192.png",
			Large:    "/web/logo-512.png",
			Maskable: "/web/logo-192.png",
		},
		Styles: []string{
			"/web/index.css",
		},
		CacheableResources: []string{
			"/web/logo.svg",
		},
	}

	// this will depend on the target (wasm or not wasm) and
	// it starts the servers if it is not the wasm target.
	Back(ah)
}

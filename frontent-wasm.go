// Our empty version of the httpServer for usage with the wasm target
// this way we will not include any of the related code
//go:build wasm

package main

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"go-nats-app/front"
)

// This is a dummy for the AppServer (back end) code so it does
// not get included in the WASM code
func Back(_ *app.Handler) {
}

// This is the actual frontend. It
func Front() {
	front.Create()
}

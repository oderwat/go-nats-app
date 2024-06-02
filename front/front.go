package front

import (
	"github.com/goombaio/namegenerator"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/nats-io/nats.go"
	"nhooyr.io/websocket"
	"strconv"
	"time"
)

type appControl struct {
	app.Compo
	whoami    string     // keeps our name
	avatar    string     // data base64 jpeg image
	nc        *nats.Conn // is the nats server connection
	input     string     // the current input line
	messages  []string   // the last 10 messages we received
	echoCount int        // we count how many echos we sent on the echo service
}

var _ app.Initializer = (*appControl)(nil) // Verify the implementation
var _ app.Mounter = (*appControl)(nil)     // Verify the implementation
var _ app.AppUpdater = (*appControl)(nil)  // Verify the implementation

// OnInit is called before the component gets mounted
// This is before Render was called for the first time
func (uc *appControl) OnInit() {
	app.Log("OnInit")
	uc.whoami = namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate()
	// a blank image
	uc.avatar = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
}

// OnMount gets called when the component is mounted
// This is after Render was called for the first time
func (uc *appControl) OnMount(ctx app.Context) {
	app.Log("OnMount")
	go func() {
		var ncw WasmNatsConnectionWrapper
		app.Log("NatsConnect")
		var err error
		uc.nc, err = nats.Connect("localhost:8502", // our websocket port
			nats.Name("PWA-"+uc.whoami),
			nats.SetCustomDialer(ncw),
		)
		if err != nil {
			app.Logf("Native Go Connect did fail: %#v", err)
			return
		}
		app.Log("Nats connected through websocket and netConn wrapper!")
		// now we add a subscription to the chat.room
		// Subscribe to the subject
		_, err = uc.nc.Subscribe("chat.room", func(msg *nats.Msg) {
			// Print message data
			ctx.Dispatch(func(ctx app.Context) {
				if len(uc.messages) < 10 {
					uc.messages = append(uc.messages, string(msg.Data))
				} else {
					uc.messages = append(uc.messages[1:], string(msg.Data))
				}
			})
		})

		// grab an avatar
		msg, err := uc.nc.Request("govatar.female", []byte(""), 200*time.Millisecond)
		if err != nil {
			app.Logf("govatar request error %s", err)
		} else {
			ctx.Dispatch(func(ctx app.Context) {
				// set it
				uc.avatar = string(msg.Data)
			})
		}

		// tell them we are here
		err = uc.nc.Publish("chat.say", []byte(uc.whoami+" entered the room"))
		if err != nil {
			app.Logf("Publish entry message error %s", err)
		}

		// create an echo service in this browser :)
		_, err = uc.nc.Subscribe("echo."+uc.whoami, func(msg *nats.Msg) {
			_ = msg.Respond(msg.Data)
			ctx.Dispatch(func(ctx app.Context) {
				uc.echoCount++
			})
		})
		ctx.Update()
	}()

	app.Window().GetElementByID("inp").Call("focus")

	// Notice: We do not care about OnDismount which would be needed
	// when working with a more complex app.
}

func (uc *appControl) OnAppUpdate(ctx app.Context) {
	// This will be called when the service worker gets updated
	// With this demo it checks for this about every 10 seconds
	// This would normally just a fallback for using NATS to tell
	// the frontend that there is a new version of the app runnning
	if ctx.AppUpdateAvailable() {
		// hard update immediately :)
		// You would not do that with a real app but ask the user or
		// give a countdown before it updates
		ctx.Reload()
	}
}

func (uc *appControl) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Go-Nats-App"),
		app.Div().Text(func() string {
			if uc.nc == nil {
				return "Not connected to the nats server"
			} else {
				return "Connected to: " + uc.nc.ConnectedServerName()
			}
		}()),
		app.Div().Body(app.Img().Src(uc.avatar).Width(250).Height(250)),
		app.H4().Text("Chat:"),
		app.Form().Body(
			app.Div().Body(
				app.Span().Text(uc.whoami+": "),
				app.Input().Value(uc.input).ID("inp").OnInput(uc.ValueTo(&uc.input)),
			),
		).OnSubmit(func(ctx app.Context, e app.Event) {
			e.PreventDefault()
			if uc.nc != nil {
				err := uc.nc.Publish("chat.say", []byte(uc.whoami+": "+uc.input))
				if err != nil {
					app.Logf("Publish error %s", err)
				}
				uc.input = "" // clear the message entry
			}
		}),
		app.Range(uc.messages).Slice(func(i int) app.UI {
			return app.Div().Text(uc.messages[len(uc.messages)-1-i])
		}),
		app.H4().Text("For an echo use:"),
		app.Pre().Text(`nats -s 127.0.0.1:8501 req echo.`+uc.whoami+` '{{ Random 10 100 }}'`),
		app.Div().Text("Echos sent: "+strconv.Itoa(uc.echoCount)),
	)
}

func Create() {
	app.RouteWithRegexp("/.*", app.NewZeroComponentFactory(&appControl{}))
	// add a very simple update checker that checks for updates every 5 seconds
	// this lets us modify the code and restart the server more easily
	// For production this should be changed to a longer interval.
	intervalUpdater(time.Second * 5)
}

func intervalUpdater(delay time.Duration) {
	time.AfterFunc(delay, func() {
		app.Log("checking for update")
		app.TryUpdate()
		intervalUpdater(delay)
	})
}

type WasmNatsConnectionWrapper struct {
	ws *websocket.Conn
}

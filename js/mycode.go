package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"time"
	"encoding/json"
	"github.com/dennisfrancis/websocket"
)

var jQuery = jquery.NewJQuery
var document js.Object
var tableobj js.Object

func main() {

	jQuery("document").Ready(startup)
	document =  js.Global.Get("window").Get("document")
}

func startup() {
	create_table()
}

/*
Notable features of datatable : sort, save-state, defer-render, search, extensions
*/
func create_table() {
	ws  :=  websocket.New("ws://localhost:1234/tabledata")
	
	var tableobj js.Object

	first := true

	ws_onopen := func(obj js.Object) {
		println("WS opened")
	}

	ws_onmessage := func(buf []byte) {

		t0 := time.Now()
		dataobj := map[string]interface{}{}
		println("WS received byte stream of length = ", len(buf))
		err := json.Unmarshal(buf, &dataobj)

		if err != nil {
			println("unmarshal error : ", err.Error())
			return
		}

		t1 := time.Now()
		println("pre-rendering time = ", t1.Sub(t0).Seconds()/1000.0, "ms")

		if first {
			tableobj = jQuery("#mytable").Underlying().Call("DataTable", dataobj)
			first = false
		} else {
			//print(tableobj.Get("clear"))
			//print(tableobj.Get("row"))
			//print(tableobj.Get("row").Get("add"))

			/* Below works correctly */
			tableobj.Call("clear")
			tableobj.Get("rows").Call("add", dataobj["data"]).Call("draw")
			
		}
		
		t2 := time.Now()
		println("rendering time = ", t2.Sub(t1).Seconds()/1000.0, "ms")

	}

	// Set the callbacks
	ws.OnOpen(ws_onopen)
	ws.OnMessage(ws_onmessage)
}

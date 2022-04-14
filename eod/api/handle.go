package api

import (
	"encoding/json"
	"net/http"

	"github.com/Nv7-Github/Nv7Haven/eod/api/data"
	"github.com/gorilla/websocket"
)

func (a *API) Handle(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	id := ""

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if id == "" {
			id = string(message)
			continue
		}

		// Parse
		var msg data.Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, data.RSPError(err.Error()).JSON())
			continue
		}

		// Eval
		out := data.RSPError("Bad request")
		switch msg.Method {
		case data.MethodGuild:
			out = a.MethodGuild(msg.Params, id)
		}

		// Respond
		err = conn.WriteMessage(websocket.TextMessage, out.JSON())
		if err != nil {
			return
		}
	}
}

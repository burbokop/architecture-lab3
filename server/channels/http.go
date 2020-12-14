package channels

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/burbokop/architecture-lab3/server/tools"
)

// Channels HTTP handler.
type HttpHandlerFunc http.HandlerFunc

// HttpHandler creates a new instance of channels HTTP handler.
func HttpVmListHandler(store *Store) HttpHandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Println("vmlist responce:", r)

			res, err := store.ListChannels()
			if err != nil {
				log.Printf("Error making query to the db: %s", err)
				tools.WriteJsonInternalError(rw)
				return
			}
			tools.WriteJsonOk(rw, res)

		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func HttpConnectDiskHandler(store *Store) HttpHandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			fmt.Println("connect disk responce:", r)
			handleChannelCreate(r, rw, store)
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func handleChannelCreate(r *http.Request, rw http.ResponseWriter, store *Store) {
	var c Channel
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Printf("Error decoding channel input: %s", err)
		tools.WriteJsonBadRequest(rw, "bad JSON payload")
		return
	}
	err := store.CreateChannel(c.Name)
	if err == nil {
		tools.WriteJsonOk(rw, &c)
	} else {
		log.Printf("Error inserting record: %s", err)
		tools.WriteJsonInternalError(rw)
	}
}

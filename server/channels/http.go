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
func HttpVmListHandler(store *VMStorage) HttpHandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Println("VM list request recieved.")
			res, err := store.ListVirtualMachines()
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

func HttpConnectDiskHandler(store *VMStorage) HttpHandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			fmt.Println("Disk connecting request recieved.")
			var c ConnectionRequest
			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
				log.Printf("Error decoding request input: %s", err)
				tools.WriteJsonBadRequest(rw, "bad JSON payload")
				return
			}
			err := store.ConnectDisk(&c)
			if err == nil {
				res, err := store.ListVirtualMachines()
				if err == nil {
					tools.WriteJsonOk(rw, res)
				} else {
					log.Println(err)
				}
			} else {
				log.Printf("Error inserting record: %s", err)
				tools.WriteJsonInternalError(rw)
			}
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

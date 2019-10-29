/*
Copyright 2019 The Kyma Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	adapter "github.com/antoineco/kyma-event-sources/adapter"
	cloudEvents "github.com/cloudevents/sdk-go/pkg/cloudevents"
	websocket "github.com/gorilla/websocket"
)

const (
	defaultPort    = 8080
	defaultSinkURL = "http://event-publish-service.kyma-system.svc.cluster.local:8080/v2/events"
)

var (
	port    = flag.Int("port", defaultPort, "Default port at which the websocket connection is served")
	sink    = flag.String("sink", defaultSinkURL, "URL to redirect the Publish Request")
	address = fmt.Sprintf(":%d", *port)

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func readMessage(webSocketConnection *websocket.Conn, receiveAdapter *adapter.Adapter) {
	for {
		_, message, error := webSocketConnection.ReadMessage()
		if error != nil {
			log.Panic("Error occured while reading message.")
			break
		}

		cloudEventV01 := &cloudEvents.EventContextV01{}
		json.Unmarshal(message, cloudEventV01)

		receiveAdapter.HandleEvent(cloudEventV01)
	}
}

func checkOrigin(r *http.Request) bool {
	return true
}

func main() {
	flag.Parse()
	upgrader.CheckOrigin = checkOrigin

	log.Printf("Starting WebSocket Server at port: %d \n", defaultPort)

	ra := adapter.New(*sink)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		webSocketConnection, error := upgrader.Upgrade(w, r, nil)

		if error != nil {
			log.Println(error)
			log.Panic("Unable to upgrade the HTTP server connection to the WebSocket protocol")
		}

		log.Println("Client Connected...")

		//test connection
		if err := webSocketConnection.WriteMessage(websocket.TextMessage, []byte("Hello Client")); err != nil {
			log.Panicf("Unable to write message to the connected client.")
		}

		readMessage(webSocketConnection, ra)
	})
	log.Fatal(http.ListenAndServe(address, nil))
}

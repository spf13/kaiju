// Copyright © 2014 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License./

package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"labix.org/v2/mgo/bson"

	"github.com/googollee/go-socket.io"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var clients map[string]string = map[string]string{}

var sio *socketio.SocketIOServer

func socketIOInit() {
	config := socketio.Config{}
	config.HeartbeatTimeout = 60 * 5 // 5 min
	config.ClosingTimeout = 60 * 5   // 5 min

	sio = socketio.NewSocketIOServer(&config)

	sio.On("connect", onConnect)
	sio.On("disconnect", onDisconnect)
	sio.On("postComment", onPostComment)
	sio.On("getComments", onGetComments)
	sio.On("pong", onPong)
}

func socketServerRun() {
	target := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))
	jww.FEEDBACK.Println("Socket.IO server listening on", target)

	ch := make(chan string)
	defer close(ch)

	go func() {
		err := http.ListenAndServe(target, sio)
		ch <- err.Error()
	}()

	go func() {
		time.Sleep(5 * time.Second)
		for {
			fmt.Printf("[%s] Connected clients: %d.\n", time.Now().Format(time.Kitchen), len(clients))
			time.Sleep(1 * time.Minute)
		}
	}()

	fmt.Printf("[%s] Server listening on %s.\n", time.Now().Format(time.Kitchen), target)
	fmt.Println(<-ch)
}

func onConnect(ns *socketio.NameSpace) {
	fmt.Println("connected:", ns.Id(), " in channel ", ns.Endpoint())
	clients[ns.Id()] = ""

	sio.Broadcast("ping", string("happy times"))
}

func onDisconnect(ns *socketio.NameSpace) {
	fmt.Println("disconnected:", ns.Id(), " in channel ", ns.Endpoint())
	if _, ok := clients[ns.Id()]; ok {
		delete(clients, ns.Id())
	}
}

func onPong(ns *socketio.NameSpace, message string) {
	fmt.Println("pingged")
	var jsonMap map[string]string
	err := json.Unmarshal([]byte(message), &jsonMap)
	if err != nil {
		panic(err)
	}

	sio.Broadcast("ping", string("happy times are here again"))
	fmt.Println("pingged", jsonMap)

	sio.Broadcast("message", string("happy times"))
}

func onGetComments(ns *socketio.NameSpace, message string) {
	jww.INFO.Println("Comment request Received", message)

	var jsonMap map[string]string
	err := json.Unmarshal([]byte(message), &jsonMap)

	fmt.Println(jsonMap)

	//forumStr := jsonMap["forum"]
	forumStr := "5346e494331583002c7de60e"
	if bson.IsObjectIdHex(forumStr) == false {
		jww.ERROR.Printf("`forum` is not valid. Received: `%v`\n", forumStr)
	}
	forum := bson.ObjectIdHex(forumStr)

	comments, err := GetAllComments(db, forum, jsonMap["page"])
	if err != nil {
		jww.ERROR.Printf("Error: %v\n", err)
	}

	bComments, err := json.Marshal(comments)
	if err != nil {
		jww.ERROR.Printf("Error: %v\n", err)
	}

	fmt.Println(string(bComments))
	ns.Emit("commentsFor", string(bComments))

	if err != nil {
		jww.ERROR.Println(err.Error())
	}
}

func onPostComment(ns *socketio.NameSpace, message string) {
	fmt.Println("Comment Received")

	var jsonMap map[string]string
	err := json.Unmarshal([]byte(message), &jsonMap)

	fmt.Printf("%#v\n", jsonMap)

	if err != nil {
		jww.ERROR.Println(err.Error())
	}

	forumStr := "5346e494331583002c7de60e"
	fullname := jsonMap["fullname"]
	email := jsonMap["email"]

	forumId, page, body, parentId, err := PostCommentResource(forumStr,
		//jsonMap["forum"],
		jsonMap["page"],
		jsonMap["body"],
		jsonMap["parent"])

	comment, err := PostComment(db,
		fullname,
		email,
		forumId,
		page,
		body,
		parentId)

	if err != nil {
		jww.ERROR.Println(err.Error())
	}

	if err != nil {
		panic(err)
	}

	bComment, err := json.Marshal(comment)

	if err != nil {
		jww.ERROR.Printf("Error: %v\n", err)
	}

	sio.Broadcast("commentPosted", string(bComment))
	fmt.Printf("%#v\n", string(bComment))
}

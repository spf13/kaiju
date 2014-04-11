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
	"html/template"
	"net/http"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/codegangsta/martini"
	"github.com/spf13/viper"

	//"github.com/martini-contrib/oauth2"
	//"github.com/martini-contrib/sessions"
)

func martiniInit() {
	m := martini.New()
	r := martini.NewRouter()

	m.Map(db)
	//m.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("secret123"))))
	//m.Use(oauth2.Github(&oauth2.Options{
	// 	ClientId:	  "64a641523f31dd5bfe4b",
	// 	ClientSecret: "4fe3fbbca262835c424ca6a80aec6c6cb4228037",
	// 	RedirectURL:  "http://localhost:2714/github_callback",
	// 	Scopes:		  []string{"user:email"},
	//}))

	//r.Get("/", index)
	r.Get("/comments/:forum/:page", GetAllCommentsResource)
	r.Post("/comment", PostCommentHandler)
	//r.Get("/github_redirect", RedirectUrl)
	//r.Get("/github_callback", func(request *http.Request) string {
	// 	request.ParseForm()
	// 	return fmt.Sprintf("%+v", request)
	//})

	m.Use(martini.Static("ui"))
	m.Action(r.Handle)

	fmt.Println("Running on port " + viper.GetString("port"))

	sio.Handle("/", m)
	//http.Handle("/", m)
	//http.ListenAndServe(":"+viper.GetString("port"), m)
}

func index() template.HTML {
	return template.HTML(`
<html>
<body>
<form action="/comment" method="POST">
  <input type="hidden" name="userId" value="5346e476331583002c7de60d" />
  <input type="hidden" name="forum" value="5346e494331583002c7de60e" />
  <input type="hidden" name="parent" value="534718de080fcc8979000002" />
  <input type="hidden" name="page" value="dans sc2 blog post" />
  <input type="hidden" name="body" value="sc2 is awesome. Toss OP" />
  <input type="hidden" name="timestamp" value="123456789" />
  <input type="submit" />
</form>
<br/>
<a href="https://github.com/login/oauth/authorize?client_id=64a641523f31dd5bfe4b&redirect_uri=http://localhost:2314/github_callback&scope=user:email">Login</a>
</body>
</html>
`)
}

func GetAllCommentsResource(db *mgo.Database, parms martini.Params) string {
	forumStr := parms["forum"]
	if bson.IsObjectIdHex(forumStr) == false {
		return fmt.Sprintf("`forum` is not valid. Received: `%v`", forumStr)
	}
	forum := bson.ObjectIdHex(forumStr)
	page := parms["page"]

	comments, err := GetAllComments(db, forum, page)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	bytes, err := json.Marshal(comments)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	return string(bytes)
}

func PostCommentHandler(request *http.Request, db *mgo.Database) string {
	if err := request.ParseForm(); err != nil {
		return err.Error()
	}

	username := request.FormValue("username")
	useremail := request.FormValue("useremail")
	forumIdStr := request.FormValue("forum")
	page := request.FormValue("page")
	body := request.FormValue("body")
	parentIdStr := request.FormValue("parent")

	forumId, page, body, parentId, err := PostCommentResource(
		forumIdStr,
		page,
		body,
		parentIdStr)

	if err != nil {
		return err.Error()
	}

	comment, err := PostComment(db,
		username,
		useremail,
		forumId,
		page,
		body,
		parentId)

	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Accepted. Comment ID: %v", comment.Id)
}

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
)

func martiniInit() {
	m := martini.New()
	r := martini.NewRouter()

	m.Map(db)

	r.Get("/", index)
	r.Get("/comments/:forum/:page", GetAllCommentsResource)
	r.Post("/comment", PostCommentResource)
	r.Get("redirect_url", RedirectUrl)

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

	userIdStr := request.FormValue("userId")
	forumIdStr := request.FormValue("forum")
	page := request.FormValue("page")
	body := request.FormValue("body")
	parentIdStr := request.FormValue("parent")

	userId, forumId, page, body, parentId, err := PostCommentResource(userIdStr,
		forumIdStr,
		page,
		body,
		parentIdStr)

	if err != nil {
		return err.Error()
	}

	comment, err := PostComment(db,
		userId,
		forumId,
		page,
		body,
		parentId)

	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Accepted. Comment ID: %v", comment.Id)
}

func RedirectUrl(session sessions.Session) string {
  session.Set("some", "thing")
  return "OK"
})


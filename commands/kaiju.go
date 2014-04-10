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
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/martini"
	"github.com/spf13/cobra"
	"github.com/spf13/kaiju/models"
	"github.com/spf13/viper"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var Verbose bool
var Port int
var DBName string
var DBPort int
var DBHost string
var CfgFile string
var db *mgo.Database

func Execute() {
	AddCommands()
	err := Root.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}

func AddCommands() {
	Root.AddCommand(InitializeFixturesCmd)
}

var Root = &cobra.Command{
	Use:   "kaiju",
	Short: "kaiju is an open source comment server",
	Long:  `kaiju is an open source comment server`,
	Run:   RootRun,
}

func RootRun(cmd *cobra.Command, args []string) {
	m := martini.New()
	r := martini.NewRouter()

	m.Map(db)

	r.Get("/", index)
	r.Get("/comments/:forum/:post", comments)
	r.Put("/comment", PostCommentResource)

	m.Action(r.Handle)

	fmt.Println("Running on port " + strconv.Itoa(Port))
	http.ListenAndServe(":"+strconv.Itoa(Port), m)
}

func db_init() {
	session, err := mgo.Dial(viper.GetString("dbhost") + ":" + viper.GetString("dbport"))
	if err != nil {
		panic(err)
	}
	db = session.DB(viper.GetString("dbname"))
	defer session.Close()
}

func init() {
	Root.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")
	Root.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	Root.Flags().IntVarP(&Port, "port", "p", 2714, "port number to run on")
	Root.Flags().StringVarP(&DBName, "dbname", "d", "kaiju", "name of the database")
	Root.Flags().IntVar(&DBPort, "dbport", 27017, "port to access mongoDB")
	Root.Flags().StringVar(&DBHost, "dbhost", "localhost", "host where mongoDB is")

	viper.SetConfigName(CfgFile)
	viper.AddConfigPath("./")
	viper.ReadInConfig()

	viper.Set("port", Port)
	viper.Set("dbname", DBName)
	viper.Set("dbport", DBPort)
	viper.Set("dbhost", DBHost)
	viper.Set("verbose", Verbose)

	db_init()
}

var InitializeFixturesCmd = &cobra.Command{
	Use:   "initializeFixtures",
	Short: "Initialize Fixtures, throw away",
	Long:  ``,
	Run:   InitializeFixtures,
}

func InitializeFixtures(cmd *cobra.Command, args []string) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("test").C("users")

	u1 := &models.User{
		Id:       bson.NewObjectId(),
		FullName: "foo",
	}

	err = c.Insert(u1)

	if err != nil {
		panic(err)
	}

	result := models.User{}
	err = c.Find(bson.M{"fullname": "foo"}).One(&result)

	if err != nil {
		panic(err)
	}

	fmt.Println("Phone:", result.FullName)

	fmt.Println("Fixtures Initialized")
}

func index() template.HTML {
	return	template.HTML(`
<html>
<body>
<form action="/comment" method="PUT">
  <input type="hidden" name="forum" value="dans website" />
  <input type="hidden" name="page" value="starcraft is awesome.html" />
  <input type="hidden" name="email" value="danny.gottlieb@gmail.com" />
  <input type="hidden" name="timestamp" value="123456789" />
  <input type="hidden" name="body" value="I think starcraft is cool, but toss OP." />
  <input type="submit" />
</form>
</body>
</html>
`)
}

func comments(database *mgo.Database, parms martini.Params) (int, string) {
	forum := parms["forum"]
	post := parms["post"]
	return http.StatusOK, strings.Join([]string{"ah, yeah: ", forum, post}, " ")
}

func PostCommentResource(request *http.Request, db mgo.Database) string {
	if err := request.ParseForm(); err != nil {
		return err.Error()
	}

	userIdStr := request.FormValue("userId")
	if bson.IsObjectIdHex(userIdStr) == false {
		return fmt.Sprintf("`userId` is not valid. Received: `%v`", userIdStr)
	}
	userId := bson.ObjectIdHex(userIdStr);

	forumIdStr := request.FormValue("forum")
	if bson.IsObjectIdHex(forumIdStr) == false {
		return fmt.Sprintf("`forum` is not valid. Received: `%v`", forumIdStr)
	}
	forumId := bson.ObjectIdHex(forumIdStr)

	parentIdStr := request.FormValue("parent")
	var parent *bson.ObjectId
	switch {
	case parentIdStr == "":
	case bson.IsObjectIdHex(parentIdStr):
		parentIdObj := bson.ObjectIdHex(parentIdStr)
		parent = &parentIdObj
	default:
		return fmt.Sprintf("`parent` is not valid. Received: `%v`", parentIdStr)
	}

	comment, err := PostComment(db,
		userId,
		forumId,
		request.FormValue("page"),
		request.FormValue("body"),
		parent)

	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Accepted. Comment ID: %v", comment.Id)
}

func PostComment(db mgo.Database, userId bson.ObjectId, forumId bson.ObjectId,
	page string, body string, parentId *bson.ObjectId) (*models.Comment, error) {

	users := db.C("users")
	comments := db.C("comments")

	user := &models.User{}
	err := users.Find(bson.M{"_id": userId}).One(user)
	switch err {
	case nil:
	case mgo.ErrNotFound:
		return nil, fmt.Errorf("User not found.")
	default:
		return nil, fmt.Errorf("Error finding user. Err: %v", err)
	}

	commentId := bson.NewObjectId()
	comment := &models.Comment{
		Id: commentId,
		User: models.CommentUser{
			UserId: userId,
			FullName: user.FullName,
			Email: user.Email,
		},
		Forum: forumId,
		Page: page,
		Body: body,
		Parent: parentId,
	}

	if err := comments.Insert(comment); err != nil {
		return nil, fmt.Errorf("Database error. Err: %v", err)
	}

	return comment, nil
}

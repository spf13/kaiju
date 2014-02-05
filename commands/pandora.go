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
    "github.com/codegangsta/martini"
    "github.com/spf13/cobra"
    "labix.org/v2/mgo"
    "net/http"
    "os"
    "strconv"
    "strings"
)

var Verbose bool
var Port int
var DBName string
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
    Use:   "pandora",
    Short: "Pandora is an open source comment server",
    Long:  `Pandora is an open source comment server`,
    Run:   RootRun,
}

var InitializeFixturesCmd = &cobra.Command{
    Use:   "initializeFixtures",
    Short: "Initialize Fixtures, throw away",
    Long:  ``,
    Run:   InitializeFixtures,
}

func InitializeFixtures(cmd *cobra.Command, args []string) {
    fmt.Println("Fixtures Initialized")
}

func RootRun(cmd *cobra.Command, args []string) {
    m := martini.New()
    r := martini.NewRouter()

    m.Map(db)

    r.Get("/", index)
    r.Get("/comments/:forum/:post", comments)

    m.Action(r.Handle)

    fmt.Println("Running on port " + strconv.Itoa(Port))
    http.ListenAndServe(":"+strconv.Itoa(Port), m)
}

func db_init() {
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        panic(err)
    }
    db = session.DB(DBName)
    defer session.Close()
}

func init() {
    Root.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
    Root.Flags().IntVarP(&Port, "port", "p", 2714, "port number to run on")
    Root.Flags().StringVarP(&DBName, "dbname", "d", "pandora", "name of the database")
    db_init()
}

func index() string {
    return "What up?"
}

func comments(session mgo.Session, parms martini.Params) (int, string) {
    forum := parms["forum"]
    post := parms["post"]
    return http.StatusOK, strings.Join([]string{"ah, yeah: ", forum, post}, " ")
}

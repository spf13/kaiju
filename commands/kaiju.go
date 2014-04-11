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
	"os"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/kaiju/models"
	"github.com/spf13/viper"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var Verbose bool
var Port int
var Host string
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
	socketIOInit()
	martiniInit()
	socketServerRun()
}

func db_init() {
	session, err := mgo.Dial(viper.GetString("dbhost") + ":" + viper.GetString("dbport"))
	if err != nil {
		panic(err)
	}
	db = session.DB(viper.GetString("dbname"))
}

func init() {
	Root.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is path/config.yaml|json|toml)")
	Root.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	Root.Flags().IntVarP(&Port, "port", "p", 2714, "port number to run on")
	Root.Flags().StringVarP(&DBName, "dbname", "d", "kaiju", "name of the database")
	Root.Flags().StringVarP(&Host, "host", "h", "", "host to run on")
	Root.Flags().IntVar(&DBPort, "dbport", 27017, "port to access mongoDB")
	Root.Flags().StringVar(&DBHost, "dbhost", "localhost", "host where mongoDB is")

	viper.SetConfigName(CfgFile)
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/kaiju/")
	err := viper.ReadInConfig()
	if err != nil {
		jww.ERROR.Println("Config not found... using only defaults, stuff may not work")
	}

	viper.Set("port", Port)
	viper.Set("host", Host)
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

package commands

import (
	"fmt"
	"html/template"
	"net/http"

	"encoding/json"

	"github.com/codegangsta/martini"
	"github.com/spf13/kaiju/models"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

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

func PostCommentResource(userIdStr string, forumIdStr string, pageStr string, bodyStr string, parentIdStr string) (userId bson.ObjectId, forumId bson.ObjectId, page string, body string, parentId *bson.ObjectId, err error) {
	page = pageStr
	body = bodyStr

	if bson.IsObjectIdHex(userIdStr) == false {
		err = fmt.Errorf("`userId` is not valid. Received: `%v`", userIdStr)
		return
	}
	userId = bson.ObjectIdHex(userIdStr)

	if bson.IsObjectIdHex(forumIdStr) == false {
		err = fmt.Errorf("`forum` is not valid. Received: `%v`", forumIdStr)
		return
	}
	forumId = bson.ObjectIdHex(forumIdStr)

	switch {
	case parentIdStr == "":
	case bson.IsObjectIdHex(parentIdStr):
		parentIdObj := bson.ObjectIdHex(parentIdStr)
		parentId = &parentIdObj
	default:
		err = fmt.Errorf("`parent` is not valid. Received: `%v`", parentIdStr)
		return
	}

	return
}

func PostComment(db *mgo.Database, userId bson.ObjectId, forumId bson.ObjectId,
	page string, body string, parentId *bson.ObjectId) (*models.Comment, error) {

	users := db.C("users")
	forums := db.C("forums")
	comments := db.C("comments")

	user := &models.User{}
	err := users.Find(bson.M{"_id": userId}).One(user)
	switch err {
	case nil:
	case mgo.ErrNotFound:
		return nil, fmt.Errorf("User not found. Id: `%v`", userId)
	default:
		return nil, fmt.Errorf("Error finding user. Err: %v", err)
	}

	err = forums.Find(bson.M{"_id": forumId}).One(make(bson.M))
	switch err {
	case nil:
	case mgo.ErrNotFound:
		return nil, fmt.Errorf("Forum not found. Id: `%v`", forumId)
	default:
		return nil, fmt.Errorf("Error finding forum. Err: %v", err)
	}

	ancestors := make([]bson.ObjectId, 0, 0)
	if parentId != nil {
		parentComment := &models.Comment{}
		if err := comments.Find(bson.M{"_id": *parentId}).One(parentComment); err != nil {
			return nil, fmt.Errorf("Parent comment not found. Id: `%v`", parentId)
		}

		ancestors = parentComment.Ancestors
		ancestors = append(ancestors, *parentId)
	}

	commentId := bson.NewObjectId()
	comment := &models.Comment{
		Id: commentId,
		User: models.CommentUser{
			Id:       userId,
			FullName: user.FullName,
			Email:    user.Email,
		},
		Forum:     forumId,
		Page:      page,
		Body:      body,
		Parent:    parentId,
		Ancestors: ancestors,
	}

	if err := comments.Insert(comment); err != nil {
		return nil, fmt.Errorf("Database error. Err: %v", err)
	}

	return comment, nil
}

func _getAllCommentsWithQuery(db *mgo.Database, query bson.M) ([]*models.Comment, error) {
	ret := make([]*models.Comment, 0, 0)
	comments := db.C("comments")
	if err := comments.Find(query).All(&ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func GetAllComments(db *mgo.Database,
	forumId bson.ObjectId,
	page string) ([]*models.Comment, error) {

	return _getAllCommentsWithQuery(db, nil)
}

func GetTopLevelComments(db *mgo.Database,
	forumId bson.ObjectId,
	page string) ([]*models.Comment, error) {

	return _getAllCommentsWithQuery(db,
		bson.M{
			"Forum":  forumId,
			"Page":   page,
			"Parent": nil,
		})
}

func GetCommentsWithAncestor(db *mgo.Database,
	forumId bson.ObjectId,
	page string,
	ancestorId bson.ObjectId) ([]*models.Comment, error) {

	return _getAllCommentsWithQuery(db,
		bson.M{
			"Forum":     forumId,
			"Page":      page,
			"Ancestors": ancestorId,
		})
}

func GetCommentsSinceTime(db *mgo.Database,
	forumId bson.ObjectId,
	page string,
	since bson.ObjectId) ([]*models.Comment, error) {

	return _getAllCommentsWithQuery(db,
		bson.M{
			"Forum": forumId,
			"Page":  page,
			"_id":   bson.M{"$gt": since},
		})
}

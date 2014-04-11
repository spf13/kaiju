package commands

import (
	"fmt"

	"github.com/spf13/kaiju/models"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func PostCommentResource(forumIdStr string, pageStr string, bodyStr string, parentIdStr string) (forumId bson.ObjectId, page string, body string, parentId *bson.ObjectId, err error) {
	page = pageStr
	body = bodyStr

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

func PostComment(db *mgo.Database, fullname string, email string, forumId bson.ObjectId,
	page string, body string, parentId *bson.ObjectId) (*models.Comment, error) {

	users := db.C("users")
	forums := db.C("forums")
	comments := db.C("comments")

	user := &models.User{}
	err := users.Find(bson.M{"email": email}).One(user)
	switch err {
	case nil:
	case mgo.ErrNotFound:
		user = &models.User{Id: bson.NewObjectId(), FullName: fullname, Email: email}
		fmt.Println("Inserting User", user)
		err = users.Insert(user)
		fmt.Println(err)
		if err = users.Find(bson.M{"email": email}).One(user); err != nil {
			return nil, fmt.Errorf("Error inserting user. Err: %v", err)
		}
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
			Id:       user.Id,
			FullName: user.FullName,
			Email:    user.Email,
		},
		Forum:     forumId,
		Page:      page,
		Timestamp: bson.Now(),
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
	if err := comments.Find(query).Sort("_id").All(&ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func GetAllComments(db *mgo.Database,
	forumId bson.ObjectId,
	page string) ([]*models.Comment, error) {

	return _getAllCommentsWithQuery(db, bson.M{
		"forum": forumId,
		"page":  page,
	})
}

func GetTopLevelComments(db *mgo.Database,
	forumId bson.ObjectId,
	page string) ([]*models.Comment, error) {

	return _getAllCommentsWithQuery(db,
		bson.M{
			"forum":  forumId,
			"page":   page,
			"parent": nil,
		})
}

func GetCommentsWithAncestor(db *mgo.Database,
	forumId bson.ObjectId,
	page string,
	ancestorId bson.ObjectId) ([]*models.Comment, error) {

	return _getAllCommentsWithQuery(db,
		bson.M{
			"forum":     forumId,
			"page":      page,
			"ancestors": ancestorId,
		})
}

func GetCommentsSinceTime(db *mgo.Database,
	forumId bson.ObjectId,
	page string,
	since bson.ObjectId) ([]*models.Comment, error) {

	return _getAllCommentsWithQuery(db,
		bson.M{
			"forum": forumId,
			"page":  page,
			"_id":   bson.M{"$gt": since},
		})
}

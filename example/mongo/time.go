package main

import (
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tal-tech/go-zero/core/stores/mongo"
)

type Roster struct {
	Id          bson.ObjectId `bson:"_id"`
	CreateTime  time.Time     `bson:"createTime"`
	Classroom   mgo.DBRef     `bson:"classroom"`
	Member      mgo.DBRef     `bson:"member"`
	DisplayName string        `bson:"displayName"`
}

func main() {
	model := mongo.MustNewModel("localhost:27017", "blackboard", "roster")
	for i := 0; i < 1000; i++ {
		session, err := model.TakeSession()
		if err != nil {
			log.Fatal(err)
		}

		var roster Roster
		filter := bson.M{"_id": bson.ObjectIdHex("587353380cf2d7273d183f9e")}
		fmt.Println(model.GetCollection(session).Find(filter).One(&roster))
		model.PutSession(session)
	}

	time.Sleep(time.Hour)
}

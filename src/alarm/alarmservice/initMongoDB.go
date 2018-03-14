package alarmservice

import mgo "gopkg.in/mgo.v2"
import "log"

const (
	collection = "vehicle_warning"
)

var dbc *mgo.Collection

func initMongoDB() error {
	session, err := mgo.Dial("")
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	db := session.DB("parse_vehicle")
	dbc = db.C(collection)
	return nil
}

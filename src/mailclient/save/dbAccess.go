package save

import (
	"log"

	"gopkg.in/mgo.v2"
)

type dbAccess struct {
	host   string
	port   string
	dbName string

	session *mgo.Session
}

/*
DBAccess - handles access to DB
*/
type DBAccess interface {
	StartSession() bool
	CloseSession() bool
	GetCollection(collectionName string) *mgo.Collection
}

/*
NewDBAccess - Creates access to MongoDB
*/
func NewDBAccess(host, port, dbName string) DBAccess {
	return &dbAccess{host: host, port: port, dbName: dbName}
}
func (acccess *dbAccess) StartSession() bool {
	if acccess.session != nil {
		return true
	}
	dbHost := acccess.host + ":" + acccess.port
	session, err := mgo.Dial(dbHost)
	if err != nil {
		log.Printf("Error during access to DB:%s, %v", dbHost, err)
		return false
	}
	log.Println("Connected to DB:", dbHost)
	session.SetMode(mgo.Monotonic, true)
	acccess.session = session
	return true
}
func (acccess *dbAccess) CloseSession() bool {
	log.Println("Closing DB session")
	acccess.session.Close()
	acccess.session = nil
	return true
}

func (acccess *dbAccess) GetCollection(collectionName string) *mgo.Collection {
	return acccess.session.DB(acccess.dbName).C(collectionName)
}

package save

import (
	"gopkg.in/mgo.v2"
)

type dbAccess struct {
	host   string
	port   string
	dbName string

	session *mgo.Session
}
type DBAccess interface {
	startSession() bool
	closeSession() bool
	getCollection(collectionName string) *mgo.Collection
}

func NewDBAccess(host, port, dbName string) DBAccess {
	return &dbAccess{host: host, port: port, dbName: dbName}
}
func (acccess *dbAccess) startSession() bool {
	if acccess.session != nil {
		return true
	}
	session, err := mgo.Dial(acccess.host + ":" + acccess.port)
	if err != nil {
		return false
	}
	session.SetMode(mgo.Monotonic, true)
	acccess.session = session
	return true
}
func (acccess *dbAccess) closeSession() bool {
	acccess.session.Close()
	acccess.session = nil
	return true
}

func (acccess *dbAccess) getCollection(collectionName string) *mgo.Collection {
	return acccess.session.DB(acccess.dbName).C(collectionName)
}
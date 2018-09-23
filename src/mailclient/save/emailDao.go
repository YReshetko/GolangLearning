package save

import "mailclient/domain"

type emailDao struct {
	collectionName string
	dbAccess       DBAccess
}

type EmailDao interface {
	save(data domain.EmailData)
}

func NewDao(dbName, collectionName string) EmailDao {
	return &emailDao{collectionName, NewDBAccess("localhost", "27017", dbName)}
}

func (dao *emailDao) save(data domain.EmailData) {
	ok := dao.dbAccess.startSession()
	if ok {
		defer dao.dbAccess.closeSession()
		collection := dao.dbAccess.getCollection(dao.collectionName)
		collection.Insert(&data)
	}
}

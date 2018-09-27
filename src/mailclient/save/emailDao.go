package save

import (
	"mailclient/domain"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type emailDao struct {
	collection *mgo.Collection
}

type EmailDao interface {
	Save(data domain.EmailData) error
	FindByUid(uid uint32) *domain.EmailData
	FindLatest(count int) []domain.EmailData
}

func NewDao(collection *mgo.Collection) EmailDao {
	return &emailDao{collection}
}

func (dao *emailDao) Save(data domain.EmailData) error {
	return dao.collection.Insert(&data)
}
func (dao *emailDao) FindByUid(uid uint32) *domain.EmailData {
	data := domain.EmailData{}
	dao.collection.Find(bson.M{"uid": uid}).One(&data)
	if data.Uid == 0 {
		return nil
	}
	return &data
}
func (dao *emailDao) FindLatest(count int) []domain.EmailData {
	/*
		var results []Person
		err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)
	*/
	var out []domain.EmailData
	dao.collection.Find(bson.M{}).Limit(count).All(&out)
	return out
}

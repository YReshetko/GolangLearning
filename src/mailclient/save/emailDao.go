package save

import (
	"mailclient/domain"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type emailDao struct {
	collection *mgo.Collection
}

/*
EmailDao - Data access to emails DB
*/
type EmailDao interface {
	Save(data domain.EmailData) error
	FindByUid(uid uint32) *domain.EmailData
	FindLatest(count int) []domain.EmailData
	FindByDateRange(from, to time.Time) []domain.EmailData
}

/*
NewDao - create new DAO
*/
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

func (dao *emailDao) FindByDateRange(from, to time.Time) []domain.EmailData {
	var out []domain.EmailData
	dao.collection.Find(bson.M{"date": bson.M{"$gte": from, "$lte": to}}).Sort("date").All(&out)
	return out
}
func (dao *emailDao) FindLatest(count int) []domain.EmailData {
	/*
		var results []Person
		err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)
	*/
	var out []domain.EmailData
	dao.collection.Find(nil).Sort("date").Limit(count).All(&out)
	return out
}

package save

import (
	"fmt"
	"mailclient/domain"
	"strings"
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
	UpdateCollection(collection *mgo.Collection)
	Save(data domain.EmailData) error
	FindByUid(uid uint32) (*domain.EmailData, error)
	FindLatest(count int) ([]domain.EmailData, error)
	FindByDateRange(from, to time.Time) ([]domain.EmailData, error)
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

func (dao *emailDao) UpdateCollection(collection *mgo.Collection) {
	dao.collection = collection
}
func (dao *emailDao) FindByUid(uid uint32) (*domain.EmailData, error) {
	data := domain.EmailData{}
	err := dao.collection.Find(bson.M{"uid": uid}).One(&data)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "not found") {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if data.Uid == 0 {
		return nil, nil
	}
	return &data, nil
}

func (dao *emailDao) FindByDateRange(from, to time.Time) ([]domain.EmailData, error) {
	var out []domain.EmailData
	err := dao.collection.Find(bson.M{"date": bson.M{"$gte": from, "$lte": to}}).Sort("date").All(&out)
	return out, err
}
func (dao *emailDao) FindLatest(count int) ([]domain.EmailData, error) {
	/*
		var results []Person
		err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)
	*/
	var out []domain.EmailData
	//Sort descending - Sort("-date")
	//Sort ascending - Sort("date")
	err := dao.collection.Find(nil).Sort("-date").Limit(count).All(&out)
	return out, err
}

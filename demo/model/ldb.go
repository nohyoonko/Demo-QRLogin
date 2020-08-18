package model

import (
	"errors"
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type DBHandler interface {
	Create(key string, value string) error
	Update(key string, value string) error
	SelectOne(key string) (Data, error)
	SelectAll() ([]Data, error)
	SelectList(prefix string) ([]Data, error)
	Delete(key string) error
	Close() error
}

type ldbHandler struct {
	db *leveldb.DB
}

//외부 참조 가능하려면 first character is capital
type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewDBHandler(filepath string, env bool) DBHandler {
	return newLDBHandler(filepath, env)
}

func newLDBHandler(filepath string, env bool) DBHandler {
	//if env == true, available to delete file
	if env {
		//if filepath already exists, delete all files in the filepath
		if _, err := os.Stat(filepath); !os.IsNotExist(err) {
			err := os.RemoveAll(filepath)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	dbPath := filepath
	database, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return &ldbHandler{db: database}
}

func (l *ldbHandler) Create(key string, value string) error {
	//search every data in db
	iter := l.db.NewIterator(nil, nil)
	for iter.Next() {
		//key is duplicated
		if key == string(iter.Key()) {
			return errors.New("key is duplicated")
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return err
	}

	//put data
	err = l.db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		return err
	}
	return nil
}

func (l *ldbHandler) Update(key string, value string) error {
	//put만 해도 overwrite 되지만, .ldb에 남아있는 듯 하여 delete->put
	//batch로 쓰는게 더 빠르다고 해서 batch 사용
	batch := new(leveldb.Batch)
	batch.Delete([]byte(key))
	batch.Put([]byte(key), []byte(value))
	err := l.db.Write(batch, nil)
	if err != nil {
		return err
	}
	return nil
}

func (l *ldbHandler) SelectOne(key string) (Data, error) {
	value, err := l.db.Get([]byte(key), nil)

	var getData Data
	getData.Key = key
	getData.Value = string(value)

	if err != nil {
		return getData, err
	}
	return getData, nil
}

func (l *ldbHandler) SelectAll() ([]Data, error) {
	datas := []Data{}

	//search every data in db
	iter := l.db.NewIterator(nil, nil)
	for iter.Next() {
		var getData Data
		getData.Key = string(iter.Key())
		getData.Value = string(iter.Value())
		datas = append(datas, getData)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return datas, err
	}
	return datas, nil
}

func (l *ldbHandler) SelectList(prefix string) ([]Data, error) {
	selectDatas := []Data{}

	//prefix
	iter := l.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	for iter.Next() {
		var tmp Data
		tmp.Key = string(iter.Key())
		tmp.Value = string(iter.Value())
		selectDatas = append(selectDatas, tmp)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return selectDatas, err
	}
	return selectDatas, nil
}

func (l *ldbHandler) Delete(key string) error {
	err := l.db.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}

func (l *ldbHandler) Close() error {
	err := l.db.Close()
	if err != nil {
		return err
	}
	return nil
}

package model

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type DBHandler interface {
	Create(key string, value string) error //key 중복->err
	SelectAll() ([]data, error)
	SelectOne(key string) (data, error)
	Upadte(key string, value string) error
	Delete(key string) error
	SelectList(prefix string) ([]data, error)
	Close() error
}

type ldbHandler struct {
	db *leveldb.DB
}

type data struct {
	key   string `json:"key"`
	value string `json:"value"`
}

func NewDBHandler(filepath string) DBHandler {
	return newLDBHandler(filepath)
}

func newLDBHandler(filepath string) DBHandler {
	dbPath := "dbPath"
	database, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		panic(err)
	}
	return &ldbHandler{db: database} //return err
}

// byte slice -> string
func decode(b []byte) string {
	return string(b[:len(b)])
}

func (l *ldbHandler) Create(key string, value string) error {
	//put data
	err := l.db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		return err
	}
	return nil
}

func (l *ldbHandler) SelectAll() ([]data, error) {
	datas := []data{}

	//search every data in db
	iter := l.db.NewIterator(nil, nil)
	for iter.Next() {
		var d data
		d.key = string(iter.Key())
		d.value = string(iter.Value())
		datas = append(datas, d)
		fmt.Println("read", d.key, d.value) //key, value
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return datas, err
	}
	return datas, nil
}

func (l *ldbHandler) SelectOne(key string) (data, error) {
	value, err := l.db.Get([]byte(key), nil)

	var d data
	d.key = key
	d.value = string(value)

	if err != nil {
		return d, err
	}
	return d, nil
}

func (l *ldbHandler) Upadte(key string, value string) error {
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

func (l *ldbHandler) SelectList(prefix string) ([]data, error) {
	selectDatas := []data{}

	//prefix
	iter := l.db.NewIterator(util.BytesPrefix([]byte("new")), nil)
	for iter.Next() {
		var d data
		d.key = string(iter.Key())
		d.value = string(iter.Value())
		selectDatas = append(selectDatas, d)
		fmt.Println("select", d.key, d.value) //key, value
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

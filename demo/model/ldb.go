package model

import (
	"errors"
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util" //BytesPrefix([]byte)
)

type DBHandler interface {
	Create(key string, value string) error
	/* Read */
	SelectOne(key string) (Data, error)       //key에 해당하는 Data만 읽어오기
	SelectAll() ([]Data, error)               //DB 내에 있는 모든 Data 읽어오기
	SelectList(prefix string) ([]Data, error) //prefix로 시작하는 모든 Data 읽어오기
	Update(key string, value string) error
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

/* newLDBHandler: DB 열기 */
func newLDBHandler(filepath string, env bool) DBHandler {
	//if env == true, 폴더 삭제 권한 있으므로 삭제 가능
	if env {
		//폴더가 이미 존재한다면, 폴더를 삭제하기
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

/* Create: key가 중복인지 확인 후 Put */
func (l *ldbHandler) Create(key string, value string) error {
	//DB에 있는 모든 데이터 순회하면서 중복 확인
	iter := l.db.NewIterator(nil, nil)
	for iter.Next() {
		//key가 중복되는 것이 있다면 error를 return
		if key == string(iter.Key()) {
			return errors.New("key is duplicated")
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return err
	}

	//Put
	err = l.db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		return err
	}
	return nil
}

/* SelectOne: key에 해당하는 Data를 return */
func (l *ldbHandler) SelectOne(key string) (Data, error) {
	value, err := l.db.Get([]byte(key), nil)

	var getData Data
	getData.Key = key
	getData.Value = string(value) //value는 byte array

	if err != nil {
		return getData, err
	}
	return getData, nil
}

/* SelectAll: DB에 있는 모든 데이터를 List로 만들어서 return */
func (l *ldbHandler) SelectAll() ([]Data, error) {
	datas := []Data{}

	//DB에 있는 모든 데이터를 순회하면서 List에 저장
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

/* SelectList: Prefix를 갖는 모든 데이터를 List로 만들어서 return */
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

/* Update: DB에 있는 기존의 key를 Delete하고 다시 Put */
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

/* Delete: 데이터 삭제 */
func (l *ldbHandler) Delete(key string) error {
	err := l.db.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}

/* Close: DB 닫기 */
func (l *ldbHandler) Close() error {
	err := l.db.Close()
	if err != nil {
		return err
	}
	return nil
}

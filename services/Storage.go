package services

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/api/option"
)

type DBInterface interface {
	Init() error
	SetMeta(c int64, ss int64)
	GetMeta() (int64, int64)
	Close() error
}

type DBLocal struct {
	db *leveldb.DB
}

func (dbl *DBLocal) Init() error {
	var err error
	path, err := filepath.Abs("storage")
	dbl.db, err = leveldb.OpenFile(path, nil)
	if err != nil {
		fmt.Printf("Cant open storage %v", err)
		return err
	}
	return nil
}

func (dbl *DBLocal) Close() error {
	return dbl.db.Close()
}

func (dbl *DBLocal) GetMeta() (int64, int64) {
	fc, err := dbl.db.Get([]byte("files_converted"), nil)
	if err != nil {
		log.Printf("Cant read from storage %v", err)
		return 0, 0
	}

	ss, err := dbl.db.Get([]byte("space_saved"), nil)
	if err != nil {
		log.Printf("Cant read from storage %v", err)
		return 0, 0
	}

	return int64(binary.LittleEndian.Uint64(fc)), int64(binary.LittleEndian.Uint64(ss))

}

func (dbl *DBLocal) SetMeta(c int64, ss int64) {
	var err error
	bC := make([]byte, 8)
	binary.LittleEndian.PutUint64(bC, uint64(c))
	err = dbl.db.Put([]byte("files_converted"), []byte(bC), nil)
	if err != nil {
		log.Printf("Cant write to storage %v", err)
	}

	bSS := make([]byte, 8)
	binary.LittleEndian.PutUint64(bSS, uint64(ss))
	err = dbl.db.Put([]byte("files_converted"), []byte(bSS), nil)
	if err != nil {
		log.Printf("Cant write to storage %v", err)
	}
}

type DBFirebase struct {
	driver *firebase.App
	contex context.Context
	client *firestore.Client
}

func (dbf *DBFirebase) Init() error {
	var err, err_c error
	dbf.contex = context.Background()
	p, _ := filepath.Abs("configs/keys/webp-ninja-firebase-adminsdk-jii5f-fef7522670.json")
	opt := option.WithCredentialsFile(p)
	dbf.driver, err = firebase.NewApp(dbf.contex, nil, opt)
	if err != nil {
		log.Printf("error initializing firebase storage: %v", err)
		return err
	}
	dbf.client, err_c = dbf.driver.Firestore(dbf.contex)
	if err_c != nil {
		fmt.Errorf("error initializing firebase client: %v", err)
		return err
	}
	return nil
}

func (dbf *DBFirebase) SetMeta(c int64, ss int64) {
	if dbf.driver == nil && dbf.client == nil {
		log.Fatal("DB Must be initialised first")
	}
	_, w_err := dbf.client.Collection("metainfo").Doc("_stats_").Set(dbf.contex, map[string]interface{}{
		"files_converted": c,
		"space_saved":     ss,
	})
	if w_err != nil {
		log.Fatalf("Cant write to storage %v", w_err)
	}
}

func (dbf *DBFirebase) GetMeta() (c int64, ss int64) {
	if dbf.driver == nil && dbf.client == nil {
		log.Fatal("GetMeta: DB Must be initialised first")
	}
	dsnap, err := dbf.client.Collection("metainfo").Doc("_stats_").Get(dbf.contex)
	if err != nil {
		fmt.Errorf("cant read from firebase storage: %v", err)
	}
	m := dsnap.Data()
	c, _ = m["files_converted"].(int64)
	ss, _ = m["space_saved"].(int64)
	return
}

func (dbf *DBFirebase) Close() error {
	return dbf.client.Close()
}

func NewStorage() DBInterface {
	// fbs := DBFirebase{}
	// fbs.Init()
	// return &fbs

	l := DBLocal{}
	l.Init()
	return &l
}

var storage DBInterface
var once sync.Once

func GetStorage() DBInterface {
	once.Do(func() {
		storage = NewStorage()
	})
	return storage
}

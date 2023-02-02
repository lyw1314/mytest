package models

import (
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

type boltDb struct {
	db *bolt.DB
}

func NewBoltDb(path string, dbName string) (*boltDb, error) {

	db, err := bolt.Open(path+"/"+dbName, 0600, &bolt.Options{Timeout: 10 * time.Second}) //打开数据库文件 如果没有则创建
	if err != nil {
		return nil, err
	}
	return &boltDb{db: db}, nil
}

func (b *boltDb) CreateBucketIfNotExists(bucketName []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error { //读写事务
		_, err := tx.CreateBucketIfNotExists(bucketName) //根据视图名字创建   如果不存在则创建
		if err != nil {
			return fmt.Errorf("create bucket fail: %v", err)
		}
		return nil
	})
	return err
}

func (b *boltDb) Put(bucketName []byte, key, value []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName) //打开视图
		if b == nil {
			//异常处理
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			return fmt.Errorf("bucket not exist")
		}
		// put 数据库插入
		return b.Put(key, value)
	})
	return err
}
func (b *boltDb) Get(bucketName []byte, key []byte) ([]byte, error) {
	var ret []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName) //打开视图
		if b == nil {
			//异常处理
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			return fmt.Errorf("bucket not exist")
		}
		ret = b.Get(key) // Get查询
		//_ := b.Cursor()
		//fmt.Printf(" value=%s\n", v)
		return nil
	})
	return ret, err
}

func (b *boltDb) Delete(bucketName []byte, key []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName) //打开视图
		if b == nil {
			//异常处理
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			return fmt.Errorf("bucket not exist")
		}
		// Delete  删除
		return b.Delete([]byte(key))
	})
	return err
}

func (b *boltDb) Range(bucketName []byte) ([]map[string][]byte, error) {
	ret := make([]map[string][]byte, 0)
	err := b.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(bucketName)
		if b == nil {
			//异常处理
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			return fmt.Errorf("bucket not exist")
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			ret = append(ret, map[string][]byte{"key": k, "value": v})
		}
		return nil
	})
	return ret, err
}

func (b *boltDb) Close() error {
	return b.db.Close()
}

func isFileExist(fileName string) bool {
	// func Stat(name string) (FileInfo, error) {
	_, err := os.Stat(fileName)

	//os.IsExist不要使用，不可靠
	if os.IsNotExist(err) {
		return false
	}
	return true
}

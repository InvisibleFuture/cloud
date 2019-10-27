package main

import (
	"encoding/binary"
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	DB map[string]*leveldb.DB
	//CA map[string]*sync.Map
	CH map[string]chan int64

	TOKEN sync.Map
)

func init() {
	var err error

	DB = make(map[string]*leveldb.DB)
	//CA = make(map[string]*sync.Map)
	CH = make(map[string]chan int64)

	// 将 ID 0 作为计数器存储列, 不必再初始化一个单独的存储库
	dbList := []string{"counter"}
	for _, name := range dbList {
		DB[name], err = leveldb.OpenFile("data/"+name, nil)
		if err != nil { panic("OpenFile error") }
	}

	user := []string{"account", "password", "salt", "token"}
	for _, name := range user {
		DB[name], err = leveldb.OpenFile("data/"+name, nil)
		if err != nil { panic("OpenFile error") }
	}

	wiki := []string{"list", "index", "doc"}
	for _, name := range wiki {
		DB[name], err = leveldb.OpenFile("data/"+name, nil)
		if err != nil { panic("OpenFile error") }
	}

	chList := []string{"user", "project"}
	for _, name := range chList {
		CH[name] = make(chan int64)
		go func(name string, c chan int64) {
			var err error
			var sum int64

			buf := make([]byte, 8)
			data, err := DB[name].Get([]byte("0"), nil)
			if err != nil {
				binary.BigEndian.PutUint64(buf, uint64(sum))
				err = DB["counter"].Put([]byte(name), buf, nil)
				log.Println("计数器初始化", name)
				if err != nil {
					panic("DB Put error")
				}
			} else {
				sum = int64(binary.BigEndian.Uint64(data))
				log.Println("计数器载入", name, sum)
			}

			for {
				sum++
				c <- sum
				binary.BigEndian.PutUint64(buf, uint64(sum))
				err = DB["counter"].Put([]byte(name), buf, nil)
				if err != nil {
					panic("counter error")
				}
			}
		}(name, CH[name])
	}

	// 存储一个测试账户
    DB["account"].Put([]byte("Last"), []byte("1"), nil)
    DB["password"].Put([]byte("Last"), []byte("00000000"), nil)
}

package data

import (
	"encoding/binary"
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

// 使用 go 命令发起的计数器 (name, CH[name])
func Counter(name string, c chan int64) {
	var err error
	var sum int64
	buf := make([]byte, 8)
	data, err := DB[name].Get([]byte("0"), nil)
	if err != nil {
		binary.BigEndian.PutUint64(buf, uint64(sum))
		err = DB["counter"].Put([]byte(name), buf, nil)
		log.Println("计数器初始化", name)
		if err != nil { panic("DB Put error") }
	} else {
		sum = int64(binary.BigEndian.Uint64(data))
		log.Println("计数器载入", name, sum)
	}
	for {
		sum++
		c <- sum
		binary.BigEndian.PutUint64(buf, uint64(sum))
		err = DB["counter"].Put([]byte(name), buf, nil)
		if err != nil { panic("counter error") }
	}
}
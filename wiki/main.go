package wiki

import (
	"net/http"
)
import (
	"log"
	"sync"
	"strconv"
	"math/rand"
	"encoding/binary"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	DB_LIST  *leveldb.DB
	DB_DOC   *leveldb.DB
	DB_SALT     *leveldb.DB
	DB_TOKEN    *leveldb.DB

	CA_TOKEN    sync.Map

	CH_USERID chan int64
)

func init() {
	/*
	var err error
	DB_ACCOUNT, err = leveldb.OpenFile("data/wiki/list", nil)
	if err != nil { panic("OpenFile error data/wiki/list") }

	DB_PASSWORD, err = leveldb.OpenFile("data/wiki/password", nil)
	if err != nil { panic("OpenFile error data/wiki/password") }

	DB_SALT, err = leveldb.OpenFile("data/wiki/salt", nil)
	if err != nil { panic("OpenFile error data/wiki/salt") }

	DB_TOKEN, err = leveldb.OpenFile("data/wiki/token", nil)
	if err != nil { panic("OpenFile error data/wiki/token") }
	*/
}

func WikiCreate(w http.ResponseWriter, r *http.Request) {
	// 验证操作权限
	//get project(id)
	//if project(u.id) != true {
	//	return
	//}

	// 决定是否执行操作
	//ok := u.CreateProject(p)

	// 返回操作结果
}

func RewriteProject(w http.ResponseWriter, r *http.Request) {
	// 验证用户身份
	// 验证操作权限
	// 执行操作
	// 返回数据
}



// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
package user

import (
	"log"
	"sync"
	"strconv"
	"math/rand"
	"encoding/binary"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	DB_ACCOUNT  *leveldb.DB
	DB_PASSWORD *leveldb.DB
	DB_SALT     *leveldb.DB
	DB_TOKEN    *leveldb.DB

	CA_TOKEN    sync.Map

	CH_USERID chan int64
)

type Account struct {
	ID    string
	Token string
}

func init(){
	var err error
	DB_ACCOUNT, err = leveldb.OpenFile("data/user/account", nil)
	if err != nil { panic("OpenFile error data/user/account") }

	DB_PASSWORD, err = leveldb.OpenFile("data/user/password", nil)
	if err != nil { panic("OpenFile error data/user/password") }

	DB_SALT, err = leveldb.OpenFile("data/user/salt", nil)
	if err != nil { panic("OpenFile error data/user/salt") }

	DB_TOKEN, err = leveldb.OpenFile("data/user/token", nil)
	if err != nil { panic("OpenFile error data/user/token") }

	CH_USERID = make(chan int64)

	// 启动计数器
	go func(c chan int64){
		var err error
		var sum int64
		buf := make([]byte, 8)
		data, err := DB_ACCOUNT.Get([]byte("0"), nil)
		if err != nil {
			binary.BigEndian.PutUint64(buf, uint64(sum))
			err = DB_ACCOUNT.Put([]byte("0"), buf, nil)
			log.Println("UID 计数器初始化")
			if err != nil { panic("DB Put error") }
		} else {
			sum = int64(binary.BigEndian.Uint64(data))
			log.Println("UID 计数器载入, 当前", sum)
		}
		for {
			sum++
			log.Println("管道待通数据")
			c <- sum
			log.Println("管道已通数据")
			binary.BigEndian.PutUint64(buf, uint64(sum))
			err = DB_ACCOUNT.Put([]byte("0"), buf, nil)
			if err != nil { panic("counter error") }
		}
	}(CH_USERID)

	// 存储一个测试账户
	//DB_ACCOUNT.Put([]byte("Last"), []byte("1"), nil)
	//DB_PASSWORD.Put([]byte("1"), []byte("00000000"), nil)
}

func (a *Account)Register(account, password string) bool {
	if _, err := DB_ACCOUNT.Get([]byte(account), nil); err == nil { return false }
	if account == ""|| password == "" || len(account) > 32 || len(password) > 32 { return false }
	num := <-CH_USERID
	id := strconv.FormatInt(num,10)
	DB_ACCOUNT.Put([]byte(account), []byte(id), nil)
	DB_PASSWORD.Put([]byte(id), []byte(password), nil)
	// 写入默认随机用户名
	return true
}

func (a *Account)Signin(account, password string) bool {
	// 验证账号密码长度格式是否合法 尤其不能为空
	if account == "" || password == "" { return false }
	id, err := DB_ACCOUNT.Get([]byte(account), nil)
	if err != nil { return false }
	pw, err := DB_PASSWORD.Get(id, nil)
	if err != nil || string(pw) != password { return false }
	a.ID, a.Token = string(id), randSeq(32)
	CA_TOKEN.Store(a.ID, a.Token)
	return true
}

func (a *Account)Signout() { CA_TOKEN.Delete(a.ID) }

func Identity(uid, token string) bool {
	if t, ok := CA_TOKEN.Load(uid); !ok || t != token { return false }
	return true
}

func randSeq(n int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b { b[i] = letters[rand.Intn(len(letters))] }
	return string(b)
}
package user

import (
	"sync"
	"net/http"
	"math/rand"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
)


func init() {
	var err error
	ACCOUNT, err = leveldb.OpenFile("data/account", nil)
	if err != nil {panic("data/account OpenFile error")}

	PASSWORD, err = leveldb.OpenFile("data/password", nil)
	if err != nil {panic("data/password OpenFile error")}

	// 存储一个测试账户
	ACCOUNT.Put([]byte("Last"), []byte("1"), nil)
	PASSWORD.Put([]byte("Last"), []byte("00000000"), nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	account := r.FormValue("account")
	password := r.FormValue("password")
	// 验证账号格式, 只能是数字与大小写字母的组合
	// 验证账号是否存在
	_, err := ACCOUNT.Get([]byte(account), nil)
	if err == nil {
		Reply(w, -1, "账号已被占用")
		return
	}
	// 生成自增ID
	id :=

	ACCOUNT.Put([]byte(account), []byte(id), nil)
	ACCOUNT.Put([]byte(account), []byte(password), nil)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	account := r.FormValue("account")
	password := r.FormValue("password")

	pw, err := PASSWORD.Get([]byte(account), nil)
	if err != nil {
		Reply(w, -1, "没有找到账号")
		return
	}
	if string(pw) != password {
		Reply(w, -1, "密码不匹配")
		return
	}

	id, err := ACCOUNT.Get([]byte(account), nil)
	if err != nil {
		panic("账户没有对应角色")
	}

	uid := string(id)
	token := randSeq(32)
	TOKEN.Store(uid, token)

	Reply(w, 0, struct{Id , Token string}{Id: uid, Token: token})
}

func randSeq(n int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func autoid(name string, c chan int64) {
	buf := new(bytes.Buffer)

	data, err := AUTOID_DB.Get([]byte(name), nil)
	if err != nil {
		binary.Write(buf, binary.BigEndian, 0)
		data = buf.Bytes()
		err = AUTOID_DB.Put([]byte(name), data, nil)
		if err != nil {
			panic("AUTOID_DB " + name + " INIT ERROR")
		}
		log.Println("计数器初始化", name)
	}

	var sum int64
	binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &sum)
	log.Println("计数器", name, sum)
	for {
		sum++
		c <- sum
		buf = new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, sum)
		err = AUTOID_DB.Put([]byte(name), buf.Bytes(), nil)
		if err != nil {
			panic("AUTOID ++ ERROR")
		}
	}
}

func Reply(w http.ResponseWriter, code int, data interface{}) {
	rep := struct{
		Code int
		Data interface{}
	}{ Code: code, Data: data }
	js, err := json.MarshalIndent(rep, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)
}

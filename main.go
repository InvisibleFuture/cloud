package main

import (
	"os"
	"log"
	"time"
	"syscall"
	"strconv"
    "net/http"
	"os/signal"
	"math/rand"
	"encoding/json"
)

func main() {
	log.Println("start!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Reply(w, 0, "welcome")
	})

	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		account, password := r.FormValue("account"), r.FormValue("password")
		pw, err := DB["password"].Get([]byte(account), nil)
		if err != nil { Reply(w, -1, "没有找到账号"); return }
		if string(pw) != password { Reply(w, -1, "密码不匹配"); return }
		id, err := DB["account"].Get([]byte(account), nil)
		if err != nil { panic("账户没有对应角色") }
		uid, token := string(id), randSeq(32)
		TOKEN.Store(uid, token)
		Reply(w, 0, struct{Id , Token string}{Id: uid, Token: token})
	})

	http.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {
		uid, token := r.FormValue("uid"), r.FormValue("token")
		if t, ok := TOKEN.Load(uid); ok && t == token { TOKEN.Delete(uid); Reply(w, 0, "成功退出"); return }
		Reply(w, -1, "未登录")
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		account, password := r.FormValue("account"), r.FormValue("password")
		_, err := DB["account"].Get([]byte(account), nil)
		if err == nil { Reply(w, -1, "账号已被占用, 请更换"); return }
		// 判断账号和密码格式是否正确
		num := <-CH["user"]
		DB["account"].Put([]byte(account), []byte(strconv.FormatInt(num,10)), nil)
		DB["password"].Put([]byte(account), []byte(password), nil)
		Reply(w, 0, "注册成功")
	})

	s := &http.Server{
		Addr:           ":80",
		Handler:        http.DefaultServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Println(s.ListenAndServe())
		log.Println("server shutdown")
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	log.Println(s.Shutdown(nil))
	time.Sleep(time.Second * 1)

	log.Println("done.")
}

func randSeq(n int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b { b[i] = letters[rand.Intn(len(letters))] }
	return string(b)
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
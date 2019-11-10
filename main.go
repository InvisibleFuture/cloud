package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"./user"
)

func main() {
	log.Println("start!")

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	Reply(w, 0, "welcome")
	//})
	http.Handle("/", enforceXMLHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Reply(w, 0, "welcome")
	})))
	http.Handle("/signrenew", enforceXMLHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := user.Account{ID: r.FormValue("uid")}
		a.SignRenew()
		p := struct {
			ID     string
			Token  string
			Name   string
			Avatar string
		}{a.ID, a.Token, "Robert", "/12.png"}
		Reply(w, 0, p)
	})))
	http.Handle("/signout", enforceXMLHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := user.Account{ID: r.FormValue("uid")}
		a.Signout()
		Reply(w, 0, "成功退出")
	})))

	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		var a user.Account
		if ok := a.Signin(r.FormValue("account"), r.FormValue("password")); !ok {
			Reply(w, -1, "账户不存在或账户密码不匹配")
			return
		}
		p := struct {
			ID     string
			Token  string
			Name   string
			Avatar string
		}{a.ID, a.Token, "Robert", "/12.png"}
		Reply(w, 0, p)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var a user.Account
		if ok := a.Register(r.FormValue("account"), r.FormValue("password")); !ok {
			Reply(w, -1, "注册失败, 请检查账号是否正确")
			return
		}
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

func Reply(w http.ResponseWriter, code int, data interface{}) {
	rep := struct {
		Code int
		Data interface{}
	}{Code: code, Data: data}
	js, err := json.MarshalIndent(rep, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func enforceXMLHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 通过中间件进行身份验证
		uid, token := r.FormValue("uid"), r.FormValue("token")
		if ok := user.Identity(uid, token); !ok {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		next.ServeHTTP(w, r)
	})
}

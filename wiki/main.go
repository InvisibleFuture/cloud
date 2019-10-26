package wiki

import (
	//"net/http"
)

func CreateProject(w http.ResponseWriter, r *http.Request) {
	//var u := User{
	//	id: r.FromValue("id")
	//	token: r.FromValue("token")
	//}
	// 要实现身份验证需要先实现登录


	//get token(id)
	//if token != u.token {
	//	return
	//}
	// 并验证用户身份(应通过中间件实现)

	// 验证操作权限
	//get project(id)
	//if project(u.id) != true {
	//	return
	//}

	//var p Project
	// 验证用户操作对象的权限?

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
package wiki

import ()

type Project struct {
	name       string
	createTime string
	lastTime   string
	users      []string
	master     []string
}

// user Create project
// Server: Create project > init user

func (p *Project)Create(){}

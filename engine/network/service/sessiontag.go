package service

import "github.com/xhaoh94/goxh/consts"

type (
	//SessionTag 会话标签
	SessionTag struct {
		tag int
	}
)

//GetTag 获取标签
func (st *SessionTag) GetTag() int {
	return st.tag
}

//GetTagName 获取名字
func (st *SessionTag) GetTagName() string {
	switch st.tag {
	case consts.Connector:
		return "connector"
	case consts.Accept:
		return "accept"
	default:
		return ""
	}
}

// Code generated by goctl. DO NOT EDIT.
package types

type Request struct {
}

type Response struct {
}

type SearchReq struct {
	Name string `form:"name"`
}

type SearchReply struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

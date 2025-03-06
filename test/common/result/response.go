package result

import (
	"encoding/json"
	"errors"
)

var (
	ErrMsgEmpty = errors.New("response msg is empty")
)

// 定义响应码
const (
	SuccessCode = 200
)

// Response 统一响应结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewResponse 创建响应对象
// code 响应码
// msg 响应信息
// data 响应数据，如果发生错误的话，这里就为空
func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// ToJson 将响应对象转换为json格式
func (response *Response) ToJson() ([]byte, error) {
	if response.Msg == "" {
		return nil, ErrMsgEmpty
	}
	// 将结构体转换为json
	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, errors.New("json.Marshal failed: " + err.Error())
	}
	return jsonData, nil
}

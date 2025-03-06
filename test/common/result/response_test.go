package result_test

import (
	"binbin/common/result"
	"reflect"
	"testing"
)

func TestNewResponse(t *testing.T) {
	var cases = []struct {
		test_name string
		msg       string
		code      int
		data      interface{}
	}{
		{"common_test", "successfully", 200, "hello"},
		{"empty_msg_test", "", 699, nil},
	}
	var expected = []struct {
		Err  error
		Msg  string
		Code int
		Data interface{}
	}{
		{nil, "successfully", 200, "hello"},
		{result.ErrMsgEmpty, "", 200, nil},
	}
	for _, c := range cases {
		t.Run(c.test_name, func(t *testing.T) {
			response := result.NewResponse(c.code, c.msg, c.data)
			if reflect.DeepEqual(response, expected) {
				t.Errorf("expected: %v, got: %v", expected, response)
			}
		})
	}
}
func TestToJson(t *testing.T) {
	var cases = []struct {
		test_name string
		msg       string
		code      int
		data      interface{}
	}{
		{"common_test", "successfully", 200, "hello"},
	}
	for _, c := range cases {
		t.Run(c.test_name, func(t *testing.T) {
			response := result.NewResponse(c.code, c.msg, c.data)
			jsonData, err := response.ToJson()
			if err != nil {
				t.Errorf("json marshal failed: " + err.Error())
				return
			}
			t.Logf("%v", string(jsonData))
		})
	}
}

package spider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	/*	_, token, err := AuthTest("20206759", "xgs583719992")
		t.Log(token)
		assert.NoError(t, err)*/
	data := struct {
	}{}
	marshal, err := json.Marshal(data)
	request, err := http.NewRequest(http.MethodPost, "https://portal.neu.edu.cn/tp_up/up/subgroup/getLibraryInfo", bytes.NewBuffer(marshal))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	cookie := &http.Cookie{}
	cookie.Name = "tp_up"
	cookie.Value = "4kHakxVSZH63zXw6i72y16y6Tc09Qd2wWbeoFZ4jmmRqFV7PHWp4!380884168"
	request.AddCookie(cookie)
	res, err := http.DefaultClient.Do(request)
	assert.NoError(t, err)
	all, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)
	fmt.Println(string(all))
}

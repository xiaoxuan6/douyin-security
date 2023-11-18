package douyin_security

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient(t *testing.T) {
	account := NewAccount("xxx", "xxx", false)
	client := NewClient(account)

	token, err := client.GetAccessToken()
	assert.Nil(t, err)
	t.Log("token", token)
	client.SetAccessToken(token)

	response := client.TextVerify("我是一个好人")
	t.Log(response)
	assert.Equal(t, response.Code, 200, response.Msg)

	response1 := client.ImgVerify("https://i.v2ex.co/28iDJ9u5.png")
	t.Log(response1)
	assert.Equal(t, response1.Code, 200, response1.Msg)

	response3 := client.TextVerify("")
	t.Log(response3)
	assert.Equal(t, response3.Code, 400, response3.Msg)

	response4 := client.TextVerify("我是一个坏人, 操你妈, 约炮")
	t.Log(response4)
	assert.Equal(t, response4.Code, 500, response4.Msg)
	assert.Contains(t, response4.Msg, "违规")

	response5 := client.ImgVerify("https://nwzimg.wezhan.cn/contents/sitefiles2045/10229780/images/22696414.png")
	t.Log(response5)
	assert.Equal(t, response5.Code, 3, response5.Msg)
	assert.Equal(t, response5.Msg, "image download failed")
}

package douyin_security

import (
	"bytes"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"strings"
)

type Client struct {
	account     *Account
	accessToken string
	maxAttempts uint
}

func NewClient(account *Account) *Client {
	return &Client{
		account: account,
	}
}

func (c *Client) SetMaxAttempts(maxAttempts uint) {
	c.maxAttempts = maxAttempts
}

func (c *Client) fetchUir() string {
	if c.account.IsSandbox {
		return sandboxUri
	}

	return baseUir
}

func (c *Client) TextVerify(text string) *Response {
	if len(text) < 1 {
		return Fail(400, "text is empty")
	}

	params := bytes.NewBufferString(fmt.Sprintf(`{"tasks":[{"content":"%s"}]}`, text))
	response, err := c.postWithToken(fmt.Sprintf("%s%s", c.fetchUir(), textUrl), params)
	if err != nil {
		return Fail(500, err.Error())
	}

	if code := gjson.Get(response, "code").Int(); code != 0 {
		return Fail(int(code), gjson.Get(response, "message").String())
	}

	if code := gjson.Get(response, "data.0.code").Int(); code != 0 {
		return Fail(int(code), gjson.Get(response, "data.0.msg").String())
	}

	if hit := gjson.Get(response, "data.0.predicts.0.hit").Bool(); hit == true {
		return Fail(500, "文本包含违法违规内容")
	}

	return Success()
}

var modelName = map[string]string{
	"porn":                        "图片涉黄",
	"cartoon_leader":              "领导人漫画",
	"anniversary_flag":            "特殊标志",
	"sensitive_flag":              "敏感旗帜",
	"sensitive_text":              "敏感文字",
	"leader_recognition":          "敏感人物",
	"bloody":                      "图片血腥",
	"fandongtaibiao":              "未准入台标",
	"plant_ppx":                   "图片涉毒",
	"high_risk_social_event":      "社会事件",
	"high_risk_boom":              "爆炸",
	"high_risk_money":             "人民币",
	"high_risk_terrorist_uniform": "极端服饰",
	"high_risk_sensitive_map":     "敏感地图",
	"great_hall":                  "大会堂",
	"cartoon_porn":                "色情动漫",
	"party_founding_memorial":     "建党纪念",
}

// ImgVerify image 支持图片链接或者图片 base64
func (c *Client) ImgVerify(image string) *Response {
	if len(image) < 1 {
		return Fail(400, "image is empty")
	}

	var format string
	if strings.HasPrefix(image, "http") {
		format = `{"app_id":"%s", "access_token":"%s", "image":"%s"}`
	} else {
		format = `{"app_id":"%s", "access_token":"%s", "image_data":"%s"}`
	}

	params := bytes.NewBufferString(fmt.Sprintf(
		format,
		c.account.AppId,
		c.fetchAccessToken(),
		image,
	))

	return c.imageVerify(params)
}

func (c *Client) imageVerify(body io.Reader) *Response {
	response, err := c.post(fmt.Sprintf("%s%s", c.fetchUir(), imageUrl), body)
	if err != nil {
		return Fail(500, err.Error())
	}

	if code := gjson.Get(response, "error").Int(); code != 0 {
		return Fail(int(code), gjson.Get(response, "message").String())
	}

	var code bool
	var msg string
	gjson.Get(response, "predicts").ForEach(func(key, value gjson.Result) bool {
		if hit := value.Get("hit").Bool(); hit == true {
			if val, ok := modelName[value.Get("model_name").String()]; ok {
				msg = val
			} else {
				msg = "图片违规"
			}
			code = hit
			return true
		}

		return false
	})

	if code == true {
		return Fail(500, fmt.Sprintf("图片违规: 涉及%s", msg))
	}

	return Success()
}

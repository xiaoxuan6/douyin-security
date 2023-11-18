package douyin_security

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

func (c *Client) SetAccessToken(accessToken string) {
	c.accessToken = accessToken
}

func (c *Client) GetAccessToken() (string, error) {
	uri := fmt.Sprintf("%s%s", c.fetchUir(), tokenUrl)
	params := bytes.NewBufferString(fmt.Sprintf(
		`{"appid": "%s", "secret": "%s", "grant_type":"client_credential"}`,
		c.account.AppId,
		c.account.AppSecret,
	))
	response, err := c.post(uri, params)

	if err != nil {
		return "", err
	}

	if errNo := gjson.Get(response, "err_no").Int(); errNo != 0 {
		return "", errors.New(gjson.Get(response, "err_tips").String())
	}

	return gjson.Get(response, "data.access_token").String(), nil
}

func (c *Client) fetchAccessToken() string {
	if len(c.accessToken) < 1 {
		accessToken, err := c.GetAccessToken()
		if err != nil {
			return ""
		}

		return accessToken
	}

	return c.accessToken
}

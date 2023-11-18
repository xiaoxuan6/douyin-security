package douyin_security

import (
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"io"
	"io/ioutil"
	"net/http"
)

func (c *Client) getMaxAttempts() uint {
	attempts := uint(3)
	if c.maxAttempts != 0 {
		attempts = c.maxAttempts
	}

	return attempts
}

func (c *Client) post(url string, body io.Reader) (string, error) {
	var content string
	err := retry.Do(
		func() error {
			response, err := http.Post(url, "application/json", body)
			defer response.Body.Close()

			if err != nil {
				return errors.New(fmt.Sprintf("请求错误：%s", err.Error()))
			}

			b, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return errors.New(fmt.Sprintf("获取 response 错误：%s", err.Error()))
			}

			content = string(b)
			return nil
		},
		retry.Attempts(c.getMaxAttempts()),
		retry.LastErrorOnly(true),
	)

	return content, err
}

func (c *Client) postWithToken(url string, body io.Reader) (string, error) {
	var content string
	err := retry.Do(
		func() error {
			req, err := http.NewRequest(http.MethodPost, url, body)
			if err != nil {
				return errors.New(fmt.Sprintf("构建 request 错误：%s", err.Error()))
			}
			req.Header.Set("X-Token", c.fetchAccessToken())
			req.Header.Set("Content-Type", "application/json")

			response, err := http.DefaultClient.Do(req)
			defer response.Body.Close()

			if err != nil {
				return errors.New(fmt.Sprintf("请求错误：%s", err.Error()))
			}

			b, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return errors.New(fmt.Sprintf("获取 response 错误：%s", err.Error()))
			}

			content = string(b)
			return nil
		},
		retry.Attempts(c.getMaxAttempts()),
		retry.LastErrorOnly(true),
	)

	return content, err
}

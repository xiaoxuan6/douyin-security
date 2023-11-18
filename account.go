package douyin_security

type Account struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	IsSandbox bool   `json:"is_sandbox"`
}

func NewAccount(appid, appSecret string, isSandbox bool) *Account {
	return &Account{
		AppId:     appid,
		AppSecret: appSecret,
		IsSandbox: isSandbox,
	}
}

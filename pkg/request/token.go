package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"git.ucloudadmin.com/uk8s/prometheus-webhook-wechat-public/pkg/models"
	"github.com/pkg/errors"
)

// const AccessTokenAPI = "https://api.weixin.qq.com/cgi-bin/stable_token"

type RequestBody struct {
	Appid      string `json:"appid"`
	Secret     string `json:"secret"`
	Grant_type string `json:"grant_type"`
}

// SendGetTokenRequest 获取token
func SendGetTokenRequest(ApiURL string) (*models.GetTokenResponse, error) {
	resp, err := http.Get(ApiURL)
	if err != nil {
		return nil, errors.Wrap(err, "error request Wechat request")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error IO Read")
	}

	var tokenResp models.GetTokenResponse
	err = json.Unmarshal([]byte(string(body)), &tokenResp)
	if err != nil {
		return nil, errors.Wrap(err, "error Unmarshal")
	}

	return &tokenResp, nil
}

// 获取AppID的access_token
func Get(appid string, secret string, url string) (string, error) {
	data := RequestBody{
		Appid:      appid,
		Secret:     secret,
		Grant_type: "client_credential",
	}

	// 将数据转换为JSON格式
	reqData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	AccessTokenAPI := "https://" + url + "/cgi-bin/stable_token"
	at, err := getPage(AccessTokenAPI, reqData)
	if err != nil {
		return "", err
	}

	if len(at.AccessToken) == 0 {
		return "", fmt.Errorf(at.Errmsg)
	}

	return at.AccessToken, nil
}

func getPage(urlPath string, reqData []byte) (models.GetTokenResponse, error) {
	req, err := http.Post(urlPath, "application/json", bytes.NewBuffer(reqData))
	if err != nil {
		return models.GetTokenResponse{}, err
	}
	defer req.Body.Close()

	var responseBody models.GetTokenResponse
	err = json.NewDecoder(req.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("JSON decode error:", err)
		return models.GetTokenResponse{}, err
	}

	return responseBody, nil
}

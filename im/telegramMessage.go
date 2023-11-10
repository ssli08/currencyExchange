package im

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type result struct {
}
type resultInfo struct {
	OK          bool   `json:"ok"`
	ErrCode     int32  `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
	Result      result `json:"result"`
}

// var msgURL string = tgUrl + botToken + "/sendMessage"

func SendMSG(tgUrl, botToken, chatGroupID, msg string) error {
	msgURL := tgUrl + botToken + "/sendMessage"
	payload, err := json.Marshal(map[string]string{"text": msg, "chat_id": chatGroupID})
	if err != nil {
		fmt.Println("marshal ", err)
		return err
	}

	resp, err := http.Post(msgURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("resp ", err)
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read ", err)
		return err
	}

	res := resultInfo{}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	if !res.OK {
		fmt.Println(res.ErrCode, res.Description)
	}

	// fmt.Println(msg)
	return nil
}

// use socket proxy(if set in system environment) to send http request
func SendMSGViaProxy(tgUrl, botToken, chatGroupID, msg string) error {
	msgURL := tgUrl + botToken + "/sendMessage"

	client := http.Client{}
	proxyServer := checkProxy()
	if proxyServer != "" {
		fmt.Printf("use proxy %s to access %s\n", proxyServer, strings.Split(msgURL, "+")[0])
		u, err := url.Parse(proxyServer)
		if err != nil {
			panic(err)
		}
		client = http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(u)},
			Timeout:   10 * time.Second,
		}
	}

	payload, err := json.Marshal(map[string]string{"text": msg, "chat_id": chatGroupID})
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, msgURL, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	res := resultInfo{}
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	if !res.OK {
		fmt.Println(res.ErrCode, res.Description)
	}

	// fmt.Println(msg)
	return nil

}
func checkProxy() string {
	// proxy format set in System Environment
	// HTTP_PROXY="socks5://127.0.0.1:1080"
	// HTTP_PROXY="http://127.0.0.1:1080"
	// HTTPS_PROXY="https://127.0.0.1:1080"
	// var proxyServer string

	if _, ok := os.LookupEnv("HTTP_PROXY"); ok {
		return os.Getenv("HTTP_PROXY")
	} else if _, ok := os.LookupEnv("HTTPS_PROXY"); ok {
		return os.Getenv("HTTPS_PROXY")
	} else {
		return ""
	}

}

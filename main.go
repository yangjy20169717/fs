package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)


const  (
	secret = "TspRQknxRCDkfLj4R33JRe"
	webhook = "https://open.feishu.cn/open-apis/bot/v2/hook/6b953bf0-70bf-44ea-a9bd-889a8b158a77"
)

type SendMessage struct {
	Timestamp	string	`json:"timestamp"`
	Sign	string	`json:"sign"`
	Msg_type	string	`json:"msg_type"`
	Content	 map[string]string	`json:"content"`
}


func GenSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" +secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}


func fs(sendMessage *SendMessage){
	bytesData, _ := json.Marshal(&sendMessage)
	res, err2 := http.Post(webhook,
		"application/json;charset=utf-8", bytes.NewBuffer(bytesData))

	if err2 != nil {
		fmt.Println("Fatal error ", err2.Error())
	}

	defer res.Body.Close()
}

func main(){

	var sendMessage SendMessage
	sign,err := GenSign(secret,time.Now().Unix())
	if err != nil {
		log.Fatal(err)
	}
	sendMessage.Timestamp = string(time.Now().Unix())
	sendMessage.Sign=sign
	sendMessage.Msg_type="text"





	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		title := c.PostForm("title")
		m := map[string]string{"text":title}
		sendMessage.Content=m
		fs(&sendMessage)
	})

	r.Run(":9090") // listen and serve on 0.0.0.0:8080

}





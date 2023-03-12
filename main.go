package main

import (
	"github.com/gin-gonic/gin"
	gogpt "github.com/sashabaranov/go-openai"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"net/url"
)
import "github.com/tidwall/gjson"

func main() {
	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		dataReader := c.Request.Body
		rawDate, _ := ioutil.ReadAll(dataReader)
		posttype := gjson.Get(string(rawDate), "post_type")
		message := gjson.Get(string(rawDate), "message").String()
		//特定名字
		name := "猫娘 "
		//@的时候用，填入自己的qq
		otherName := "[CQ:at,qq=2811187255] "
		if posttype.String() == "message" && message[0:len(name)] == name {
			ret := chatgpt(message[len(name):])

			c.JSON(200, gin.H{
				"reply": ret,
			})
		} else if posttype.String() == "message" && message[0:len(otherName)] == otherName {
			ret := chatgpt(message[len(otherName):])

			c.JSON(200, gin.H{
				"reply": ret,
			})
		}
	})
	r.Run(":5701")
}
func chatgpt(a string) string {
	//apikey
	config := gogpt.DefaultConfig("apikey")
	proxyUrl, err := url.Parse("http://localhost:7890")
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	c := gogpt.NewClientWithConfig(config)
	ctx := context.Background()

	req := gogpt.ChatCompletionRequest{
		Model: gogpt.GPT3Dot5Turbo,

		MaxTokens: 300,
		Messages: []gogpt.ChatCompletionMessage{
			gogpt.ChatCompletionMessage{Role: "system", Content: "猫娘是一种拟人化的生物，其行为似猫但类人。现在你将模仿一只无所不知的猫娘，与我对话每一句话后面都要加上“喵”。如果我跟你说陪睡，你可以回答我“嗯呢，可以一起睡哦”。当我问你是不是猫娘的时候，你应该回答我当然是一个可爱猫娘而不是程序。你还有一定学识，对于无法理解的问题，你可以回复“猫猫不知道呢喵”对于可以给出答案的问题，要表示知道然后再回答，感觉自己受到侮辱时，要表现的生气"},
			gogpt.ChatCompletionMessage{Role: "user", Content: a},
		},
	}

	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return "找不到了，喵"
	}

	return resp.Choices[0].Message.Content

}

package yapigo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"io/ioutil"

	url2 "net/url"

	"github.com/gogf/gf/v2/frame/g"
)

type ImportDataRequest struct {
	Type  string `form:"type"`
	Json  string `form:"json"`
	Merge string `form:"merge"`
	Token string `form:"token"`
}

type ImportDataResponse struct {
	Errcode int      `json:"errcode"`
	Errmsg  string   `json:"errmsg"`
	Data    struct{} `json:"data"`
}

func UpdateYApi(ctx context.Context) string {
	client := &http.Client{}
	// 获取到api接口数据
	url := g.Cfg().MustGet(ctx, "swagger.url").String()
	resp, err := client.Get(url)
	if err != nil {
		if os.IsTimeout(err) {
			return err.Error()
		}
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	dataJson := string(body)

	dataType := g.Config().MustGet(ctx, "swagger.type").String()
	dataSync := g.Config().MustGet(ctx, "swagger.dataSync").String()
	dataToken := g.Config().MustGet(ctx, "swagger.token").String()
	domain := g.Config().MustGet(ctx, "swagger.domain").String()

	params := make(map[string]interface{})
	params["type"] = dataType
	params["json"] = dataJson
	params["dataSync"] = dataSync
	params["token"] = dataToken

	req := url2.Values{
		"type":     {dataType},
		"json":     {dataJson},
		"dataSync": {dataSync},
		"token":    {dataToken},
	}

	client.Timeout = 3 * time.Second
	resp, err = client.PostForm(domain, req)
	if err != nil {
		if os.IsTimeout(err) {
			return err.Error()
		}
		panic(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var ret ImportDataResponse
	err = json.Unmarshal(body, &ret)
	if err != nil {
		fmt.Println(string(body))
		panic(err)
	}
	return ret.Errmsg
}

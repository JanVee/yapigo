package yapigo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
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

type ApiJsonResponse struct {
	OpenApi    interface{}            `json:"openapi"`
	Components interface{}            `json:"components"`
	Info       interface{}            `json:"info"`
	Paths      map[string]interface{} `json:"paths"`
}

func UpdateYApi(ctx context.Context) string {
	client := &http.Client{}
	// 获取到api接口数据
	url := g.Cfg().MustGet(ctx, "swagger.url").String()
	dataType := g.Config().MustGet(ctx, "swagger.type").String()
	dataSync := g.Config().MustGet(ctx, "swagger.dataSync").String()
	dataToken := g.Config().MustGet(ctx, "swagger.token").String()
	domain := g.Config().MustGet(ctx, "swagger.domain").String()

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

	var retApi ApiJsonResponse
	err = json.Unmarshal(body, &retApi)
	if err != nil {
		panic(err)
	}

	// 解析body
	if !g.IsEmpty(retApi.Paths) {
		for _, pathRow := range retApi.Paths {
			for method, row := range gconv.Map(pathRow) {
				if method == "post" {
					row1 := gconv.Map(row)
					row2 := make([]g.Map, 0)
					row2 = append(row2, g.Map{
						"name":        "Authorization",
						"in":          "header",
						"description": "",
						"required":    true,
						"type":        "string",
					})
					row1["parameters"] = row2
					row = row1
				} else if method == "get" {
					row1 := gconv.Map(row)
					if !g.IsEmpty(row1["parameters"]) {
						row2 := gconv.SliceAny(row1["parameters"])
						// 强行写入类型
						for _, v := range row2 {
							m := gconv.Map(v)
							schema := gconv.Map(m["schema"])
							m["description"] = "{" + gconv.String(schema["format"]) + "}" + gconv.String(m["description"])
							v = m
						}

						row2 = append(row2, g.Map{
							"name":        "Authorization",
							"in":          "header",
							"description": "",
							"required":    true,
							"type":        "string",
						})
						row1["parameters"] = row2
					} else {
						row2 := make([]g.Map, 0)
						row2 = append(row2, g.Map{
							"name":        "Authorization",
							"in":          "header",
							"description": "",
							"required":    true,
							"type":        "string",
						})
						row1["parameters"] = row2
					}
					row = row1
				}
			}
		}
	}

	dataJson, _ := json.Marshal(retApi)
	dataJsonStr := string(dataJson)
	req := url2.Values{
		"type":  {dataType},
		"json":  {dataJsonStr},
		"merge": {dataSync},
		"token": {dataToken},
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

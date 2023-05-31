package yapigo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/JanVee/yapigo/model"
	"github.com/gogf/gf/v2/util/gconv"
	"net/http"
	"os"
	"time"

	"io/ioutil"

	url2 "net/url"

	"github.com/gogf/gf/v2/frame/g"
)

func AddOrUpdateToYApi(ctx context.Context) string {
	client := &http.Client{}
	// 获取到api接口数据
	host := g.Cfg().MustGet(ctx, "swagger.host").String()
	yApiHost := g.Cfg().MustGet(ctx, "swagger.yApiHost").String()
	dataToken := g.Config().MustGet(ctx, "swagger.token").String()

	resp, err := client.Get(host + "/api.json")
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

	var retApi model.ApiJsonResponse
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
		"type":  {"swagger"},
		"json":  {dataJsonStr},
		"merge": {"mergin"},
		"token": {dataToken},
	}

	client.Timeout = 3 * time.Second
	resp, err = client.PostForm(yApiHost+"/api/open/import_data", req)
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

	var ret model.ImportDataResponse
	err = json.Unmarshal(body, &ret)
	if err != nil {
		fmt.Println(string(body))
		panic(err)
	}

	return ret.Errmsg
}

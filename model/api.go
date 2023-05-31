package model

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

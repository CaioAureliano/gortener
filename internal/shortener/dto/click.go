package dto

import (
	"net/http"
	"strings"

	"github.com/CaioAureliano/gortener/internal/shortener/model"
)

type ClickedRequest struct {
	Source   string
	Device   string
	Browser  string
	Language string
	System   string
}

const (
	acceptLanguageHeaderKey = "Accept-Language"
	userAgentHeaderKey      = "User-Agent"
)

var directives = map[string][]struct {
	field string
	key   string
}{
	"browser": {
		{field: "opera", key: "opr"},
		{field: "chrome", key: "chrome"},
		{field: "firefox", key: "firefox"},
		{field: "safari", key: "safari"},
		{field: "internet explorer", key: "msie"},
	},
	"system": {
		{field: "linux", key: "linux"},
		{field: "mac", key: "mac"},
		{field: "windows", key: "windows"},
	},
	"device": {
		{field: "mobile", key: "mobile"},
	},
}

func (c ClickedRequest) Set(req *http.Request) (click model.Click) {
	headers := req.Header

	click.Source = req.Referer()

	if acceptLanguage, ok := headers[acceptLanguageHeaderKey]; ok {
		languages := strings.Split(acceptLanguage[0], ",")
		click.Language = languages[0]
	}

	if userAgents, ok := headers[userAgentHeaderKey]; ok {
		h := userAgents[0]

		for k, d := range directives {
			value := ""

			for _, v := range d {
				if strings.Contains(strings.ToLower(h), v.key) {
					value = v.field
					break
				}
			}

			if value == "" && k != "device" {
				value = "other"
			}

			switch k {
			case "browser":
				click.Browser = value
			case "system":
				click.System = value
			case "device":
				if value == "" {
					value = "desktop"
				}
				click.Device = value
			}
		}
	}

	return
}

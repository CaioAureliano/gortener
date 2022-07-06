package dto

import (
	"regexp"
	"strings"
)

type UrlRequest struct {
	Url string `json:"url" validate:"required"`
}

const (
	regexValidURL = `[(http(s)?):\/\/(www\.)?a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
)

func (u *UrlRequest) IsValid() bool {
	isValid, _ := regexp.MatchString(regexValidURL, u.Url)
	return isValid
}

func (u *UrlRequest) AppendProtocolIfNotExists() {
	if !strings.Contains(u.Url, "http") {
		u.Url = "http://" + u.Url
	}
}

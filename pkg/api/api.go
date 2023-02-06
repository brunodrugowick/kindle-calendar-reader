package api

import (
	"net/http"
)

type Api interface {
	HandleRequests(w http.ResponseWriter, r *http.Request)
	GetPath() string
}

func ParseFormAndGetFromRequest(r *http.Request, queryOrAttr string) (got string, err error) {
	err = r.ParseForm()
	if err != nil {
		return got, err
	}
	got = r.Form.Get(queryOrAttr)
	return got, nil
}

package api

import (
	"net/http"
)

type Api interface {
	HandleRequests(w http.ResponseWriter, r *http.Request)
	GetPath() string
}

func ParseFormAndGetFromRequest(r *http.Request, queryOrAttrs ...string) (map[string]string, error) {
	got := make(map[string]string, len(queryOrAttrs))

	err := r.ParseForm()
	if err != nil {
		return got, err
	}

	for _, queryOrAttr := range queryOrAttrs {
		got[queryOrAttr] = r.Form.Get(queryOrAttr)
	}

	return got, nil
}

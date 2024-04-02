package request_util

import (
	"encoding/json"
	"go_todo/src/model"
	"go_todo/src/types"
	"net/http"
)

func RequestUser(r *http.Request) *model.User {
	return r.Context().Value(types.UserKey{}).(*model.User)
}

type ReturnJsonOptions struct {
	Content interface{}
	Status int
}

func ReturnJson(w http.ResponseWriter, options ReturnJsonOptions) {
	if (options.Status == 0) {
		options.Status = http.StatusOK
	}

	if (options.Content == nil) {
		options.Content = map[string]interface{}{ "status": 1 }
	}

	w.WriteHeader(options.Status)

	encorder := json.NewEncoder(w)
	encorder.SetIndent("", "  ")
	encorder.Encode(options.Content)
}

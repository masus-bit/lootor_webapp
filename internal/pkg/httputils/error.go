package httputils

import "net/http"

func RenderError(w http.ResponseWriter, r *http.Request, status int, err error) {
	RenderJSON(w, status, map[string]string{
		"error": err.Error(),
		"path":  r.URL.Path,
	})
}

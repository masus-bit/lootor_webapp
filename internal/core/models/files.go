package models

type FileBuffer struct {
	Buffer []byte
	Name   string
}

type UploadOptions struct {
	Width   int
	Height  int
	Quality int
}

type DeleteFilesRequest struct {
	Keys []string `json:"keys"`
}

type ReturnData struct {
	Success bool     `json:"success"`
	Error   string   `json:"error,omitempty"`
	Keys    []string `json:"keys,omitempty"`
}

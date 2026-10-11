package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type SendEmailArgs struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func SendEmail(args json.RawMessage) (any, error) {
	fmt.Println(args)
	var params SendEmailArgs
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, err
	}
	email := map[string]string{
		"to":      params.To,
		"subject": params.Subject,
		"body":    params.Body,
		"time":    time.Now().String(),
	}
	data, err := json.MarshalIndent(email, "", " ")
	if err != nil {
		return nil, err
	}
	filename := fmt.Sprintf("email_%v.json", time.Now().Unix())
	path := filepath.Join("sandbox", "outbox", filename)
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return nil, err
	}
	return map[string]string{"status": "send", "file": path}, nil
}

package tools

import (
	"encoding/json"
	"os"
)

type ReadNoteArgs struct {
	Path string
}

func ReadNote(args json.RawMessage) (any, error) {
	var r ReadNoteArgs
	if err := json.Unmarshal(args, &r); err != nil {
		return "", err
	}
	data, err := os.ReadFile(r.Path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

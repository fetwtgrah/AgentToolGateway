package tools

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type DeleteNoteArgs struct {
	Path string `json:"path"`
}

func DeleteNote(args json.RawMessage) (any, error) {
	var params DeleteNoteArgs
	if err := json.Unmarshal(args, &params); err != nil {
		return nil, err
	}
	base := filepath.Join("sandbox", "notes")
	target := filepath.Join(base, params.Path)

	absBase, err := filepath.Abs(base)
	if err != nil {
		return nil, err
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return nil, err
	}

	rel, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return nil, err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return nil, errors.New("permission denied")
	}
	err = os.Remove(absTarget)
	if err != nil {
		return nil, err
	}
	return map[string]string{"status": "delete", "path": params.Path}, nil
}

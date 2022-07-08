package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_extractData(t *testing.T) {
	tests := []struct {
		name string
		link string
		org  string
		repo string
	}{
		{"normal link to repo", "https://github.com/aerfio/gitclone", "aerfio", "gitclone"},
		{"link to some file", "https://github.com/gpakosz/.tmux/blob/master/.tmux.conf.local", "gpakosz", ".tmux"},
		{"git clone link via ssh", "git@github.com:aerfio/gitclone.git", "aerfio", "gitclone"},
		{"git clone link via https", "https://github.com/aerfio/gitclone.git", "aerfio", "gitclone"},
		{"link to raw file", "https://raw.githubusercontent.com/aerfio/gitclone/main/go.sum", "aerfio", "gitclone"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var org string
			var repo string

			assert.NotPanics(t, func() {
				org, repo = extractData(tt.link)
			})

			assert.Equal(t, tt.org, org)
			assert.Equal(t, tt.repo, repo)
		})
	}
}

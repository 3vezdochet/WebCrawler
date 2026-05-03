package parser

import (
	"io"
	"slices"
	"strings"
	"testing"
)

func TestExtractLinks(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		htmlBody string
		want     []string
	}{
		{
			name:     "Стандартные ссылки",
			baseURL:  "https://golang.org",
			htmlBody: `<a href="/doc">Doc</a><a href="https://google.com">Google</a>`,
			want:     []string{"https://golang.org/doc", "https://google.com"},
		},
		{
			name:     "Отсутствие ссылок",
			baseURL:  "https://golang.org",
			htmlBody: `<html><body><h1>No links here</h1></body></html>`,
			want:     []string{}, // Ждем пустой срез, а не nil
		},
		{
			name:     "Пустой атрибут href",
			baseURL:  "https://golang.org",
			htmlBody: `<a href="">Empty</a><a href=" ">Space</a>`,
			want:     []string{},
		},
		{
			name:     "Игнорирование якорей и фрагментов",
			baseURL:  "https://golang.org",
			htmlBody: `<a href="#top">Top</a><a href="/page#section">Section</a>`,
			// Канонично: либо игнорим, либо чистим от фрагментов
			want: []string{"https://golang.org/page"},
		},
		{
			name:     "Специфические протоколы (mailto, javascript)",
			baseURL:  "https://golang.org",
			htmlBody: `<a href="mailto:admin@golang.org">Mail</a><a href="javascript:void(0)">JS</a>`,
			want:     []string{}, // Краулер не должен пытаться "переходить" по почте
		},
		{
			name:     "Битая HTML верстка",
			baseURL:  "https://golang.org",
			htmlBody: `<a href="/valid">Valid <a href="/nested">Nested`,
			want:     []string{"https://golang.org/valid", "https://golang.org/nested"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := io.NopCloser(strings.NewReader(tt.htmlBody))

			got := ExtractLinks(body, tt.baseURL)

			if !slices.Equal(got, tt.want) {
				t.Errorf("ExtractLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

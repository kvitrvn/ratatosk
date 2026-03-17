package domain

import "testing"

func TestNewFeed(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid https", "https://news.ycombinator.com/rss", false},
		{"valid http", "http://feeds.example.com/rss", false},
		{"no scheme", "news.ycombinator.com/rss", true},
		{"ftp scheme", "ftp://feeds.example.com/rss", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feed, err := NewFeed(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewFeed(%q) expected error, got nil", tt.url)
				}
				return
			}
			if err != nil {
				t.Errorf("NewFeed(%q) unexpected error: %v", tt.url, err)
			}
			if feed.URL != tt.url {
				t.Errorf("feed.URL = %q, want %q", feed.URL, tt.url)
			}
			if feed.CreatedAt.IsZero() {
				t.Error("feed.CreatedAt should not be zero")
			}
		})
	}
}

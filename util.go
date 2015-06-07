package main

import "net/url"

func mustParseURL(urlString string) string {
	u, _ := url.Parse(urlString)
	return u.String()
}

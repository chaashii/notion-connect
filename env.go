package main

import "os"

// NOTION_API_KEY := os.Getenv("NOTION_API_KEY")

func GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return ""
}

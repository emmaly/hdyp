package main

import (
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Setting is a setting is a setting
type Setting struct {
	String string
	Int    int
	Bool   bool
}

func settingKey(c context.Context, key string) *datastore.Key {
	return datastore.NewKey(c, "Setting", strings.ToLower(key), 0, nil)
}

// GetSettingString gets string from setting
func GetSettingString(c context.Context, key string) string {
	var dst Setting
	datastore.Get(c, settingKey(c, key), &dst)
	return dst.String
}

// SetSettingString sets a string setting
func SetSettingString(c context.Context, key string, value string) (*datastore.Key, error) {
	return datastore.Put(c, settingKey(c, key), &Setting{String: value})
}

// GetSettingInt gets int from setting
func GetSettingInt(c context.Context, key string) int {
	var dst Setting
	datastore.Get(c, settingKey(c, key), &dst)
	return dst.Int
}

// SetSettingInt sets an int setting
func SetSettingInt(c context.Context, key string, value int) (*datastore.Key, error) {
	return datastore.Put(c, settingKey(c, key), &Setting{Int: value})
}

// GetSettingBool gets int from setting
func GetSettingBool(c context.Context, key string) bool {
	var dst Setting
	datastore.Get(c, settingKey(c, key), &dst)
	return dst.Bool
}

// SetSettingBool sets an int setting
func SetSettingBool(c context.Context, key string, value bool) (*datastore.Key, error) {
	return datastore.Put(c, settingKey(c, key), &Setting{Bool: value})
}

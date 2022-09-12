package storage

import "fmt"

var DB = make(map[string]string)

func PutMapStorage(key string, value string) (err error) {
	if _, ok := DB[key]; ok {
		return fmt.Errorf("key is already in database")
	}
	DB[key] = string(value)
	return nil
}

func GetMapStorage(key string) (value string, err error) {
	var ok bool
	if value, ok = DB[key]; !ok {
		return "", fmt.Errorf("key %v not found", key)
	}
	return value, nil

}

package authsample

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// Users holds a list of users.
type Users map[string]struct{}

// Check checks if a key could retrieve a user from a list of users.
func (u Users) Check(key string) (bool, string, string) {
	splits := strings.Split(key, ".")
	if len(splits) != 2 {
		return false, "", ""
	}

	userName := splits[0]
	_, ok := u[userName]
	if !ok {
		return false, "", ""
	}
	return ok, userName, splits[1]
}

// LoadUsers load users data from a JSON file.
func LoadUsers(jsonFile string) (Users, error) {
	var users Users
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	return users, nil
}

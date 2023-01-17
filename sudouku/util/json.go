package util

import "encoding/json"

func MarshalIndent(in interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(in, "", "    ")
	if err != nil {
		return "", err
	}
	json := string(jsonBytes)
	return json, nil
}

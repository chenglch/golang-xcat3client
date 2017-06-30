package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadJsonFile(filepath string) (data map[string]interface{}, err error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func WriteJsonFile(filepath string, data []byte) (err error) {
	var out bytes.Buffer
	json.Indent(&out, data, "", "\t")
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	out.WriteTo(f)
	return
}

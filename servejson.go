// Copyright 2015,2016,2017,2018 SeukWon Kang (kasworld@gmail.com)

package weblib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ServeJSON2HTTP(obj interface{}, w http.ResponseWriter) error {
	data, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "fail to marshal json", 404)
		return fmt.Errorf("fail to marshal json %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "fail to write", 404)
		return fmt.Errorf("fail to write %v", err)
	}
	return nil
}

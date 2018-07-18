// text_consumer.go
//
// Copyright (c) 2018 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package file

import (
	"encoding/json"
	"fmt"
)

// TextUnmarshaler interface implementation
func (m *UploadFileCreatedBody) UnmarshalText(text []byte) error {
	return m.UnmarshalJSON(text)
}

func (m *UploadFileCreatedBody) UnmarshalJSON(text []byte) (err error) {

	aux := make(map[string]interface{})
	err = json.Unmarshal(text, &aux)
	if err != nil {
		fmt.Println(string(text))
		return
	}

	m.Success = aux["success"].(bool)
	m.ID = aux["id"].(string)
	return

}

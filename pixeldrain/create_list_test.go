// create_list_test.go
//
// Copyright (c) 2018 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php

package pixeldrain

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/jkawamoto/go-pixeldrain/client"
	"github.com/jkawamoto/go-pixeldrain/client/list"
	"github.com/jkawamoto/go-pixeldrain/models"
)

type mockListServer struct {
	ID          string
	Description string
	Title       string
	Files       []*list.FilesItems0
}

func (m *mockListServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	res.Header().Add("Content-type", "application/json")
	if req.URL.Path != "/list" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(models.StandardError{Message: swag.String("received a wrong request")})
		return
	}

	raw, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(models.StandardError{
			Message: swag.String(fmt.Sprintln("failed to parse the request:", err)),
		})
		return
	}

	var body list.CreateFileListBody
	err = body.UnmarshalBinary(raw)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(models.StandardError{
			Message: swag.String(fmt.Sprintln("failed to parse the request:", err)),
		})
		return
	}

	m.Description = body.Description
	m.Title = *body.Title
	m.Files = body.Files

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(&list.CreateFileListCreatedBody{
		ID:      m.ID,
		Success: true,
	})
	return

}

func TestCreateList(t *testing.T) {

	mock := &mockListServer{
		ID: "sample-id",
	}
	server := httptest.NewServer(mock)
	defer server.Close()

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal("Cannot parse a URL:", err)
	}
	cli := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:     u.Host,
		BasePath: "/",
		Schemes:  []string{"http"},
	})

	title := "sample-title"
	description := "sample-description"
	files := []string{"file1:file1.txt", "file2:file2.dat"}
	id, err := CreateList(context.Background(), cli, title, description, files)
	if err != nil {
		t.Fatal("CreateList returned an error:", err)
	}

	if id != mock.ID {
		t.Errorf("received id was %v but expected %v", id, mock.ID)
	}
	if len(mock.Files) != len(files) {
		t.Errorf("the mock server received %v items but sent %v items", len(mock.Files), len(files))
	} else {
		for i, f := range ParseListItems(files) {
			if item := mock.Files[i]; item.ID != f.ID || item.Description != f.Description {
				t.Errorf("item %v was %+v but expected %+v", i, item, f)
			}
		}
	}

}

func TestParseListItems(t *testing.T) {

	cases := []struct {
		Input    []string
		Expected []*list.FilesItems0
	}{
		{
			Input: []string{"ID"},
			Expected: []*list.FilesItems0{
				{ID: "ID"},
			},
		},
		{
			Input: []string{"ID:description"},
			Expected: []*list.FilesItems0{
				{ID: "ID", Description: "description"},
			},
		},
		{
			Input: []string{"ID:description", "id2"},
			Expected: []*list.FilesItems0{
				{ID: "ID", Description: "description"},
				{ID: "id2"},
			},
		},
		{
			Input: []string{"ID:description", "id2:desc2"},
			Expected: []*list.FilesItems0{
				{ID: "ID", Description: "description"},
				{ID: "id2", Description: "desc2"},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case-%v", i), func(t *testing.T) {
			res := ParseListItems(c.Input)
			if len(res) != len(c.Expected) {
				t.Errorf("got %v items but want %v", len(res), len(c.Expected))
			} else {
				for j, e := range c.Expected {
					if res[j].ID != e.ID || res[j].Description != e.Description {
						t.Errorf("item %v: %+v but want ID = %v, Description = %q", j, res[j], e.ID, e.Description)
					}
				}
			}
		})
	}

}

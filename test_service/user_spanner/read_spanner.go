package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

type SpannerParams struct {
	ProjectId  string `json:"projectId,omitempty"`
	InstanceId string `json:"instanceId,omitempty"`
	DatabaseId string `json:"databaseId,omitempty"`
}

func (s SpannerParams) URI() string {
	return fmt.Sprintf("%s/databases/%s", s.Parent(), s.DatabaseId)
}

func (s SpannerParams) Parent() string {
	return fmt.Sprintf("projects/%s/instances/%s", s.ProjectId, s.InstanceId)
}

// need to have a struct that looks like this saved in your
// ~/.protoc-gen-persist-db.json
// {
// 		"projectId": "my-google-cloud-project-id",
// 		"instanceId": "my-google-cloud-instance-id",
// 		"databaseId": "my-google-cloud-database-id"
// }
func ReadSpannerParams() (out SpannerParams) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	f, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, "/.protoc-gen-persist-db.json"))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(f, &out); err != nil {
		panic(err)
	}
	return
}

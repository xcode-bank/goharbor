// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRepos(t *testing.T) {

	assert := assert.New(t)
	apiTest := newHarborAPI()
	projectID := "1"
	keyword := "hello-world"

	fmt.Println("Testing Repos Get API")
	//-------------------case 1 : response code = 200------------------------//
	fmt.Println("case 1 : response code = 200")
	code, repositories, err := apiTest.GetRepos(*admin, projectID, keyword)
	if err != nil {
		t.Errorf("failed to get repositories: %v", err)
	} else {
		assert.Equal(int(200), code, "response code should be 200")
		if repos, ok := repositories.([]repoResp); ok {
			assert.Equal(int(1), len(repos), "the length of repositories should be 1")
			assert.Equal(repos[0].Name, "library/hello-world", "unexpected repository name")
		} else {
			t.Error("unexpected reponse")
		}
	}

	//-------------------case 2 : response code = 404------------------------//
	fmt.Println("case 2 : response code = 404:project  not found")
	projectID = "111"
	httpStatusCode, _, err := apiTest.GetRepos(*admin, projectID, keyword)
	if err != nil {
		t.Error("Error whihle get repos by projectID", err.Error())
		t.Log(err)
	} else {
		assert.Equal(int(404), httpStatusCode, "httpStatusCode should be 404")
	}

	//-------------------case 3 : response code = 400------------------------//
	fmt.Println("case 3 : response code = 400,invalid project_id")
	projectID = "ccc"
	httpStatusCode, _, err = apiTest.GetRepos(*admin, projectID, keyword)
	if err != nil {
		t.Error("Error whihle get repos by projectID", err.Error())
		t.Log(err)
	} else {
		assert.Equal(int(400), httpStatusCode, "httpStatusCode should be 400")
	}

	fmt.Printf("\n")
}

func TestGetReposTags(t *testing.T) {

	assert := assert.New(t)
	apiTest := newHarborAPI()

	//-------------------case 1 : response code = 404------------------------//
	fmt.Println("case 1 : response code = 404,repo not found")
	repository := "errorRepos"
	code, _, err := apiTest.GetReposTags(*admin, repository)
	if err != nil {
		t.Errorf("failed to get tags of repository %s: %v", repository, err)
	} else {
		assert.Equal(int(404), code, "httpStatusCode should be 404")
	}
	//-------------------case 2 : response code = 200------------------------//
	fmt.Println("case 2 : response code = 200")
	repository = "library/hello-world"
	code, tags, err := apiTest.GetReposTags(*admin, repository)
	if err != nil {
		t.Errorf("failed to get tags of repository %s: %v", repository, err)
	} else {
		assert.Equal(int(200), code, "httpStatusCode should be 200")
		if tg, ok := tags.([]tagResp); ok {
			assert.Equal(1, len(tg), fmt.Sprintf("there should be only one tag, but now %v", tg))
			assert.Equal(tg[0].Tag, "latest", "the tag should be latest")
		} else {
			t.Error("unexpected response")
		}
	}

	fmt.Printf("\n")
}

func TestGetReposManifests(t *testing.T) {
	var httpStatusCode int
	var err error
	var repoName string
	var tag string

	assert := assert.New(t)
	apiTest := newHarborAPI()

	fmt.Println("Testing ReposManifests Get API")
	//-------------------case 1 : response code = 200------------------------//
	fmt.Println("case 1 : response code = 200")
	repoName = "library/hello-world"
	tag = "latest"
	httpStatusCode, err = apiTest.GetReposManifests(*admin, repoName, tag)
	if err != nil {
		t.Error("Error whihle get reposManifests by repoName and tag", err.Error())
		t.Log(err)
	} else {
		assert.Equal(int(200), httpStatusCode, "httpStatusCode should be 200")
	}
	//-------------------case 2 : response code = 404------------------------//
	fmt.Println("case 2 : response code = 404:tags error,manifest unknown")
	tag = "l"
	httpStatusCode, err = apiTest.GetReposManifests(*admin, repoName, tag)
	if err != nil {
		t.Error("Error whihle get reposManifests by repoName and tag", err.Error())
		t.Log(err)
	} else {
		assert.Equal(int(404), httpStatusCode, "httpStatusCode should be 404")
	}

	//-------------------case 3 : response code = 404------------------------//
	fmt.Println("case 3 : response code = 404,repo not found")
	repoName = "111"
	httpStatusCode, err = apiTest.GetReposManifests(*admin, repoName, tag)
	if err != nil {
		t.Error("Error whihle get reposManifests by repoName and tag", err.Error())
		t.Log(err)
	} else {
		assert.Equal(int(404), httpStatusCode, "httpStatusCode should be 404")
	}

	fmt.Printf("\n")
}

func TestGetReposTop(t *testing.T) {

	assert := assert.New(t)
	apiTest := newHarborAPI()

	fmt.Println("Testing ReposTop Get API")
	//-------------------case 1 : response code = 400------------------------//
	fmt.Println("case 1 : response code = 400,invalid count")
	count := "cc"
	code, _, err := apiTest.GetReposTop(*admin, count)
	if err != nil {
		t.Errorf("failed to get the most popular repositories: %v", err)
	} else {
		assert.Equal(int(400), code, "response code should be 400")
	}

	//-------------------case 2 : response code = 200------------------------//
	fmt.Println("case 2 : response code = 200")
	count = "1"
	code, repos, err := apiTest.GetReposTop(*admin, count)
	if err != nil {
		t.Errorf("failed to get the most popular repositories: %v", err)
	} else {
		assert.Equal(int(200), code, "response code should be 200")
		if r, ok := repos.([]*repoResp); ok {
			assert.Equal(int(1), len(r), "the length should be 1")
			assert.Equal(r[0].Name, "library/docker", "the name of repository should be library/docker")
		} else {
			t.Error("unexpected response")
		}
	}

	fmt.Printf("\n")
}

package cr

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// StartRepoBuild invokes the cr.StartRepoBuild API synchronously
// api document: https://help.aliyun.com/api/cr/startrepobuild.html
func (client *Client) StartRepoBuild(request *StartRepoBuildRequest) (response *StartRepoBuildResponse, err error) {
	response = CreateStartRepoBuildResponse()
	err = client.DoAction(request, response)
	return
}

// StartRepoBuildWithChan invokes the cr.StartRepoBuild API asynchronously
// api document: https://help.aliyun.com/api/cr/startrepobuild.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) StartRepoBuildWithChan(request *StartRepoBuildRequest) (<-chan *StartRepoBuildResponse, <-chan error) {
	responseChan := make(chan *StartRepoBuildResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.StartRepoBuild(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// StartRepoBuildWithCallback invokes the cr.StartRepoBuild API asynchronously
// api document: https://help.aliyun.com/api/cr/startrepobuild.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) StartRepoBuildWithCallback(request *StartRepoBuildRequest, callback func(response *StartRepoBuildResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *StartRepoBuildResponse
		var err error
		defer close(result)
		response, err = client.StartRepoBuild(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// StartRepoBuildRequest is the request struct for api StartRepoBuild
type StartRepoBuildRequest struct {
	*requests.RoaRequest
	RepoNamespace string `position:"Path" name:"RepoNamespace"`
	RepoName      string `position:"Path" name:"RepoName"`
}

// StartRepoBuildResponse is the response struct for api StartRepoBuild
type StartRepoBuildResponse struct {
	*responses.BaseResponse
}

// CreateStartRepoBuildRequest creates a request to invoke StartRepoBuild API
func CreateStartRepoBuildRequest() (request *StartRepoBuildRequest) {
	request = &StartRepoBuildRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("cr", "2016-06-07", "StartRepoBuild", "/repos/[RepoNamespace]/[RepoName]/build", "cr", "openAPI")
	request.Method = requests.PUT
	return
}

// CreateStartRepoBuildResponse creates a response to parse from StartRepoBuild response
func CreateStartRepoBuildResponse() (response *StartRepoBuildResponse) {
	response = &StartRepoBuildResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}

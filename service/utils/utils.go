/*
   Copyright (c) 2016 VMware, Inc. All Rights Reserved.
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package utils

import (
	"github.com/vmware/harbor/utils/log"
	"net/http"
	"os"
)

func VerifySecret(r *http.Request) bool {
	secret := os.Getenv("UI_SECRET")
	c, err := r.Cookie("uisecret")
	if err != nil {
		log.Errorf("Failed to get secret cookie, error: %v", err)
	}
	return c != nil && c.Value == secret
}

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

*** Settings ***
Documentation  Sign in and out
Resource  ../../resources/Util.robot
Default Tags  regression

*** Test Cases ***
Test Case - Sign in and out
    Init Chrome Driver
    Sign In Harbor  %{HARBOR_ADMIN}  %{HARBOR_PASSWORD}
    Logout Harbor
    Close Browser

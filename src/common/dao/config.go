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

package dao

import (
	"github.com/vmware/harbor/src/common/models"
)

// AuthModeCanBeModified determines whether auth mode can be
// modified or not. Auth mode can modified when there is only admin
// user in database.
func AuthModeCanBeModified() (bool, error) {
	c, err := GetOrmer().QueryTable(&models.User{}).Count()
	if err != nil {
		return false, err
	}
	// admin and anonymous
	return c == 2, nil
}

// GetConfigEntries Get configuration from database
func GetConfigEntries() ([]*models.ConfigEntry, error) {
	o := GetOrmer()
	var p []*models.ConfigEntry
	sql:="select * from properties"
	n,err := o.Raw(sql,[]interface{}{}).QueryRows(&p)

	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, nil
	}
	return p,nil
}

// SaveConfigEntries Save configuration to database.
func SaveConfigEntries(entries []models.ConfigEntry) error{
	o := GetOrmer()
	tempEntry:=models.ConfigEntry{}
	for _, entry := range entries{
		tempEntry.Key = entry.Key
		tempEntry.Value = entry.Value
		created, _, error := o.ReadOrCreate(&tempEntry,"k")
		if error != nil {
			return error
		}
		if !created {
			entry.ID = tempEntry.ID
			_ ,err := o.Update(&entry,"v")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
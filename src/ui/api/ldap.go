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
	"net/http"
	"strings"

	"github.com/vmware/harbor/src/common/models"
	ldapUtils "github.com/vmware/harbor/src/common/utils/ldap"
	"github.com/vmware/harbor/src/common/utils/log"
	"github.com/vmware/harbor/src/ui/auth"
)

// LdapAPI handles requesst to /api/ldap/ping /api/ldap/user/search /api/ldap/user/import
type LdapAPI struct {
	BaseController
}

const metaChars = "&|!=~*<>()"

// Prepare ...
func (l *LdapAPI) Prepare() {
	l.BaseController.Prepare()
	if !l.SecurityCtx.IsAuthenticated() {
		l.HandleUnauthorized()
		return
	}
	if !l.SecurityCtx.IsSysAdmin() {
		l.HandleForbidden(l.SecurityCtx.GetUsername())
		return
	}
}

// Ping ...
func (l *LdapAPI) Ping() {
	var ldapConfs models.LdapConf
	var err error
	var ldapSession *ldapUtils.Session

	l.Ctx.Input.CopyBody(1 << 32)

	if string(l.Ctx.Input.RequestBody) == "" {
		ldapSession, err = ldapUtils.LoadSystemLdapConfig()
		if err != nil {
			log.Errorf("Can't load system configuration, error: %v", err)
			l.RenderError(http.StatusInternalServerError, fmt.Sprintf("can't load system configuration: %v", err))
			return
		}
		err = ldapSession.ConnectionTest()
	} else {
		l.DecodeJSONReqAndValidate(&ldapConfs)
		err = ldapUtils.ConnectionTestWithConfig(ldapConfs)
	}

	if err != nil {
		log.Errorf("ldap connect fail, error: %v", err)
		l.RenderError(http.StatusBadRequest, fmt.Sprintf("ldap connect fail: %v", err))
		return
	}
}

// Search ...
func (l *LdapAPI) Search() {
	var err error
	var ldapUsers []models.LdapUser
	var ldapConfs models.LdapConf
	var ldapSession *ldapUtils.Session
	l.Ctx.Input.CopyBody(1 << 32)
	if string(l.Ctx.Input.RequestBody) == "" {
		ldapSession, err = ldapUtils.LoadSystemLdapConfig()
		if err != nil {
			log.Errorf("can't load system configuration, error: %v", err)
			l.RenderError(http.StatusInternalServerError, fmt.Sprintf("can't load system configuration: %v", err))
			return
		}
	} else {
		l.DecodeJSONReqAndValidate(&ldapConfs)
		ldapSession, err = ldapUtils.CreateWithUIConfig(ldapConfs)
	}

	if err = ldapSession.Open(); err != nil {
		log.Errorf("can't Open ldap session, error: %v", err)
		l.RenderError(http.StatusInternalServerError, fmt.Sprintf("can't open ldap session: %v", err))
		return
	}
	defer ldapSession.Close()

	searchName := l.GetString("username")

	if searchName != "" {
		for _, c := range metaChars {
			if strings.ContainsRune(searchName, c) {
				log.Errorf("the search username contains meta char: %q", c)
				l.RenderError(http.StatusBadRequest, fmt.Sprintf("the search username contains meta char: %q", c))
				return
			}
		}
	}

	ldapUsers, err = ldapSession.SearchUser(searchName)

	if err != nil {
		log.Errorf("Ldap search fail, error: %v", err)
		l.RenderError(http.StatusBadRequest, fmt.Sprintf("ldap search fail: %v", err))
		return
	}

	l.Data["json"] = ldapUsers
	l.ServeJSON()

}

// ImportUser ...
func (l *LdapAPI) ImportUser() {
	var ldapImportUsers models.LdapImportUser
	var ldapFailedImportUsers []models.LdapFailedImportUser
	var ldapConfs models.LdapConf

	l.DecodeJSONReqAndValidate(&ldapImportUsers)

	ldapFailedImportUsers, err := importUsers(ldapConfs, ldapImportUsers.LdapUIDList)

	if err != nil {
		log.Errorf("Ldap import user fail, error: %v", err)
		l.RenderError(http.StatusBadRequest, fmt.Sprintf("ldap import user fail: %v", err))
		return
	}

	if len(ldapFailedImportUsers) > 0 {
		log.Errorf("Import ldap user have internal error")
		l.RenderError(http.StatusInternalServerError, fmt.Sprintf("import ldap user have internal error"))
		l.Data["json"] = ldapFailedImportUsers
		l.ServeJSON()
		return
	}

}

func importUsers(ldapConfs models.LdapConf, ldapImportUsers []string) ([]models.LdapFailedImportUser, error) {
	var failedImportUser []models.LdapFailedImportUser
	var u models.LdapFailedImportUser

	ldapSession, err := ldapUtils.LoadSystemLdapConfig()
	if err != nil {
		log.Errorf("can't load system configuration, error: %v", err)
		return nil, err
	}

	if err = ldapSession.Open(); err != nil {
		log.Errorf("Can't connect to ldap, error: %v", err)
	}
	defer ldapSession.Close()

	for _, tempUID := range ldapImportUsers {
		u.UID = tempUID
		u.Error = ""

		if u.UID == "" {
			u.Error = "empty_uid"
			failedImportUser = append(failedImportUser, u)
			continue
		}

		for _, c := range metaChars {
			if strings.ContainsRune(u.UID, c) {
				u.Error = "invaild_username"
				break
			}
		}

		if u.Error != "" {
			failedImportUser = append(failedImportUser, u)
			continue
		}

		ldapUsers, err := ldapSession.SearchUser(u.UID)
		if err != nil {
			u.UID = tempUID
			u.Error = "failed_search_user"
			failedImportUser = append(failedImportUser, u)
			log.Errorf("Invalid ldap search request for %s, error: %v", tempUID, err)
			continue
		}

		if ldapUsers == nil || len(ldapUsers) <= 0 {
			u.UID = tempUID
			u.Error = "unknown_user"
			failedImportUser = append(failedImportUser, u)
			continue
		}

		var user models.User

		user.Username = ldapUsers[0].Username
		user.Realname = ldapUsers[0].Realname
		user.Email = ldapUsers[0].Email
		err = auth.OnBoardUser(&user)

		if err != nil || user.UserID <= 0 {
			u.UID = tempUID
			u.Error = err.Error()
			failedImportUser = append(failedImportUser, u)
			log.Errorf("Can't import user %s, error: %s", tempUID, u.Error)
		}

	}

	return failedImportUser, nil
}

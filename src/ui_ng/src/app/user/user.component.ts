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
import { Component, OnInit, ViewChild, OnDestroy, ChangeDetectionStrategy, ChangeDetectorRef } from '@angular/core';
import 'rxjs/add/operator/toPromise';
import { Subscription } from 'rxjs/Subscription';

import { UserService } from './user.service';
import { User } from './user';
import { NewUserModalComponent } from './new-user-modal.component';
import { TranslateService } from '@ngx-translate/core';
import { ConfirmationDialogService } from '../shared/confirmation-dialog/confirmation-dialog.service';
import { ConfirmationMessage } from '../shared/confirmation-dialog/confirmation-message';
import { ConfirmationState, ConfirmationTargets, ConfirmationButtons } from '../shared/shared.const';
import { MessageHandlerService } from '../shared/message-handler/message-handler.service';

import { SessionService } from '../shared/session.service';
import { AppConfigService } from '../app-config.service';
import {BatchInfo, BathInfoChanges} from '../shared/confirmation-dialog/confirmation-batch-message';

/**
 * NOTES:
 *   Pagination for this component is a temporary workaround solution. It will be replaced in future release.
 *
 * @export
 * @class UserComponent
 * @implements {OnInit}
 * @implements {OnDestroy}
 */

@Component({
  selector: 'harbor-user',
  templateUrl: 'user.component.html',
  styleUrls: ['user.component.css'],
  providers: [UserService],
  changeDetection: ChangeDetectionStrategy.OnPush
})

export class UserComponent implements OnInit, OnDestroy {
  users: User[] = [];
  originalUsers: Promise<User[]>;
  private onGoing: boolean = true;
  private adminMenuText: string = "";
  private adminColumn: string = "";
  private deletionSubscription: Subscription;
  selectedRow: User[] = [];
  ISADMNISTRATOR: string = "USER.ENABLE_ADMIN_ACTION";
  batchDelectionInfos: BatchInfo[] = [];

  currentTerm: string;
  totalCount: number = 0;
  currentPage: number = 1;

  @ViewChild(NewUserModalComponent)
  newUserDialog: NewUserModalComponent;

  timerHandler: any;

  constructor(
    private userService: UserService,
    private translate: TranslateService,
    private deletionDialogService: ConfirmationDialogService,
    private msgHandler: MessageHandlerService,
    private session: SessionService,
    private appConfigService: AppConfigService,
    private ref: ChangeDetectorRef) {
    this.deletionSubscription = deletionDialogService.confirmationConfirm$.subscribe(confirmed => {
      if (confirmed &&
        confirmed.source === ConfirmationTargets.USER &&
        confirmed.state === ConfirmationState.CONFIRMED) {
        this.delUser(confirmed.data);
      }
    });
  }

  isMySelf(uid: number): boolean {
    let currentUser = this.session.getCurrentUser();
    if (currentUser) {
      if (currentUser.user_id === uid) {
        return true;
      }
    }

    return false;
  }

  get onlySelf(): boolean {
    return this.selectedRow.length === 1 && this.isMySelf(this.selectedRow[0].user_id);
  }

  private isMatchFilterTerm(terms: string, testedItem: string): boolean {
    return testedItem.toLowerCase().indexOf(terms.toLowerCase()) !== -1;
  }

  public get canCreateUser(): boolean {
    let appConfig = this.appConfigService.getConfig();
    if (appConfig) {
      return !(appConfig.auth_mode === 'ldap_auth' || appConfig.auth_mode === 'uaa_auth');
    } else {
      return true;
    }
  }

  public get ifSameRole(): boolean {
    let usersRole: number[] = [];
    this.selectedRow.forEach(user => {
      if (user.user_id === 0 || this.isMySelf(user.user_id)) {
        return false;
      }
      usersRole.push(user.has_admin_role);
    })
    if (usersRole.length && usersRole.every(num => num === 0)) {
      this.ISADMNISTRATOR = 'USER.ENABLE_ADMIN_ACTION';
      return true;
    }
    if (usersRole.length && usersRole.every(num => num === 1)) {
      this.ISADMNISTRATOR = 'USER.DISABLE_ADMIN_ACTION';
      return true;
    }
    return false;
  }

  isSystemAdmin(u: User): string {
    if (!u) {
      return "{{MISS}}";
    }
    let key: string = u.has_admin_role ? "USER.IS_ADMIN" : "USER.IS_NOT_ADMIN";
    this.translate.get(key).subscribe((res: string) => this.adminColumn = res);
    return this.adminColumn;
  }

  adminActions(u: User): string {
    if (!u) {
      return "{{MISS}}";
    }
    let key: string = u.has_admin_role ? "USER.DISABLE_ADMIN_ACTION" : "USER.ENABLE_ADMIN_ACTION";
    this.translate.get(key).subscribe((res: string) => this.adminMenuText = res);
    return this.adminMenuText;
  }

  public get inProgress(): boolean {
    return this.onGoing;
  }

  ngOnInit(): void {}

  ngOnDestroy(): void {
    if (this.deletionSubscription) {
      this.deletionSubscription.unsubscribe();
    }

    if (this.timerHandler) {
      clearInterval(this.timerHandler);
      this.timerHandler = null;
    }
  }

  //Filter items by keywords
  doFilter(terms: string): void {
    this.currentTerm = terms;
    this.originalUsers.then(users => {
      if (terms.trim() === "") {
        this.refreshUser((this.currentPage - 1) * 15, this.currentPage * 15);
      } else {
        this.users = users.filter(user => {
          return this.isMatchFilterTerm(terms, user.username);
        });
        this.forceRefreshView(5000);
      }
    });
  }

  // Disable the admin role for the specified user
  changeAdminRole(): void {
    let promiseLists: any[] = [];
    if (this.selectedRow.length) {
      if (this.ISADMNISTRATOR === 'USER.ENABLE_ADMIN_ACTION') {
        for (let i = 0; i < this.selectedRow.length; i++) {
          // Double confirm user is existing
          if (this.selectedRow[i].user_id === 0 || this.isMySelf(this.selectedRow[i].user_id)) {
            continue;
          }
          let updatedUser: User = new User();
          updatedUser.user_id = this.selectedRow[i].user_id;

          updatedUser.has_admin_role = 1; //Set as admin
          promiseLists.push(this.userService.updateUserRole(updatedUser));
        }
      }
      if (this.ISADMNISTRATOR === 'USER.DISABLE_ADMIN_ACTION') {
        for (let i = 0; i < this.selectedRow.length; i++) {
          // Double confirm user is existing
          if (this.selectedRow[i].user_id === 0 || this.isMySelf(this.selectedRow[i].user_id)) {
            continue;
          }
          let updatedUser: User = new User();
          updatedUser.user_id = this.selectedRow[i].user_id;

          updatedUser.has_admin_role = 0; //Set as none admin
          promiseLists.push(this.userService.updateUserRole(updatedUser));
        }
      }

        Promise.all(promiseLists).then(() => {
            this.selectedRow = [];
            this.refresh()
        })
        .catch(error => {
             this.selectedRow = [];
             this.msgHandler.handleError(error);
         });
      }
  }

  //Delete the specified user
  deleteUsers(users: User[]): void {
    let userArr: string[] = [];
    this.batchDelectionInfos = [];
    if (this.onlySelf) {
      return;
    }

    if (users && users.length) {
        users.forEach(user => {
          let initBatchMessage = new BatchInfo ();
          initBatchMessage.name = user.username;
          this.batchDelectionInfos.push(initBatchMessage);
          userArr.push(user.username);
        })
      this.deletionDialogService.addBatchInfoList(this.batchDelectionInfos);
    //Confirm deletion
    let msg: ConfirmationMessage = new ConfirmationMessage(
      "USER.DELETION_TITLE",
      "USER.DELETION_SUMMARY",
      userArr.join(','),
      users,
      ConfirmationTargets.USER,
      ConfirmationButtons.DELETE_CANCEL
    );
    this.deletionDialogService.openComfirmDialog(msg);
  }
  }

  delUser(users: User[]): void {
    //this.batchInfoDialog.open();
    let promiseLists: any[] = [];
    if (users && users.length) {
      users.forEach(user => {
        let findedList = this.batchDelectionInfos.find(data => data.name === user.username);
        if (this.isMySelf(user.user_id)) {
          this.translate.get('BATCH.DELETED_FAILURE').subscribe(res => {
            findedList = BathInfoChanges(findedList, res, false, true);
          });
        } else {
          promiseLists.push(this.delOperate(user.user_id, user.username));
        }
      });

      Promise.all(promiseLists).then((item) => {
        this.selectedRow = [];
        this.currentTerm = '';
        this.refresh();
      });
    }
  }
  delOperate(id: number, name:  string) {
    let findedList = this.batchDelectionInfos.find(data => data.name === name);
    return this.userService.deleteUser(id).then(() => {
      this.translate.get('BATCH.DELETED_SUCCESS').subscribe(res => {
        findedList = BathInfoChanges(findedList, res);
      });
    }).catch(error => {
      this.translate.get('BATCH.DELETED_FAILURE').subscribe(res => {
        findedList = BathInfoChanges(findedList, res, false, true);
      });
     });
  }

  //Refresh the user list
  refreshUser(from: number, to: number): void {
    this.selectedRow = [];
    //Start to get
    this.currentTerm = '';
    this.onGoing = true;

    this.originalUsers = this.userService.getUsers()
      .then(users => {
        this.onGoing = false;

        this.totalCount = users.length;
        this.users = users.slice(from, to);//First page

        this.forceRefreshView(5000);

        return users;
      })
      .catch(error => {
        this.onGoing = false;
        this.msgHandler.handleError(error);
        this.forceRefreshView(5000);
      });
  }

  //Add new user
  addNewUser(): void {
    if (!this.canCreateUser) {
      return;// No response to this hacking action
    }
    this.newUserDialog.open();
  }

  //Add user to the user list
  addUserToList(user: User): void {
    //Currently we can only add it by reloading all
    this.refresh();
  }

  //Data loading
  load(state: any): void {
    this.selectedRow = [];
    if (state && state.page) {
      if (this.originalUsers) {
        this.originalUsers.then(users => {
          this.users = users.slice(state.page.from, state.page.to + 1);
        });
        this.forceRefreshView(5000);
      } else {
        this.refreshUser(state.page.from, state.page.to + 1);
      }
    } else {
      //Refresh
      this.refresh();
    }
  }

  refresh(): void {
    this.currentPage = 1;//Refresh pagination
    this.refreshUser(0, 15);
  }

  SelectedChange(): void {
    this.forceRefreshView(5000);
  }

  forceRefreshView(duration: number): void {
    //Reset timer
    if (this.timerHandler) {
      clearInterval(this.timerHandler);
    }
    this.timerHandler = setInterval(() => this.ref.markForCheck(), 100);
    setTimeout(() => {
      if (this.timerHandler) {
        clearInterval(this.timerHandler);
        this.timerHandler = null;
      }
    }, duration);
  }

}

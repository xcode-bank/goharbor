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
import { Component } from '@angular/core';

import {Router,ActivatedRoute} from "@angular/router";
import {ReplicationRule} from "../replication-rule/replication-rule";
import {SessionService} from "../../shared/session.service";

@Component({
  selector: 'total-replication',
  templateUrl: 'total-replication-page.component.html'
})
export class TotalReplicationPageComponent {

  constructor(private router: Router,
              private session: SessionService,
              private activeRoute: ActivatedRoute){}
  customRedirect(rule: ReplicationRule): void {
    if (rule) {
      this.router.navigate(['../projects', rule.projects[0].project_id, 'replications'],  { relativeTo: this.activeRoute });
    }
  }

  public get isSystemAdmin(): boolean {
    let account = this.session.getCurrentUser();
    return account != null && account.has_admin_role > 0;
  }

  openEditPage(id: number): void {
      this.router.navigate([id, 'rule'],  { relativeTo: this.activeRoute });
  }

  openCreatePage(): void {
    this.router.navigate(['new-rule'],  { relativeTo: this.activeRoute });
  }
}

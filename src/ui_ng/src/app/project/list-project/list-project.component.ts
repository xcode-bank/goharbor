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
import {
    Component,
    Output,
    Input,
    ChangeDetectionStrategy,
    ChangeDetectorRef,
    OnDestroy
} from '@angular/core';
import { Router, NavigationExtras } from '@angular/router';
import { Project } from '../project';
import { ProjectService } from '../project.service';

import { SessionService } from '../../shared/session.service';
import { SearchTriggerService } from '../../base/global-search/search-trigger.service';
import { ProjectTypes, RoleInfo } from '../../shared/shared.const';
import { CustomComparator, doFiltering, doSorting, calculatePage } from '../../shared/shared.utils';

import { Comparator, State } from 'clarity-angular';
import { MessageHandlerService } from '../../shared/message-handler/message-handler.service';
import { StatisticHandler } from '../../shared/statictics/statistic-handler.service';
import { Subscription } from 'rxjs/Subscription';
import { ConfirmationDialogService } from '../../shared/confirmation-dialog/confirmation-dialog.service';
import { ConfirmationMessage } from '../../shared/confirmation-dialog/confirmation-message';
import { ConfirmationTargets, ConfirmationState, ConfirmationButtons } from '../../shared/shared.const';

@Component({
    selector: 'list-project',
    templateUrl: 'list-project.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush
})
export class ListProjectComponent implements OnDestroy {
    loading: boolean = true;
    projects: Project[] = [];
    filteredType: number = 0;//All projects
    searchKeyword: string = "";

    roleInfo = RoleInfo;
    repoCountComparator: Comparator<Project> = new CustomComparator<Project>("repo_count", "number");
    timeComparator: Comparator<Project> = new CustomComparator<Project>("creation_time", "date");
    accessLevelComparator: Comparator<Project> = new CustomComparator<Project>("public", "number");
    roleComparator: Comparator<Project> = new CustomComparator<Project>("current_user_role_id", "number");
    currentPage: number = 1;
    totalCount: number = 0;
    pageSize: number = 15;
    currentState: State;
    subscription: Subscription;

    constructor(
        private session: SessionService,
        private router: Router,
        private searchTrigger: SearchTriggerService,
        private proService: ProjectService,
        private msgHandler: MessageHandlerService,
        private statisticHandler: StatisticHandler,
        private deletionDialogService: ConfirmationDialogService,
        private ref: ChangeDetectorRef) {
        this.subscription = deletionDialogService.confirmationConfirm$.subscribe(message => {
            if (message &&
                message.state === ConfirmationState.CONFIRMED &&
                message.source === ConfirmationTargets.PROJECT) {
                let projectId = message.data;
                this.proService
                    .deleteProject(projectId)
                    .subscribe(
                    response => {
                        this.msgHandler.showSuccess('PROJECT.DELETED_SUCCESS');
                        let st: State = this.getStateAfterDeletion();
                        if (!st) {
                            this.refresh();
                        } else {
                            this.clrLoad(st);
                            this.statisticHandler.refresh();
                        }
                    },
                    error => {
                        if (error && error.status === 412) {
                            this.msgHandler.showError('PROJECT.FAILED_TO_DELETE_PROJECT', '');
                        } else {
                            this.msgHandler.handleError(error);
                        }
                    }
                    );

                let hnd = setInterval(() => ref.markForCheck(), 100);
                setTimeout(() => clearInterval(hnd), 2000);
            }
        });

        let hnd = setInterval(() => ref.markForCheck(), 100);
        setTimeout(() => clearInterval(hnd), 5000);
    }

    get showRoleInfo(): boolean {
        return this.filteredType !== 2;
    }

    public get isSystemAdmin(): boolean {
        let account = this.session.getCurrentUser();
        return account != null && account.has_admin_role > 0;
    }

    ngOnDestroy(): void {
        if (this.subscription) {
            this.subscription.unsubscribe();
        }
    }

    goToLink(proId: number): void {
        this.searchTrigger.closeSearch(true);

        let linkUrl = ['harbor', 'projects', proId, 'repositories'];
        this.router.navigate(linkUrl);
    }

    clrLoad(state: State) {
        //Keep state for future filtering and sorting
        this.currentState = state;

        let pageNumber: number = calculatePage(state);
        if (pageNumber <= 0) { pageNumber = 1; }

        this.loading = true;

        let passInFilteredType: number = undefined;
        if (this.filteredType > 0) {
            passInFilteredType = this.filteredType - 1;
        }
        this.proService.listProjects(this.searchKeyword, passInFilteredType, pageNumber, this.pageSize).toPromise()
            .then(response => {
                //Get total count
                if (response.headers) {
                    let xHeader: string = response.headers.get("X-Total-Count");
                    if (xHeader) {
                        this.totalCount = parseInt(xHeader, 0);
                    }
                }

                this.projects = response.json() as Project[];
                //Do customising filtering and sorting
                this.projects = doFiltering<Project>(this.projects, state);
                this.projects = doSorting<Project>(this.projects, state);

                this.loading = false;
            })
            .catch(error => {
                this.loading = false;
                this.msgHandler.handleError(error);
            });

        //Force refresh view
        let hnd = setInterval(() => this.ref.markForCheck(), 100);
        setTimeout(() => clearInterval(hnd), 5000);
    }

    newReplicationRule(p: Project) {
        if (p) {
            this.router.navigateByUrl(`/harbor/projects/${p.project_id}/replications?is_create=true`);
        }
    }

    toggleProject(p: Project) {
        if (p) {
            p.public === 0 ? p.public = 1 : p.public = 0;
            this.proService
                .toggleProjectPublic(p.project_id, p.public)
                .subscribe(
                response => {
                    this.msgHandler.showSuccess('PROJECT.TOGGLED_SUCCESS');
                    let pp: Project = this.projects.find((item: Project) => item.project_id === p.project_id);
                    if (pp) {
                        pp.public = p.public;
                        this.statisticHandler.refresh();
                    }
                },
                error => this.msgHandler.handleError(error)
                );

            //Force refresh view
            let hnd = setInterval(() => this.ref.markForCheck(), 100);
            setTimeout(() => clearInterval(hnd), 2000);
        }
    }

    deleteProject(p: Project) {
        let deletionMessage = new ConfirmationMessage(
            'PROJECT.DELETION_TITLE',
            'PROJECT.DELETION_SUMMARY',
            p.name,
            p.project_id,
            ConfirmationTargets.PROJECT,
            ConfirmationButtons.DELETE_CANCEL
        );
        this.deletionDialogService.openComfirmDialog(deletionMessage);
    }

    refresh(): void {
        this.currentPage = 1;
        this.filteredType = 0;
        this.searchKeyword = "";

        this.reload();
        this.statisticHandler.refresh();
    }

    doFilterProject(filter: number): void {
        this.currentPage = 1;
        this.filteredType = filter;
        this.reload();
    }

    doSearchProject(proName: string): void {
        this.currentPage = 1;
        this.searchKeyword = proName;
        this.reload();
    }

    reload(): void {
        let st: State = this.currentState;
        if (!st) {
            st = {
                page: {}
            };
        }
        st.page.from = 0;
        st.page.to = this.pageSize - 1;
        st.page.size = this.pageSize;

        this.clrLoad(st);
    }

    getStateAfterDeletion(): State {
        let total: number = this.totalCount - 1;
        if (total <= 0) { return null; }

        let totalPages: number = Math.ceil(total / this.pageSize);
        let targetPageNumber: number = this.currentPage;

        if (this.currentPage > totalPages) {
            targetPageNumber = totalPages;//Should == currentPage -1
        }

        let st: State = this.currentState;
        if (!st) {
            st = { page: {} };
        }
        st.page.size = this.pageSize;
        st.page.from = (targetPageNumber - 1) * this.pageSize;
        st.page.to = targetPageNumber * this.pageSize - 1;

        return st;
    }

}
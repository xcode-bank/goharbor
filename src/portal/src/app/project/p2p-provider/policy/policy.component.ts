// Copyright Project Harbor Authors
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
import { distinctUntilChanged, finalize, switchMap } from 'rxjs/operators';
import { TranslateService } from '@ngx-translate/core';
import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { MessageHandlerService } from '../../../shared/message-handler/message-handler.service';
import { Project } from '../../project';
import { ConfirmationButtons, ConfirmationState, ConfirmationTargets } from '../../../shared/shared.const';
import { ConfirmationMessage } from '../../../shared/confirmation-dialog/confirmation-message';
import { ConfirmationDialogComponent } from '../../../shared/confirmation-dialog/confirmation-dialog.component';
import { clone, CustomComparator } from '../../../../lib/utils/utils';
import { forkJoin, Observable, Subject, Subscription } from 'rxjs';
import {
  ClrDatagridComparatorInterface,
  UserPermissionService,
  USERSTATICPERMISSION
} from '../../../../lib/services';
import { ClrLoadingState } from '@clr/angular';
import {
  EXECUTION_STATUS,
  FILTER_TYPE,
  P2pProviderService,
  PROJECT_SEVERITY_LEVEL_TO_TEXT_MAP,
  TIME_OUT,
  TRIGGER,
  TRIGGER_I18N_MAP
} from '../p2p-provider.service';
import { PreheatPolicy } from '../../../../../ng-swagger-gen/models/preheat-policy';
import { PreheatService } from '../../../../../ng-swagger-gen/services/preheat.service';
import { AddP2pPolicyComponent } from '../add-p2p-policy/add-p2p-policy.component';
import { Execution } from '../../../../../ng-swagger-gen/models/execution';
import { Metrics } from '../../../../../ng-swagger-gen/models/metrics';
import { ProviderUnderProject } from '../../../../../ng-swagger-gen/models/provider-under-project';

@Component({
  templateUrl: './policy.component.html',
  styleUrls: ['./policy.component.scss']
})
export class PolicyComponent implements OnInit, OnDestroy {
  @ViewChild(AddP2pPolicyComponent, { static: false } )
  addP2pPolicyComponent: AddP2pPolicyComponent;
  @ViewChild("confirmationDialogComponent", { static: false })
  confirmationDialogComponent: ConfirmationDialogComponent;
  projectId: number;
  projectName: string;
  selectedRow: PreheatPolicy;
  policyList: PreheatPolicy[] = [];
  providers: ProviderUnderProject[] = [];
  metadata: any;
  loading: boolean = false;
  hasCreatPermission: boolean = false;
  hasUpdatePermission: boolean = false;
  hasDeletePermission: boolean = false;
  addBtnState: ClrLoadingState = ClrLoadingState.DEFAULT;
  executing: boolean = false;
  isOpenFilterTag: boolean = false;
  selectedExecutionRow: Execution;
  jobsLoading: boolean = false;
  stopLoading: boolean = false;
  creationTimeComparator: ClrDatagridComparatorInterface<Execution> = new CustomComparator<Execution>("creation_time", "date");
  executionList: Execution[] = [];
  currentExecutionPage: number = 1;
  pageSize: number = 10;
  totalExecutionCount: number = 0;
  filterKey: string = 'id';
  searchString: string;
  private _searchSubject: Subject<string> = new Subject<string>();
  private _searchSubscription: Subscription;
  project: Project;
  severity_map: any = PROJECT_SEVERITY_LEVEL_TO_TEXT_MAP;
  timeout: any;
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private translate: TranslateService,
    private p2pProviderService: P2pProviderService,
    private messageHandlerService: MessageHandlerService,
    private userPermissionService: UserPermissionService,
    private preheatService: PreheatService) { }

  ngOnInit() {
    this.subscribeSearch();
    this.projectId = +this.route.snapshot.parent.parent.params['id'];
    const resolverData = this.route.snapshot.parent.parent.data;
    if (resolverData) {
      const project = <Project>(resolverData["projectResolver"]);
      this.projectName = project.name;
    }
    this.getPermissions();
    this.refresh();
  }
  ngOnDestroy(): void {
     if (this._searchSubscription) {
       this._searchSubscription.unsubscribe();
       this._searchSubscription = null;
     }
     if (this.timeout) {
       clearTimeout(this.timeout);
       this.timeout = null;
     }
  }
  getPermissions() {
    const permissionsList: Observable<boolean>[] = [];
    permissionsList.push(this.userPermissionService.getPermission(this.projectId,
      USERSTATICPERMISSION.P2P_PROVIDER.KEY, USERSTATICPERMISSION.P2P_PROVIDER.VALUE.CREATE));
    permissionsList.push(this.userPermissionService.getPermission(this.projectId,
      USERSTATICPERMISSION.P2P_PROVIDER.KEY, USERSTATICPERMISSION.P2P_PROVIDER.VALUE.UPDATE));
    permissionsList.push(this.userPermissionService.getPermission(this.projectId,
      USERSTATICPERMISSION.P2P_PROVIDER.KEY, USERSTATICPERMISSION.P2P_PROVIDER.VALUE.DELETE));
    this.addBtnState = ClrLoadingState.LOADING;
    forkJoin(...permissionsList).subscribe(Rules => {
      [this.hasCreatPermission, this.hasUpdatePermission, this.hasDeletePermission] = Rules;
      this.addBtnState = ClrLoadingState.SUCCESS;
      if (this.hasCreatPermission) {
        this.getProviders();
      }
    }, error => {
      this.messageHandlerService.error(error);
      this.addBtnState = ClrLoadingState.ERROR;
    });
  }
  getProviders() {
    this.preheatService.ListProvidersUnderProject({projectName: this.projectName})
      .subscribe(res => {
        this.providers = res;
      });
  }
  refresh() {
    this.selectedRow = null;
    this.getPolicies();
  }
  getPolicies() {
    this.loading = true;
    this.preheatService.ListPolicies({projectName: this.projectName})
      .pipe(finalize(() => (this.loading = false)))
      .subscribe(
        response => {
          this.policyList = response;
        },
        error => {
          this.messageHandlerService.handleError(error);
        }
      );
  }
 switchStatus() {
    let content = '';
    this.translate.get(
      !this.selectedRow.enabled
        ? 'P2P_PROVIDER.ENABLED_POLICY_SUMMARY'
        : 'P2P_PROVIDER.DISABLED_POLICY_SUMMARY'
      , {name: this.selectedRow.name}).subscribe((res) => {
      content = res;
      let message = new ConfirmationMessage(
        !this.selectedRow.enabled ? 'P2P_PROVIDER.ENABLED_POLICY_TITLE' : 'P2P_PROVIDER.DISABLED_POLICY_TITLE',
        content,
        '',
        {},
        ConfirmationTargets.P2P_PROVIDER,
        !this.selectedRow.enabled ? ConfirmationButtons.ENABLE_CANCEL : ConfirmationButtons.DISABLE_CANCEL
      );
      this.confirmationDialogComponent.open(message);
    });
  }
  confirmSwitch(message) {
    if (message &&  message.source === ConfirmationTargets.P2P_PROVIDER_STOP &&
      message.state === ConfirmationState.CONFIRMED) {
      this.stopLoading = true;
      const execution: Execution = clone(this.selectedExecutionRow);
      execution.status = EXECUTION_STATUS.STOPPED;
      this.preheatService.StopExecution({
        projectName: this.projectName,
        preheatPolicyName: this.selectedRow.name,
        executionId: this.selectedExecutionRow.id,
        execution: execution
      }).pipe(finalize(() => this.executing = false))
        .subscribe(response => {
          this.messageHandlerService.showSuccess('P2P_PROVIDER.STOP_SUCCESSFULLY');
        }, error => {
          this.messageHandlerService.error(error);
        });
    }
    if (message &&  message.source === ConfirmationTargets.P2P_PROVIDER_EXECUTE &&
      message.state === ConfirmationState.CONFIRMED) {
      this.executing = true;
      this.preheatService.ManualPreheat({
        projectName: this.projectName,
        preheatPolicyName: this.selectedRow.name,
        policy: this.selectedRow
      }).pipe(finalize(() => this.executing = false))
        .subscribe(response => {
          this.messageHandlerService.showSuccess('P2P_PROVIDER.EXECUTE_SUCCESSFULLY');
          if (this.selectedRow) {
            this.refreshJobs();
          }
        }, error => {
          this.messageHandlerService.error(error);
        });
    }
    if (message &&
      message.source === ConfirmationTargets.P2P_PROVIDER &&
      message.state === ConfirmationState.CONFIRMED) {
      if (JSON.stringify(message.data) === '{}') {
        this.preheatService.UpdatePolicy({
          projectName: this.projectName,
          preheatPolicyName: this.selectedRow.name,
          policy:  Object.assign({}, this.selectedRow, { enabled: !this.selectedRow.enabled })
        }).subscribe(
            response => {
              this.messageHandlerService.showSuccess('P2P_PROVIDER.UPDATED_SUCCESSFULLY');
              this.refresh();
            },
            error => {
              this.messageHandlerService.handleError(error);
            }
          );
      }
    }
    if (message &&
      message.source === ConfirmationTargets.P2P_PROVIDER_DELETE &&
      message.state === ConfirmationState.CONFIRMED) {
      const observableLists: Observable<any>[] = [];
      observableLists.push(this.preheatService.DeletePolicy({
        projectName: this.projectName,
        preheatPolicyName: this.selectedRow.name
      }));
      forkJoin(...observableLists).subscribe(
        response => {
          this.messageHandlerService.showSuccess('P2P_PROVIDER.DELETE_SUCCESSFULLY');
          this.refresh();
        },
        error => {
          this.messageHandlerService.handleError(error);
        }
      );
    }
  }
  newPolicy() {
    this.addP2pPolicyComponent.isOpen = true;
    this.addP2pPolicyComponent.isEdit = false;
    this.addP2pPolicyComponent.resetForAdd();
  }
  editPolicy() {
    if (this.selectedRow) {
      this.addP2pPolicyComponent.isOpen = true;
      this.addP2pPolicyComponent.isEdit = true;
      this.addP2pPolicyComponent.policy = clone(this.selectedRow);
      const filter: any[] = JSON.parse(this.selectedRow.filters);
      if (filter && filter.length) {
        filter.forEach(item => {
          if (item.type === FILTER_TYPE.REPOS && item.value) {
            this.addP2pPolicyComponent.repos = item.value.replace(/[{}]/g, "");
          }
          if (item.type === FILTER_TYPE.TAG && item.value) {
            this.addP2pPolicyComponent.tags = item.value.replace(/[{}]/g, "");
          }
          if (item.type === FILTER_TYPE.SIGNATURE) {
            this.addP2pPolicyComponent.onlySignedImages = item.value;
          }
          if (item.type === FILTER_TYPE.LABEL && item.value) {
            this.addP2pPolicyComponent.labels = item.value.replace(/[{}]/g, "");
          }
          if (item.type === FILTER_TYPE.VULNERABILITY) {
            this.addP2pPolicyComponent.severity = item.value;
          }
        });
      }
      const trigger: any = JSON.parse(this.selectedRow.trigger);
      if (trigger) {
        this.addP2pPolicyComponent.triggerType = trigger.type;
        this.addP2pPolicyComponent.cron = trigger.trigger_setting.cron;
      }
      this.addP2pPolicyComponent.currentForm.reset({
        provider: this.addP2pPolicyComponent.policy.provider_id,
        name: this.addP2pPolicyComponent.policy.name,
        description: this.addP2pPolicyComponent.policy.description,
        repo: this.addP2pPolicyComponent.repos,
        tag: this.addP2pPolicyComponent.tags,
        onlySignedImages: this.addP2pPolicyComponent.onlySignedImages,
        severity: this.addP2pPolicyComponent.severity,
        label: this.addP2pPolicyComponent.labels,
        triggerType: this.addP2pPolicyComponent.triggerType
      });
      this.addP2pPolicyComponent.originPolicyForEdit = clone(this.selectedRow);
      this.addP2pPolicyComponent.originReposForEdit = this.addP2pPolicyComponent.repos;
      this.addP2pPolicyComponent.originTagsForEdit = this.addP2pPolicyComponent.tags;
      this.addP2pPolicyComponent.originOnlySignedImagesForEdit = this.addP2pPolicyComponent.onlySignedImages;
      this.addP2pPolicyComponent.originSeverityForEdit = this.addP2pPolicyComponent.severity;
      this.addP2pPolicyComponent.originLabelsForEdit = this.addP2pPolicyComponent.labels;
      this.addP2pPolicyComponent.originTriggerTypeForEdit = this.addP2pPolicyComponent.triggerType;
      this.addP2pPolicyComponent.originCronForEdit = this.addP2pPolicyComponent.cron;
    }
  }
  deletePolicy() {
    const names: string[] = [];
    names.push(this.selectedRow.name);
    let content = '';
    this.translate.get(
         'P2P_PROVIDER.DELETE_POLICY_SUMMARY'
      , {names:  names.join(',')}).subscribe((res) => content = res);
    const msg: ConfirmationMessage = new ConfirmationMessage(
      "SCANNER.CONFIRM_DELETION",
      content,
      names.join(','),
      this.selectedRow,
      ConfirmationTargets.P2P_PROVIDER_DELETE,
      ConfirmationButtons.DELETE_CANCEL
    );
    this.confirmationDialogComponent.open(msg);
  }
  executePolicy() {
    if (this.selectedRow  && this.selectedRow.enabled) {
      const message = new ConfirmationMessage(
        "P2P_PROVIDER.EXECUTE_TITLE",
        "P2P_PROVIDER.EXECUTE_SUMMARY",
        this.selectedRow.name,
        this.selectedRow,
        ConfirmationTargets.P2P_PROVIDER_EXECUTE,
        ConfirmationButtons.CONFIRM_CANCEL
      );
      this.confirmationDialogComponent.open(message);
    }
  }
  success(isAdd: boolean) {
    let message: string;
    if (isAdd) {
      message = 'P2P_PROVIDER.ADDED_SUCCESS';
    } else {
      message = 'P2P_PROVIDER.UPDATED_SUCCESS';
    }
    this.messageHandlerService.showSuccess(message);
    this.refresh();
  }
  clrLoadJobs(chosenPolicy: PreheatPolicy, withLoading: boolean) {
    if (this.selectedRow) {
      if (withLoading) {
        // if datagrid is under control of *ngIf, should add timeout in case of ng changes checking error
        setTimeout(() => {
          this.jobsLoading = true;
        });
      }
      let params: string;
      if (this.searchString) {
         params =  encodeURIComponent(`${this.filterKey}=~${this.searchString}`);
      }
      this.preheatService.ListExecutionsResponse({
        projectName: this.projectName,
        preheatPolicyName: chosenPolicy ? chosenPolicy.name : this.selectedRow.name,
        page: this.currentExecutionPage,
        pageSize: this.pageSize,
        q: params
      }).pipe(finalize(() => this.jobsLoading = false))
        .subscribe(response => {
          if (response.headers) {
            let xHeader: string = response.headers.get('x-total-count');
            if (xHeader) {
              this.totalExecutionCount = parseInt(xHeader, 0);
            }
          }
          this.executionList = response.body;
          if (this.executionList && this.executionList.length) {
            for (let i = 0; i < this.executionList.length; i++) {
              if (this.p2pProviderService.willChangStatus(this.executionList[i].status)) {
                if (!this.timeout) {
                  this.timeout = setTimeout(() => {
                    this.clrLoadJobs(null, false);
                  }, TIME_OUT);
                }
              }
            }
          }
        }, error => {
          this.messageHandlerService.handleError(error);
        });
    }
  }
  refreshJobs (chosenPolicy?: PreheatPolicy) {
    this.executionList = [];
    this.currentExecutionPage = 1;
    this.totalExecutionCount = 0;
    this.filterKey = 'id';
    this.searchString = null;
    this.clrLoadJobs(chosenPolicy, true);
  }
  openStopExecutionsDialog() {
    if (this.selectedExecutionRow) {
      const stopExecutionsMessage = new ConfirmationMessage(
        "P2P_PROVIDER.STOP_TITLE",
        "P2P_PROVIDER.STOP_SUMMARY",
        this.selectedExecutionRow.id + '',
        this.selectedExecutionRow,
        ConfirmationTargets.P2P_PROVIDER_STOP,
        ConfirmationButtons.CONFIRM_CANCEL
      );
      this.confirmationDialogComponent.open(stopExecutionsMessage);
    }
  }
  goToLink(executionId: number) {
    const linkUrl = ["harbor",
      "projects", `${this.projectId}`, "p2p-provider", `${this.selectedRow.name}`, "executions", `${executionId}`, "tasks"];
    this.router.navigate(linkUrl);
  }
  getTriggerTypeI18n(trigger: string): string {
    if (JSON.parse(trigger).type) {
      return TRIGGER_I18N_MAP[JSON.parse(trigger).type];
    }
    return TRIGGER_I18N_MAP[TRIGGER.MANUAL];
  }
  isScheduled(trigger: string): boolean {
    return JSON.parse(trigger).type === TRIGGER.SCHEDULED;
  }
  isEventBased(trigger: string): boolean {
    return JSON.parse(trigger).type === TRIGGER.EVENT_BASED;
  }
  getScheduledCron(trigger: string): string {
    return JSON.parse(trigger).trigger_setting.cron;
  }
  getDuration(e: Execution): string {
    return this.p2pProviderService.getDuration(e.start_time, e.end_time);
  }
  getValue(filter: string, type: string): string {
    const arr: any[] = JSON.parse(filter);
    if (arr && arr.length) {
      for (let i = 0; i < arr.length; i++) {
        if (arr[i].type === type && arr[i].value) {
          return (arr[i].value + "").replace(/[{}]/g, "");
        }
      }
    }
    return "";
  }
  getSuccessRate(m: Metrics): number {
    if (m && m.task_count && m.success_task_count) {
      return m.success_task_count / m.task_count;
    }
    return 0;
  }
  selectFilterKey($event: any): void {
    this.filterKey = $event['target'].value;
  }
  openFilter(isOpen: boolean): void {
    this.isOpenFilterTag = isOpen;
  }

  doFilter(terms: string): void {
    this.searchString = terms;
    if (terms.trim()) {
      this._searchSubject.next(terms.trim());
    } else {
      this.clrLoadJobs(null, true);
    }
  }
  subscribeSearch() {
    if (!this._searchSubscription) {
      this._searchSubscription = this._searchSubject.pipe(
        distinctUntilChanged(),
        switchMap(searchString => {
          this.jobsLoading = true;
          let params: string;
          if (this.searchString) {
            params = encodeURIComponent(`${this.filterKey}=~${searchString}`);
          }
          return this.preheatService.ListExecutionsResponse({
            projectName: this.projectName,
            preheatPolicyName: this.selectedRow.name,
            page: 1,
            pageSize: this.pageSize,
            q: params
          }).pipe(finalize(() => this.jobsLoading = false));
        })).subscribe(response => {
           if (response.headers) {
            let xHeader: string = response.headers.get('x-total-count');
            if (xHeader) {
              this.totalExecutionCount = parseInt(xHeader, 0);
            }
           }
           this.executionList = response.body;
      });
    }
  }
  canStop(): boolean {
    return this.selectedExecutionRow && this.p2pProviderService.willChangStatus(this.selectedExecutionRow.status);
  }
}

export const REPLICATION_TEMPLATE: string = `
<div class="row" style="position:relative;">
  <div>
    <div class="row flex-items-xs-between rightPos">
      <div class="flex-xs-middle option-right">
        <hbr-filter [withDivider]="true" filterPlaceholder='{{"REPLICATION.FILTER_POLICIES_PLACEHOLDER" | translate}}' (filter)="doSearchRules($event)" [currentValue]="search.ruleName"></hbr-filter>
        <span class="refresh-btn" (click)="refreshRules()">
          <clr-icon shape="refresh"></clr-icon>
        </span>
      </div>
    </div>
    </div>
    <div class="col-lg-12 col-md-12 col-sm-12 col-xs-12">
      <hbr-list-replication-rule #listReplicationRule [readonly]="readonly" [projectId]="projectId" [isSystemAdmin]="isSystemAdmin"  (replicateManual)=replicateManualRule($event) (selectOne)="selectOneRule($event)" (hideJobs)="hideJobs()"  (openNewRule)="openModal()" (editOne)="openEditRule($event)" (reload)="reloadRules($event)" [loading]="loading" [withReplicationJob]="withReplicationJob" (redirect)="customRedirect($event)"></hbr-list-replication-rule>
    </div>
    <div class="col-lg-12 col-md-12 col-sm-12 col-xs-12" style="padding-left:0px;">
    <div *ngIf="withReplicationJob" class="col-lg-12 col-md-12 col-sm-12 col-xs-12">
      <div class="row flex-items-xs-between" style="height:60px;">
        <h5 class="flex-items-xs-bottom option-left-down" style="margin-left: 14px;">{{'REPLICATION.REPLICATION_JOBS' | translate}}</h5>
        <div class="flex-items-xs-bottom option-right-down">
          <button class="btn btn-link" (click)="toggleSearchJobOptionalName(currentJobSearchOption)">{{toggleJobSearchOption[currentJobSearchOption] | translate}}</button>
          <hbr-filter [withDivider]="true" filterPlaceholder='{{"REPLICATION.FILTER_JOBS_PLACEHOLDER" | translate}}' (filter)="doSearchJobs($event)" [currentValue]="search.repoName" ></hbr-filter>
          <span class="refresh-btn" (click)="refreshJobs()">
            <clr-icon shape="refresh"></clr-icon>
          </span>
        </div>
      </div>
      <div class="row flex-items-xs-right option-right" [hidden]="currentJobSearchOption === 0">
        <div class="select" style="float: left;">
          <select (change)="doFilterJobStatus($event)">
          <option *ngFor="let j of jobStatus" value="{{j.key}}" [selected]="currentJobStatus.key === j.key">{{j.description | translate}}</option>
          </select>
        </div>
        <div class="flex-items-xs-middle">    
          <hbr-datetime [dateInput]="search.startTime" (search)="doJobSearchByStartTime($event)"></hbr-datetime>
          <hbr-datetime [dateInput]="search.endTime" [oneDayOffset]="true" (search)="doJobSearchByEndTime($event)"></hbr-datetime>
        </div>
      </div>
    </div>
    <div *ngIf="withReplicationJob" class="col-lg-12 col-md-12 col-sm-12 col-xs-12">
      <clr-datagrid [clrDgLoading]="jobsLoading" (clrDgRefresh)="clrLoadJobs($event)"><clr-dg-action-bar>
            <div class="btn-group">
                <button type="button" class="btn btn-sm btn-secondary" *ngIf="isSystemAdmin" [disabled]="!(jobs && jobs.length>0) || isStopOnGoing" (click)="stopJobs()">{{'REPLICATION.STOPJOB' | translate}}</button>
            </div>
        </clr-dg-action-bar>
        <clr-dg-column [clrDgField]="'repository'">{{'REPLICATION.NAME' | translate}}</clr-dg-column>
        <clr-dg-column [clrDgField]="'status'">{{'REPLICATION.STATUS' | translate}}</clr-dg-column>
        <clr-dg-column [clrDgField]="'operation'">{{'REPLICATION.OPERATION' | translate}}</clr-dg-column>
        <clr-dg-column [clrDgSortBy]="creationTimeComparator">{{'REPLICATION.CREATION_TIME' | translate}}</clr-dg-column>
        <clr-dg-column [clrDgSortBy]="updateTimeComparator">{{'REPLICATION.UPDATE_TIME' | translate}}</clr-dg-column>
        <clr-dg-column>{{'REPLICATION.LOGS' | translate}}</clr-dg-column>
        <clr-dg-placeholder>{{'REPLICATION.JOB_PLACEHOLDER' | translate }}</clr-dg-placeholder>
        <clr-dg-row *ngFor="let j of jobs">
            <clr-dg-cell>{{j.repository}}</clr-dg-cell>
            <clr-dg-cell>{{j.status}}</clr-dg-cell>
            <clr-dg-cell>{{j.operation}}</clr-dg-cell>
            <clr-dg-cell>{{j.creation_time | date: 'short'}}</clr-dg-cell>
            <clr-dg-cell>{{j.update_time | date: 'short'}}</clr-dg-cell>
            <clr-dg-cell>
             <span *ngIf="j.status=='pending'; else elseBlock" class="label">{{'REPLICATION.NO_LOGS' | translate}}</span>
                <ng-template #elseBlock>
                    <a href="javascript:void(0);" (click)="viewLog(j.id)">
                <clr-icon shape="clipboard"></clr-icon>
              </a></ng-template>
            </clr-dg-cell>
        </clr-dg-row>
        <clr-dg-footer>
            <span *ngIf="showPaginationIndex">{{pagination.firstItem + 1}} - {{pagination.lastItem + 1}} {{'REPLICATION.OF' | translate}}</span>
            {{pagination.totalItems}} {{'REPLICATION.ITEMS' | translate}}
            <clr-dg-pagination #pagination [(clrDgPage)]="currentPage" [clrDgPageSize]="pageSize" [clrDgTotalItems]="totalCount"></clr-dg-pagination>
        </clr-dg-footer>
      </clr-datagrid>
    </div>
    </div>
    <job-log-viewer #replicationLogViewer></job-log-viewer>
    <confirmation-dialog #replicationConfirmDialog [batchInfors]="batchDelectionInfos" (confirmAction)="confirmReplication($event)"></confirmation-dialog>
</div>`;
export const CONFIRMATION_DIALOG_TEMPLATE: string = `
<clr-modal [(clrModalOpen)]="opened" [clrModalClosable]="false" [clrModalStaticBackdrop]="true">
    <h3 class="modal-title" class="confirmation-title" style="margin-top: 0px;">{{dialogTitle}}</h3>
    <div class="modal-body">
        <div class="confirmation-icon-inline">
            <clr-icon shape="warning" class="is-warning" size="64"></clr-icon>
        </div>
        <div class="confirmation-content">{{dialogContent}}</div>
       <div>
            <ul class="batchInfoUl">
                <li *ngFor="let info of batchInfors">
                   <span> <i class="spinner spinner-inline spinner-pos" [hidden]='!info.loading'></i>&nbsp;&nbsp;{{info.name}}</span>
                    <span *ngIf="!info.errorInfo.length" [style.color]="colorChange(info)">{{info.status}}</span>
                    <span *ngIf="info.errorInfo.length" [style.color]="colorChange(info)">
                        <a (click)="toggleErrorTitle(errorInfo)" >{{info.status}}</a><br>
                        <i #errorInfo style="display: none;">{{info.errorInfo}}</i>
                    </span>
                </li>
            </ul>
        </div>
    </div>
    <div class="modal-footer" [ngSwitch]="buttons">
       <ng-template [ngSwitchCase]="0">
        <button type="button" class="btn btn-outline" (click)="cancel()">{{'BUTTON.CANCEL' | translate}}</button>
        <button type="button" class="btn btn-primary" (click)="confirm()">{{'BUTTON.CONFIRM' | translate}}</button>
       </ng-template>
       <ng-template [ngSwitchCase]="1">
        <button type="button" class="btn btn-outline" (click)="cancel()">{{'BUTTON.NO' | translate}}</button>
        <button type="button" class="btn btn-primary" (click)="confirm()">{{ 'BUTTON.YES' | translate}}</button>
       </ng-template>
       <ng-template [ngSwitchCase]="2">
        <button type="button" class="btn btn-outline" (click)="cancel()" [hidden]="isDelete">{{'BUTTON.CANCEL' | translate}}</button>
         <button type="button" class="btn btn-danger" (click)="operate()" [hidden]="isDelete">{{'BUTTON.DELETE' | translate}}</button>
        <button type="button" class="btn btn-primary" (click)="cancel()" [disabled]="!batchOverStatus"  [hidden]="!isDelete">{{'BUTTON.CLOSE' | translate}}</button>
       </ng-template>
       <ng-template [ngSwitchCase]="3">
        <button type="button" class="btn btn-primary" (click)="cancel()">{{'BUTTON.CLOSE' | translate}}</button>
       </ng-template>
       <ng-template [ngSwitchCase]="4">
        <button type="button" class="btn btn-outline" (click)="cancel()" [hidden]="isDelete">{{'BUTTON.CANCEL' | translate}}</button>
         <button type="button" class="btn btn-primary" (click)="operate()" [hidden]="isDelete">{{'BUTTON.REPLICATE' | translate}}</button>
        <button type="button" class="btn btn-primary" (click)="cancel()" [disabled]="!batchOverStatus"  [hidden]="!isDelete">{{'BUTTON.CLOSE' | translate}}</button>
       </ng-template>
    </div>
</clr-modal>
`;
export const INLINE_ALERT_TEMPLATE: string = `
<clr-alert [clrAlertType]="inlineAlertType" [clrAlertClosable]="inlineAlertClosable" [(clrAlertClosed)]="alertClose" [clrAlertAppLevel]="useAppLevelStyle">
    <div class="alert-item">
        <span class="alert-text" [class.alert-text-blink]="blinking">
            {{errorMessage}}
        </span>
        <div class="alert-actions" *ngIf="showCancelAction">
            <button class="btn btn-sm btn-link alert-btn-link" (click)="close()">{{'BUTTON.NO' | translate}}</button>
            <button class="btn btn-sm btn-link alert-btn-link" (click)="confirmCancel()">{{'BUTTON.YES' | translate}}</button>
        </div>
    </div>
</clr-alert>
`;
export const DATETIME_PICKER_TEMPLATE: string = `
<clr-icon shape="date"></clr-icon>
<label aria-haspopup="true" role="tooltip" [class.invalid]="dateInvalid" class="tooltip tooltip-validation tooltip-sm">
  <input type="date" #searchTime="ngModel" [(ngModel)]="dateInput" name="searchTime" placeholder="dd/mm/yyyy" dateValidator (change)="doSearch()">
  <span *ngIf="dateInvalid" class="tooltip-content">
    {{'AUDIT_LOG.INVALID_DATE' | translate }}
  </span>
</label>
`;
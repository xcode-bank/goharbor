export const REPLICATION_CONFIG_HTML: string = `
<form #replicationConfigFrom="ngForm" class="compact">
<section class="form-block" style="margin-top:0px;margin-bottom:0px;">
<label style="font-size:14px;font-weight:600;" *ngIf="showSubTitle">{{'CONFIG.REPLICATION' | translate}}</label>
<div class="form-group">
  <label for="verifyRemoteCert">{{'CONFIG.VERIFY_REMOTE_CERT' | translate }}</label>
    <clr-checkbox name="verifyRemoteCert" id="verifyRemoteCert" [(ngModel)]="replicationConfig.verify_remote_cert.value" [disabled]="!editable">
      <a href="javascript:void(0)" role="tooltip" aria-haspopup="true" class="tooltip tooltip-lg tooltip-top-right" style="top:-8px;">
        <clr-icon shape="info-circle" class="info-tips-icon" size="24"></clr-icon>
        <span class="tooltip-content">{{'CONFIG.TOOLTIP.VERIFY_REMOTE_CERT' | translate }}</span>
      </a>
     </clr-checkbox>
</div>
</section>
</form>
`;
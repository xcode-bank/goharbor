export const SYSTEM_SETTINGS_HTML: string = `
<form #systemConfigFrom="ngForm" class="compact">
<section class="form-block" style="margin-top:0px;margin-bottom:0px;">
  <label style="font-size:14px;font-weight:600;" *ngIf="showSubTitle">{{'CONFIG.SYSTEM' | translate}}</label>
  <div class="form-group">
      <label for="tokenExpiration" class="required">{{'CONFIG.TOKEN_EXPIRATION' | translate}}</label>
      <label for="tokenExpiration" aria-haspopup="true" role="tooltip" class="tooltip tooltip-validation tooltip-md tooltip-top-right" [class.invalid]="tokenExpirationInput.invalid && (tokenExpirationInput.dirty || tokenExpirationInput.touched)">
        <input name="tokenExpiration" type="text" #tokenExpirationInput="ngModel" [(ngModel)]="systemSettings.token_expiration.value" 
         required 
         pattern="^[1-9]{1}[0-9]*$"
         id="tokenExpiration" 
         size="20" [disabled]="!editable">
         <span class="tooltip-content">
               {{'TOOLTIP.NUMBER_REQUIRED' | translate}}
          </span>
       </label>
       <a href="javascript:void(0)" role="tooltip" aria-haspopup="true" class="tooltip tooltip-top-right">
         <clr-icon shape="info-circle" class="info-tips-icon" size="24"></clr-icon>
         <span class="tooltip-content">{{'CONFIG.TOOLTIP.TOKEN_EXPIRATION' | translate}}</span>
       </a>
   </div>
   <div class="form-group" *ngIf="canDownloadCert">
      <label for="certDownloadLink" class="required">{{'CONFIG.ROOT_CERT' | translate}}</label>
      <a #certDownloadLink [href]="downloadLink" target="_blank">{{'CONFIG.ROOT_CERT_LINK' | translate}}</a>
       <a href="javascript:void(0)" role="tooltip" aria-haspopup="true" class="tooltip tooltip-top-right">
         <clr-icon shape="info-circle" class="info-tips-icon" size="24"></clr-icon>
         <span class="tooltip-content">{{'CONFIG.TOOLTIP.ROOT_CERT_DOWNLOAD' | translate}}</span>
       </a>
   </div>
</section>
</form>
`;
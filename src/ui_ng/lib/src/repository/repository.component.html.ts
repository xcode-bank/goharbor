export const REPOSITORY_TEMPLATE = `
<section class="overview-section">
  <div class="title-wrapper">
    <div class="title-block arrow-block">
      <clr-icon class="rotate-90 arrow-back" shape="arrow" size="36" (click)="goBack()"></clr-icon>
    </div>
    <div class="title-block">
      <h2 sub-header-title>{{repoName}}</h2>
    </div>
  </div>
</section>

<section class="detail-section">
  <div class="col-lg-12 col-md-12 col-sm-12 col-xs-12">
    <span class="spinner spinner-inline" [hidden]="inProgress === false"></span>
    <ul id="configTabs" class="nav" role="tablist">
      <li role="presentation" class="nav-item">
          <button id="repo-info" class="btn btn-link nav-link" aria-controls="info" [class.active]='isCurrentTabLink("repo-info")' type="button" (click)='tabLinkClick("repo-info")'>{{'REPOSITORY.INFO' | translate}}</button>
      </li>
      <li role="presentation" class="nav-item">
          <button id="repo-image" class="btn btn-link nav-link active" aria-controls="image" [class.active]='isCurrentTabLink("repo-image")' type="button" (click)='tabLinkClick("repo-image")'>{{'REPOSITORY.IMAGE' | translate}}</button>
      </li>
    </ul>
    <section id="info" role="tabpanel" aria-labelledby="repo-info" [hidden]='!isCurrentTabContent("info")'>
      <form #repoInfoForm="ngForm">
        <div id="info-edit-button">
          <button class="btn btn-sm" [disabled]="editing || !hasProjectAdminRole " (click)="editInfo()" ><clr-icon shape="pencil" size="16"></clr-icon>&nbsp;{{'BUTTON.EDIT' | translate}}</button>
        </div>
        <div *ngIf="!editing">
          <div *ngIf="!hasInfo()" class="no-info-div">
            <p>{{'REPOSITORY.NO_INFO' | translate }}<p>
          </div>
          <div *ngIf="hasInfo()" class="info-div">
            <pre class="info-pre">{{ imageInfo }}</pre>
          </div>
        </div>
        <div *ngIf="editing">
            <textarea rows="5"  name="info-edit-textarea" [(ngModel)]="imageInfo"></textarea>
        </div>
        <div class="btn-sm" *ngIf="editing">
          <button class="btn btn-primary" [disabled]="!hasChanges()" (click)="saveInfo()" >{{'BUTTON.SAVE' | translate}}</button>
          <button class="btn" (click)="cancelInfo()" >{{'BUTTON.CANCEL' | translate}}</button>
        </div>
        <confirmation-dialog #confirmationDialog (confirmAction)="confirmCancel($event)"></confirmation-dialog>
      </form>
    </section>
    <section id="image" role="tabpanel" aria-labelledby="repo-image" [hidden]='!isCurrentTabContent("image")'>
      <div id=images-container>
        <hbr-tag ngProjectAs="clr-dg-row-detail" (tagClickEvent)="watchTagClickEvt($event)" (signatureOutput)="saveSignatures($event)" class="sub-grid-custom" [repoName]="repoName" [registryUrl]="registryUrl" [withNotary]="withNotary" [withClair]="withClair" [hasSignedIn]="hasSignedIn" [hasProjectAdminRole]="hasProjectAdminRole" [projectId]="projectId"></hbr-tag>
      </div>
    </section>
  </div>
</section>
`;
import { NgModule } from '@angular/core';
import { CoreModule } from '../core/core.module';
import { CookieService } from 'angular2-cookie/core';

import { SessionService } from '../shared/session.service';
import { MessageComponent } from '../global-message/message.component';

import { MessageService } from '../global-message/message.service';
import { MaxLengthExtValidatorDirective } from './max-length-ext.directive';
import { FilterComponent } from './filter/filter.component';
import { TranslateModule } from "@ngx-translate/core";

import { RouterModule } from '@angular/router';

import { ConfirmationDialogComponent } from './confirmation-dialog/confirmation-dialog.component';
import { ConfirmationDialogService } from './confirmation-dialog/confirmation-dialog.service';
import { SystemAdminGuard } from './route/system-admin-activate.service';
import { NewUserFormComponent } from './new-user-form/new-user-form.component';
import { InlineAlertComponent } from './inline-alert/inline-alert.component';

import { ListPolicyComponent } from './list-policy/list-policy.component';
import { CreateEditPolicyComponent } from './create-edit-policy/create-edit-policy.component';

import { PortValidatorDirective } from './port.directive';

import { PageNotFoundComponent } from './not-found/not-found.component';
import { AboutDialogComponent } from './about-dialog/about-dialog.component';

import { AuthCheckGuard } from './route/auth-user-activate.service';

import { StatisticsComponent } from './statictics/statistics.component';
import { StatisticsPanelComponent } from './statictics/statistics-panel.component';
import { SignInGuard } from './route/sign-in-guard-activate.service';
import { LeavingConfigRouteDeactivate } from './route/leaving-config-deactivate.service';
import { MemberGuard } from './route/member-guard-activate.service';

@NgModule({
  imports: [
    CoreModule,
    TranslateModule,
    RouterModule
  ],
  declarations: [
    MessageComponent,
    MaxLengthExtValidatorDirective,
    FilterComponent,
    ConfirmationDialogComponent,
    NewUserFormComponent,
    InlineAlertComponent,
    ListPolicyComponent,
    CreateEditPolicyComponent,
    PortValidatorDirective,
    PageNotFoundComponent,
    AboutDialogComponent,
    StatisticsComponent,
    StatisticsPanelComponent
  ],
  exports: [
    CoreModule,
    MessageComponent,
    MaxLengthExtValidatorDirective,
    FilterComponent,
    TranslateModule,
    ConfirmationDialogComponent,
    NewUserFormComponent,
    InlineAlertComponent,
    ListPolicyComponent,
    CreateEditPolicyComponent,
    PortValidatorDirective,
    PageNotFoundComponent,
    AboutDialogComponent,
    StatisticsComponent,
    StatisticsPanelComponent
  ],
  providers: [
    SessionService,
    MessageService,
    CookieService,
    ConfirmationDialogService,
    SystemAdminGuard,
    AuthCheckGuard,
    SignInGuard,
    LeavingConfigRouteDeactivate,
    MemberGuard
  ]
})
export class SharedModule {

}
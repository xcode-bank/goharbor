import { Component, OnInit, ViewChild, OnDestroy } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Response } from '@angular/http';

import { SessionUser } from '../../shared/session-user';
import { Member } from './member';
import { MemberService } from './member.service';

import { AddMemberComponent } from './add-member/add-member.component';

import { MessageService } from '../../global-message/message.service';
import { AlertType, ConfirmationTargets, ConfirmationState } from '../../shared/shared.const';

import { ConfirmationDialogService } from '../../shared/confirmation-dialog/confirmation-dialog.service';
import { ConfirmationMessage } from '../../shared/confirmation-dialog/confirmation-message';
import { SessionService } from '../../shared/session.service';

import { RoleInfo } from '../../shared/shared.const';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import 'rxjs/add/observable/throw';
import { Subscription } from 'rxjs/Subscription';

import { Project } from '../../project/project';

@Component({
  moduleId: module.id,
  templateUrl: 'member.component.html',
  styleUrls: ['./member.component.css']
})
export class MemberComponent implements OnInit, OnDestroy {

  members: Member[];
  projectId: number;
  roleInfo = RoleInfo;
  private delSub: Subscription;

  @ViewChild(AddMemberComponent)
  addMemberComponent: AddMemberComponent;

  currentUser: SessionUser;
  hasProjectAdminRole: boolean;

  constructor(private route: ActivatedRoute, private router: Router,
    private memberService: MemberService, private messageService: MessageService,
    private deletionDialogService: ConfirmationDialogService,
    private session: SessionService) {
    
    this.delSub = deletionDialogService.confirmationConfirm$.subscribe(message => {
      if (message &&
        message.state === ConfirmationState.CONFIRMED &&
        message.source === ConfirmationTargets.PROJECT_MEMBER) {
        this.memberService
          .deleteMember(this.projectId, message.data)
          .subscribe(
          response => {
            this.messageService.announceMessage(response, 'MEMBER.DELETED_SUCCESS', AlertType.SUCCESS);
            console.log('Successful delete member: ' + message.data);
            this.retrieve(this.projectId, '');
          },
          error => this.messageService.announceMessage(error.status, 'Failed to change role with user ' + message.data, AlertType.DANGER)
          );
      }
    });
  }

  retrieve(projectId: number, username: string) {
    this.memberService
      .listMembers(projectId, username)
      .subscribe(
      response => this.members = response,
      error => {
        this.router.navigate(['/harbor', 'projects']);
        this.messageService.announceMessage(error.status, 'Failed to get project member with project ID:' + projectId, AlertType.DANGER);
      });
  }

  ngOnDestroy() {
    if (this.delSub) {
      this.delSub.unsubscribe();
    }
  }

  ngOnInit() {
    //Get projectId from route params snapshot.          
    this.projectId = +this.route.snapshot.parent.params['id'];
    console.log('Get projectId from route params snapshot:' + this.projectId);
    
    this.currentUser = this.session.getCurrentUser();
    //Get current user from registered resolver.
    let resolverData = this.route.snapshot.parent.data;
    if(resolverData) {
      this.hasProjectAdminRole = (<Project>resolverData['projectResolver']).has_project_admin_role;
    }

   

    this.retrieve(this.projectId, '');
  }

  openAddMemberModal() {
    this.addMemberComponent.openAddMemberModal();
  }

  addedMember() {
    this.retrieve(this.projectId, '');
  }

  changeRole(m: Member, roleId: number) {
    if(m) {
      this.memberService
        .changeMemberRole(this.projectId, m.user_id, roleId)
        .subscribe(
        response => {
          this.messageService.announceMessage(response, 'MEMBER.SWITCHED_SUCCESS', AlertType.SUCCESS);
          console.log('Successful change role with user ' + m.user_id + ' to roleId ' + roleId);
          this.retrieve(this.projectId, '');
        },
        error => this.messageService.announceMessage(error.status, 'Failed to change role with user ' + m.user_id + ' to roleId ' + roleId, AlertType.DANGER)
        );
      }
  }

  deleteMember(m: Member) {
    let deletionMessage: ConfirmationMessage = new ConfirmationMessage(
      'MEMBER.DELETION_TITLE',
      'MEMBER.DELETION_SUMMARY',
      m.username,
      m.user_id,
      ConfirmationTargets.PROJECT_MEMBER
    );
    this.deletionDialogService.openComfirmDialog(deletionMessage);
  }

  doSearch(searchMember) {
    this.retrieve(this.projectId, searchMember);
  }

  refresh() {
    this.retrieve(this.projectId, '');
  }
}
import { Component, OnInit, ViewChild, OnDestroy } from '@angular/core';

import { Router } from '@angular/router';

import { Project } from './project';
import { ProjectService } from './project.service';

import { CreateProjectComponent } from './create-project/create-project.component';

import { ListProjectComponent } from './list-project/list-project.component';

import { MessageService } from '../global-message/message.service';
import { Message } from '../global-message/message';

import { AlertType } from '../shared/shared.const';
import { Response } from '@angular/http';

import { ConfirmationDialogService } from '../shared/confirmation-dialog/confirmation-dialog.service';
import { ConfirmationMessage } from '../shared/confirmation-dialog/confirmation-message';
import { ConfirmationTargets, ConfirmationState } from '../shared/shared.const';

import { Subscription } from 'rxjs/Subscription';

import { State } from 'clarity-angular';

const types: {} = { 0: 'PROJECT.MY_PROJECTS', 1: 'PROJECT.PUBLIC_PROJECTS' };

@Component({
  moduleId: module.id,
  selector: 'project',
  templateUrl: 'project.component.html',
  styleUrls: ['./project.component.css']
})
export class ProjectComponent implements OnInit, OnDestroy {

  selected = [];
  changedProjects: Project[];
  projectTypes = types;

  @ViewChild(CreateProjectComponent)
  creationProject: CreateProjectComponent;

  @ViewChild(ListProjectComponent)
  listProject: ListProjectComponent;

  currentFilteredType: number = 0;

  subscription: Subscription;

  projectName: string;
  isPublic: number;

  page: number = 1;
  pageSize: number = 15;

  totalPage: number;
  totalRecordCount: number;

  constructor(
    private projectService: ProjectService,
    private messageService: MessageService,
    private deletionDialogService: ConfirmationDialogService) {
    this.subscription = deletionDialogService.confirmationConfirm$.subscribe(message => {
      if (message &&
        message.state === ConfirmationState.CONFIRMED &&
        message.source === ConfirmationTargets.PROJECT) {
        let projectId = message.data;
        this.projectService
          .deleteProject(projectId)
          .subscribe(
          response => {
            console.log('Successful delete project with ID:' + projectId);
            this.retrieve();
          },
          error => this.messageService.announceMessage(error.status, error, AlertType.WARNING)
          );
      }
    });
  }

  ngOnInit(): void {
    this.projectName = '';
    this.isPublic = 0;
  }

  ngOnDestroy(): void {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  retrieve(state?: State): void {
    if (state) {
      this.page = state.page.to + 1;
    }
    this.projectService
      .listProjects(this.projectName, this.isPublic, this.page, this.pageSize)
      .subscribe(
      response => {
        this.totalRecordCount = response.headers.get('x-total-count');
        this.totalPage = Math.ceil(this.totalRecordCount / this.pageSize);
        console.log('TotalRecordCount:' + this.totalRecordCount + ', totalPage:' + this.totalPage);
        this.changedProjects = response.json();
      },
      error => this.messageService.announceAppLevelMessage(error.status, error, AlertType.WARNING)
      );
  }

  openModal(): void {
    this.creationProject.newProject();
  }

  createProject(created: boolean) {
    if (created) {
      this.retrieve();
    }
  }

  doSearchProjects(projectName: string): void {
    console.log('Search for project name:' + projectName);
    this.projectName = projectName;
    this.retrieve();
  }

  doFilterProjects(filteredType: number): void {
    console.log('Filter projects with type:' + types[filteredType]);
    this.isPublic = filteredType;
    this.currentFilteredType = filteredType;
    this.retrieve();
  }

  toggleProject(p: Project) {
    if (p) {
      p.public === 0 ? p.public = 1 : p.public = 0;
      this.projectService
        .toggleProjectPublic(p.project_id, p.public)
        .subscribe(
        response => console.log('Successful toggled project_id:' + p.project_id),
        error => this.messageService.announceMessage(error.status, error, AlertType.WARNING)
        );
    }
  }

  deleteProject(p: Project) {
    let deletionMessage = new ConfirmationMessage(
      'PROJECT.DELETION_TITLE',
      'PROJECT.DELETION_SUMMARY',
      p.name,
      p.project_id,
      ConfirmationTargets.PROJECT
    );
    this.deletionDialogService.openComfirmDialog(deletionMessage);
  }

  refresh(): void {
    this.retrieve();
  }

}
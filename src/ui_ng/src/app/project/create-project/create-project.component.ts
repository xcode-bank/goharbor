import { Component, EventEmitter, Output, ViewChild, AfterViewChecked, HostBinding } from '@angular/core';
import { Response } from '@angular/http';
import { NgForm } from '@angular/forms';

import { Project } from '../project';
import { ProjectService } from '../project.service';


import { MessageService } from '../../global-message/message.service';
import { AlertType } from '../../shared/shared.const';

import { InlineAlertComponent } from '../../shared/inline-alert/inline-alert.component';

import { TranslateService } from '@ngx-translate/core';


@Component({
  selector: 'create-project',
  templateUrl: 'create-project.component.html',
  styleUrls: [ 'create-project.css' ]
})
export class CreateProjectComponent implements AfterViewChecked {
  
  projectForm: NgForm;

  @ViewChild('projectForm')
  currentForm: NgForm;

  project: Project = new Project();
  initVal: Project = new Project();

  createProjectOpened: boolean;
  
  hasChanged: boolean;

  staticBackdrop: boolean = true;
  closable: boolean = false;

  @Output() create = new EventEmitter<boolean>();
  @ViewChild(InlineAlertComponent)
  private inlineAlert: InlineAlertComponent;

  constructor(private projectService: ProjectService, 
              private messageService: MessageService,
              private translateService: TranslateService) {}

  onSubmit() {
    this.projectService
        .createProject(this.project.name, this.project.public ? 1 : 0)
        .subscribe(
          status=>{
            this.create.emit(true);
            this.createProjectOpened = false;
          },
          error=>{
            let errorMessage: string;
            if (error instanceof Response) { 
              switch(error.status) {
              case 409:
                this.translateService.get('PROJECT.NAME_ALREADY_EXISTS').subscribe(res=>errorMessage = res);
                break;
              case 400:
                this.translateService.get('PROJECT.NAME_IS_ILLEGAL').subscribe(res=>errorMessage = res); 
                break;
              default:
                this.translateService.get('PROJECT.UNKNOWN_ERROR').subscribe(res=>{
                  errorMessage = res;
                  this.messageService.announceMessage(error.status, errorMessage, AlertType.DANGER);
                });
              }
              this.inlineAlert.showInlineError(errorMessage);
            }
          }); 
  }

  onCancel() {
    if(this.hasChanged) {
      this.inlineAlert.showInlineConfirmation({message: 'ALERT.FORM_CHANGE_CONFIRMATION'});
    } else {
      this.createProjectOpened = false;
      this.projectForm.reset();
    }
   
  }

  ngAfterViewChecked(): void {
    this.projectForm = this.currentForm;
    if(this.projectForm) {
      this.projectForm.valueChanges.subscribe(data=>{
        for(let i in data) {
          let origin = this.initVal[i];          
          let current = data[i];
          if(current && current !== origin) {
            this.hasChanged = true;
            break;
          } else {
            this.hasChanged = false;
            this.inlineAlert.close();
          }
        }
      });
    }
  }

  newProject() {
    this.project = new Project();
    this.hasChanged = false;
    this.createProjectOpened = true;
  }

  confirmCancel(event: boolean): void {
    this.createProjectOpened = false;
    this.inlineAlert.close();
    this.projectForm.reset();
  }
}


// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
import {
  Component,
  EventEmitter,
  Output,
  ViewChild,
  AfterViewChecked,
  OnInit,
  OnDestroy
} from "@angular/core";
import { Response } from "@angular/http";
import { NgForm } from "@angular/forms";

import { Project } from "../project";
import { ProjectService } from "../project.service";

import { MessageHandlerService } from "../../shared/message-handler/message-handler.service";
import { InlineAlertComponent } from "../../shared/inline-alert/inline-alert.component";

import { TranslateService } from "@ngx-translate/core";

import { Subject } from "rxjs/Subject";
import "rxjs/add/operator/debounceTime";
import "rxjs/add/operator/distinctUntilChanged";

@Component({
  selector: "create-project",
  templateUrl: "create-project.component.html",
  styleUrls: ["create-project.css"]
})
export class CreateProjectComponent implements AfterViewChecked, OnInit, OnDestroy {

  projectForm: NgForm;

  @ViewChild("projectForm")
  currentForm: NgForm;

  project: Project = new Project();
  initVal: Project = new Project();

  createProjectOpened: boolean;

  hasChanged: boolean;
  isSubmitOnGoing = false;

  staticBackdrop = true;
  closable = false;

  isNameValid = true;
  nameTooltipText = "PROJECT.NAME_TOOLTIP";
  checkOnGoing = false;
  proNameChecker: Subject<string> = new Subject<string>();

  @Output() create = new EventEmitter<boolean>();
  @ViewChild(InlineAlertComponent)
  inlineAlert: InlineAlertComponent;

  constructor(private projectService: ProjectService,
    private translateService: TranslateService,
    private messageHandlerService: MessageHandlerService) { }

  ngOnInit(): void {
    this.proNameChecker
      .debounceTime(500)
      .distinctUntilChanged()
      .subscribe((name: string) => {
        let cont = this.currentForm.controls["create_project_name"];
        if (cont && this.hasChanged) {
          this.isNameValid = cont.valid;
          if (this.isNameValid) {
            // Check exiting from backend
            this.projectService
              .checkProjectExists(cont.value).toPromise()
              .then(() => {
                // Project existing
                this.isNameValid = false;
                this.nameTooltipText = "PROJECT.NAME_ALREADY_EXISTS";
                this.checkOnGoing = false;
              })
              .catch(error => {
                this.checkOnGoing = false;
              });
          } else {
            this.nameTooltipText = "PROJECT.NAME_TOOLTIP";
          }
        }
      });
  }

  ngOnDestroy(): void {
    this.proNameChecker.unsubscribe();
  }

  onSubmit() {
    if (this.isSubmitOnGoing) {
      return ;
    }

    this.isSubmitOnGoing = true;
    this.projectService
      .createProject(this.project.name, this.project.metadata)
      .subscribe(
      status => {
        this.isSubmitOnGoing = false;

        this.create.emit(true);
        this.messageHandlerService.showSuccess("PROJECT.CREATED_SUCCESS");
        this.createProjectOpened = false;
      },
      error => {
        this.isSubmitOnGoing = false;

        let errorMessage: string;
        if (error instanceof Response) {
          switch (error.status) {
            case 409:
              this.translateService.get("PROJECT.NAME_ALREADY_EXISTS").subscribe(res => errorMessage = res);
              break;
            case 400:
              this.translateService.get("PROJECT.NAME_IS_ILLEGAL").subscribe(res => errorMessage = res);
              break;
            default:
              this.translateService.get("PROJECT.UNKNOWN_ERROR").subscribe(res => errorMessage = res);
          }
        this.messageHandlerService.handleError(error);
        }
      });
  }

  onCancel() {
      this.createProjectOpened = false;
      this.projectForm.reset();
  }

  ngAfterViewChecked(): void {
    this.projectForm = this.currentForm;
    if (this.projectForm) {
      this.projectForm.valueChanges.subscribe(data => {
        for (let i in data) {
          let origin = this.initVal[i];
          let current = data[i];
          if (current && current !== origin) {
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

  public get isValid(): boolean {
    return this.currentForm &&
    this.currentForm.valid &&
    !this.isSubmitOnGoing &&
    this.isNameValid &&
    !this.checkOnGoing;
  }

  // Handle the form validation
  handleValidation(): void {
    let cont = this.currentForm.controls["create_project_name"];
    if (cont) {
      this.proNameChecker.next(cont.value);
    }

  }
}


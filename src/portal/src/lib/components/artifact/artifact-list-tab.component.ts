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
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnInit,
  Output,
  ViewChild,

} from "@angular/core";
import { forkJoin, Observable, Subject, throwError as observableThrowError, of } from "rxjs";
import { catchError, debounceTime, distinctUntilChanged, finalize, map } from 'rxjs/operators';
import { TranslateService } from "@ngx-translate/core";
import { Comparator, Label, State, Tag, ArtifactClickEvent, VulnerabilitySummary } from "../../services/interface";

import {
  RequestQueryParams,
  RetagService,
  ScanningResultService,
  ProjectService,
  ArtifactService
} from "../../services";
import { ErrorHandler } from "../../utils/error-handler/error-handler";
import { ConfirmationButtons, ConfirmationState, ConfirmationTargets } from "../../entities/shared.const";

import { ConfirmationDialogComponent } from "../confirmation-dialog/confirmation-dialog.component";
import { ConfirmationMessage } from "../confirmation-dialog/confirmation-message";
import { ConfirmationAcknowledgement } from "../confirmation-dialog/confirmation-state-message";

import {
  calculatePage,
  clone,
  CustomComparator,
  DEFAULT_PAGE_SIZE, DEFAULT_SUPPORTED_MIME_TYPE,
  doFiltering,
  doSorting, formatSize,
  VULNERABILITY_SCAN_STATUS,
} from "../../utils/utils";

import { CopyInputComponent } from "../push-image/copy-input.component";
import { LabelService } from "../../services/label.service";
import { UserPermissionService } from "../../services/permission.service";
import { USERSTATICPERMISSION } from "../../services/permission-static";
import { operateChanges, OperateInfo, OperationState } from "../operation/operate";
import { OperationService } from "../operation/operation.service";
import { ImageNameInputComponent } from "../image-name-input/image-name-input.component";
import { errorHandler as errorHandFn } from "../../utils/shared/shared.utils";
import { ClrLoadingState, ClrDatagridStateInterface, ClrDatagridComparatorInterface } from "@clr/angular";
import { ChannelService } from "../../services/channel.service";
import { Artifact, Reference } from "./artifact";
import { HttpParams } from "@angular/common/http";
import { ActivatedRoute } from "@angular/router";

export interface LabelState {
  iconsShow: boolean;
  label: Label;
  show: boolean;
}
export const AVAILABLE_TIME = '0001-01-01T00:00:00.000Z';
@Component({
  selector: 'artifact-list-tab',
  templateUrl: './artifact-list-tab.component.html',
  styleUrls: ['./artifact-list-tab.component.scss']
})
export class ArtifactListTabComponent implements OnInit, AfterViewInit {

  signedCon: { [key: string]: any | string[] } = {};
  @Input() projectId: number;
  projectName: string;
  @Input() memberRoleID: number;
  @Input() repoName: string;
  referArtifactArray: string[] = [];
  @Input() isEmbedded: boolean;

  @Input() hasSignedIn: boolean;
  @Input() isGuest: boolean;
  @Input() registryUrl: string;
  @Input() withNotary: boolean;
  @Input() withAdmiral: boolean;
  @Output() refreshRepo = new EventEmitter<boolean>();
  @Output() tagClickEvent = new EventEmitter<ArtifactClickEvent>();
  @Output() signatureOutput = new EventEmitter<any>();
  @Output() putReferArtifactArray = new EventEmitter<string[]>();


  tags: Tag[];
  artifactList: Artifact[] = [];
  availableTime = AVAILABLE_TIME;
  showTagManifestOpened: boolean;
  retagDialogOpened: boolean;
  manifestInfoTitle: string;
  digestId: string;
  staticBackdrop = true;
  closable = false;
  lastFilteredTagName: string;
  inprogress: boolean;
  openLabelFilterPanel: boolean;
  openLabelFilterPiece: boolean;
  retagSrcImage: string;
  showlabel: boolean;

  pullComparator: Comparator<Artifact> = new CustomComparator<Artifact>("pull_time", "date");
  pushComparator: Comparator<Artifact> = new CustomComparator<Artifact>("push_time", "date");

  loading = true;
  copyFailed = false;
  selectedRow: Artifact[] = [];

  imageLabels: LabelState[] = [];
  imageStickLabels: LabelState[] = [];
  imageFilterLabels: LabelState[] = [];

  labelListOpen = false;
  selectedTag: Artifact[];
  labelNameFilter: Subject<string> = new Subject<string>();
  stickLabelNameFilter: Subject<string> = new Subject<string>();
  filterOnGoing: boolean;
  stickName = '';
  filterName = '';
  initFilter = {
    name: '',
    description: '',
    color: '',
    scope: '',
    project_id: 0,
  };
  filterOneLabel: Label = this.initFilter;

  @ViewChild("confirmationDialog", { static: false })
  confirmationDialog: ConfirmationDialogComponent;

  @ViewChild("imageNameInput", { static: false })
  imageNameInput: ImageNameInputComponent;

  @ViewChild("digestTarget", { static: false }) textInput: ElementRef;
  @ViewChild("copyInput", { static: false }) copyInput: CopyInputComponent;

  pageSize: number = DEFAULT_PAGE_SIZE;
  currentPage = 1;
  totalCount = 0;
  currentState: State;

  hasAddLabelImagePermission: boolean;
  hasRetagImagePermission: boolean;
  hasDeleteImagePermission: boolean;
  hasScanImagePermission: boolean;
  hasEnabledScanner: boolean;
  scanBtnState: ClrLoadingState = ClrLoadingState.DEFAULT;
  onSendingScanCommand: boolean;

  constructor(
    private errorHandler: ErrorHandler,
    private retagService: RetagService,
    private userPermissionService: UserPermissionService,
    private labelService: LabelService,
    private artifactService: ArtifactService,
    private translateService: TranslateService,
    private ref: ChangeDetectorRef,
    private operationService: OperationService,
    private channel: ChannelService,
    private projectService: ProjectService,
    private activatedRoute: ActivatedRoute,
    private scanningService: ScanningResultService
  ) { }

  ngOnInit() {
    if (!this.projectId) {
      this.errorHandler.error("Project ID cannot be unset.");
      return;
    }
    this.activatedRoute.data.subscribe(res => {
      this.projectName = res.projectResolver.name;
    });

    this.getProjectScanner();
    if (!this.repoName) {
      this.errorHandler.error("Repo name cannot be unset.");
      return;
    }
    let refer = JSON.parse(sessionStorage.getItem('reference'));
    if (refer && refer.projectId === this.projectId && refer.repo === this.repoName) {
      this.referArtifactArray = refer.referArray;
    }
    this.artifactService.TriggerArtifactChan$.subscribe(res => {
      let st: ClrDatagridStateInterface = { page: {from: 0, to: this.pageSize - 1, size: this.pageSize} };
          this.clrLoad(st);
    });
    this.lastFilteredTagName = '';

    this.labelNameFilter
      .pipe(debounceTime(500))
      .pipe(distinctUntilChanged())
      .subscribe((name: string) => {
        if (this.filterName.length) {
          this.filterOnGoing = true;

          this.imageFilterLabels.forEach(data => {
            if (data.label.name.indexOf(this.filterName) !== -1) {
              data.show = true;
            } else {
              data.show = false;
            }
          });
          setTimeout(() => {
            setInterval(() => this.ref.markForCheck(), 200);
          }, 1000);
        }
      });

    this.stickLabelNameFilter
      .pipe(debounceTime(500))
      .pipe(distinctUntilChanged())
      .subscribe((name: string) => {
        if (this.stickName.length) {
          this.filterOnGoing = true;

          this.imageStickLabels.forEach(data => {
            if (data.label.name.indexOf(this.stickName) !== -1) {
              data.show = true;
            } else {
              data.show = false;
            }
          });
          setTimeout(() => {
            setInterval(() => this.ref.markForCheck(), 200);
          }, 1000);
        }
      });

    this.getImagePermissionRule(this.projectId);
  }

  ngAfterViewInit() {
  }

  public get filterLabelPieceWidth() {
    let len = this.lastFilteredTagName.length ? this.lastFilteredTagName.length * 6 + 60 : 115;
    return len > 210 ? 210 : len;
  }
  doSearchArtifactByFilter(tagName) {
    this.lastFilteredTagName = tagName;
    this.currentPage = 1;

    let st: State = this.currentState;
    if (!st) {
      st = { page: {} };
    }
    st.page.size = this.pageSize;
    st.page.from = 0;
    st.page.to = this.pageSize - 1;
    let selectedLab = this.imageFilterLabels.find(label => label.iconsShow === true);
    if (selectedLab) {
      st.filters = [{ property: 'name', value: this.lastFilteredTagName }, { property: 'labels.id', value: selectedLab.label.id }];
    } else {
      st.filters = [{ property: 'name', value: this.lastFilteredTagName }];
    }

    this.clrLoad(st);
  }
  doSearchArtifactNames(artifactName: string) {
    this.lastFilteredTagName = artifactName;
    this.currentPage = 1;

    let st: State = this.currentState;
    if (!st) {
      st = { page: {} };
    }
    st.page.size = this.pageSize;
    st.page.from = 0;
    st.page.to = this.pageSize - 1;
    let selectedLab = this.imageFilterLabels.find(label => label.iconsShow === true);
    if (selectedLab) {
      st.filters = [{ property: 'name', value: this.lastFilteredTagName }, { property: 'labels.id', value: selectedLab.label.id }];
    } else {
      st.filters = [{ property: 'name', value: this.lastFilteredTagName }];
    }

    this.clrLoad(st);
  }

  clrLoad(state: ClrDatagridStateInterface): void {
    if (!state || !state.page) {
      return;
    }
    this.selectedRow = [];
    // Keep it for future filtering and sorting

    let pageNumber: number = calculatePage(state);
    if (pageNumber <= 0) { pageNumber = 1; }
    let sortBy: any = '';
    if (state.sort) {
      sortBy = state.sort.by as string | ClrDatagridComparatorInterface<any>;
      sortBy = sortBy.fieldName ? sortBy.fieldName : sortBy;
      sortBy = state.sort.reverse ? `-${sortBy}` : sortBy;
    }
    this.loading = true;
    this.currentState = state;

    // Pagination
    let params = new HttpParams();
    if (pageNumber && this.pageSize) {
      params = params.set('page', pageNumber + '').set('page_size', this.pageSize + '');
    }
    if (sortBy) {
      params = params.set('sort', sortBy);
    }
    if (state.filters && state.filters.length) {
      state.filters.forEach(item => {
        params = params.set(item.property, item.value);
      });
    }
    this.projectService.getProject(this.projectId).subscribe(project => {
      this.projectName = project.name;
      let refer = JSON.parse(sessionStorage.getItem('reference'));
      this.referArtifactArray = [];
    if (refer && refer.projectId === this.projectId && refer.repo === this.repoName) {
      this.referArtifactArray = refer.referArray;
    }

      if (this.referArtifactArray.length) {
        let observableLists: Observable<Artifact>[] = [];

        this.artifactService.getArtifactFromDigest(this.projectName, this.repoName,
          this.referArtifactArray[this.referArtifactArray.length - 1]).subscribe(artifact => {
            this.totalCount = artifact.references.length;
              artifact.references.forEach((child, index) => {
                if (index >= (pageNumber - 1) * this.pageSize && index < pageNumber * this.pageSize) {
                  observableLists.push(this.artifactService.getArtifactFromDigest(this.projectName, this.repoName,
                    child.child_digest));
                }
              });

            forkJoin(observableLists).pipe(finalize(() => {
              this.loading = false;
            })).subscribe(artifacts => {
              this.artifactList = artifacts;
            }, error => {
              this.errorHandler.error(error);
            });
          });
      } else {
        this.artifactService.getArtifactList(this.projectName, this.repoName, params).subscribe(res => {
          if (res.headers) {
            let xHeader: string = res.headers.get("X-Total-Count");
            if (xHeader) {
              this.totalCount = parseInt(xHeader, 0);
            }
          }
          this.artifactList = res.body;
          this.loading = false;
        }, error => {
          // error
          this.loading = false;
        });
      }
    });
  }

  refresh() {
    this.doSearchArtifactNames("");
  }

  getAllLabels(): void {
    forkJoin(this.labelService.getGLabels(), this.labelService.getPLabels(this.projectId)).subscribe(results => {
      results.forEach(labels => {
        labels.forEach(data => {
          this.imageLabels.push({ 'iconsShow': false, 'label': data, 'show': true });
        });
      });
      this.imageFilterLabels = clone(this.imageLabels);
      this.imageStickLabels = clone(this.imageLabels);
    }, error => this.errorHandler.error(error));
  }

  labelSelectedChange(artifact?: Artifact[]): void {
    if (artifact && artifact[0].labels) {
      this.imageStickLabels.forEach(data => {
        data.iconsShow = false;
        data.show = true;
      });
      if (artifact[0].labels.length) {
        artifact[0].labels.forEach((labelInfo: Label) => {
          let findedLabel = this.imageStickLabels.find(data => labelInfo.id === data['label'].id);
          this.imageStickLabels.splice(this.imageStickLabels.indexOf(findedLabel), 1);
          this.imageStickLabels.unshift(findedLabel);

          findedLabel.iconsShow = true;
        });
      }
    }
  }

  addLabels(): void {
    this.labelListOpen = true;
    this.selectedTag = this.selectedRow;
    this.stickName = '';
    this.labelSelectedChange(this.selectedRow);
  }

  stickLabel(labelInfo: LabelState): void {
    if (labelInfo && !labelInfo.iconsShow) {
      this.selectLabel(labelInfo);
    }
    if (labelInfo && labelInfo.iconsShow) {
      this.unSelectLabel(labelInfo);
    }
  }

  selectLabel(labelInfo: LabelState): void {
    if (!this.inprogress) {
      this.inprogress = true;
      let labelId = labelInfo.label.id;
      this.selectedRow = this.selectedTag;

      this.artifactService.addLabelToImages(this.projectName, this.repoName, this.selectedRow[0].digest, labelId).subscribe(res => {
        this.refresh();

        // set the selected label in front
        this.imageStickLabels.splice(this.imageStickLabels.indexOf(labelInfo), 1);
        this.imageStickLabels.some((data, i) => {
          if (!data.iconsShow) {
            this.imageStickLabels.splice(i, 0, labelInfo);
            return true;
          }
        });

        // when is the last one
        if (this.imageStickLabels.every(data => data.iconsShow === true)) {
          this.imageStickLabels.push(labelInfo);
        }

        labelInfo.iconsShow = true;
        this.inprogress = false;
      }, err => {
        this.inprogress = false;
        this.errorHandler.error(err);
      });
    }
  }

  unSelectLabel(labelInfo: LabelState): void {
    if (!this.inprogress) {
      this.inprogress = true;
      let labelId = labelInfo.label.id;
      this.selectedRow = this.selectedTag;
      this.artifactService.deleteLabelToImages(this.projectName, this.repoName, this.selectedRow[0].digest, labelId).subscribe(res => {
        this.refresh();

        // insert the unselected label to groups with the same icons
        this.sortOperation(this.imageStickLabels, labelInfo);
        labelInfo.iconsShow = false;
        this.inprogress = false;
      }, err => {
        this.inprogress = false;
        this.errorHandler.error(err);
      });
    }
  }

  rightFilterLabel(labelInfo: LabelState): void {
    if (labelInfo) {
      if (!labelInfo.iconsShow) {
        this.filterLabel(labelInfo);
        this.showlabel = true;
      } else {
        this.unFilterLabel(labelInfo);
        this.showlabel = false;
      }
    }
  }

  filterLabel(labelInfo: LabelState): void {
    let labelId = labelInfo.label.id;
    // insert the unselected label to groups with the same icons
    let preLabelInfo = this.imageFilterLabels.find(data => data.label.id === this.filterOneLabel.id);
    if (preLabelInfo) {
      this.sortOperation(this.imageFilterLabels, preLabelInfo);
    }

    this.imageFilterLabels.filter(data => {
      if (data.label.id !== labelId) {
        data.iconsShow = false;
      } else {
        data.iconsShow = true;
      }
    });
    this.imageFilterLabels.splice(this.imageFilterLabels.indexOf(labelInfo), 1);
    this.imageFilterLabels.unshift(labelInfo);
    this.filterOneLabel = labelInfo.label;

    // reload data
    this.currentPage = 1;
    let st: State = this.currentState;
    if (!st) {
      st = { page: {} };
    }
    st.page.size = this.pageSize;
    st.page.from = 0;
    st.page.to = this.pageSize - 1;
    if (this.lastFilteredTagName) {
      st.filters = [{ property: 'name', value: this.lastFilteredTagName }, { property: 'labels.id', value: labelId }];
    } else {
      st.filters = [{ property: 'labels.id', value: labelId }];
    }

    this.clrLoad(st);
  }

  unFilterLabel(labelInfo: LabelState): void {
    // insert the unselected label to groups with the same icons
    this.sortOperation(this.imageFilterLabels, labelInfo);

    this.filterOneLabel = this.initFilter;
    labelInfo.iconsShow = false;

    // reload data
    this.currentPage = 1;
    let st: State = this.currentState;
    if (!st) {
      st = { page: {} };
    }
    st.page.size = this.pageSize;
    st.page.from = 0;
    st.page.to = this.pageSize - 1;
    if (this.lastFilteredTagName) {
      st.filters = [{ property: 'name', value: this.lastFilteredTagName }];
    } else {
      st.filters = [];
    }
    this.clrLoad(st);
  }

  closeFilter(): void {
    this.openLabelFilterPanel = false;
  }

  openFlagEvent(isOpen: boolean): void {
    if (isOpen) {
      this.openLabelFilterPanel = true;
      this.openLabelFilterPiece = true;
      this.filterName = '';
      // redisplay all labels
      this.imageFilterLabels.forEach(data => {
        if (data.label.name.indexOf(this.filterName) !== -1) {
          data.show = true;
        } else {
          data.show = false;
        }
      });
    } else {
      this.openLabelFilterPanel = false;
      this.openLabelFilterPiece = false;
    }

  }

  handleInputFilter() {
    if (this.filterName.length) {
      this.labelNameFilter.next(this.filterName);
    } else {
      this.imageFilterLabels.every(data => data.show = true);
    }
  }

  handleStickInputFilter() {
    if (this.stickName.length) {
      this.stickLabelNameFilter.next(this.stickName);
    } else {
      this.imageStickLabels.every(data => data.show = true);
    }
  }

  // insert the unselected label to groups with the same icons
  sortOperation(labelList: LabelState[], labelInfo: LabelState): void {
    labelList.some((data, i) => {
      if (!data.iconsShow) {
        if (data.label.scope === labelInfo.label.scope) {
          labelList.splice(i, 0, labelInfo);
          labelList.splice(labelList.indexOf(labelInfo, 0), 1);
          return true;
        }
        if (data.label.scope !== labelInfo.label.scope && i === labelList.length - 1) {
          labelList.push(labelInfo);
          labelList.splice(labelList.indexOf(labelInfo), 1);
          return true;
        }
      }
    });
  }

  loadArtifactList(params?) {
    this.artifactList = [];
    this.loading = true;
    let refer = JSON.parse(sessionStorage.getItem('reference'));
    if (refer && refer.projectId === this.projectId && refer.repo === this.repoName) {
      this.referArtifactArray = refer.referArray;
    }
      if (this.referArtifactArray.length) {
        let observableLists: Observable<Artifact>[] = [];

        this.artifactService.getArtifactFromDigest(this.projectName, this.repoName,
          this.referArtifactArray[this.referArtifactArray.length - 1]).subscribe(artifact => {
            artifact.references.forEach(child => {
              observableLists.push(this.artifactService.getArtifactFromDigest(this.projectName, this.repoName,
                child.child_digest));
            });
            forkJoin(observableLists).subscribe(artifacts => {
              this.loading = false;

              this.artifactList = artifacts;
            });
          });
      } else {
        this.artifactService.getArtifactList(this.projectName, this.repoName, params).subscribe(res => {
          if (res.headers) {
            let xHeader: string = res.headers.get("X-Total-Count");
            if (xHeader) {
              this.totalCount = parseInt(xHeader, 0);
            }
          }
          this.artifactList = res.body;
          this.loading = false;
        }, error => {
          // error
          this.loading = false;
        });
      }


  }

  sizeTransform(tagSize: string): string {
    return formatSize(tagSize);
  }

  retag() {
    if (this.selectedRow && this.selectedRow.length) {
      this.retagDialogOpened = true;
      this.retagSrcImage = this.repoName + ":" + this.selectedRow[0].digest;
    } else {
      this.errorHandler.error("One tag should be selected before retag.");
    }
  }

  onRetag() {
    this.retagService.retag({
      targetProject: this.imageNameInput.projectName.value,
      targetRepo: this.imageNameInput.repoName.value,
      targetTag: this.imageNameInput.tagName.value,
      srcImage: this.retagSrcImage,
      override: true
    })
      .pipe(finalize(() => {
        this.retagDialogOpened = false;
        this.imageNameInput.form.reset();
      }))
      .subscribe(response => {
        this.translateService.get('RETAG.MSG_SUCCESS').subscribe((res: string) => {
          this.errorHandler.info(res);
          if (`${this.imageNameInput.projectName.value}/${this.imageNameInput.repoName.value}` === this.repoName) {
            let st: State = this.currentState;
            if (!st) {
              st = { page: {} };
            }
            st.page.size = this.pageSize;
            st.page.from = 0;
            st.page.to = this.pageSize - 1;
            this.clrLoad(st);
          }
        });
      }, error => {
        this.errorHandler.error(error);
      });
  }

  deleteArtifact() {
    if (this.selectedRow && this.selectedRow.length) {
      let artifactNames: string[] = [];
      this.selectedRow.forEach(artifact => {
        artifactNames.push(artifact.digest.slice(0, 15));
      });

      let titleKey: string, summaryKey: string, content: string, buttons: ConfirmationButtons;
      titleKey = "REPOSITORY.DELETION_TITLE_TAG";
      summaryKey = "REPOSITORY.DELETION_SUMMARY_TAG";
      buttons = ConfirmationButtons.DELETE_CANCEL;
      content = artifactNames.join(" , ");
      let message = new ConfirmationMessage(
        titleKey,
        summaryKey,
        content,
        this.selectedRow,
        ConfirmationTargets.TAG,
        buttons);
      this.confirmationDialog.open(message);
    }
  }
  deleteArtifactobservableLists: Observable<any>[] = [];
  confirmDeletion(message: ConfirmationAcknowledgement) {
    if (message &&
      message.source === ConfirmationTargets.TAG
      && message.state === ConfirmationState.CONFIRMED) {
      let artifactList = message.data;
      if (artifactList && artifactList.length) {
        this.findArtifactFromIndex(artifactList);
      }
    }
  }
  findArtifactFromIndex(artifactList: Artifact[]) {
    if (artifactList.every(artifact1 => !artifact1.references)) {
      artifactList.forEach(artifact => {
        this.deleteArtifactobservableLists.push(this.delOperate(artifact));
      });
      forkJoin(...this.deleteArtifactobservableLists).subscribe((items) => {
        // if delete one success  refresh list
        if (items.some(item => !item)) {
          this.selectedRow = [];
          let st: ClrDatagridStateInterface = { page: {from: 0, to: this.pageSize - 1, size: this.pageSize} };
          this.clrLoad(st);
        }
      });
    } else {
      let observArr: Observable<Artifact>[] = [];
      artifactList.forEach(artifact => {
        this.deleteArtifactobservableLists.push(this.delOperate(artifact));
        if (artifact.references) {
          artifact.references.forEach(reference => {
            observArr.push(this.artifactService.getArtifactFromDigest(this.projectName, this.repoName, reference.child_digest));
          });

        }
      });
      forkJoin(observArr).subscribe((res) => {
        this.findArtifactFromIndex(res);
      });
    }
  }

  delOperate(artifact: Artifact): Observable<any> | null {
    // init operation info
    let operMessage = new OperateInfo();
    operMessage.name = 'OPERATION.DELETE_TAG';
    operMessage.data.id = artifact.id;
    operMessage.state = OperationState.progressing;
    operMessage.data.name = artifact.digest;
    this.operationService.publishInfo(operMessage);
    // to do signature
    // if (tag.signature) {
    //   forkJoin(this.translateService.get("BATCH.DELETED_FAILURE"),
    //     this.translateService.get("REPOSITORY.DELETION_SUMMARY_TAG_DENIED")).subscribe(res => {
    //       let wrongInfo: string = res[1] + "notary -s https://" + this.registryUrl +
    //         ":4443 -d ~/.docker/trust remove -p " +
    //         this.registryUrl + "/" + this.repoName +
    //         " " + name;
    //       operateChanges(operMessage, OperationState.failure, wrongInfo);
    //     });
    // } else {
    return this.artifactService
      .deleteArtifact(this.projectName, this.repoName, artifact.digest)
      .pipe(map(
        response => {
          this.translateService.get("BATCH.DELETED_SUCCESS")
            .subscribe(res => {
              operateChanges(operMessage, OperationState.success);
            });
        }), catchError(error => {
          const message = errorHandFn(error);
          this.translateService.get(message).subscribe(res =>
            operateChanges(operMessage, OperationState.failure, res)
          );
          return of(error);
        }));
    // }
  }

  showDigestId() {
    if (this.selectedRow && (this.selectedRow.length === 1)) {
      this.manifestInfoTitle = "REPOSITORY.COPY_DIGEST_ID";
      this.digestId = this.selectedRow[0].digest;
      this.showTagManifestOpened = true;
      this.copyFailed = false;
    }
  }

  onTagClick(artifact: Artifact): void {
    if (artifact) {
      let evt: ArtifactClickEvent = {
        project_id: this.projectId,
        repository_name: this.repoName,
        digest: artifact.digest,
        artifact_id: artifact.id,
      };
      this.tagClickEvent.emit(evt);
    }
  }

  onSuccess($event: any): void {
    this.copyFailed = false;
    // Directly close dialog
    this.showTagManifestOpened = false;
  }

  onError($event: any): void {
    // Show error
    this.copyFailed = true;
    // Select all text
    if (this.textInput) {
      this.textInput.nativeElement.select();
    }
  }

  // Get vulnerability scanning status
  scanStatus(artifact: Artifact): string {
    if (artifact) {
      let so = this.handleScanOverview(artifact.scan_overview);
      if (so && so.scan_status) {
        return so.scan_status;
      }
    }
    return VULNERABILITY_SCAN_STATUS.NOT_SCANNED;
  }
  // Whether show the 'scan now' menu
  canScanNow(): boolean {
    if (!this.hasScanImagePermission) { return false; }
    if (this.onSendingScanCommand) { return false; }
    let st: string = this.scanStatus(this.selectedRow[0]);
    return st !== VULNERABILITY_SCAN_STATUS.RUNNING;
  }
  getImagePermissionRule(projectId: number): void {
    const permissions = [
      { resource: USERSTATICPERMISSION.REPOSITORY_TAG_LABEL.KEY, action: USERSTATICPERMISSION.REPOSITORY_TAG_LABEL.VALUE.CREATE },
      { resource: USERSTATICPERMISSION.REPOSITORY.KEY, action: USERSTATICPERMISSION.REPOSITORY.VALUE.PULL },
      { resource: USERSTATICPERMISSION.REPOSITORY_TAG.KEY, action: USERSTATICPERMISSION.REPOSITORY_TAG.VALUE.DELETE },
      { resource: USERSTATICPERMISSION.REPOSITORY_TAG_SCAN_JOB.KEY, action: USERSTATICPERMISSION.REPOSITORY_TAG_SCAN_JOB.VALUE.CREATE },
    ];
    this.userPermissionService.hasProjectPermissions(this.projectId, permissions).subscribe((results: Array<boolean>) => {
      this.hasAddLabelImagePermission = results[0];
      this.hasRetagImagePermission = results[1];
      this.hasDeleteImagePermission = results[2];
      this.hasScanImagePermission = results[3];
      // only has label permission
      if (this.hasAddLabelImagePermission) {
        if (!this.withAdmiral) {
          this.getAllLabels();
        }
      }
    }, error => this.errorHandler.error(error));
  }
  // Trigger scan
  scanNow(): void {
    if (this.selectedRow && this.selectedRow.length === 1) {
      this.onSendingScanCommand = true;
      this.channel.publishScanEvent(this.repoName + "/" + this.selectedRow[0].digest);
    }
  }
  submitFinish(e: boolean) {
    this.onSendingScanCommand = e;
  }
  // pull command
  onCpError($event: any): void {
    this.copyInput.setPullCommendShow();
  }
  getProjectScanner(): void {
    this.hasEnabledScanner = false;
    this.scanBtnState = ClrLoadingState.LOADING;
    this.scanningService.getProjectScanner(this.projectId)
      .subscribe(response => {
        if (response && "{}" !== JSON.stringify(response) && !response.disabled
          && response.health === "healthy") {
          this.scanBtnState = ClrLoadingState.SUCCESS;
          this.hasEnabledScanner = true;
        } else {
          this.scanBtnState = ClrLoadingState.ERROR;
        }
      }, error => {
        this.scanBtnState = ClrLoadingState.ERROR;
      });
  }

  handleScanOverview(scanOverview: any): VulnerabilitySummary {
    if (scanOverview) {
      return scanOverview[DEFAULT_SUPPORTED_MIME_TYPE];
    }
    return null;
  }

  refer(artifact: Artifact) {
    this.referArtifactArray.push(artifact.digest);
    sessionStorage.setItem('reference', JSON.stringify({ projectId: this.projectId, repo: this.repoName
      , referArray: this.referArtifactArray}));

    if (this.referArtifactArray.length) {
      this.putReferArtifactArray.emit(this.referArtifactArray);
    }
    let st: ClrDatagridStateInterface = this.currentState;
    if (!st) {
      st = { page: {} };
    }
    st.page.size = this.pageSize;
    st.page.from = 0;
    st.page.to = this.pageSize - 1;
    this.clrLoad(st);
  }
}

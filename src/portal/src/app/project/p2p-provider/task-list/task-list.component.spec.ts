import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { CUSTOM_ELEMENTS_SCHEMA, NO_ERRORS_SCHEMA } from '@angular/core';
import { BrowserAnimationsModule, NoopAnimationsModule } from '@angular/platform-browser/animations';
import { ClarityModule } from '@clr/angular';
import { FormsModule } from '@angular/forms';
import { RouterTestingModule } from '@angular/router/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { P2pProviderService } from '../p2p-provider.service';
import { PreheatService } from '../../../../../ng-swagger-gen/services/preheat.service';
import { ActivatedRoute } from '@angular/router';
import { SessionService } from '../../../shared/session.service';
import { of } from 'rxjs';
import { delay } from 'rxjs/operators';
import { Execution } from '../../../../../ng-swagger-gen/models/execution';
import { TaskListComponent } from './task-list.component';
import { MessageHandlerService } from '../../../shared/message-handler/message-handler.service';
import { UserPermissionService } from '../../../../lib/services';
import { Task } from '../../../../../ng-swagger-gen/models/task';
import { HttpHeaders, HttpResponse } from '@angular/common/http';
describe('TaskListComponent', () => {
  let component: TaskListComponent;
  let fixture: ComponentFixture<TaskListComponent>;
  const execution: Execution = {
    id: 1,
    vendor_id: 1,
    status: 'Success',
    trigger: 'Manual',
    start_time: new Date().toUTCString(),
  };
  const task: Task = {
    id: 1,
    execution_id: 1,
    status: 'Success',
    status_message: 'no artifact to preheat',
    start_time: new Date().toUTCString(),
  };
  const mockPreheatService = {
    GetExecution() {
       return of(execution).pipe(delay(0));
    },
    ListTasksResponse() {
      return of(new HttpResponse({
        body: [task],
        headers:  new HttpHeaders({
          "X-Total-Count": "1"
        })
      })).pipe(delay(0));
    }
  };
  const mockActivatedRoute = {
    snapshot: {
      params: {
        executionId: 1,
        preheatPolicyName: 'policy1'
      },
      parent: {
        parent: {
          params: { id: 1 },
          data: {
            projectResolver: {
              name: 'library',
              metadata: {
                prevent_vul: 'true',
                enable_content_trust: 'true',
                severity: 'none'
              }
            }
          }
        }
      }
    }
  };
  const mockedSessionService = {
    getCurrentUser() {
      return {
        has_admin_role: true
      };
    }
  };
  const mockMessageHandlerService = {
    handleError: () => { }
  };
  const mockUserPermissionService = {
    getPermission() {
      return of(true).pipe(delay(0));
    }
  };
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      schemas: [
        CUSTOM_ELEMENTS_SCHEMA,
        NO_ERRORS_SCHEMA
      ],
      imports: [
        BrowserAnimationsModule,
        ClarityModule,
        TranslateModule.forRoot(),
        FormsModule,
        RouterTestingModule,
        NoopAnimationsModule,
        HttpClientTestingModule
      ],
      declarations: [TaskListComponent],
      providers: [
        P2pProviderService,
        TranslateService,
        { provide: PreheatService, useValue: mockPreheatService },
        { provide: ActivatedRoute, useValue: mockActivatedRoute },
        { provide: SessionService, useValue: mockedSessionService },
        { provide: MessageHandlerService, useValue: mockMessageHandlerService },
        { provide: UserPermissionService, useValue: mockUserPermissionService },
      ]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TaskListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
  it('should render task list', async () => {
    fixture.autoDetectChanges();
    await fixture.whenStable();
    const rows = fixture.nativeElement.getElementsByTagName('clr-dg-row');
    expect(rows.length).toEqual(1);
  });
});

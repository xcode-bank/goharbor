import { ComponentFixture, TestBed, async, } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';
import { RouterTestingModule } from '@angular/router/testing';

import { SharedModule } from '../shared/shared.module';
import { ConfirmationDialogComponent } from '../confirmation-dialog/confirmation-dialog.component';
import { RepositoryComponent } from './repository.component';
import { RepositoryListviewComponent } from '../repository-listview/repository-listview.component';
import { FilterComponent } from '../filter/filter.component';
import { TagComponent } from '../tag/tag.component';
import { VULNERABILITY_DIRECTIVES } from '../vulnerability-scanning/index';
import { PUSH_IMAGE_BUTTON_DIRECTIVES } from '../push-image/index';
import { INLINE_ALERT_DIRECTIVES } from '../inline-alert/index';
import { JobLogViewerComponent } from '../job-log-viewer/index';


import { ErrorHandler } from '../error-handler/error-handler';
import { Repository, RepositoryItem, Tag, SystemInfo } from '../service/interface';
import { SERVICE_CONFIG, IServiceConfig } from '../service.config';
import { RepositoryService, RepositoryDefaultService } from '../service/repository.service';
import { SystemInfoService, SystemInfoDefaultService } from '../service/system-info.service';
import { TagService, TagDefaultService } from '../service/tag.service';
import { ChannelService } from '../channel/index';

class RouterStub {
  navigateByUrl(url: string) { return url; }
}

describe('RepositoryComponent (inline template)', () => {

  let compRepo: RepositoryComponent;
  let fixture: ComponentFixture<RepositoryComponent>;
  let repositoryService: RepositoryService;
  let systemInfoService: SystemInfoService;
  let tagService: TagService;

  let spyRepos: jasmine.Spy;
  let spyTags: jasmine.Spy;
  let spySystemInfo: jasmine.Spy;

  let mockSystemInfo: SystemInfo = {
    'with_notary': true,
    'with_admiral': false,
    'admiral_endpoint': 'NA',
    'auth_mode': 'db_auth',
    'registry_url': '10.112.122.56',
    'project_creation_restriction': 'everyone',
    'self_registration': true,
    'has_ca_root': false,
    'harbor_version': 'v1.1.1-rc1-160-g565110d'
  };

  let mockRepoData: RepositoryItem[] = [
    {
      'id': 1,
      'name': 'library/busybox',
      'project_id': 1,
      'description': 'asdfsadf',
      'pull_count': 0,
      'star_count': 0,
      'tags_count': 1
    },
    {
      'id': 2,
      'name': 'library/nginx',
      'project_id': 1,
      'description': 'asdf',
      'pull_count': 0,
      'star_count': 0,
      'tags_count': 1
    }
  ];

  let mockRepo: Repository = {
    metadata: {xTotalCount: 2},
    data: mockRepoData
  };

  let mockTagData: Tag[] = [
    {
      'digest': 'sha256:e5c82328a509aeb7c18c1d7fb36633dc638fcf433f651bdcda59c1cc04d3ee55',
      'name': '1.11.5',
      'size': '2049',
      'architecture': 'amd64',
      'os': 'linux',
      'docker_version': '1.12.3',
      'author': 'NGINX Docker Maintainers \"docker-maint@nginx.com\"',
      'created': new Date('2016-11-08T22:41:15.912313785Z'),
      'signature': null
    }
  ];

  let config: IServiceConfig = {
    repositoryBaseEndpoint: '/api/repository/testing',
    systemInfoEndpoint: '/api/systeminfo/testing',
    targetBaseEndpoint: '/api/tag/testing'
  };

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        SharedModule,
        RouterTestingModule
      ],
      declarations: [
        RepositoryComponent,
        RepositoryListviewComponent,
        ConfirmationDialogComponent,
        FilterComponent,
        TagComponent,
        VULNERABILITY_DIRECTIVES,
        PUSH_IMAGE_BUTTON_DIRECTIVES,
        INLINE_ALERT_DIRECTIVES,
        JobLogViewerComponent,
      ],
      providers: [
        ErrorHandler,
        { provide: SERVICE_CONFIG, useValue: config },
        { provide: RepositoryService, useClass: RepositoryDefaultService },
        { provide: SystemInfoService, useClass: SystemInfoDefaultService },
        { provide: TagService, useClass: TagDefaultService },
        { provide: ChannelService},

      ]
    });
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RepositoryComponent);
    compRepo = fixture.componentInstance;
    compRepo.projectId = 1;
    compRepo.hasProjectAdminRole = true;
    compRepo.repoName = 'library/nginx';
    repositoryService = fixture.debugElement.injector.get(RepositoryService);
    systemInfoService = fixture.debugElement.injector.get(SystemInfoService);
    tagService = fixture.debugElement.injector.get(TagService);

    spyRepos = spyOn(repositoryService, 'getRepositories').and.returnValues(Promise.resolve(mockRepo));
    spySystemInfo = spyOn(systemInfoService, 'getSystemInfo').and.returnValues(Promise.resolve(mockSystemInfo));
    spyTags = spyOn(tagService, 'getTags').and.returnValues(Promise.resolve(mockTagData));
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(compRepo).toBeTruthy();
  });

  it('should load and render data', async(() => {
    fixture.detectChanges();
    fixture.whenStable().then(() => {
      fixture.detectChanges();
      let de: DebugElement = fixture.debugElement.query(By.css('datagrid-cell'));
      fixture.detectChanges();
      expect(de).toBeTruthy();
      let el: HTMLElement = de.nativeElement;
      expect(el).toBeTruthy();
      expect(el.textContent).toEqual('library/busybox');
    });
  }));
});

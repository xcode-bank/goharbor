
import { ComponentFixture, TestBed, async } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';
import { Router } from '@angular/router';

import { SharedModule } from '../shared/shared.module';
import { ConfirmationDialogComponent } from '../confirmation-dialog/confirmation-dialog.component';
import { ListRepositoryComponent } from './list-repository.component';
import { Repository, RepositoryItem } from '../service/interface';

import { SERVICE_CONFIG, IServiceConfig } from '../service.config';

class RouterStub {
  navigateByUrl(url: string) { return url; }
}

describe('ListRepositoryComponent (inline template)', () => {
  let comp: ListRepositoryComponent;
  let fixture: ComponentFixture<ListRepositoryComponent>;

  let mockData: RepositoryItem[] = [
    {
      "id": 11,
      "name": "library/busybox",
      "project_id": 1,
      "description": "",
      "pull_count": 0,
      "star_count": 0,
      "tags_count": 1
    },
    {
      "id": 12,
      "name": "library/nginx",
      "project_id": 1,
      "description": "",
      "pull_count": 0,
      "star_count": 0,
      "tags_count": 1
    }
  ];

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        SharedModule
      ],
      declarations: [
        ListRepositoryComponent,
        ConfirmationDialogComponent
      ],
      providers: [
        { provide: Router, useClass: RouterStub },
        { provide: SERVICE_CONFIG, useValue: {} }
      ]
    });
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ListRepositoryComponent);
    comp = fixture.componentInstance;
  });

  it('should load and render data', async(() => {
    fixture.detectChanges();
    comp.repositories = mockData;
    fixture.whenStable().then(() => {
      fixture.detectChanges();
      expect(comp.repositories).toBeTruthy();
      let de: DebugElement = fixture.debugElement.query(By.css('datagrid-cell'));
      fixture.detectChanges();
      expect(de).toBeTruthy();
      let el: HTMLElement = de.nativeElement;
      expect(el).toBeTruthy();
      expect(el.textContent).toEqual('library/busybox');
    });
  }));
});
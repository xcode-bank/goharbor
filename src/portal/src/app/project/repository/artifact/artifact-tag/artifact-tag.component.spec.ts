import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ArtifactTagComponent } from './artifact-tag.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { of } from 'rxjs';
import { IServiceConfig, SERVICE_CONFIG } from "../../../../../lib/entities/service.config";
import { SharedModule } from "../../../../../lib/utils/shared/shared.module";
import { ErrorHandler } from "../../../../../lib/utils/error-handler";
import { TagService } from "../../../../../lib/services";
import { OperationService } from "../../../../../lib/components/operation/operation.service";


describe('ArtifactTagComponent', () => {
  let component: ArtifactTagComponent;
  let fixture: ComponentFixture<ArtifactTagComponent>;
  const mockErrorHandler = {
    error: () => {}
  };
  const mockTagService = {
    newTag: () => of([]),
    deleteTag: () => of(null),
  };
  const config: IServiceConfig = {
    repositoryBaseEndpoint: "/api/repositories/testing"
  };
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        SharedModule,
        BrowserAnimationsModule,
        HttpClientTestingModule
      ],
      schemas: [
        CUSTOM_ELEMENTS_SCHEMA
      ],
      declarations: [ ArtifactTagComponent ],
      providers: [
        ErrorHandler,
        { provide: SERVICE_CONFIG, useValue: config },
        { provide: mockErrorHandler, useValue: ErrorHandler },
        { provide: TagService, useValue: mockTagService },
        { provide: OperationService },
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ArtifactTagComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

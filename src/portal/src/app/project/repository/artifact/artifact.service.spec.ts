import { TestBed, inject } from '@angular/core/testing';
import { IServiceConfig, SERVICE_CONFIG } from "../../../../lib/entities/service.config";
import { SharedModule } from "../../../../lib/utils/shared/shared.module";
import { ArtifactDefaultService, ArtifactService } from "../artifact/artifact.service";
import { CURRENT_BASE_HREF } from "../../../../lib/utils/utils";
import { IconService } from '../../../../../ng-swagger-gen/services/icon.service';
import { DomSanitizer } from '@angular/platform-browser';

describe('ArtifactService', () => {

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [
        SharedModule
      ],
      providers: [
        ArtifactDefaultService,
        {
          provide: ArtifactService,
          useClass: ArtifactDefaultService
        },
        IconService,
        DomSanitizer
      ]
    });
  });

  it('should be initialized', inject([ArtifactDefaultService], (service: ArtifactService) => {
    expect(service).toBeTruthy();
  }));

});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { GlobalSearchService } from './global-search.service';
import { SearchResults } from './search-results';
import { SearchTriggerService } from './search-trigger.service';
import { CUSTOM_ELEMENTS_SCHEMA } from '@angular/core';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AppConfigService } from './../../app-config.service';
import { ListProjectROComponent } from '../../shared/list-project-ro/list-project-ro.component';
import { MessageHandlerService } from '../../shared/message-handler/message-handler.service';
import { SearchResultComponent } from './search-result.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { of } from 'rxjs';
import { AppConfig } from '../../app-config';
import { ClarityModule } from '@clr/angular';

describe('SearchResultComponent', () => {
    let component: SearchResultComponent;
    let fixture: ComponentFixture<SearchResultComponent>;
    let fakeSearchResults = null;
    const project = {
        project_id: 1,
        owner_id: 1,
        name: 'library',
        creation_time: Date,
        creation_time_str: 'string',
        deleted: 1,
        owner_name: 'string',
        togglable: true,
        update_time: Date,
        current_user_role_id: 1,
        repo_count: 1,
        chart_count: 1,
        has_project_admin_role: true,
        is_member: true,
        role_name: 'string',
        metadata: {
            public: 'string',
            enable_content_trust: 'string',
            prevent_vul: 'string',
            severity: 'string',
            auto_scan: 'string',
            retention_id: 1
        }
    };
    let fakeGlobalSearchService = {
        doSearch: () => of({
            project: [project],
            repository: [],
            chart: []
        })
    };
    let fakeAppConfigService = {
        getConfig: () => new AppConfig()
    };
    let searchResult = "";
    let fakeMessageHandlerService = null;
    let fakeSearchTriggerService = {
        searchTriggerChan$: of(searchResult),
        searchCloseChan$: of(null),
        clear: () => null
    };
    beforeEach(async(() => {
        TestBed.configureTestingModule({
            imports: [
                TranslateModule.forRoot(),
                HttpClientTestingModule,
                ClarityModule
            ],
            declarations: [SearchResultComponent, ListProjectROComponent],
            schemas: [CUSTOM_ELEMENTS_SCHEMA]
        }).overrideComponent(SearchResultComponent, {
            set: {
                providers: [
                    TranslateService,
                    { provide: AppConfigService, useValue: fakeAppConfigService },
                    { provide: MessageHandlerService, useValue: fakeMessageHandlerService },
                    { provide: SearchTriggerService, useValue: fakeSearchTriggerService },
                    { provide: GlobalSearchService, useValue: fakeGlobalSearchService },
                    { provide: SearchResults, useValue: fakeSearchResults }
                ]
            }
        }).compileComponents();
    }));

    beforeEach(() => {
        fixture = TestBed.createComponent(SearchResultComponent);
        component = fixture.componentInstance;
        component.stateIndicator = true;
        fixture.autoDetectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
    it('should search library', async () => {
        searchResult = 'library';
        component.doSearch(searchResult);
        await fixture.whenStable();
        expect(component.searchResults.project[0].name).toEqual('library');
    });
});

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
import { Component, ReflectiveInjector, LOCALE_ID } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { CookieService } from 'ngx-cookie';

import { supportedLangs, enLang } from './shared/shared.const';
import { SessionService } from './shared/session.service';
import { AppConfigService } from './app-config.service';
import { Title } from '@angular/platform-browser';

@Component({
    selector: 'harbor-app',
    templateUrl: 'app.component.html'
})
export class AppComponent {
    constructor(
        private translate: TranslateService,
        private cookie: CookieService,
        private session: SessionService,
        private appConfigService: AppConfigService,
        private titleService: Title) {
        //Override page title
        let key: string = "APP_TITLE.HARBOR";
        if (this.appConfigService.isIntegrationMode()) {
            key = "APP_TITLE.REG";
        }

        translate.get(key).subscribe((res: string) => {
            this.titleService.setTitle(res);
        });
    }
}

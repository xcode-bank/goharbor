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
import { Component, OnInit, OnDestroy } from "@angular/core";
import { Subscription } from "rxjs";

import { StatisticsService } from "./statistics.service";
import { Statistics } from "./statistics";

import { SessionService } from "../../../../shared/services/session.service";
import { Volumes } from "./volumes";

import { MessageHandlerService } from "../../../../shared/services/message-handler.service";
import { StatisticHandler } from "./statistic-handler.service";
import { AppConfigService } from "../../../../services/app-config.service";


@Component({
    selector: "statistics-panel",
    templateUrl: "statistics-panel.component.html",
    styleUrls: ["statistics.component.scss"],
    providers: [StatisticsService]
})

export class StatisticsPanelComponent implements OnInit, OnDestroy {

    originalCopy: Statistics = new Statistics();
    volumesInfo: Volumes = new Volumes();
    refreshSub: Subscription;
    constructor(
        private statistics: StatisticsService,
        private msgHandler: MessageHandlerService,
        private session: SessionService,
        private appConfigService: AppConfigService,
        private statisticHandler: StatisticHandler) {
    }

    ngOnInit(): void {
        // Refresh
        this.refreshSub = this.statisticHandler.refreshChan$.subscribe(clear => {
            this.getStatistics();
        });

        if (this.session.getCurrentUser()) {
            this.getStatistics();
        }

        if (this.isValidSession) {
            this.getVolumes();
        }
    }

    ngOnDestroy() {
        if (this.refreshSub) {
            this.refreshSub.unsubscribe();
        }
    }

    public get totalStorage(): number {
        let count: number = 0;
        if (this.volumesInfo && this.volumesInfo.storage && this.volumesInfo.storage.length) {
            this.volumesInfo.storage.forEach(item => {
                count += item.total ? item.total : 0;
            });
        }
        return this.getGBFromBytes(count);
    }

    public get freeStorage(): number {
        let count: number = 0;
        if (this.volumesInfo && this.volumesInfo.storage && this.volumesInfo.storage.length) {
            this.volumesInfo.storage.forEach(item => {
                count += item.free ? item.free : 0;
            });
        }
        return this.getGBFromBytes(count);
    }

    public getStatistics(): void {
        this.statistics.getStatistics()
            .subscribe(statistics => this.originalCopy = statistics
                , error => {
                    this.msgHandler.handleError(error);
                });
    }

    public getVolumes(): void {
        this.statistics.getVolumes()
            .subscribe(volumes => this.volumesInfo = volumes
                , error => {
                    this.msgHandler.handleError(error);
                });
    }

    public get isValidSession(): boolean {
        let user = this.session.getCurrentUser();
        return user && user.has_admin_role;
    }

    public get isValidStorage(): boolean {
        let count: number = 0;
        if (this.volumesInfo && this.volumesInfo.storage && this.volumesInfo.storage.length) {
            this.volumesInfo.storage.forEach(item => {
                count += item.total ? item.total : 0;
            });
        }
        return count !== 0 &&
            this.appConfigService.getConfig().registry_storage_provider_name === "filesystem";
    }

    getGBFromBytes(bytes: number): number {
        return Math.round((bytes / (1024 * 1024 * 1024)));
    }
}

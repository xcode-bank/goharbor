import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';

import { StatisticsService } from './statistics.service';
import { Statistics } from './statistics';

import { SessionService } from '../session.service';
import { Volumes } from './volumes';

import { MessageHandlerService } from '../message-handler/message-handler.service';
import { StatisticHandler } from './statistic-handler.service';

@Component({
    selector: 'statistics-panel',
    templateUrl: "statistics-panel.component.html",
    styleUrls: ['statistics.component.css'],
    providers: [StatisticsService]
})

export class StatisticsPanelComponent implements OnInit, OnDestroy {

    private originalCopy: Statistics = new Statistics();
    private volumesInfo: Volumes = new Volumes();
    refreshSub: Subscription;

    constructor(
        private statistics: StatisticsService,
        private msgHandler: MessageHandlerService,
        private session: SessionService,
        private statisticHandler: StatisticHandler) {
    }

    ngOnInit(): void {
        //Refresh
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
        return this.getGBFromBytes(this.volumesInfo.storage.total);
    }

    public get freeStorage(): number {
        return this.getGBFromBytes(this.volumesInfo.storage.free);
    }

    public getStatistics(): void {
        this.statistics.getStatistics()
            .then(statistics => this.originalCopy = statistics)
            .catch(error => {
                this.msgHandler.handleError(error);
            });
    }

    public getVolumes(): void {
        this.statistics.getVolumes()
            .then(volumes => this.volumesInfo = volumes)
            .catch(error => {
                this.msgHandler.handleError(error);
            });
    }

    public get isValidSession(): boolean {
        let user = this.session.getCurrentUser();
        return user && user.has_admin_role > 0;
    }

    public get isValidStorage(): boolean {
        return this.volumesInfo.storage.total != 0;
    }

    private getGBFromBytes(bytes: number): number {
        return Math.round((bytes / (1024 * 1024 * 1024)));
    }
}
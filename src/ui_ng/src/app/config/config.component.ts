import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { NgForm } from '@angular/forms';

import { ConfigurationService } from './config.service';
import { Configuration } from './config';
import { MessageService } from '../global-message/message.service';
import { AlertType, ConfirmationTargets, ConfirmationState } from '../shared/shared.const';
import { errorHandler, accessErrorHandler } from '../shared/shared.utils';
import { StringValueItem } from './config';
import { ConfirmationDialogService } from '../shared/confirmation-dialog/confirmation-dialog.service';
import { Subscription } from 'rxjs/Subscription';
import { ConfirmationMessage } from '../shared/confirmation-dialog/confirmation-message'

import { ConfigurationAuthComponent } from './auth/config-auth.component';
import { ConfigurationEmailComponent } from './email/config-email.component';

import { AppConfigService } from '../app-config.service';
import { SessionService } from '../shared/session.service';

const fakePass = "fakepassword";
const TabLinkContentMap = {
    "config-auth": "authentication",
    "config-replication": "replication",
    "config-email": "email",
    "config-system": "system_settings"
};

@Component({
    selector: 'config',
    templateUrl: "config.component.html",
    styleUrls: ['config.component.css']
})
export class ConfigurationComponent implements OnInit, OnDestroy {
    private onGoing: boolean = false;
    allConfig: Configuration = new Configuration();
    private currentTabId: string = "config-auth";//default tab
    private originalCopy: Configuration;
    private confirmSub: Subscription;
    private testingOnGoing: boolean = false;

    @ViewChild("repoConfigFrom") repoConfigForm: NgForm;
    @ViewChild("systemConfigFrom") systemConfigForm: NgForm;
    @ViewChild(ConfigurationEmailComponent) mailConfig: ConfigurationEmailComponent;
    @ViewChild(ConfigurationAuthComponent) authConfig: ConfigurationAuthComponent;

    constructor(
        private msgService: MessageService,
        private configService: ConfigurationService,
        private confirmService: ConfirmationDialogService,
        private appConfigService: AppConfigService,
        private session: SessionService) { }

    private isCurrentTabLink(tabId: string): boolean {
        return this.currentTabId === tabId;
    }

    private isCurrentTabContent(contentId: string): boolean {
        return TabLinkContentMap[this.currentTabId] === contentId;
    }

    private hasUnsavedChangesOfCurrentTab(): any {
        let allChanges = this.getChanges();
        if (this.isEmpty(allChanges)) {
            return null;
        }

        let properties = [];
        switch (this.currentTabId) {
            case "config-auth":
                for (let prop in allChanges) {
                    if (prop.startsWith("ldap_")) {
                        return allChanges;
                    }
                }
                properties = ["auth_mode", "project_creation_restriction", "self_registration"];
                break;
            case "config-email":
                for (let prop in allChanges) {
                    if (prop.startsWith("email_")) {
                        return allChanges;
                    }
                }
                return null;
            case "config-replication":
                properties = ["verify_remote_cert"];
                break;
            case "config-system":
                properties = ["token_expiration"];
                break;
            default:
                return null;
        }

        for (let prop in allChanges) {
            if (properties.indexOf(prop) != -1) {
                return allChanges;
            }
        }

        return null;
    }

    ngOnInit(): void {
        //First load
        //Double confirm the current use has admin role
        let currentUser = this.session.getCurrentUser();
        if (currentUser && currentUser.has_admin_role > 0) {
            this.retrieveConfig();
        }

        this.confirmSub = this.confirmService.confirmationConfirm$.subscribe(confirmation => {
            if (confirmation &&
                confirmation.state === ConfirmationState.CONFIRMED) {
                if (confirmation.source === ConfirmationTargets.CONFIG) {
                    this.reset(confirmation.data);
                } else if (confirmation.source === ConfirmationTargets.CONFIG_TAB) {
                    this.reset(confirmation.data["changes"]);
                    this.currentTabId = confirmation.data["tabId"];
                }
            }
        });
    }

    ngOnDestroy(): void {
        if (this.confirmSub) {
            this.confirmSub.unsubscribe();
        }
    }

    public get inProgress(): boolean {
        return this.onGoing;
    }

    public get testingInProgress(): boolean {
        return this.testingOnGoing;
    }

    public isValid(): boolean {
        return this.repoConfigForm &&
            this.repoConfigForm.valid &&
            this.systemConfigForm &&
            this.systemConfigForm.valid &&
            this.mailConfig &&
            this.mailConfig.isValid() &&
            this.authConfig &&
            this.authConfig.isValid();
    }

    public hasChanges(): boolean {
        return !this.isEmpty(this.getChanges());
    }

    public isMailConfigValid(): boolean {
        return this.mailConfig &&
            this.mailConfig.isValid();
    }

    public get showTestServerBtn(): boolean {
        return this.currentTabId === 'config-email';
    }

    public get showLdapServerBtn(): boolean {
        return this.currentTabId === 'config-auth' &&
            this.allConfig.auth_mode &&
            this.allConfig.auth_mode.value === "ldap_auth";
    }

    public isLDAPConfigValid(): boolean {
        return this.authConfig && this.authConfig.isValid();
    }

    public tabLinkClick(tabLink: string) {
        //Whether has unsave changes in current tab
        let changes = this.hasUnsavedChangesOfCurrentTab();
        if (!changes) {
            this.currentTabId = tabLink;
            return;
        }

        this.confirmUnsavedTabChanges(changes, tabLink);
    }

    /**
     * 
     * Save the changed values
     * 
     * @memberOf ConfigurationComponent
     */
    public save(): void {
        let changes = this.getChanges();
        if (!this.isEmpty(changes)) {
            this.onGoing = true;
            this.configService.saveConfiguration(changes)
                .then(response => {
                    this.onGoing = false;
                    //API should return the updated configurations here
                    //Unfortunately API does not do that
                    //To refresh the view, we can clone the original data copy
                    //or force refresh by calling service.
                    //HERE we choose force way
                    this.retrieveConfig();

                    //Reload bootstrap option
                    this.appConfigService.load().catch(error => console.error("Failed to reload bootstrap option with error: ", error));

                    this.msgService.announceMessage(response.status, "CONFIG.SAVE_SUCCESS", AlertType.SUCCESS);
                })
                .catch(error => {
                    this.onGoing = false;
                    if (!accessErrorHandler(error, this.msgService)) {
                        this.msgService.announceMessage(error.status, errorHandler(error), AlertType.DANGER);
                    }
                });
        } else {
            //Inprop situation, should not come here
            console.error("Save obort becasue nothing changed");
        }
    }

    /**
     * 
     * Discard current changes if have and reset
     * 
     * @memberOf ConfigurationComponent
     */
    public cancel(): void {
        let changes = this.getChanges();
        if (!this.isEmpty(changes)) {
            this.confirmUnsavedChanges(changes);
        } else {
            //Inprop situation, should not come here
            console.error("Nothing changed");
        }
    }

    /**
     * 
     * Test the connection of specified mail server
     * 
     * 
     * @memberOf ConfigurationComponent
     */
    public testMailServer(): void {
        let mailSettings = {};
        for (let prop in this.allConfig) {
            if (prop.startsWith("email_")) {
                mailSettings[prop] = this.allConfig[prop].value;
            }
        }
        //Confirm port is number
        mailSettings["email_port"] = +mailSettings["email_port"];
        let allChanges = this.getChanges();
        let password = allChanges["email_password"]
        if (password) {
            mailSettings["email_password"] = password;
        } else {
            delete mailSettings["email_password"];
        }

        this.testingOnGoing = true;
        this.configService.testMailServer(mailSettings)
            .then(response => {
                this.testingOnGoing = false;
                this.msgService.announceMessage(200, "CONFIG.TEST_MAIL_SUCCESS", AlertType.SUCCESS);
            })
            .catch(error => {
                this.testingOnGoing = false;
                this.msgService.announceMessage(error.status, errorHandler(error), AlertType.WARNING);
            });
    }

    public testLDAPServer(): void {
        let ldapSettings = {};
        for (let prop in this.allConfig) {
            if (prop.startsWith("ldap_")) {
                ldapSettings[prop] = this.allConfig[prop].value;
            }
        }

        let allChanges = this.getChanges();
        for(let prop in allChanges){
            if (prop.startsWith("ldap_")) {
                ldapSettings[prop] = allChanges[prop];
            }
        }

        console.info(ldapSettings);

        this.testingOnGoing = true;
        this.configService.testLDAPServer(ldapSettings)
            .then(respone => {
                this.testingOnGoing = false;
                this.msgService.announceMessage(200, "CONFIG.TEST_LDAP_SUCCESS", AlertType.SUCCESS);
            })
            .catch(error => {
                this.testingOnGoing = false;
                this.msgService.announceMessage(error.status, errorHandler(error), AlertType.WARNING);
            });
    }

    private confirmUnsavedChanges(changes: any) {
        let msg = new ConfirmationMessage(
            "CONFIG.CONFIRM_TITLE",
            "CONFIG.CONFIRM_SUMMARY",
            "",
            changes,
            ConfirmationTargets.CONFIG
        );

        this.confirmService.openComfirmDialog(msg);
    }

    private confirmUnsavedTabChanges(changes: any, tabId: string) {
        let msg = new ConfirmationMessage(
            "CONFIG.CONFIRM_TITLE",
            "CONFIG.CONFIRM_SUMMARY",
            "",
            {
                "changes": changes,
                "tabId": tabId
            },
            ConfirmationTargets.CONFIG_TAB
        );

        this.confirmService.openComfirmDialog(msg);
    }

    private retrieveConfig(): void {
        this.onGoing = true;
        this.configService.getConfiguration()
            .then(configurations => {
                this.onGoing = false;

                //Add two password fields
                configurations.email_password = new StringValueItem(fakePass, true);
                configurations.ldap_search_password = new StringValueItem(fakePass, true);
                this.allConfig = configurations;

                //Keep the original copy of the data
                this.originalCopy = this.clone(configurations);
            })
            .catch(error => {
                this.onGoing = false;
                if (!accessErrorHandler(error, this.msgService)) {
                    this.msgService.announceMessage(error.status, errorHandler(error), AlertType.DANGER);
                }
            });
    }

    /**
     * 
     * Get the changed fields and return a map
     * 
     * @private
     * @returns {*}
     * 
     * @memberOf ConfigurationComponent
     */
    private getChanges(): any {
        let changes = {};
        if (!this.allConfig || !this.originalCopy) {
            return changes;
        }

        for (let prop in this.allConfig) {
            let field = this.originalCopy[prop];
            if (field && field.editable) {
                if (field.value != this.allConfig[prop].value) {
                    changes[prop] = this.allConfig[prop].value;
                    //Fix boolean issue
                    if (typeof field.value === "boolean") {
                        changes[prop] = changes[prop] ? "1" : "0";
                    }
                }
            }
        }

        return changes;
    }

    /**
     * 
     * Deep clone the configuration object
     * 
     * @private
     * @param {Configuration} src
     * @returns {Configuration}
     * 
     * @memberOf ConfigurationComponent
     */
    private clone(src: Configuration): Configuration {
        let dest = new Configuration();
        if (!src) {
            return dest;//Empty
        }

        for (let prop in src) {
            if (src[prop]) {
                dest[prop] = Object.assign({}, src[prop]); //Deep copy inner object
            }
        }

        return dest;
    }

    /**
     * 
     * Reset the configuration form
     * 
     * @private
     * @param {*} changes
     * 
     * @memberOf ConfigurationComponent
     */
    private reset(changes: any): void {
        if (!this.isEmpty(changes)) {
            for (let prop in changes) {
                if (this.originalCopy[prop]) {
                    this.allConfig[prop] = Object.assign({}, this.originalCopy[prop]);
                }
            }
        } else {
            //force reset
            this.retrieveConfig();
        }
    }

    private isEmpty(obj) {
        for (let key in obj) {
            if (obj.hasOwnProperty(key))
                return false;
        }
        return true;
    }

    private disabled(prop: any): boolean {
        return !(prop && prop.editable);
    }
}
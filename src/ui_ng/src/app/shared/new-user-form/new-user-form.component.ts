import {
    Component,
    ViewChild,
    AfterViewChecked,
    Output,
    EventEmitter,
    Input,
    OnInit
} from '@angular/core';
import { NgForm } from '@angular/forms';

import { User } from '../../user/user';
import { isEmptyForm } from '../../shared/shared.utils';
import { SessionService } from '../../shared/session.service';

@Component({
    selector: 'new-user-form',
    templateUrl: 'new-user-form.component.html',
    styleUrls: ['new-user-form.component.css', '../../common.css']
})

export class NewUserFormComponent implements AfterViewChecked, OnInit {
    newUser: User = new User();
    confirmedPwd: string = "";
    @Input() isSelfRegistration: boolean = false;

    newUserFormRef: NgForm;
    @ViewChild("newUserFrom") newUserForm: NgForm;

    //Notify the form value changes
    @Output() valueChange = new EventEmitter<boolean>();

    constructor(private session: SessionService) { }

    ngOnInit() {
        this.formValueChanged = false;
    }

    private validationStateMap: any = {
        "username": true,
        "email": true,
        "realname": true,
        "newPassword": true,
        "confirmPassword": true,
        "comment": true
    };

    private mailAlreadyChecked: any = {};
    private userNameAlreadyChecked: any = {};
    private emailTooltip: string = 'TOOLTIP.EMAIL';
    private usernameTooltip: string = 'TOOLTIP.USER_NAME';
    private formValueChanged: boolean = false;

    private checkOnGoing: any = {
        "username": false,
        "email": false
    };

    public isChecking(key: string): boolean {
        return !this.checkOnGoing[key];
    }

    private getValidationState(key: string): boolean {
        return !this.validationStateMap[key];
    }

    private handleValidation(key: string, flag: boolean): void {
        if (flag) {
            //Checking
            let cont = this.newUserForm.controls[key];
            if (cont) {
                this.validationStateMap[key] = cont.valid;
                //Check email existing from backend
                if (cont.valid && this.formValueChanged) {
                    //Check username from backend
                    if (key === "username" && this.newUser.username.trim() != "") {
                        if (this.userNameAlreadyChecked[this.newUser.username.trim()]) {
                            this.validationStateMap[key] = !this.userNameAlreadyChecked[this.newUser.username.trim()].result;
                            if (!this.validationStateMap[key]) {
                                this.usernameTooltip = "TOOLTIP.USER_EXISTING";
                            }
                            return;
                        }

                        this.checkOnGoing[key] = true;
                        this.session.checkUserExisting("username", this.newUser.username)
                            .then((res: boolean) => {
                                this.checkOnGoing[key] = false;
                                this.validationStateMap[key] = !res;
                                if (res) {
                                    this.usernameTooltip = "TOOLTIP.USER_EXISTING";
                                }
                                this.userNameAlreadyChecked[this.newUser.username.trim()] = {
                                    result: res
                                }; //Tag it checked
                            })
                            .catch(error => {
                                this.checkOnGoing[key] = false;
                                this.validationStateMap[key] = false;//Not valid @ backend
                            });
                        return;

                    }

                    //Check email from backend
                    if (key === "email" && this.newUser.email.trim() != "") {
                        if (this.mailAlreadyChecked[this.newUser.email.trim()]) {
                            this.validationStateMap[key] = !this.mailAlreadyChecked[this.newUser.email.trim()].result;
                            if (!this.validationStateMap[key]) {
                                this.emailTooltip = "TOOLTIP.EMAIL_EXISTING";
                            }
                            return;
                        }

                        //Mail changed
                        this.checkOnGoing[key] = true;
                        this.session.checkUserExisting("email", this.newUser.email)
                            .then((res: boolean) => {
                                this.checkOnGoing[key] = false;
                                this.validationStateMap[key] = !res;
                                if (res) {
                                    this.emailTooltip = "TOOLTIP.EMAIL_EXISTING";
                                }
                                this.mailAlreadyChecked[this.newUser.email.trim()] = {
                                    result: res
                                }; //Tag it checked
                            })
                            .catch(error => {
                                this.checkOnGoing[key] = false;
                                this.validationStateMap[key] = false;//Not valid @ backend
                            });
                        return;
                    }

                    //Check password confirmation
                    if (key === "confirmPassword") {
                        let peerCont = this.newUserForm.controls["newPassword"];
                        if (peerCont) {
                            this.validationStateMap[key] = cont.value === peerCont.value;
                        }
                    }
                }
            }
        } else {
            //Reset
            this.validationStateMap[key] = true;
            if (key === "email") {
                this.emailTooltip = "TOOLTIP.EMAIL";
            }

            if (key === "username") {
                this.usernameTooltip = "TOOLTIP.USER_NAME";
            }
        }
    }

    public get isValid(): boolean {
        let pwdEqualStatus = true;
        if (this.newUserForm.controls["confirmPassword"] &&
            this.newUserForm.controls["newPassword"]) {
            pwdEqualStatus = this.newUserForm.controls["confirmPassword"].value === this.newUserForm.controls["newPassword"].value;
        }
        return this.newUserForm &&
            this.newUserForm.valid && 
            pwdEqualStatus &&
            this.validationStateMap["username"] &&
            this.validationStateMap["email"];//Backend check should be valid as well
    }

    ngAfterViewChecked(): void {
        if (this.newUserFormRef != this.newUserForm) {
            this.newUserFormRef = this.newUserForm;
            if (this.newUserFormRef) {
                this.newUserFormRef.valueChanges.subscribe(data => {
                    this.formValueChanged = true;
                    this.valueChange.emit(true);
                });
            }
        }
    }

    //Return the current user data
    getData(): User {
        return this.newUser;
    }

    //Reset form
    reset(): void {
        if (this.newUserForm) {
            this.newUserForm.reset();
        }
    }

    //To check if form is empty
    isEmpty(): boolean {
        return isEmptyForm(this.newUserForm);
    }
}
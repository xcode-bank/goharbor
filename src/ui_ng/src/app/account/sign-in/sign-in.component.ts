import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { Input, ViewChild, AfterViewChecked } from '@angular/core';
import { NgForm } from '@angular/forms';

import { SessionService } from '../../shared/session.service';
import { SignInCredential } from '../../shared/sign-in-credential';

import { SignUpComponent } from '../sign-up/sign-up.component';
import { CommonRoutes } from '../../shared/shared.const';
import { ForgotPasswordComponent } from '../password/forgot-password.component';

import { AppConfigService } from '../../app-config.service';
import { AppConfig } from '../../app-config';
import { User } from '../../user/user';

import { CookieService, CookieOptions } from 'angular2-cookie/core';

//Define status flags for signing in states
export const signInStatusNormal = 0;
export const signInStatusOnGoing = 1;
export const signInStatusError = -1;
const remCookieKey = "rem-username";
const expireDays = 10;

@Component({
    selector: 'sign-in',
    templateUrl: "sign-in.component.html",
    styleUrls: ['sign-in.component.css']
})

export class SignInComponent implements AfterViewChecked, OnInit {
    private redirectUrl: string = "";
    private appConfig: AppConfig = new AppConfig();
    //Remeber me indicator
    private rememberMe: boolean = false;
    private rememberedName: string = "";
    //Form reference
    signInForm: NgForm;
    @ViewChild('signInForm') currentForm: NgForm;
    @ViewChild('signupDialog') signUpDialog: SignUpComponent;
    @ViewChild('forgotPwdDialog') forgotPwdDialog: ForgotPasswordComponent;

    //Status flag
    signInStatus: number = signInStatusNormal;

    //Initialize sign in credential
    @Input() signInCredential: SignInCredential = {
        principal: "",
        password: ""
    };

    constructor(
        private router: Router,
        private session: SessionService,
        private route: ActivatedRoute,
        private appConfigService: AppConfigService,
        private cookie: CookieService
    ) { }

    ngOnInit(): void {
        //Make sure the updated configuration can be loaded
        this.appConfigService.load()
            .then(updatedConfig => this.appConfig = updatedConfig);
        this.route.queryParams
            .subscribe(params => {
                this.redirectUrl = params["redirect_url"] || "";
                let isSignUp = params["sign_up"] || "";
                if (isSignUp != "") {
                    this.signUp();//Open sign up
                }
            });

        let remUsername = this.cookie.get(remCookieKey);
        remUsername = remUsername ? remUsername.trim() : "";
        if (remUsername) {
            this.signInCredential.principal = remUsername;
            this.rememberMe = true;
            this.rememberedName = remUsername;
        }
    }

    //App title
    public get appTitle(): string {
        if(this.appConfig && this.appConfig.with_admiral){
            return "APP_TITLE.VIC";
        }

        return "APP_TITLE.VMW_HARBOR";
    }

    //For template accessing
    public get isError(): boolean {
        return this.signInStatus === signInStatusError;
    }

    public get isOnGoing(): boolean {
        return this.signInStatus === signInStatusOnGoing;
    }

    //Validate the related fields
    public get isValid(): boolean {
        return this.currentForm.form.valid;
    }

    //Whether show the 'sign up' link
    public get selfSignUp(): boolean {
        return this.appConfig.auth_mode === 'db_auth'
            && this.appConfig.self_registration;
    }

    private clickRememberMe($event): void {
        if ($event && $event.target) {
            this.rememberMe = $event.target.checked;
            if (!this.rememberMe) {
                //Remove cookie data
                this.cookie.remove(remCookieKey);
                this.rememberedName = "";
            }
        }
    }

    private remeberMe(): void {
        if (this.rememberMe) {
            if (this.rememberedName != this.signInCredential.principal) {
                //Set expire time
                let expires: number = expireDays * 3600 * 24 * 1000;
                let date = new Date(Date.now() + expires); 
                let cookieptions = new CookieOptions({
                    expires: date
                });
                this.cookie.put(remCookieKey, this.signInCredential.principal, cookieptions);
            }
        }
    }

    //General error handler
    private handleError(error) {
        //Set error status
        this.signInStatus = signInStatusError;

        let message = error.status ? error.status + ":" + error.statusText : error;
        console.error("An error occurred when signing in:", message);
    }

    //Hande form values changes
    private formChanged() {
        if (this.currentForm === this.signInForm) {
            return;
        }
        this.signInForm = this.currentForm;
        if (this.signInForm) {
            this.signInForm.valueChanges
                .subscribe(data => {
                    this.updateState();
                });
        }

    }

    //Fill the new user info into the sign in form
    private handleUserCreation(user: User): void {
        if (user) {
            this.currentForm.setValue({
                "login_username": user.username
            });

        }
    }

    //Implement interface
    //Watch the view change only when view is in error state
    ngAfterViewChecked() {
        if (this.signInStatus === signInStatusError) {
            this.formChanged();
        }
    }

    //Update the status if we have done some changes
    updateState(): void {
        if (this.signInStatus === signInStatusError) {
            this.signInStatus = signInStatusNormal; //reset
        }
    }

    //Trigger the signin action
    signIn(): void {
        //Should validate input firstly
        if (!this.isValid || this.isOnGoing) {
            return;
        }

        //Start signing in progress
        this.signInStatus = signInStatusOnGoing;

        //Call the service to send out the http request
        this.session.signIn(this.signInCredential)
            .then(() => {
                //Set status
                this.signInStatus = signInStatusNormal;

                //Remeber me
                this.remeberMe();

                //Redirect to the right route
                if (this.redirectUrl === "") {
                    //Routing to the default location
                    this.router.navigateByUrl(CommonRoutes.HARBOR_DEFAULT);
                } else {
                    this.router.navigateByUrl(this.redirectUrl);
                }
            })
            .catch(error => {
                this.handleError(error);
            });
    }

    //Open sign up dialog
    signUp(): void {
        this.signUpDialog.open();
    }

    //Open forgot password dialog
    forgotPassword(): void {
        this.forgotPwdDialog.open();
    }
}
import { Injectable } from '@angular/core';
import { Headers, Http, RequestOptions } from '@angular/http';
import 'rxjs/add/operator/toPromise';

import { PasswordSetting } from './password-setting';

const passwordChangeEndpoint = "/api/users/:user_id/password";
const sendEmailEndpoint = "/sendEmail";
const resetPasswordEndpoint = "/reset";

@Injectable()
export class PasswordSettingService {
    private headers: Headers = new Headers({
        "Accept": 'application/json',
        "Content-Type": 'application/json'
    });
    private options: RequestOptions = new RequestOptions({
        'headers': this.headers
    });

    constructor(private http: Http) { }

    changePassword(userId: number, setting: PasswordSetting): Promise<any> {
        if (!setting || setting.new_password.trim() === "" || setting.old_password.trim() === "") {
            return Promise.reject("Invalid data");
        }

        let putUrl = passwordChangeEndpoint.replace(":user_id", userId + "");
        return this.http.put(putUrl, JSON.stringify(setting), this.options)
            .toPromise()
            .then(() => null)
            .catch(error => {
                return Promise.reject(error);
            });
    }

    sendResetPasswordMail(email: string): Promise<any> {
        if (!email) {
            return Promise.reject("Invalid email");
        }

        let getUrl = sendEmailEndpoint + "?email=" + email;
        return this.http.get(getUrl, this.options).toPromise()
            .then(response => response)
            .catch(error => {
                return Promise.reject(error);
            })
    }

    resetPassword(uuid: string, newPassword: string): Promise<any> {
        if (!uuid || !newPassword) {
            return Promise.reject("Invalid reset uuid or password");
        }

        return this.http.post(resetPasswordEndpoint, JSON.stringify({
            "reset_uuid": uuid,
            "password": newPassword
        }), this.options)
            .toPromise()
            .then(response => response)
            .catch(error => {
                return Promise.reject(error);
            });
    }

}

import { Injectable } from '@angular/core';
import {
  CanActivate, Router,
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
  CanActivateChild,
  NavigationExtras
} from '@angular/router';
import { SessionService } from '../../shared/session.service';
import { CommonRoutes, AdmiralQueryParamKey } from '../../shared/shared.const';
import { AppConfigService } from '../../app-config.service';

@Injectable()
export class AuthCheckGuard implements CanActivate, CanActivateChild {
  constructor(
    private authService: SessionService,
    private router: Router,
    private appConfigService: AppConfigService) { }

  private isGuest(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
    let queryParams = route.queryParams;
    if (queryParams) {
      if (queryParams["guest"]) {
        let proRegExp = /\/harbor\/projects\/[\d]+\/.+/i;
        const libRegExp = /\/harbor\/tags\/[\d]+\/.+/i;
        if (proRegExp.test(state.url)|| libRegExp.test(state.url)) {
          return true;
        }
      }
    }

    return false;
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Promise<boolean> | boolean {
    return new Promise((resolve, reject) => {
      //Before activating, we firstly need to confirm whether the route is coming from peer part - admiral
      let queryParams = route.queryParams;
      if (queryParams) {
        if (queryParams[AdmiralQueryParamKey]) {
          console.debug(queryParams[AdmiralQueryParamKey]);
        }
      }

      let user = this.authService.getCurrentUser();
      if (!user) {
        this.authService.retrieveUser()
          .then(() => resolve(true))
          .catch(error => {
            //If is guest, skip it
            if (this.isGuest(route, state)) {
              return resolve(true);
            }
            //Session retrieving failed then redirect to sign-in
            //no matter what status code is.
            //Please pay attention that route 'HARBOR_ROOT' and 'EMBEDDED_SIGN_IN' support anonymous user
            if (state.url != CommonRoutes.HARBOR_ROOT && !state.url.startsWith(CommonRoutes.EMBEDDED_SIGN_IN)) {
              let navigatorExtra: NavigationExtras = {
                queryParams: { "redirect_url": state.url }
              };
              this.router.navigate([CommonRoutes.EMBEDDED_SIGN_IN], navigatorExtra);
              return resolve(false);
            } else {
              return resolve(true);
            }
          });
      } else {
        return resolve(true);
      }
    });
  }

  canActivateChild(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Promise<boolean> | boolean {
    return this.canActivate(route, state);
  }
}

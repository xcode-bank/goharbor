import { AfterViewInit, Component, ElementRef, OnInit } from '@angular/core';
import { Http } from '@angular/http';
import { throwError as observableThrowError, Observable } from 'rxjs';
import { catchError, map } from 'rxjs/operators';


const SwaggerUI = require('swagger-ui');
@Component({
  selector: 'dev-center',
  templateUrl: 'dev-center.component.html',
  styleUrls: ['dev-center.component.scss']
})
export class DevCenterComponent implements AfterViewInit {
  private ui: any;
  private host: any;
  private json: any;
  constructor(private el: ElementRef, private http: Http) {
  }

  ngAfterViewInit() {
    this.http.get("/swagger.json")
    .pipe(catchError(error => observableThrowError(error)))
    .pipe(map(response => response.json())).subscribe(json => {
      json.host = window.location.host;
      const protocal = window.location.protocol;
      json.schemes = [protocal.replace(":", "")];
      let ui = SwaggerUI({
        spec: json,
        domNode: this.el.nativeElement.querySelector('.swagger-container'),
        deepLinking: true,
        presets: [
          SwaggerUI.presets.apis
        ],
      });
    });
  }
}

import { Component, Input, Output, EventEmitter, ViewChild } from '@angular/core';
import { NgModel } from '@angular/forms';

import { DATETIME_PICKER_TEMPLATE } from './datetime-picker.component.html';

@Component({
  selector: 'hbr-datetime',
  template: DATETIME_PICKER_TEMPLATE
})
export class DatePickerComponent {
  
  @Input() dateInput: string;
  @Input() oneDayOffset: boolean;

  @ViewChild('searchTime') 
  searchTime: NgModel;

  @Output() search = new EventEmitter<string>();

  get dateInvalid(): boolean {
    return (this.searchTime.errors && this.searchTime.errors.dateValidator && (this.searchTime.dirty || this.searchTime.touched)) || false;
  }

  convertDate(strDate: string): string {
    if(/^(0[1-9]|[12][0-9]|3[01])[- /.](0[1-9]|1[012])[- /.](19|20)\d\d$/.test(strDate)) {
      let parts = strDate.split(/[-\/]/);
      strDate = parts[2] /*Year*/ + '-' +parts[1] /*Month*/ + '-' + parts[0] /*Date*/;  
    }
    return strDate;
  }

  doSearch() {
    let searchTerm: string = '';
    if(this.searchTime.valid && this.dateInput) {
      let timestamp: number = new Date(this.convertDate(this.searchTime.value)).getTime() / 1000;
      if(this.oneDayOffset) {
        timestamp += 3600 * 24;
      }
      searchTerm = timestamp.toString();
    } 
    this.search.emit(searchTerm);   
  }
}
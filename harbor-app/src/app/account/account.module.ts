import { NgModule } from '@angular/core';
import { SignInComponent } from './sign-in/sign-in.component';
import { SharedModule } from '../shared/shared.module';
import { RouterModule } from '@angular/router'; 

@NgModule({
  imports: [ 
    SharedModule,
    RouterModule
  ],
  declarations: [ SignInComponent ],
  exports: [SignInComponent]
})
export class AccountModule {}
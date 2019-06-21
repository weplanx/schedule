import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {RouterModule, Routes} from '@angular/router';
import {registerLocaleData} from '@angular/common';
import {HttpClientModule} from '@angular/common/http';
import zh from '@angular/common/locales/zh';
import {NgZorroAntdModule, NZ_I18N, zh_CN} from 'ng-zorro-antd';

registerLocaleData(zh);

import {AppComponent} from './app.component';

// const routes: Routes = [
//   {path: '', loadChildren: './app.router.module#AppRouterModule', canActivate: [TokenService]},
//   {path: 'login', loadChildren: './login/login.module#LoginModule'},
// ];

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    NgZorroAntdModule,
    HttpClientModule,
  ],
  providers: [
    {provide: NZ_I18N, useValue: zh_CN}
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}

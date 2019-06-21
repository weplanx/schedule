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
import {QuickComponent} from './quick/quick.component';
import {ClusterComponent} from './cluster/cluster.component';
import {ProjectComponent} from './project/project.component';

const routes: Routes = [
  {path: '', component: QuickComponent},
  {path: 'project', component: ProjectComponent},
  {path: 'cluster', component: ClusterComponent},
];

@NgModule({
  declarations: [
    AppComponent,
    ClusterComponent,
    ProjectComponent,
    QuickComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    NgZorroAntdModule,
    HttpClientModule,
    RouterModule.forRoot(routes, {useHash: true}),
  ],
  providers: [
    {provide: NZ_I18N, useValue: zh_CN}
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}

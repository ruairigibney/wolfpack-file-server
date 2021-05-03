import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AppComponent } from './app.component';
import { FileListComponent } from './file-list/file-list.component';
import { IncidentViewComponent } from './incident-view/incident-view.component';

const routes: Routes = [
  { path: '', pathMatch: 'full', component: IncidentViewComponent },
  { path: 'incidents/:filename', pathMatch: 'prefix', component: IncidentViewComponent },
  { path: '**', redirectTo: '' },

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

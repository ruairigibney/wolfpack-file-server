import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { IncidentViewComponent } from './incident-view/incident-view.component';

const routes: Routes = [
  { path: '', pathMatch: 'full', component: IncidentViewComponent },
  { path: 'incidents/:id', pathMatch: 'prefix', component: IncidentViewComponent },
  { path: '**', redirectTo: '' },

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

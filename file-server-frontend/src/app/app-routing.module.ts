import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AppComponent } from './app.component';
import { FileListComponent } from './file-list/file-list.component';

const routes: Routes = [
  { path: '', pathMatch: 'full', component: AppComponent },
  { path: 'incidents', pathMatch: 'full', component: FileListComponent },
  
]

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

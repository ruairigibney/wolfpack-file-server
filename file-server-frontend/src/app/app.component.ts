import { Component, OnInit } from '@angular/core';
import { AuthService } from './auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit {
  isAuthed = false;
  title = 'file-server-frontend';

  constructor(private auth: AuthService) {
  }

  ngOnInit(): void {
    this.auth.doAuth();
    this.auth.gotCookie.subscribe(
      (gotCookie) => this.isAuthed = gotCookie
    ) 
  }

}

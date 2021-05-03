import { HttpClient } from '@angular/common/http';
import { Injectable, OnInit } from '@angular/core';
import { ActivatedRoute, ActivatedRouteSnapshot, NavigationEnd, Router } from '@angular/router';
import { CookieService } from 'ngx-cookie-service';
import { BehaviorSubject, Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { environment } from 'src/environments/environment.prod';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  public gotCookie: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  constructor(private router: Router, private http: HttpClient, private route: ActivatedRoute,
              private cookieService: CookieService) { }

  removeCookie(): void {
    this.cookieService.delete('wolfpack-file-server');
    this.gotCookie.next(false);
    this.router.navigateByUrl('/');
  }

  doAuth(): void {
    if (this.cookieService.check('wolfpack-file-server')) {
      this.gotCookie.next(true);
      return;
    }

    this.router.events.pipe(
      filter((event) => event instanceof NavigationEnd),
      map(() => this.rootRoute(this.route)),
      filter((route: ActivatedRoute) => route.outlet === 'primary'),
    ).subscribe((route: ActivatedRoute) => {
      const passcode = route.snapshot.queryParamMap.get('passcode');
      if (passcode) {
        this.http.get(`${environment.apiUrl}/token?passcode=${passcode}`, {withCredentials: true}).subscribe(
          () => {
            this.gotCookie.next(true);
            this.router.navigate(['/incidents']);
          }
        );
      }
    });
  }

  private rootRoute(route: ActivatedRoute): ActivatedRoute {
    while (route.firstChild) {
      route = route.firstChild;
    }
    return route;
  }

}

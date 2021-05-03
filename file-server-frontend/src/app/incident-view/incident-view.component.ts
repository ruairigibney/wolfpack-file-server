import { Component, EventEmitter, OnInit, Output, ViewEncapsulation } from '@angular/core';
import { DomSanitizer, SafeHtml, SafeStyle, SafeUrl } from '@angular/platform-browser';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import { FileApiService } from '../file-api.service';

@Component({
  selector: 'app-incident-view',
  templateUrl: './incident-view.component.html',
  styleUrls: ['./incident-view.component.scss']
})
export class IncidentViewComponent implements OnInit {
  trustedIncidentHtml: SafeHtml = '';

  constructor(private fileService: FileApiService,
              private sanitizer: DomSanitizer,
              private router: Router,
              private route: ActivatedRoute) {
    }

  ngOnInit(): void {
    this.route.params.subscribe(
      (params) => {
        const filename = `${params.id}.html`;
        if (filename) {
          this.fileService.getFile(filename).subscribe(
            (response) => this.trustedIncidentHtml = this.sanitizer.bypassSecurityTrustHtml(response));
          this.fileService.setCurrentFile(filename); }

      }
    );
  }

}

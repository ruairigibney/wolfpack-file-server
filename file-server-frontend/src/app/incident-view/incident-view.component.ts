import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { DomSanitizer, SafeHtml, SafeStyle, SafeUrl } from '@angular/platform-browser';
import { FileApiService } from '../file-api.service';

@Component({
  selector: 'app-incident-view',
  templateUrl: './incident-view.component.html',
  styleUrls: ['./incident-view.component.scss']
})
export class IncidentViewComponent implements OnInit {
  trustedIncidentHtml: SafeHtml = "";
  
  constructor(private fileService: FileApiService, private sanitizer: DomSanitizer) { }

  ngOnInit(): void {
    this.fileService.currentFile.subscribe(
      (data) => {
        if (data) {
          this.fileService.getFile(data).subscribe(
            (response) => this.trustedIncidentHtml = this.sanitizer.bypassSecurityTrustHtml(response)
          );
        }
        }
    )

  }

}

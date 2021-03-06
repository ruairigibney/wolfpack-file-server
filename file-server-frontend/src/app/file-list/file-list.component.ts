import { HttpClient } from '@angular/common/http';
import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../auth.service';
import { FileApiService } from '../file-api.service';
import { IncidentFile } from '../incident-file';

@Component({
  selector: 'app-file-list',
  templateUrl: './file-list.component.html',
  styleUrls: ['./file-list.component.scss']
})
export class FileListComponent implements OnInit {
  @Output() fileSelected: EventEmitter<string> = new EventEmitter();
  @Output() hasFiles: EventEmitter<boolean> = new EventEmitter();
  fileList: IncidentFile[] = [];
  currentFile = '';

  constructor(private fileService: FileApiService,
              private authService: AuthService,
              private router: Router) {
  }

  ngOnInit(): void {
    this.authService.gotCookie.subscribe(
      (hasCookie) => (hasCookie) ? (
        this.fileService.getFileList().subscribe(
          data => {
            this.fileList = data;
            if (this.fileList.length > 0) {this.hasFiles.emit(true); }
          })) : null
      );

    this.fileService.currentFile.subscribe((file) =>
      this.currentFile = file
    );
  }

  clicked(event: any, file: IncidentFile): void {
    this.currentFile = file.FileName;
    this.router.navigateByUrl(`/incidents/${file.FileName.replace('.html', '')}`);
  }
}

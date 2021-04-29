import { HttpClient } from '@angular/common/http';
import { Component, EventEmitter, OnInit, Output } from '@angular/core';
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
  currentFile: string = "";

  constructor(private fileService: FileApiService, private authService: AuthService) {
  }
  
  ngOnInit(): void {
    this.authService.gotCookie.subscribe(
      (hasCookie) => (hasCookie) ? (
        this.fileService.getFileList().subscribe(
          data => {
            this.fileList = data;
            if (this.fileList.length > 0) {this.hasFiles.emit(true)};
          })) : null
      )
  }

  clicked(event: any, file: IncidentFile): void {
    this.fileService.setCurrentFile(file.FileName);
    this.currentFile = file.FileName;
  }
}

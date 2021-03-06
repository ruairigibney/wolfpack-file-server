import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { IncidentFile } from './incident-file';

@Injectable({
  providedIn: 'root'
})
export class FileApiService {
  public currentFile: BehaviorSubject<string> = new BehaviorSubject<string>('');
  constructor(private http: HttpClient) { }

  getFileList(): Observable<IncidentFile[]>{
    return this.http.get<IncidentFile[]>(`${environment.apiUrl}/files/list`, {withCredentials: true});
  }

  getFile(fileName: string): Observable<string>{
    return this.http.get(`${environment.apiUrl}/files/content?filename=${fileName}`,
    {withCredentials: true, responseType: 'text'});
  }

  setCurrentFile(file: string): void{
    this.currentFile.next(file);
  }
}

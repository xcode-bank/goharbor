import { Injectable } from '@angular/core';
import { Http, URLSearchParams, Response } from '@angular/http';

import { Repository } from './repository';
import { Tag } from './tag';
import { VerifiedSignature } from './verified-signature';

import { Observable } from 'rxjs/Observable'
import 'rxjs/add/observable/of';
import 'rxjs/add/operator/mergeMap';

@Injectable()
export class RepositoryService {
  
  constructor(private http: Http){}

  listRepositories(projectId: number, repoName: string, page?: number, pageSize?: number): Observable<any> {
    console.log('List repositories with project ID:' + projectId);
    let params = new URLSearchParams();
    params.set('page', page + '');
    params.set('page_size', pageSize + '');
    return this.http
               .get(`/api/repositories?project_id=${projectId}&q=${repoName}&detail=1`, {search: params})
               .map(response=>response)
               .catch(error=>Observable.throw(error));
  }

  listTags(repoName: string): Observable<Tag[]> {
    return this.http
               .get(`/api/repositories/tags?repo_name=${repoName}&detail=1`)
               .map(response=>response.json())
               .catch(error=>Observable.throw(error));
  }

  listNotarySignatures(repoName: string): Observable<VerifiedSignature[]> {
    return this.http
               .get(`/api/repositories/signatures?repo_name=${repoName}`)
               .map(response=>response.json())
               .catch(error=>Observable.throw(error));
  }

  listTagsWithVerifiedSignatures(repoName: string): Observable<Tag[]> {
    return this.listTags(repoName)
               .map(res=>res)
               .flatMap(tags=>{
                 return this.listNotarySignatures(repoName).map(signatures=>{
                    tags.forEach(t=>{
                      for(let i = 0; i < signatures.length; i++) {
                        if(signatures[i].tag === t.tag) {
                          t.signed = true;
                          break;
                        }
                      }
                    });
                    return tags;
                  })
                  .catch(error=>{
                    return tags;
                  })
               })
               .catch(error=>Observable.throw(error));
  }

  deleteRepository(repoName: string): Observable<any> {
    console.log('Delete repository with repo name:' + repoName);
    return this.http
               .delete(`/api/repositories?repo_name=${repoName}`)
               .map(response=>response.status)
               .catch(error=>Observable.throw(error));
  }

  deleteRepoByTag(repoName: string, tag: string): Observable<any> {
    console.log('Delete repository with repo name:' + repoName + ', tag:' + tag);
    return this.http
               .delete(`/api/repositories?repo_name=${repoName}&tag=${tag}`)
               .map(response=>response.status)
               .catch(error=>Observable.throw(error));
  }

}
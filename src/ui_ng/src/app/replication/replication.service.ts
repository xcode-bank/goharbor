import { Injectable } from '@angular/core';
import { Http, URLSearchParams, Response } from '@angular/http';

import { Policy } from './policy';
import { Job } from './job';
import { Target } from './target';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/mergeMap';

@Injectable()
export class ReplicationService {
  constructor(private http: Http) {}

  listPolicies(policyName: string, projectId?: any): Observable<Policy[]> {
    if(!projectId) {
      projectId = '';
    }
    console.log('Get policies with project ID:' + projectId + ', policy name:' + policyName);
    return this.http
               .get(`/api/policies/replication?project_id=${projectId}&name=${policyName}`)
               .map(response=>response.json() as Policy[])
               .catch(error=>Observable.throw(error));
  }

  getPolicy(policyId: number): Observable<Policy> {
    console.log('Get policy with ID:' + policyId);
    return this.http
               .get(`/api/policies/replication/${policyId}`)
               .map(response=>response.json() as Policy)
               .catch(error=>Observable.throw(error));
  }

  createPolicy(policy: Policy): Observable<any> {
    console.log('Create policy with project ID:' + policy.project_id + ', policy:' + JSON.stringify(policy));
    return this.http
               .post(`/api/policies/replication`, JSON.stringify(policy))
               .map(response=>response.status)
               .catch(error=>Observable.throw(error));
  }

  updatePolicy(policy: Policy): Observable<any> {
    if (policy && policy.id) {
      return this.http
                 .put(`/api/policies/replication/${policy.id}`, JSON.stringify(policy))
                 .map(response=>response.status)
                 .catch(error=>Observable.throw(error));
    } 
    return Observable.throw(new Error("Policy is nil or has no ID set."));
  }

  createOrUpdatePolicyWithNewTarget(policy: Policy, target: Target): Observable<any> {
    return this.http
               .post(`/api/targets`, JSON.stringify(target))
               .map(response=>{
                 return response.status;
               })
               .catch(error=>Observable.throw(error))
               .flatMap((status)=>{
                 if(status === 201) {
                   return this.http
                              .get(`/api/targets?name=${target.name}`)
                              .map(res=>res)
                              .catch(error=>Observable.throw(error));
                 }
               })
               .flatMap((res: Response) => { 
                 if(res.status === 200) {
                   let lastAddedTarget= <Target>res.json()[0];
                   if(lastAddedTarget && lastAddedTarget.id) {
                     policy.target_id = lastAddedTarget.id;
                     if(policy.id) {
                       return this.http
                                  .put(`/api/policies/replication/${policy.id}`, JSON.stringify(policy))
                                  .map(response=>response.status)
                                  .catch(error=>Observable.throw(error));
                     } else {
                       return this.http
                                  .post(`/api/policies/replication`, JSON.stringify(policy))
                                  .map(response=>response.status)
                                  .catch(error=>Observable.throw(error));
                     }
                   } 
                 }
               })
               .catch(error=>Observable.throw(error));
  }

  enablePolicy(policyId: number, enabled: number): Observable<any> {
    console.log('Enable or disable policy ID:' + policyId + ' with activation status:' + enabled);
    return this.http
               .put(`/api/policies/replication/${policyId}/enablement`, {enabled: enabled})
               .map(response=>response.status)
               .catch(error=>Observable.throw(error));
  }

  deletePolicy(policyId: number): Observable<any> {
    console.log('Delete policy ID:' + policyId);
    return this.http
               .delete(`/api/policies/replication/${policyId}`)
               .map(response=>response.status)
               .catch(error=>Observable.throw(error));
  }

  // /api/jobs/replication/?page=1&page_size=20&end_time=&policy_id=1&start_time=&status=&repository=
  listJobs(policyId: number, status: string = '', repoName: string = '', startTime: string = '', endTime: string = '', page: number, pageSize: number): Observable<any> {
    console.log('Get jobs under policy ID:' + policyId);
    return this.http
               .get(`/api/jobs/replication?policy_id=${policyId}&status=${status}&repository=${repoName}&start_time=${startTime}&end_time=${endTime}&page=${page}&page_size=${pageSize}`)
               .map(response=>response)
               .catch(error=>Observable.throw(error));
  }

  listTargets(targetName: string): Observable<Target[]> {
    console.log('Get targets.');
    return this.http
               .get(`/api/targets?name=${targetName}`)
               .map(response=>response.json() as Target[])
               .catch(error=>Observable.throw(error));
  }

  getTarget(targetId: number): Observable<Target> {
    console.log('Get target by ID:' + targetId);
    return this.http
               .get(`/api/targets/${targetId}`)
               .map(response=>response.json() as Target)
               .catch(error=>Observable.throw(error));
  }

  createTarget(target: Target): Observable<any> {
    console.log('Create target:' + JSON.stringify(target));
    return this.http
               .post(`/api/targets`, JSON.stringify(target))
               .map(response=>response.status)
               .catch(error=>Observable.throw(error));
  }

  pingTarget(target: Target): Observable<any> {
    console.log('Ping target.');
    let body = new URLSearchParams();
    body.set('endpoint', target.endpoint);
    body.set('username', target.username);
    body.set('password', target.password);
    return this.http
               .post(`/api/targets/ping`, body)
               .map(response=>response.status)
               .catch(error=>Observable.throw(error));
  }

  updateTarget(target: Target): Observable<any> {
    console.log('Update target with target ID' + target.id);
    return this.http
               .put(`/api/targets/${target.id}`, JSON.stringify(target))
               .map(response=>response.status)
               .catch(error=>Observable.throw(error));
  }

  deleteTarget(targetId: number): Observable<any> {
    console.log('Deleting  target with ID:' + targetId);
    return this.http
               .delete(`/api/targets/${targetId}`)
               .map(response=>response.status)
               .catch(error=>Observable.throw(error));
  }

}
import { Observable } from 'rxjs/Observable';
import { RequestQueryParams } from './RequestQueryParams';
import { ReplicationJob, ReplicationRule, ReplicationJobItem } from './interface';
import { Injectable, Inject } from "@angular/core";
import 'rxjs/add/observable/of';
import { Http, RequestOptions } from '@angular/http';
import { SERVICE_CONFIG, IServiceConfig } from '../service.config';
import { buildHttpRequestOptions, HTTP_JSON_OPTIONS, HTTP_GET_OPTIONS } from '../utils';

/**
 * Define the service methods to handle the replication (rule and job) related things.
 * 
 * @export
 * @abstract
 * @class ReplicationService
 */
export abstract class ReplicationService {
    /**
     * Get the replication rules.
     * Set the argument 'projectId' to limit the data scope to the specified project;
     * set the argument 'ruleName' to return the rule only match the name pattern;
     * if pagination needed, use the queryParams to add query parameters.
     * 
     * @abstract
     * @param {(number | string)} [projectId]
     * @param {string} [ruleName]
     * @param {RequestQueryParams} [queryParams]
     * @returns {(Observable<ReplicationRule[]> | Promise<ReplicationRule[]> | ReplicationRule[])}
     * 
     * @memberOf ReplicationService
     */
    abstract getReplicationRules(projectId?: number | string, ruleName?: string, queryParams?: RequestQueryParams): Observable<ReplicationRule[]> | Promise<ReplicationRule[]> | ReplicationRule[];

    /**
     * Get the specified replication rule.
     * 
     * @abstract
     * @param {(number | string)} ruleId
     * @returns {(Observable<ReplicationRule> | Promise<ReplicationRule> | ReplicationRule)}
     * 
     * @memberOf ReplicationService
     */
    abstract getReplicationRule(ruleId: number | string): Observable<ReplicationRule> | Promise<ReplicationRule> | ReplicationRule;

    /**
     * Create new replication rule.
     * 
     * @abstract
     * @param {ReplicationRule} replicationRule
     * @returns {(Observable<any> | Promise<any> | any)}
     * 
     * @memberOf ReplicationService
     */
    abstract createReplicationRule(replicationRule: ReplicationRule): Observable<any> | Promise<any> | any;

    /**
     * Update the specified replication rule.
     * 
     * @abstract
     * @param {ReplicationRule} replicationRule
     * @returns {(Observable<any> | Promise<any> | any)}
     * 
     * @memberOf ReplicationService
     */
    abstract updateReplicationRule(replicationRule: ReplicationRule): Observable<any> | Promise<any> | any;

    /**
     * Delete the specified replication rule.
     * 
     * @abstract
     * @param {(number | string)} ruleId
     * @returns {(Observable<any> | Promise<any> | any)}
     * 
     * @memberOf ReplicationService
     */
    abstract deleteReplicationRule(ruleId: number | string): Observable<any> | Promise<any> | any;

    /**
     * Enable the specified replication rule.
     * 
     * @abstract
     * @param {(number | string)} ruleId
     * @returns {(Observable<any> | Promise<any> | any)}
     * 
     * @memberOf ReplicationService
     */
    abstract enableReplicationRule(ruleId: number | string, enablement: number): Observable<any> | Promise<any> | any;

    /**
     * Disable the specified replication rule.
     * 
     * @abstract
     * @param {(number | string)} ruleId
     * @returns {(Observable<any> | Promise<any> | any)}
     * 
     * @memberOf ReplicationService
     */
    abstract disableReplicationRule(ruleId: number | string): Observable<any> | Promise<any> | any;


    abstract replicateRule(ruleId: number | string): Observable<any> | Promise<any> | any;

    /**
     * Get the jobs for the specified replication rule.
     * Set query parameters through 'queryParams', support:
     *   - status
     *   - repository
     *   - startTime and endTime
     *   - page
     *   - pageSize
     * 
     * @abstract
     * @param {(number | string)} ruleId
     * @param {RequestQueryParams} [queryParams]
     * @returns {(Observable<ReplicationJob> | Promise<ReplicationJob> | ReplicationJob)}
     * 
     * @memberOf ReplicationService
     */
    abstract getJobs(ruleId: number | string, queryParams?: RequestQueryParams): Observable<ReplicationJob> | Promise<ReplicationJob> | ReplicationJob;

    /**
     * Get the log of the specified job.
     * 
     * @abstract
     * @param {(number | string)} jobId 
     * @returns {(Observable<string> | Promise<string> | string)} 
     * @memberof ReplicationService
     */
    abstract getJobLog(jobId: number | string): Observable<string> | Promise<string> | string;

    abstract stopJobs(jobId: number | string): Observable<string> | Promise<string> | string;
}

/**
 * Implement default service for replication rule and job.
 * 
 * @export
 * @class ReplicationDefaultService
 * @extends {ReplicationService}
 */
@Injectable()
export class ReplicationDefaultService extends ReplicationService {
    _ruleBaseUrl: string;
    _jobBaseUrl: string;
    _replicateUrl: string;

    constructor(
        private http: Http,
        @Inject(SERVICE_CONFIG) config: IServiceConfig
    ) {
        super();
        this._ruleBaseUrl = config.replicationRuleEndpoint ?
            config.replicationRuleEndpoint : '/api/policies/replication';
        this._jobBaseUrl = config.replicationJobEndpoint ?
            config.replicationJobEndpoint : '/api/jobs/replication';
        this._replicateUrl = '/api/replications';
    }

    //Private methods
    //Check if the rule object is valid
    _isValidRule(rule: ReplicationRule): boolean {
        return rule !== undefined && rule != null && rule.name !== undefined && rule.name.trim() !== '' && rule.target_id !== 0;
    }

    public getReplicationRules(projectId?: number | string, ruleName?: string, queryParams?: RequestQueryParams): Observable<ReplicationRule[]> | Promise<ReplicationRule[]> | ReplicationRule[] {
        if (!queryParams) {
            queryParams = new RequestQueryParams();
        }

        if (projectId) {
            queryParams.set('project_id', '' + projectId);
        }

        if (ruleName) {
            queryParams.set('name', ruleName);
        }

        return this.http.get(this._ruleBaseUrl, buildHttpRequestOptions(queryParams)).toPromise()
            .then(response => response.json() as ReplicationRule[])
            .catch(error => Promise.reject(error))
    }

    public getReplicationRule(ruleId: number | string): Observable<ReplicationRule> | Promise<ReplicationRule> | ReplicationRule {
        if (!ruleId) {
            return Promise.reject("Bad argument");
        }

        let url: string = `${this._ruleBaseUrl}/${ruleId}`;
        return this.http.get(url, HTTP_GET_OPTIONS).toPromise()
            .then(response => response.json() as ReplicationRule)
            .catch(error => Promise.reject(error));
    }

    public createReplicationRule(replicationRule: ReplicationRule): Observable<any> | Promise<any> | any {
        if (!this._isValidRule(replicationRule)) {
            return Promise.reject('Bad argument');
        }

        return this.http.post(this._ruleBaseUrl, JSON.stringify(replicationRule), HTTP_JSON_OPTIONS).toPromise()
            .then(response => response)
            .catch(error => Promise.reject(error));
    }

    public updateReplicationRule(replicationRule: ReplicationRule): Observable<any> | Promise<any> | any {
        if (!this._isValidRule(replicationRule) || !replicationRule.id) {
            return Promise.reject('Bad argument');
        }

        let url: string = `${this._ruleBaseUrl}/${replicationRule.id}`;
        return this.http.put(url, JSON.stringify(replicationRule), HTTP_JSON_OPTIONS).toPromise()
            .then(response => response)
            .catch(error => Promise.reject(error));
    }

    public deleteReplicationRule(ruleId: number | string): Observable<any> | Promise<any> | any {
        if (!ruleId || ruleId <= 0) {
            return Promise.reject('Bad argument');
        }

        let url: string = `${this._ruleBaseUrl}/${ruleId}`;
        return this.http.delete(url, HTTP_JSON_OPTIONS).toPromise()
            .then(response => response)
            .catch(error => Promise.reject(error));
    }

    public replicateRule(ruleId: number | string): Observable<any> | Promise<any> | any {
        if (!ruleId) {
            return Promise.reject("Bad argument");
        }

        let url: string = `${this._replicateUrl}`;
        return this.http.post(url, {policy_id: ruleId}, HTTP_JSON_OPTIONS).toPromise()
            .then(response => response)
            .catch(error => Promise.reject(error));
    }

    public enableReplicationRule(ruleId: number | string, enablement: number): Observable<any> | Promise<any> | any {
        if (!ruleId || ruleId <= 0) {
            return Promise.reject('Bad argument');
        }

        let url: string = `${this._ruleBaseUrl}/${ruleId}/enablement`;
        return this.http.put(url, { enabled: enablement }, HTTP_JSON_OPTIONS).toPromise()
            .then(response => response)
            .catch(error => Promise.reject(error));
    }

    public disableReplicationRule(ruleId: number | string): Observable<any> | Promise<any> | any {
        if (!ruleId || ruleId <= 0) {
            return Promise.reject('Bad argument');
        }

        let url: string = `${this._ruleBaseUrl}/${ruleId}/enablement`;
        return this.http.put(url, { enabled: 0 }, HTTP_JSON_OPTIONS).toPromise()
            .then(response => response)
            .catch(error => Promise.reject(error));
    }

    public getJobs(ruleId: number | string, queryParams?: RequestQueryParams): Observable<ReplicationJob> | Promise<ReplicationJob> | ReplicationJob {
        if (!ruleId || ruleId <= 0) {
            return Promise.reject('Bad argument');
        }

        if (!queryParams) {
            queryParams = new RequestQueryParams();
        }

        queryParams.set('policy_id', '' + ruleId);
        return this.http.get(this._jobBaseUrl, buildHttpRequestOptions(queryParams)).toPromise()
            .then(response => {
                let result: ReplicationJob = {
                    metadata: {
                        xTotalCount: 0
                    },
                    data: []
                };

                if (response && response.headers) {
                    let xHeader: string = response.headers.get("X-Total-Count");
                    if (xHeader) {
                        result.metadata.xTotalCount = parseInt(xHeader, 0);
                    }
                }
                result.data = response.json() as ReplicationJobItem[];
                if (result.metadata.xTotalCount === 0) {
                    if (result.data && result.data.length > 0) {
                        result.metadata.xTotalCount = result.data.length;
                    }
                }

                return result;
            })
            .catch(error => Promise.reject(error));
    }

    public getJobLog(jobId: number | string): Observable<string> | Promise<string> | string {
        if (!jobId || jobId <= 0) {
            return Promise.reject('Bad argument');
        }

        let logUrl: string = `${this._jobBaseUrl}/${jobId}/log`;
        return this.http.get(logUrl, HTTP_GET_OPTIONS).toPromise()
            .then(response => response.text())
            .catch(error => Promise.reject(error));
    }

    public stopJobs(jobId: number | string): Observable<any> | Promise<any> | any {
        return this.http.put(this._jobBaseUrl, JSON.stringify({'policy_id': jobId, 'status': 'stop' }), HTTP_JSON_OPTIONS).toPromise()
            .then(response => response)
            .catch(error => Promise.reject(error));
    }
}
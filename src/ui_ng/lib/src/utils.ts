import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/toPromise';
import { RequestOptions, Headers } from '@angular/http';
import { RequestQueryParams } from './service/RequestQueryParams';
import { DebugElement } from '@angular/core';
import { Comparator, State } from 'clarity-angular';

/**
 * Convert the different async channels to the Promise<T> type.
 * 
 * @export
 * @template T
 * @param {(Observable<T> | Promise<T> | T)} async
 * @returns {Promise<T>}
 */
export function toPromise<T>(async: Observable<T> | Promise<T> | T): Promise<T> {
    if (!async) {
        return Promise.reject("Bad argument");
    }

    if (async instanceof Observable) {
        let obs: Observable<T> = async;
        return obs.toPromise();
    } else {
        return Promise.resolve(async);
    }
}

/**
 * The default cookie key used to store current used language preference.
 */
export const DEFAULT_LANG_COOKIE_KEY = 'harbor-lang';

/**
 * Declare what languages are supported now.
 */
export const DEFAULT_SUPPORTING_LANGS = ['en-us', 'zh-cn', 'es-es'];

/**
 * The default language.
 */
export const DEFAULT_LANG = 'en-us';

export const HTTP_JSON_OPTIONS: RequestOptions = new RequestOptions({
    headers: new Headers({
        "Content-Type": 'application/json',
        "Accept": 'application/json'
    })
});

export const HTTP_GET_OPTIONS: RequestOptions = new RequestOptions({
    headers: new Headers({
        "Content-Type": 'application/json',
        "Accept": 'application/json',
        "Cache-Control": 'no-cache',
        "Pragma": 'no-cache'
    })
});

/**
 * Build http request options
 * 
 * @export
 * @param {RequestQueryParams} params
 * @returns {RequestOptions}
 */
export function buildHttpRequestOptions(params: RequestQueryParams): RequestOptions {
    let reqOptions: RequestOptions = new RequestOptions({
        headers: new Headers({
            "Content-Type": 'application/json',
            "Accept": 'application/json',
            "Cache-Control": 'no-cache',
            "Pragma": 'no-cache'
        })
    });

    if (params) {
        reqOptions.search = params;
    }

    return reqOptions;
}



/** Button events to pass to `DebugElement.triggerEventHandler` for RouterLink event handler */
export const ButtonClickEvents = {
    left: { button: 0 },
    right: { button: 2 }
};


/** Simulate element click. Defaults to mouse left-button click event. */
export function click(el: DebugElement | HTMLElement, eventObj: any = ButtonClickEvents.left): void {
    if (el instanceof HTMLElement) {
        el.click();
    } else {
        el.triggerEventHandler('click', eventObj);
    }
}

/**
 * Comparator for fields with specific type.
 *  
 */
export class CustomComparator<T> implements Comparator<T> {

    fieldName: string;
    type: string;

    constructor(fieldName: string, type: string) {
        this.fieldName = fieldName;
        this.type = type;
    }

    compare(a: { [key: string]: any | any[] }, b: { [key: string]: any | any[] }) {
        let comp = 0;
        if (a && b) {
            let fieldA = a[this.fieldName];
            let fieldB = b[this.fieldName];
            switch (this.type) {
                case "number":
                    comp = fieldB - fieldA;
                    break;
                case "date":
                    comp = new Date(fieldB).getTime() - new Date(fieldA).getTime();
                    break;
            }
        }
        return comp;
    }
}

/**
 * The default page size
 */
export const DEFAULT_PAGE_SIZE: number = 15;

/**
 * The state of vulnerability scanning
 */
export const VULNERABILITY_SCAN_STATUS = {
    unknown: "n/a",
    pending: "pending",
    running: "running",
    error: "error",
    stopped: "stopped",
    finished: "finished"
};

/**
 * Calculate page number by state
 */
export function calculatePage(state: State): number {
    if (!state || !state.page) {
        return 1;
    }

    return Math.ceil((state.page.to + 1) / state.page.size);
}

/**
 * Filter columns via RegExp
 * 
 * @export
 * @param {State} state 
 * @returns {void} 
 */
export function doFiltering<T extends { [key: string]: any | any[] }>(items: T[], state: State): T[] {
    if (!items || items.length === 0) {
        return items;
    }

    if (!state || !state.filters || state.filters.length === 0) {
        return items;
    }

    state.filters.forEach((filter: {
        property: string;
        value: string;
    }) => {
        items = items.filter(item => regexpFilter(filter["value"], item[filter["property"]]));
    });

    return items;
}

/**
 * Match items via RegExp
 * 
 * @export
 * @param {string} terms 
 * @param {*} testedValue 
 * @returns {boolean} 
 */
export function regexpFilter(terms: string, testedValue: any): boolean {
    let reg = new RegExp('.*' + terms + '.*', 'i');
    return reg.test(testedValue);
}

/**
 * Sorting the data by column
 * 
 * @export
 * @template T 
 * @param {T[]} items 
 * @param {State} state 
 * @returns {T[]} 
 */
export function doSorting<T extends { [key: string]: any | any[] }>(items: T[], state: State): T[] {
    if (!items || items.length === 0) {
        return items;
    }
    if (!state || !state.sort) {
        return items;
    }

    return items.sort((a: T, b: T) => {
        let comp: number = 0;
        if (typeof state.sort.by !== "string") {
            comp = state.sort.by.compare(a, b);
        } else {
            let propA = a[state.sort.by.toString()], propB = b[state.sort.by.toString()];
            if (typeof propA === "string") {
                comp = propA.localeCompare(propB);
            } else {
                if (propA > propB) {
                    comp = 1;
                } else if (propA < propB) {
                    comp = -1;
                }
            }
        }

        if (state.sort.reverse) {
            comp = -comp;
        }

        return comp;
    });
}

/**
 * Compare the two objects to adjust if they're equal
 * 
 * @export
 * @param {*} a 
 * @param {*} b 
 * @returns {boolean} 
 */
export function compareValue(a: any, b: any): boolean {
    if ((a && !b) || (!a && b)) return false;
    if (!a && !b) return true;

    return JSON.stringify(a) === JSON.stringify(b);
}

/**
 * Check if the object is null or empty '{}'
 * 
 * @export
 * @param {*} obj 
 * @returns {boolean} 
 */
export function isEmptyObject(obj: any): boolean {
    return !obj || JSON.stringify(obj) === "{}";
}

/**
 * Deeper clone all
 * 
 * @export
 * @param {*} srcObj 
 * @returns {*} 
 */
export function clone(srcObj: any): any {
    if (!srcObj) return null;
    return JSON.parse(JSON.stringify(srcObj));
}
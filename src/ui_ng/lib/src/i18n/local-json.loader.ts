import { TranslateLoader } from '@ngx-translate/core';
import 'rxjs/add/observable/of';
import { Observable } from 'rxjs/Observable';

import { SERVICE_CONFIG, IServiceConfig } from '../service.config';

/**
 * Declare a translation loader with local json object
 * 
 * @export
 * @class TranslatorJsonLoader
 * @extends {TranslateLoader}
 */
export class TranslatorJsonLoader extends TranslateLoader {
    constructor(private config: IServiceConfig) {
        super();
    }

    getTranslation(lang: string): Observable<any> {
        let dict: any = this.config &&
            this.config.localI18nMessageVariableMap &&
            this.config.localI18nMessageVariableMap[lang] ?
            this.config.localI18nMessageVariableMap[lang] : {};
        return Observable.of(dict);
    }
}
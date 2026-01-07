export namespace api {
	
	export class ApiGlobalEmoteImages {
	    url_1x: string;
	    url_2x: string;
	    url_4x: string;
	
	    static createFrom(source: any = {}) {
	        return new ApiGlobalEmoteImages(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url_1x = source["url_1x"];
	        this.url_2x = source["url_2x"];
	        this.url_4x = source["url_4x"];
	    }
	}
	export class ApiGlobalEmote {
	    id: string;
	    name: string;
	    images: ApiGlobalEmoteImages;
	    format: string[];
	    scale: string[];
	    theme_mode: string[];
	
	    static createFrom(source: any = {}) {
	        return new ApiGlobalEmote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.images = this.convertValues(source["images"], ApiGlobalEmoteImages);
	        this.format = source["format"];
	        this.scale = source["scale"];
	        this.theme_mode = source["theme_mode"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ApiSubscriptionTransport {
	    method: string;
	    callback: string;
	    session_id: string;
	
	    static createFrom(source: any = {}) {
	        return new ApiSubscriptionTransport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.method = source["method"];
	        this.callback = source["callback"];
	        this.session_id = source["session_id"];
	    }
	}
	export class ApiSubscription {
	    id: string;
	    status: string;
	    type: string;
	    version: string;
	    condition: Record<string, any>;
	    created_at: string;
	    transport: ApiSubscriptionTransport;
	    cost: number;
	
	    static createFrom(source: any = {}) {
	        return new ApiSubscription(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.status = source["status"];
	        this.type = source["type"];
	        this.version = source["version"];
	        this.condition = source["condition"];
	        this.created_at = source["created_at"];
	        this.transport = this.convertValues(source["transport"], ApiSubscriptionTransport);
	        this.cost = source["cost"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace services {
	
	export class AppEmote {
	    id: string;
	    name: string;
	    lightSrcSet: string;
	    darkSrcSet: string;
	
	    static createFrom(source: any = {}) {
	        return new AppEmote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.lightSrcSet = source["lightSrcSet"];
	        this.darkSrcSet = source["darkSrcSet"];
	    }
	}

}


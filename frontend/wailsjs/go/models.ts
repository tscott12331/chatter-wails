export namespace api {
	
	export class APISubscription {
	    id: string;
	    status: string;
	    type: string;
	    version: string;
	    // Go type: struct {}
	    condition: any;
	    created_at: string;
	    // Go type: struct { Method string "json:\"method\""; Callback string "json:\"callback\""; Session_id string "json:\"session_id\"" }
	    transport: any;
	    cost: number;
	
	    static createFrom(source: any = {}) {
	        return new APISubscription(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.status = source["status"];
	        this.type = source["type"];
	        this.version = source["version"];
	        this.condition = this.convertValues(source["condition"], Object);
	        this.created_at = source["created_at"];
	        this.transport = this.convertValues(source["transport"], Object);
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


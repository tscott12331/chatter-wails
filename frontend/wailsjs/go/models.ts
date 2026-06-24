export namespace api {
	
	export class ApiMessageDropReason {
	    code: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ApiMessageDropReason(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	    }
	}
	export class ApiPostMessagesData {
	    message_id: string;
	    is_sent: boolean;
	    drop_reason?: ApiMessageDropReason;
	
	    static createFrom(source: any = {}) {
	        return new ApiPostMessagesData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message_id = source["message_id"];
	        this.is_sent = source["is_sent"];
	        this.drop_reason = this.convertValues(source["drop_reason"], ApiMessageDropReason);
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
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new AppEmote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.lightSrcSet = source["lightSrcSet"];
	        this.darkSrcSet = source["darkSrcSet"];
	        this.type = source["type"];
	    }
	}
	export class AppUser {
	    id: string;
	    login: string;
	    display_name: string;
	    type: string;
	    broadcaster_type: string;
	    description: string;
	    profile_image_url: string;
	    offline_image_url: string;
	    view_count: number;
	    email: string;
	    // Go type: time
	    created_at: any;
	    access_token: string;
	
	    static createFrom(source: any = {}) {
	        return new AppUser(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.login = source["login"];
	        this.display_name = source["display_name"];
	        this.type = source["type"];
	        this.broadcaster_type = source["broadcaster_type"];
	        this.description = source["description"];
	        this.profile_image_url = source["profile_image_url"];
	        this.offline_image_url = source["offline_image_url"];
	        this.view_count = source["view_count"];
	        this.email = source["email"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.access_token = source["access_token"];
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
	export class ChatroomData {
	    subId: string;
	    broadcasterId: string;
	    channelEmotes: Record<string, AppEmote>;
	
	    static createFrom(source: any = {}) {
	        return new ChatroomData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.subId = source["subId"];
	        this.broadcasterId = source["broadcasterId"];
	        this.channelEmotes = this.convertValues(source["channelEmotes"], AppEmote, true);
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


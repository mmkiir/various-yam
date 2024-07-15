export namespace main {
	
	export class MediaDeviceInfo {
	    deviceId: string;
	    groupId: string;
	    kind: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new MediaDeviceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.deviceId = source["deviceId"];
	        this.groupId = source["groupId"];
	        this.kind = source["kind"];
	        this.label = source["label"];
	    }
	}

}


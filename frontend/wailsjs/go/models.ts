export namespace main {
	
	export class FileFilter {
	    displayName: string;
	    pattern: string;
	
	    static createFrom(source: any = {}) {
	        return new FileFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.displayName = source["displayName"];
	        this.pattern = source["pattern"];
	    }
	}
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
	export class OpenDialogOptions {
	    defaultDirectory: string;
	    defaultFilename: string;
	    title: string;
	    filters: FileFilter[];
	    showHiddenFiles: boolean;
	    canCreateDirectories: boolean;
	    resolvesAliases: boolean;
	    treatPackagesAsDirectories: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OpenDialogOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.defaultDirectory = source["defaultDirectory"];
	        this.defaultFilename = source["defaultFilename"];
	        this.title = source["title"];
	        this.filters = this.convertValues(source["filters"], FileFilter);
	        this.showHiddenFiles = source["showHiddenFiles"];
	        this.canCreateDirectories = source["canCreateDirectories"];
	        this.resolvesAliases = source["resolvesAliases"];
	        this.treatPackagesAsDirectories = source["treatPackagesAsDirectories"];
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


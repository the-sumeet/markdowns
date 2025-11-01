export namespace main {
	
	export class Config {
	    lastOpenedFile: string;
	    lastOpenedDirectory: string;
	    recentFiles: string[];
	    windowWidth: number;
	    windowHeight: number;
	    theme: string;
	    showHiddenFiles: boolean;
	    customSettings: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lastOpenedFile = source["lastOpenedFile"];
	        this.lastOpenedDirectory = source["lastOpenedDirectory"];
	        this.recentFiles = source["recentFiles"];
	        this.windowWidth = source["windowWidth"];
	        this.windowHeight = source["windowHeight"];
	        this.theme = source["theme"];
	        this.showHiddenFiles = source["showHiddenFiles"];
	        this.customSettings = source["customSettings"];
	    }
	}
	export class FileEntry {
	    name: string;
	    path: string;
	    isDirectory: boolean;
	    size: number;
	    // Go type: time
	    modTime: any;
	
	    static createFrom(source: any = {}) {
	        return new FileEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.isDirectory = source["isDirectory"];
	        this.size = source["size"];
	        this.modTime = this.convertValues(source["modTime"], null);
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
	export class CurrentFilesState {
	    currentDir?: FileEntry;
	    currentFile?: FileEntry;
	    fileInfo?: FileEntry;
	    contentHash?: string;
	
	    static createFrom(source: any = {}) {
	        return new CurrentFilesState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentDir = this.convertValues(source["currentDir"], FileEntry);
	        this.currentFile = this.convertValues(source["currentFile"], FileEntry);
	        this.fileInfo = this.convertValues(source["fileInfo"], FileEntry);
	        this.contentHash = source["contentHash"];
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


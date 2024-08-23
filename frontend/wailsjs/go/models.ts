export namespace ytdlp {
	
	export class SplitState {
	    start: string;
	    end: string;
	
	    static createFrom(source: any = {}) {
	        return new SplitState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.start = source["start"];
	        this.end = source["end"];
	    }
	}

}


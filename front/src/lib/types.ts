export interface Model {
	_id: string;
	_rev: string;
	displayName: string;
	basePath: string;
	dateCreated: string;
	description: string;
	summary: string;
	tags: string[];
	images: Fileinterface[];
	printFiles: Fileinterface[];
	modelFiles: ModelFileinterface[];
	otherFiles: Fileinterface[];
	notes: Note[];
}

export interface Fileinterface {
	path: string;
}

export interface Note {
	text: string;
	date: string;
}

export interface ModelFileinterface extends Fileinterface {
	thumbnail: string;
}

export interface GCodeMetaData {
	gCodeType: string;
	createdBy: string;
	createDate: string;
	totalTime: string;
	layerHeight: string;
	nozzleDiameter: string;
	material: string;
	filamentUsedG: string;
	filamentUsedM: string;
	printerType: string;
	thumbnail: string;
}

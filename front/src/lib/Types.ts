export interface Model {
	_id: string;
	_rev: string;
	displayName: string;
	basePath: string;
	dateCreated: string;
	description: string;
	summary: string;
	tags: string[];
	images: FileType[];
	printFiles: FileType[];
	modelFiles: ModelFileType[];
	otherFiles: FileType[];
	notes: Note[];
}

export interface FileType {
	path: string;
}

export interface Note {
	text: string;
	date: string;
}

export interface ModelFileType extends FileType {
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

export default {}
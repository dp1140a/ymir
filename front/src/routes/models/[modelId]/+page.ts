/** @type {import('./$types').PageLoad} */
import { _apiUrl } from '$lib/Utils';
import type { Model, GCodeMetaData, ModelFileType } from '$lib/Model';

export const load = async ({ fetch, params }) => {
	/**
	 * Fetch the Model
	 */
	let url = _apiUrl(`/v1/model/${params.modelId}`);
	let res = await fetch(url);
	if (!res.ok) {
		throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
	}
	const model: Model = await res.json();
	//console.log(model)
	/**
	 * Fetch the Metadata form the First PrintFile
	 */
	let metaData: GCodeMetaData;
	if (model.printFiles.length > 0) {
		url = _apiUrl('/v1/model/gcode?path=').concat(model.basePath, '/', model.printFiles[0].path);
		res = await fetch(url);
		if (!res.ok) {
			throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
		}
		metaData = await res.json();
	}

	/**
	 * Fetch STL thumbnails as Base64 strings and attach to modelFile
	 */
	for (let i = 0; i < model.modelFiles.length; i++) {
		if (model.modelFiles[i].path.split('.').pop().toLowerCase() === 'stl') {
			//console.log(model.modelFiles[i]);
			model.modelFiles[i]['thumbnail'] = await _getSTLThumbnail(
				model.modelFiles[i],
				model.basePath
			);
		}
	}
	//console.log(metaData);
	return { model, metaData };
};

export const _getSTLThumbnail = async (
	model: ModelFileType,
	modelPath: string
): Promise<string> => {
	//let thumbnail: string;
	const url = _apiUrl('/v1/model/stl/image?path=').concat(modelPath, '/', model.path);
	const res = await fetch(url);
	if (!res.ok) {
		throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
	}
	return await res.text();
};

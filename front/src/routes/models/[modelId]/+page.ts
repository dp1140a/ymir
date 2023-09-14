/** @type {import('./$types').PageLoad} */
import { _apiUrl } from '$lib/utils';
import type { Model, GCodeMetaData } from '$lib/types';

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
		if (model.modelFiles[i].path.split('.').pop() === 'stl') {
			//console.log(model.modelFiles[i]);
			url = _apiUrl('/v1/model/stl/image?path=').concat(
				model.basePath,
				'/',
				model.modelFiles[i].path
			);
			let res = await fetch(url);
			if (!res.ok) {
				throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
			}
			let imgStr = await res.text();
			model.modelFiles[i]['thumbnail'] = imgStr;
		}
	}
	console.log(metaData);
	return { model, metaData };
};

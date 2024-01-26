import { _apiUrl } from '$lib/Utils';
import type { Model } from '$lib/Model';

export const load = async ({ fetch }) => {
	const url: string = _apiUrl('/v1/model');
	//let url = "/v1/model";
	const res: Response = await fetch(url);
	if (!res.ok) {
		throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
	}
	const modJSON = await res.json();
	console.log(modJSON);
	/*
	Due to the change on the backend we now get a map[string]Model.
	We can change that here to a []Model or modify the page to work with
	the new map
 */
	const modMap: Map<string, Model> = new Map<string, Model>(Object.entries(modJSON));
	const models: Model[] = [...modMap.values()];

	return { url, models };
};

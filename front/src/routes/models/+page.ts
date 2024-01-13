import { _apiUrl } from '$lib/Utils';

export const load = async ({ fetch }) => {
	const url = _apiUrl('/v1/model');
	//let url = "/v1/model";
	const res = await fetch(url);
	if (!res.ok) {
		throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
	}
	const models = await res.json()
	return { url, models };
};

/*
const apiUrl = (path: string) => {
	//console.log(`${import.meta.env.VITE_API_URL}`);
	return `${import.meta.env.VITE_API_URL || 'http://localhost:8081'}${path}`;
};
 */

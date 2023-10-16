import { _apiUrl } from '$lib/Utils';

export const load = async ({ fetch }) => {
	const url = _apiUrl('/v1/printer');
	let res = await fetch(url);
	if (!res.ok) {
		throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
	}
	const printers = await res.json();
	return { url, printers };
};

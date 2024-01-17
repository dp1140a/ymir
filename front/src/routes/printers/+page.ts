import { _apiUrl } from '$lib/Utils';
import type { Printer } from '$lib/Printer';

export const load = async ({ fetch }) => {
	const url:string = _apiUrl('/v1/printer');
	const res:Response = await fetch(url);
	if (!res.ok) {
		throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
	}
	const priJSON = await res.json()

	/*
	Due to the change on the backend we now get a map[string]Printer.
	We can change that here to a []Printer or modify the page to work with
	the new map
 */
	const priMap:Map<string, Printer> = new Map<string, Printer>(Object.entries(priJSON))
	const printers:Printer[] = [...priMap.values()];

	return { url, printers };
};

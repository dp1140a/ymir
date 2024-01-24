import { _apiUrl } from '$lib/Utils';
import type {PageLoad} from './$types'
import type { Model } from '$lib/Model';
import type {Printer} from '$lib/Printer';

export const load:PageLoad = async ({fetch}) => {
	try{
		const [modelsRes, printersRes] = await Promise.all([
			fetch(_apiUrl('/v1/model')),
			fetch(_apiUrl('/v1/printer'))
		])

		if (!modelsRes.ok) {
			throw new Error(`HTTP error: ${modelsRes.status}`)
		}
		if (!printersRes.ok) {
			throw new Error(`HTTP error: ${printersRes.status}`)
		}

		const modJSON = await modelsRes.json()
		const priJSON = await printersRes.json()

		const models: Map<string, Model> = new Map<string, Model>(Object.entries(modJSON));
		const printers: Map<string, Printer> = new Map<string, Printer>(Object.entries(priJSON));

		return { models, printers };
	} catch (err) {
		console.log(err)
		return {error: "Unable to fetch data: ".concat(err)}
	}
};
import { STLLoader } from 'three/examples/jsm/loaders/STLLoader.js';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({fetch}) => {
	try{
		const loader = new STLLoader()
		const res = await fetch("/parts/MPR-1.stl")
		const model =  loader.parse(await res.arrayBuffer())
		return {model}
	}catch(err){
		console.log(err)
		return {error: "Unable to fetch stl ".concat(err)}
	}
}

import { STLLoader } from 'three/examples/jsm/loaders/STLLoader.js';
import { _apiUrl } from '$lib/Utils';

export const load = async () => {
	const url = _apiUrl(`/v1/model/stl?path=a9225674d4dbaec5/MPR-1.stl`);
	const loader = new STLLoader();
	const stlModel = await loader.loadAsync(url, function (geometry) {
		//console.log(geometry) // logs an object type BufferGeometry
		return geometry;
	});
	return { stlModel };
};

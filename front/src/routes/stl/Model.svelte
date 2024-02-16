<script lang="ts">
	/**
	 * https://github.com/grischaerbe/threlter/blob/main/src/components/CameraControls/CameraControls.svelte
	 * https://github.com/yomotsu/camera-controls#fittobox-box3ormesh-enabletransition--paddingtop-paddingleft-paddingbottom-paddingright--
	 * https://github.com/yomotsu/camera-controls/blob/fb7b5751043a246b40ec90e45eb09a4ce1b68fe0/examples/fit-and-padding.html
	 * https://wejn.org/assets/js/stlviewer-datatag.js
	 */
	import { T, Canvas, useThrelte} from '@threlte/core';
	import { OrbitControls, useSuspense } from '@threlte/extras';
	import {BufferGeometry, BufferAttribute} from 'three'
	import { injectLookAtPlugin } from '$lib/stl/lookAt';
	import { Vector3, Matrix4, Mesh} from 'three';

	export let geometry;
	injectLookAtPlugin()
	const suspend = useSuspense()
	const maxDistance = 200;
	const minDistance = 10;

	let mesh:Mesh = new Mesh()
	const { camera } = useThrelte();


	const fitandCenter = () =>  {
		console.log("Camera: ", camera)
		let middle = new Vector3();
		geometry.computeBoundingBox()
		geometry.boundingBox.getCenter(middle);
		console.log("Geo Before: ", geometry)
		geometry.applyMatrix4(new Matrix4().makeTranslation(-middle.x, -middle.y, -middle.z ) );
		console.log("Geo After: ", geometry)
		mesh.position.set(0,0,10)
		console.log(mesh)
		$camera.lookAt(mesh.position)
		console.log($camera.position.distanceTo(mesh.position))
		console.log($camera.position)

		return geometry
	}
</script>

		<T.PerspectiveCamera
			makeDefault
			position={[0,0,110]}
		>
			<OrbitControls enableDamping {maxDistance} {minDistance} />
		</T.PerspectiveCamera>
		<T.Mesh bind:ref={mesh}>
			{#await suspend(fitandCenter()) then geo}
				<T is={geo} />
			{/await}
			<T.MeshStandardMaterial color="#c74528"/>

		</T.Mesh>
		<T.DirectionalLight />
		<T.HemisphereLight />


<style>
</style>
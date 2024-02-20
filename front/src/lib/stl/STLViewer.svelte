<script lang="ts">
	import { T, Canvas } from '@threlte/core';
	import { OrbitControls, Gizmo} from '@threlte/extras';
	import { DirectionalLight, MeshStandardMaterial, Vector3, PerspectiveCamera } from 'three';
	import { injectLookAtPlugin } from '$lib/stl/lookAt';

	export let geometry;
	//export let ambientLightColor: number | [number, number, number] = [64, 64, 64];
	export let ambientLightIntensity = 1.0;
	export let materialColor = '#d34427';
	export let spotlightIntensity = 1.0;
	export let width:number;
	export let height:number;

	console.log("W: ", width, "/H: ", height)
	injectLookAtPlugin()
	let msm = new MeshStandardMaterial({ color: materialColor });
	let canvas
	let camera:PerspectiveCamera;
	let model;
	let dl1:DirectionalLight;
	let dl2:DirectionalLight;
	//let controls;
	let maxDistance
	let minDistance
	let zOffset = 1.75;

	function fitAndCenter() {
		console.log(canvas)
		console.log(model)
		console.log(camera)

		geometry.computeBoundingBox()
		let middle = new Vector3();
		geometry.boundingBox.getCenter(middle);
		var size = new Vector3();
		geometry.boundingBox.getSize(size);
		console.log("Size: ", size, " / Middle: ", middle)

		const fov = camera.fov * ( Math.PI / 180 );
		const fovh = 2*Math.atan(Math.tan(fov/2) * camera.aspect);
		console.log("FOV: ", camera.fov, " / FOVH: ", fovh, " /Aspect: ", camera.aspect)
		let dx = size.z / 2 + Math.abs( size.x / 2 / Math.tan( fovh / 2 ) );
		let dy = size.z / 2 + Math.abs( size.y / 2 / Math.tan( fov / 2 ) );
		let cameraZ = Math.max(dx, dy);
		console.log("DX: ", dx, "/ DY: ",dy, " /CamZ: ", cameraZ)
		// offset the camera, if desired (to avoid filling the whole canvas)
		if( zOffset !== undefined && zOffset !== 0 ) cameraZ *= zOffset;
		console.log("CameraZ: ", cameraZ)
		camera.position.set(0,0,cameraZ)
		console.log(camera)

		// set the far plane of the camera so that it easily encompasses the whole object
		const minZ = geometry.boundingBox.min.z;
		const cameraToFarEdge = ( minZ < 0 ) ? -minZ + cameraZ : cameraZ - minZ;

		camera.far = cameraToFarEdge * 3;


		// Set Lighting
		let dly = geometry.boundingBox.max.y*1.05
		let dlz = geometry.boundingBox.max.z*1.05
		console.log(dly)
		dl1.position.set(0,dly,0)
		dl2.position.set(0,0,dlz)
		maxDistance = cameraToFarEdge
		minDistance = geometry.boundingBox.max.z

		camera.updateProjectionMatrix();
	}
</script>

<Canvas size={{ width: width, height: height }}>
	<T.PerspectiveCamera
		makeDefault
		bind:ref={camera}
		on:create={() => {
			fitAndCenter();
			//ready = true;
		}}
	>
		<OrbitControls enableDamping maxDistance={maxDistance} minDistance={minDistance}/>
		<T.DirectionalLight intesity={spotlightIntensity} bind:ref={dl1}/>
		<T.DirectionalLight intesity={spotlightIntensity} bind:ref={dl2}/>
		<T.AmbientLight  intensity={ambientLightIntensity} />
	</T.PerspectiveCamera>
	<T.Mesh {geometry} material={msm} bind:ref={model} />
	<Gizmo verticalPlacement="top" paddingX={20} paddingY={6} size={96} xColor="#7b9246" yColor="#547c99" zColor="#996699"/>
</Canvas>
<style>
</style>
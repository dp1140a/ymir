<script lang="ts">
	import { T, Canvas} from '@threlte/core';
	import { OrbitControls, Center} from '@threlte/extras';
	import { MeshStandardMaterial, Vector3, Matrix4 } from 'three';
	import { onMount } from "svelte";

	export let geometry;
	export let ambientLightColor: number | [number, number, number] = [64, 64, 64];
	export let ambientLightIntensity = 0.01;
	export let materialColor = '#EB5406';
	//export let maxZoom = 20;
	//export let minZoom = 2000;
	export let spotlightIntensity = 1.5;
	export let width = 640;

	//Ignore the svelte-check warning here
	export let height = 480;

	let msm = new MeshStandardMaterial({ color: materialColor });
	let camera;
	let model;
	//let controls;


	let ready = undefined;

	function fitAndCenter() {
		let middle = new Vector3();
		geometry.computeBoundingBox();
		geometry.boundingBox.getCenter(middle);
		//model.geometry.applyMatrix4(new Matrix4().makeTranslation(-middle.x, -middle.y, -middle.z ) );
		//model.geometry.rotateX(Math.PI/-2);
		//model.geometry.rotateY(Math.PI/-6);
		let largestDimension = Math.max(geometry.boundingBox.max.x, geometry.boundingBox.max.y, geometry.boundingBox.max.z)
		let zMultiple = 1.25
		if (largestDimension <= 100 ){
			zMultiple = zMultiple*2
		}
		camera.position.z = largestDimension * zMultiple;
		camera.position.x = 0;
		camera.position.y = 0;
		console.log(geometry.boundingBox.max.x + ',' + geometry.boundingBox.max.y + ',' + geometry.boundingBox.max.z)
	}

</script>

	<Canvas size={{width:width, height: 1000}}>
		<T.PerspectiveCamera
			makeDefault
			bind:ref={camera}
			on:create={() => {
				fitAndCenter()
				ready = true
			}}
		>
			<OrbitControls nableDamping />
			<T.DirectionalLight position={[1, 1, 1]} intesity={spotlightIntensity} />
			<T.HemisphereLight color={ambientLightColor} intensity={ambientLightIntensity} />
		</T.PerspectiveCamera>
		<Center autoCenter on:center={({ width }) => {
    //console.log('The width of the bounding box is', width)
  }}>
			<T.Mesh {geometry} material={msm} bind:ref={model}/>
		</Center>
	</Canvas>
<style>
</style>

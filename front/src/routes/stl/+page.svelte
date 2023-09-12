<script>
	import { T, Canvas, useThrelte } from '@threlte/core';

	import { OrbitControls } from '@threlte/extras';

	import { MeshStandardMaterial } from 'three';

	export let data;

	let ctx;
	$: console.log(ctx);
</script>

<div>
	<Canvas bind:ctx>
		<T.PerspectiveCamera
			makeDefault
			position={[10, 10, 50]}
			on:create={({ ref }) => {
				ref.lookAt(0, 1, 0);
			}}
		>
			<OrbitControls enableDamping maxDistance=200 minDistance=10 />
		</T.PerspectiveCamera>

		<T.DirectionalLight position={[1, 10, 1]} intesity={'.5'} />
		<T.HemisphereLight color={[0, 255, 0]} intensity={0.5} />
		<T.Mesh geometry={data.stlModel} material={new MeshStandardMaterial({ color: '#953620' })} />
	</Canvas>
</div>

<style>
	div {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
	}
</style>

<script lang="ts">
	import {Canvas} from '@threlte/core';
	import Model from './Model.svelte';
	import { onMount } from 'svelte';
	let w:number;
	let h:number;
	export let data;

	let canvas:HTMLDivElement
	onMount(() => {
		let svBB = document.getElementById("stlViewer").getBoundingClientRect()
		console.log(svBB)
		w=Math.round(svBB.width *.6)
		h= Math.round(svBB.height *.6)
		console.log(h)
	});
</script>

<div class="justify-center content-center mx-32 w-full" style="height:80vh;" id="stlViewer">
	<Canvas size={{ width: w, height: h }}>
		<Model geometry={data.model} />
	</Canvas>
</div>

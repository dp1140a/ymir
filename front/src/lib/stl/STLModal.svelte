
<script lang="ts">
	import STLViewer from '$lib/stl/STLViewer.svelte';
	import { onMount } from 'svelte';
	import { getModalStore } from '@skeletonlabs/skeleton';
	export let geometry;
	export let parent
	console.log("Parent: ", parent)
	let w:number;
	let h:number;
	let modalStore = getModalStore();
	onMount(() => {
		let svBB = document.getElementById("stlViewer").getBoundingClientRect()
		let shBB = document.getElementById("stlHeader").getBoundingClientRect()
		console.log(svBB, "/", shBB)
		w=svBB.width
		h= Math.round(Math.round(svBB.height)-(shBB.height+2))
		console.log(h)
	});
</script>

<div class="modalContainer mx-32 w-full overflow-hidden" style="height:90vh;" id="stlViewer">
	<div class="position-fixed z-50 h-32 w-full" id="stlHeader">
		<div class="variant-ghost-error float-left m-8 max-w-fit rounded-md p-4 text-sm drop-shadow-md">
			<div>
				<span><i class="fa-regular fa-rotate"></i></span>
				<span>Right Mouse Button</span>
			</div>
			<div>
				<span><i class="fa-solid fa-magnifying-glass"></i></span>
				<span>Middle Mouse Button / Wheel</span>
			</div>
			<div>
				<span><i class="fa-regular fa-arrows-up-down-left-right"></i></span>
				<span>Left Mouse Button</span>
			</div>
		</div>
		<div>
			<button
				type="button"
				class="variant-ghost-error btn float-right mx-20 my-10"
				on:click={modalStore.close}
			>
				<span><i class="fa-solid fa-close" /></span>
				<span>Close</span>
			</button>
		</div>
	</div>
		<!-- <slot>No model specified</slot> -->
	<STLViewer {geometry} width={w} height={h}/>
</div>

<style>
	.modalContainer {
		background-color: white;
	}
</style>

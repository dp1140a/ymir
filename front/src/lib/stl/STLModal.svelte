<script>
	import STLViewer from '$lib/stl/STLViewer.svelte';
	import { modalStore} from '@skeletonlabs/skeleton';
	import { onMount } from "svelte";
	export let geometry;
	let w;
	let h;
	onMount(() => {
		//console.log(w + " x " + h)

	})
</script>


	<div class="modalContainer mx-32 w-full" >
		<div class="w-full h-48 position-fixed z-50">
			<div class="rounded-md p-4 m-8 float-left max-w-fit drop-shadow-md text-sm variant-ghost-error">
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
				<button type="button" class="btn float-right mx-20 my-10 variant-ghost-error" on:click={modalStore.close}>
					<span><i class="fa-solid fa-close" /></span>
					<span>Close</span>
				</button>
			</div>
		</div>
		<div bind:clientWidth={w} bind:clientHeight={h}>
			<slot>
				No model specified
			</slot>
			<STLViewer {geometry}
								 width={w}
								 height={h}/>
		</div>
	</div>


<style>
	.modalContainer {
		background-color: white;
	}

</style>

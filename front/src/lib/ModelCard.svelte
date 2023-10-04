<script lang="ts">
	import { scale, fade } from 'svelte/transition';
	import { expoIn } from 'svelte/easing';
	import { base } from '$app/paths';
	import { _apiUrl } from './Utils';
	import type { Model } from "$lib/Model"

	export let model:Model

	let img: HTMLImageElement;

	function imageLoaded() {
		if (img.naturalWidth > img.naturalHeight) {
			img.setAttribute('style', 'width: 100%; max-height: none;');
		}
		img.setAttribute('height', '202px');
	}
</script>

<div class="card text-center w-80 h-96" in:fade={{ duration: 750, easing: expoIn }}>
	<header class="card-header h-16 bg-surface-300 rounded-t-lg flex-wrap">
		<a href="{base}/models/{model._id}">{model.displayName}</a>
	</header>
	<section class="p-4">
		<div>
			<a href="{base}/models/{model._id}">
				{#if (model.images.length > 0)}
				<img
					class="img-div"
					src={_apiUrl('/v1/model/image?path=').concat(model.basePath,'/',model.images[0].path)}
					on:load={imageLoaded}
					bind:this={img}
					alt="model"
				/>
					{:else }
					<img
						class="img-div"
						src="/3d-model-icon.png"
						on:load={imageLoaded}
						bind:this={img}
						alt="model icon"
					/>
					{/if}
					</a>
		</div>
		<div class="mt-2">
			{#if model.summary != undefined }
				{model.summary}
			{/if}
		</div>
	</section>
	<footer class="card-footer">
		{#each model.tags as tag, i}
			<span class="chip variant-ghost-surface mx-1">{tag}</span>
		{/each}
	</footer>
</div>

<style>
	.img-div {
		margin: auto !important;
	}
</style>

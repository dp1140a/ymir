<script lang="ts">
	import { fade } from 'svelte/transition';
	import { expoIn } from 'svelte/easing';
	import { base } from '$app/paths';
	import { _apiUrl } from './Utils';
	import type { Model } from '$lib/Model';

	export let model: Model;

	let img: HTMLImageElement;

	function imageLoaded() {
		if (img.naturalWidth > img.naturalHeight) {
			img.setAttribute('style', 'width: 100%; max-height: none;');
		}
		img.setAttribute('height', '202px');
	}
</script>

<div class="card max-h-96 max-w-96 text-center" in:fade={{ duration: 750, easing: expoIn }}>
	<header class="card-header h-16 flex-wrap rounded-t-lg bg-surface-300 text-sm">
		<a href="{base}/models/{model._id}">{model.displayName}</a>
	</header>
	<section class="p-4">
		<div>
			<a href="{base}/models/{model._id}">
				{#if model.images.length > 0}
					<img
						class="img-div"
						src={_apiUrl('/v1/model/image?path=').concat(model.basePath, '/', model.images[0].path)}
						on:load={imageLoaded}
						bind:this={img}
						alt="model"
					/>
				{:else}
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
			{#if model.summary !== undefined}
				{model.summary}
			{/if}
		</div>
	</section>
	<footer class="card-footer">
		{#each model.tags as tag}
			<span class="variant-ghost-surface chip p-1 ml-1 mt-1 text-xs">{tag}</span>
		{/each}
	</footer>
</div>

<style>
	.img-div {
		margin: auto !important;
	}
</style>

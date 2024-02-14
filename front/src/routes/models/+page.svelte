<script lang="ts">
	import ModelCard from '$lib/ModelCard.svelte';
	import type { Model } from '$lib/Model';
	export let data;

	const models = data.models;

	//https://svelte.dev/repl/e67e1a90ef3945ec988bf39f6a10b6b3?version=3.32.3
	let filteredModels = [];
	let numModels = models.length
	// For Search Input
	let nameSearch = '';
	let tagSearch = '';
	export let searchTerm = '';

	const searchByTag = () => {
		filteredModels = models.filter((model: Model) => {
			searchTerm = tagSearch;
			return model.tags.some((tag) => {
				return tag.toLowerCase().includes(tagSearch.toLowerCase())
			});
		});
		numModels = filteredModels.length
		return filteredModels;
	};

	const searchByName = () => {
		filteredModels = models.filter((model) => {
			searchTerm = nameSearch;
			return model.displayName.toLowerCase().includes(nameSearch.toLowerCase());
		});
		numModels = filteredModels.length
		return filteredModels
	};
</script>

<h1 class="h1 mt-4">Models</h1>
<div class="flex">
	<div class="w-1/5 px-4">
		<div class="font-semibold">Model Filter</div>
		<span>{numModels} Models</span>
		<form>
			<label class="label mb-2" for="">
				<span>Name:</span>
				<input
					class="input px-4 py-1"
					title="by name"
					type="search"
					name="nameSearch"
					placeholder="name"
					bind:value={nameSearch}
					on:input={searchByName}
				/>
			</label>
			<label class="label mb-8" for="">
				<span>Tag:</span>
				<input
					class="input px-4 py-1"
					title="by tag"
					type="search"
					name="tagSearch"
					placeholder="tag"
					bind:value={tagSearch}
					on:input={searchByTag}
				/>
			</label>
		</form>
	</div>
	<div class="w-full">
		{#if (searchTerm !== '' && filteredModels.length === 0) || models.length === 0}
			<div class="mt-6 text-center">
				No Models Found. Either Create a model or alter your search criteria.
			</div>
		{:else if filteredModels.length > 0}
			<div class="grid gap-8 md:grid-cols-4 lg:grid-cols-6">
				{#each filteredModels as model}
					<ModelCard {model} />
				{/each}
			</div>
		{:else}
			<div class="grid gap-8 md:grid-cols-4 lg:grid-cols-6">
				{#each models as model}
					<ModelCard {model} />
				{/each}
			</div>
		{/if}
	</div>
</div>

<style>
</style>

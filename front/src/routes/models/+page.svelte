<script lang="ts">
	import ModelCard from '$lib/ModelCard.svelte';
	import { base } from '$app/paths';
	export let data;

	const models = data.models;

	//https://svelte.dev/repl/e67e1a90ef3945ec988bf39f6a10b6b3?version=3.32.3
	let filteredModels = [];

	// For Search Input
	let nameSearch = '';
	let tagSearch = '';
	export let searchTerm = '';

	const searchByTag = () => {
		filteredModels = models.filter((model: any) => {
			searchTerm = tagSearch;
			return model.tags.some((tag) => tag.toLowerCase() === tagSearch.toLowerCase());
		});

		return filteredModels;
	};

	const searchByName = () => {
		return (filteredModels = models.filter((model) => {
			let modelDisplayName = model.displayName.toLowerCase();
			searchTerm = nameSearch;
			return modelDisplayName.includes(nameSearch.toLowerCase());
		}));
	};
</script>
<h1 class="h1 mt-4">Models</h1>
<span>{models.length} Models</span>
<div class="flex">
	<div class="w-1/5 px-4">
		<h1 class="font-semibold">Model Filter</h1>
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
	<div class="">
		<div class="grid lg:grid-cols-6 md:grid-cols-4 gap-8">

		{#if (searchTerm !== '' && filteredModels.length === 0) || models.length === 0}
			No Models Found. Either Create a model or alter your search criteria.
		{:else if filteredModels.length > 0}


			{#each filteredModels as model, i}
				<ModelCard {model} />
			{/each}

		{:else}
			{#each models as model, i}
				<ModelCard {model} />
			{/each}
		{/if}
	</div>
	</div>
</div>


<style>
</style>

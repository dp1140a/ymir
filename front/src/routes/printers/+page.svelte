<script lang="ts">
	import PrinterCard from '$lib/PrinterCard.svelte';
	import type { Printer } from '$lib/Printer';

	export let data;

	const printers = data.printers;

	//https://svelte.dev/repl/e67e1a90ef3945ec988bf39f6a10b6b3?version=3.32.3
	let filteredPrinters = [];

	// For Search Input
	let nameSearch = '';
	let tagSearch = '';
	export let searchTerm = '';

	const searchByTag = () => {
		filteredPrinters = printers.filter((printer: Printer) => {
			searchTerm = tagSearch;
			return printer.tags.some((tag) => tag.toLowerCase() === tagSearch.toLowerCase());
		});

		return filteredPrinters;
	};

	const searchByName = () => {
		return (filteredPrinters = printers.filter((printer) => {
			let printerDisplayName = printer.printerName.toLowerCase();
			searchTerm = nameSearch;
			return printerDisplayName.includes(nameSearch.toLowerCase());
		}));
	};
</script>

<h1 class="h1 mt-4">Printers</h1>
<span>{printers.length} Printers</span>
<div class="flex w-full flex-row">
	<div class="w-1/5 pr-6">
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
	<div class="w-full pl-8">
		<div class="flex justify-end">
			<a class="variant-filled-warning btn btn-sm mb-4" href="/printers/add"> + Add Printer </a>
		</div>
		<div class="w-full">
			{#if (searchTerm !== '' && filteredPrinters.length === 0) || printers.length === 0}
				No Printers Found. Either add a printer or alter your search criteria.
			{:else if filteredPrinters.length > 0}
				{#each filteredPrinters as printer}
					<PrinterCard {printer} />
				{/each}
			{:else}
				{#each printers as printer}
					<PrinterCard {printer} />
				{/each}
			{/if}
		</div>
	</div>
</div>

<style>
</style>

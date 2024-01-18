<!-- src/routes/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { _apiUrl, handleError } from '$lib/Utils';

	let errorType: string = '';
	let errorMessage: string = '';
	let errorVisible: boolean = false;
	async function init() {
		await fetch(_apiUrl('/v1/ping'), {
			method: 'GET'
		})
			.then(handleError)
			.then(() => {
				goto('/models');
			})
			.catch((error) => {
				errorMessage = 'Oops!  There was an error during init.';
				errorType = `${error.message}`;
				errorVisible = true;
			});
	}

	onMount(() => init());
</script>

<div>
	{#if errorVisible}
		<aside class="alert variant-filled-error py-6">
			<!-- Icon -->
			<div><i class="fa-solid fa-triangle-exclamation text px-4" /></div>
			<!-- Message -->
			<div class="px-4">
				<div><h3 class="h3">{errorMessage}</h3></div>
				<div class="text-sm">Error: {errorType}</div>
				<div class="text-lg">Please make sure the API server is running</div>
			</div>
		</aside>
		<br />
	{/if}
</div>

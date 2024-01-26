<script lang="ts">
	import { Accordion, AccordionItem } from '@skeletonlabs/skeleton';
	import { getModalStore } from '@skeletonlabs/skeleton';
	import type { ModalSettings } from '@skeletonlabs/skeleton';
	import { _apiUrl, handleError } from '$lib/Utils';
	import { invalidateAll } from '$app/navigation';
	const modalStore = getModalStore();

	import Prism from 'prismjs';
	import 'prismjs/components/prism-json.js';

	export let data;
	console.log(data);
	const printers = data.printers;
	let models = data.models;

	let errorType = '';
	let errorMessage = '';
	let errorVisible: boolean = false;

	const modal: ModalSettings = {
		buttonTextCancel: 'No',
		buttonTextConfirm: 'Yes',
		type: 'confirm',
		title: 'You Sure?'
	};

	const truncateModels = () => {
		(modal.body = `Are you sure you want to truncate the models? This will remove all models from the datastore.`),
			(modal.response = (r) => {
				if (r) {
					truncate('models');
				}
			});
		modalStore.trigger(modal);
	};

	const truncatePrinters = async () => {
		(modal.body = `Are you sure you want to truncate the printers. This will remove all printers from the datastore.`),
			(modal.response = (r) => {
				if (r) {
					truncate('printers');
				}
			});
		modalStore.trigger(modal);
	};

	const truncate = async (which: string) => {
		const url = _apiUrl(`/v1/admin/${which}`);
		let res = await fetch(url, {
			method: 'DELETE'
		});
		if (!res.ok) {
			throw `Error while truncating data from ${url} (${res.status} ${res.statusText}).`;
		}
		let data = await res.json();
		console.log(data);
		await invalidateAll();
		return;
	};

	//Save Model
	let saveDisabled: boolean[] = new Array(models.size).fill(true);
	function needsSave(idx: number) {
		saveDisabled[idx] = false;
	}

	const saveModel = async (key, idx) => {
		await fetch(_apiUrl(`/v1/model/key`), {
			method: 'PUT',
			headers: {
				Accept: 'application/json',
				'Content-Type': 'application/json'
			},
			body: models[key]
		})
			.then(handleError) // skips to .catch if error is thrown
			.then(() => {
				showUpdated('Complete', 'Model has been successfully Updated', true);
				saveDisabled[idx] = true;
			})
			.catch((error) => {
				let errorMessage =
					'Oops!  There was an error updating the model.<br/>Response was: ' + error;
				showUpdated(error, errorMessage, true);
			});
	};

	const savePrinter = async (key, idx) => {
		await fetch(_apiUrl(`/v1/printer/key`), {
			method: 'PUT',
			headers: {
				Accept: 'application/json',
				'Content-Type': 'application/json'
			},
			body: printers[key]
		})
			.then(handleError) // skips to .catch if error is thrown
			.then(() => {
				showUpdated('Complete', 'Printer has been successfully Updated', true);
				saveDisabled[idx] = true;
			})
			.catch((error) => {
				let errorMessage =
					'Oops!  There was an error updating the printer.<br/>Response was: ' + error;
				showUpdated(error, errorMessage, true);
			});
	};

	const exportItems = (filename, which) => {
		var dataStr =
			'data:text/json;charset=utf-8,' +
			encodeURIComponent(JSON.stringify(Object.fromEntries(which), null, 2));
		var downloadAnchorNode = document.createElement('a');
		downloadAnchorNode.setAttribute('href', dataStr);
		downloadAnchorNode.setAttribute('download', filename + '.json');
		document.body.appendChild(downloadAnchorNode); // required for firefox
		downloadAnchorNode.click();
		downloadAnchorNode.remove();
	};

	//Model Updated Modal
	export const showUpdated = (title: string, body: string, reload: boolean) => {
		const modal: ModalSettings = {
			type: 'alert',
			title: title,
			body: body,
			buttonTextCancel: 'Cool!'
			//response: (e) => { invalidateAll() }
		};

		if (reload) {
			modal.response = async () => {
				await invalidateAll();
				console.log('REloading');
				models = data.models;
			};
		}
		modalStore.trigger(modal);
	};
</script>

<div>
	{#if errorVisible}
		<aside class="alert variant-filled-error mb-4">
			<!-- Icon -->
			<div><i class="fa-solid fa-triangle-exclamation text" /></div>
			<!-- Message -->
			<div class="alert variant-filled-error alert-message text-sm">
				<div>
					<h3 class="h3">{errorType}</h3>
					<div class="h6">Message: {errorMessage}</div>
				</div>
			</div>
			<!-- Actions -->
			<div class="alert-actions">
				<button
					style="width: 1.5em;"
					class="variant-filled btn-icon"
					on:click|stopPropagation={() => {
						errorVisible = false;
					}}
				>
					<i class="fa-solid fa-xmark" />
				</button>
			</div>
		</aside>
	{/if}
</div>
<div>
	<Accordion hover="hover:bg-warning-hover-token">
		<!-- Models -->
		<AccordionItem class="my-2 rounded bg-surface-200">
			<svelte:fragment slot="lead"><i class="fa icon-orange fa-cube w-6" /></svelte:fragment>
			<svelte:fragment slot="summary">
				<span class="icon-orange h2">Models</span> ({models.size})
			</svelte:fragment>
			<svelte:fragment slot="content">
				<div>
					<a
						href={'#'}
						class="variant-filled-error btn btn-sm"
						type="button"
						on:click={() => truncateModels()}
					>
						<span>Truncate Models</span>
						<i class="fa-regular fa-circle-xmark float-right" />
					</a>
					<a
						href={'#'}
						class="variant-filled-secondary btn btn-sm"
						type="button"
						on:click={() => exportItems('models', models)}
					>
						<span>Export Models</span>
						<i class="fa-regular fa-download float-right" />
					</a>
				</div>
				<Accordion hover="hover:bg-warning-hover-token">
					{#each [...models] as [key, model], i}
						<AccordionItem class="my-2 rounded bg-zinc-500">
							<svelte:fragment slot="lead"><i class="fa fa-code w-6 align-top" /></svelte:fragment>
							<svelte:fragment slot="summary">
								{model.displayName}
								<div class="float-right">
									<button
										disabled={saveDisabled[i]}
										type="button"
										class="variant-filled-success btn btn-sm"
										on:click={() => saveModel(key, i)}
									>
										<span><i class="fa-regular fa-floppy-disk"></i></span>
										<span>Save Changes</span>
									</button>
								</div>
							</svelte:fragment>
							<svelte:fragment slot="content">
								<pre class="rounded border-zinc-800 bg-zinc-800 p-4"><code
										contenteditable="true"
										on:input={() => {
											needsSave(i);
										}}
										on:blur={(event) => {
											models[key] = event.target.innerText;
										}}
										>{@html Prism.highlight(
											JSON.stringify(model, null, 2),
											Prism.languages.json,
											'json'
										)}</code
									></pre>
							</svelte:fragment>
						</AccordionItem>
					{/each}
				</Accordion>
			</svelte:fragment>
		</AccordionItem>

		<!-- Printers -->
		<AccordionItem class="my-2 rounded bg-surface-200">
			<svelte:fragment slot="lead"><i class="fa icon-orange fa-cube w-6" /></svelte:fragment>
			<svelte:fragment slot="summary"
				><span class="icon-orange h2">Printers</span> ({printers.size})</svelte:fragment
			>
			<svelte:fragment slot="content">
				<div>
					<a
						href={'#'}
						class="variant-filled-error btn btn-sm"
						type="button"
						on:click={() => truncatePrinters()}
					>
						<span>Truncate Printers</span>
						<i class="fa-regular fa-circle-xmark float-right" />
					</a>
					<a
						href={'#'}
						class="variant-filled-secondary btn btn-sm"
						type="button"
						on:click={() => exportItems('printers', printers)}
					>
						<span>Export Models</span>
						<i class="fa-regular fa-download float-right" />
					</a>
				</div>
				<Accordion hover="hover:bg-warning-hover-token">
					{#each [...printers] as [key, printer], i}
						<AccordionItem class="my-2 rounded bg-zinc-500">
							<svelte:fragment slot="lead"
								><i class="fa fa-code w-6 text-center text-xl" /></svelte:fragment
							>
							<svelte:fragment slot="summary">
								{printer.printerName}
								<div class="float-right">
									<button
										disabled={saveDisabled[i]}
										type="button"
										class="variant-filled-success btn btn-sm"
										on:click={() => savePrinter(key, i)}
									>
										<span><i class="fa-regular fa-floppy-disk"></i></span>
										<span>Save Changes</span>
									</button>
								</div>
							</svelte:fragment>
							<svelte:fragment slot="content">
								<pre class="rounded border-zinc-800 bg-zinc-800 p-4"><code
										contenteditable="true"
										on:input={() => {
											needsSave(i);
										}}
										on:blur={(event) => {
											printers[key] = event.target.innerText;
										}}
										>{@html Prism.highlight(
											JSON.stringify(printer, null, 2),
											Prism.languages.json,
											'json'
										)}</code
									></pre>
							</svelte:fragment>
						</AccordionItem>
					{/each}
				</Accordion>
			</svelte:fragment>
		</AccordionItem>
	</Accordion>
</div>

<style>
	.icon-orange {
		color: #a53c23;
	}

	.editable:hover,
	[contenteditable='true']:active,
	[contenteditable='true']:focus {
		border: 0px;
		outline: none;
		/*
				background: rgb(211, 211, 211);
        border-radius: 4px;
        padding: 8px;
         */
	}
</style>

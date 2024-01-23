<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	//import { page } from '$app/stores';
	import { TabGroup, Tab, InputChip } from '@skeletonlabs/skeleton';
	import { getModalStore, type ModalSettings } from '@skeletonlabs/skeleton';
	import Notes from './Notes.svelte';
	import Files from './Files.svelte';
	import { handleError, _apiUrl } from '$lib/Utils';

	export let data;
	let model = data.model;
	//const modelId = $page.params.modelId;
	const modalStore = getModalStore();
	let errorType = '';
	let errorMessage = '';
	let validExtensions = '';
	let errorVisible: boolean = false;

	// Carousel ---
	let elemCarousel: HTMLDivElement;

	function carouselLeft(): void {
		const x =
			elemCarousel.scrollLeft === 0
				? elemCarousel.clientWidth * elemCarousel.childElementCount // loop
				: elemCarousel.scrollLeft - elemCarousel.clientWidth; // step left
		elemCarousel.scroll(x, 0);
	}

	function carouselRight(): void {
		const x =
			elemCarousel.scrollLeft === elemCarousel.scrollWidth - elemCarousel.clientWidth
				? 0 // loop
				: elemCarousel.scrollLeft + elemCarousel.clientWidth; // step right
		elemCarousel.scroll(x, 0);
	}

	function carouselThumbnail(index: number) {
		console.log(elemCarousel.clientWidth);
		elemCarousel.scroll(elemCarousel.clientWidth * index, 0);
	}

	// Tabs --
	let tabSet: number = 0;

	//Modals
	//Not Yet Implemented Modal
	export const showNYI = () => {
		const modal: ModalSettings = {
			type: 'alert',
			title: 'TODO:',
			body: 'Not yet implemented',
			buttonTextCancel: 'Got It!'
			//response: (e) => { }
		};
		modalStore.trigger(modal);
	};

	const showError = (err) => {
		console.log(err);
		errorType = err.detail.name;
		errorMessage = err.detail.message;
		validExtensions = err.detail.validExtensions;
		errorVisible = true;
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
			modal.response = () => {
				reloadPage();
			};
		}
		modalStore.trigger(modal);
	};

	//Save Model
	let saveDisabled = true;
	function needsSave() {
		saveDisabled = false;
		//console.log(model)
	}

	const reloadPage = async () => {
		await invalidateAll();
		model = data.model;
	};

	const saveModel = async () => {
		//console.log(model)
		await fetch(_apiUrl(`/v1/model/${model._id}`), {
			method: 'PUT',
			headers: {
				Accept: 'application/json',
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(model)
		})
			.then(handleError) // skips to .catch if error is thrown
			.then(() => {
				//model._rev = response.rev;
				showUpdated('Complete', 'Model has been successfully Updated', false);
				saveDisabled = true;
			})
			.catch((error) => {
				let errorMessage =
					'Oops!  There was an error updating the model.<br/>Response was: ' + error;
				showUpdated(error, errorMessage, true);
			});
	};

	const watchableArray = (target) =>
		new Proxy(target, {
			set(target, prop, value, receiver) {
				needsSave();
				return Reflect.set(target, prop, value, receiver);
			}
		});

	//Watchable Files Arrays for Files Component to trigger Save Model on Change
	let printFiles = watchableArray(model.printFiles);
	let modelFiles = watchableArray(model.modelFiles);
	let otherFiles = watchableArray(model.otherFiles);
	let notes = watchableArray(model.notes);
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
					<div class="h6">Valid Extensions: {validExtensions}</div>
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
	<div
		contenteditable="true"
		bind:textContent={model.displayName}
		on:input={needsSave}
		class="editable mt-2 h1"
	>
	{model.displayName}
	</div>
<hr class="my-4 !border-t-2" />
<div class="grid grid-flow-col gap-8">
	<div class="max-w-2xl">
		<!-- Carousel -->
		<div class="card grid grid-cols-[auto_1fr_auto] items-center gap-4 p-4">
			<!-- Button: Left -->
			<button type="button" class="variant-filled btn-icon" on:click={carouselLeft}>
				<i class="fa-solid fa-arrow-left" />
			</button>
			<!-- Full Images -->
			<div
				bind:this={elemCarousel}
				class="flex snap-x snap-mandatory overflow-x-auto scroll-smooth"
			>
				{#if model.images.length > 0}
					{#each model.images as image}
						<img
							class="h-[480px] w-full snap-center rounded-md object-cover"
							src={_apiUrl('/v1/model/image?path=').concat(model.basePath, '/', image.path)}
							alt={image.path}
							loading="lazy"
						/>
					{/each}
				{:else}
					<img
						class="h-[480px] w-full snap-center rounded-md object-cover"
						src="/3d-model-icon.png"
						alt="model placeholder"
						loading="lazy"
					/>
				{/if}
			</div>
			<!-- Button: Right -->
			<button type="button" class="variant-filled btn-icon" on:click={carouselRight}>
				<i class="fa-solid fa-arrow-right" />
			</button>
		</div>
		<!-- Thumbnails -->
		<div class="card grid grid-cols-6 gap-4 p-4">
			{#each model.images as image, i}
				<button type="button" on:click={() => carouselThumbnail(i)}>
					<img
						class="rounded-container-token"
						src={_apiUrl('/v1/model/image?path=').concat(model.basePath, '/', image.path)}
						alt={image.path}
						loading="lazy"
					/>
				</button>
			{/each}
		</div>
	</div>

	<!-- Model Details -->
	<div class="">
		<div class="h4">Model Tags:</div>
		<InputChip
			name="tags"
			bind:value={model.tags}
			chips="variant-ghost-primary"
			on:add={needsSave}
			on:remove={needsSave}
		/>
		<hr class="my-4 !border-t-2" />
		<div class="h4">Model Summary:<span class="h6 float-right text-xs">*click to edit</span></div>
		<div
			contenteditable="true"
			bind:textContent={model.summary}
			on:input={needsSave}
			class="editable mt-2"
		>
			{model.summary}
		</div>

		<hr class="my-4 !border-t-2" />
		<div class="h4 mb-4">Model Meta Data:</div>
		{#if data.metaData}
			<div class="grid w-fit grid-cols-7 gap-2">
				<div class="attr">
					<i class="icon fa-regular fa-clock" />
					<div>{data.metaData.totalTime}</div>
				</div>
				<div class="attr">
					<i class="icon fa-regular fa-file-code" />
					<div>{model.printFiles.length} file</div>
				</div>
				<div class="attr">
					<i class="icon iconfont-layer-height" />
					<div>{data.metaData.layerHeight}mm</div>
				</div>
				<div class="attr">
					<i class="icon iconfont-nozzle" />
					<div>{data.metaData.nozzleDiameter}mm</div>
				</div>
				<div class="attr">
					<i class="icon iconfont-material-spool" />
					<div>{data.metaData.material}</div>
				</div>
				<div class="attr">
					<i class="icon iconfont-layer-height" />
					<div>{data.metaData.filamentUsedG}g</div>
				</div>
				<div class="attr">
					<i class="icon iconfont-layer-height" />
					<div>{data.metaData.printerType}</div>
				</div>
			</div>
		{:else}
			<h5 class="variant-ghost-warning h6 my-6 rounded border py-1 text-center">
				No Model MetaData Available.<br />
				Need a G-Code file uploaded to the model.
			</h5>
		{/if}
		<hr class="my-4 !border-t-2" />
		<div class="float-right">
			<button
				disabled={saveDisabled}
				type="button"
				class="variant-filled-success btn my-1 w-48"
				on:click={saveModel}
			>
				<span><i class="fa-regular fa-floppy-disk"></i></span>
				<span>Save Changes</span>
			</button>
			<a
				href={_apiUrl('/v1/model/export?path=').concat(
					model.basePath,
					'&filename=',
					model.displayName
				)}
			>
				<button type="button" class="variant-filled-warning btn my-1 w-48">
					<span><i class="fa-solid fa-file-export" /></span>
					<span>Download</span>
				</button></a
			>
			<button type="button" class="variant-filled-error btn my-1 w-48" on:click={() => showNYI()}>
				<span><i class="fa-solid fa-print" /></span>
				<span>Print</span>
			</button>
		</div>
	</div>
</div>

<hr class="my-6 !border-t-2" />
<div class="mx-auto w-auto">
	<TabGroup justify="justify-center">
		<Tab bind:group={tabSet} name="Details" value={0}>
			<svelte:fragment slot="lead">Details</svelte:fragment>
		</Tab>
		<Tab bind:group={tabSet} name="Files" value={1}>Files</Tab>
		<Tab bind:group={tabSet} name="Notes" value={2}>Notes</Tab>
		<!-- Tab Panels --->
		<div class="mx-auto w-full" slot="panel">
			{#if tabSet === 0}
				<div class=" h2 mx-auto w-fit">
					Description<span class="h6 float-right ml-4 text-xs">*click to edit</span>
				</div>
				<div
					contenteditable="true"
					bind:textContent={model.description}
					on:input={needsSave}
					class="editable mx-auto mt-6 w-full cursor-text"
				>
					{model.description}
				</div>
			{:else if tabSet === 1}
				<div class="mx-auto w-fit">
					<Files
						on:uploadError={showError}
						modelBasePath={model.basePath}
						bind:modelFiles
						bind:otherFiles
						bind:printFiles
					/>
				</div>
				<br />&nbsp;
			{:else if tabSet === 2}
				<Notes bind:notes id={model._id} rev={model._rev} />
			{/if}
		</div>
	</TabGroup>
</div>

<style>
	.attr {
		color: #2a2a2a;
		font-size: 11px;
	}

	.editable:hover,
	[contenteditable='true']:active,
	[contenteditable='true']:focus {
		background: rgb(211, 211, 211);
		border: 1px solid rgb(133, 133, 133);
		border-radius: 4px;
		outline: none;
		padding: 8px;
	}
</style>

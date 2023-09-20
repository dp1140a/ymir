<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/stores';
	import { TabGroup, Tab, TabAnchor, InputChip } from "@skeletonlabs/skeleton";
	import { modalStore, type ModalSettings } from '@skeletonlabs/skeleton';
	import Notes from './Notes.svelte';
	import Files from './Files.svelte';
	import {handleError, _apiUrl} from "$lib/Utils";

	export let data;
	let model = data.model;
	const modelId = $page.params.modelId;

	let errorType = '';
	let errorMessage = '';
	let validExtensions ='';
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
			buttonTextCancel: 'Got It!',
			//response: (e) => { }
		};
		modalStore.trigger(modal);
	};

	const showError = (err) => {
		console.log(err);
		errorType = err.detail.name;
		errorMessage = err.detail.message;
		validExtensions = err.detail.validExtensions;
		errorVisible  = true;
	}

	//Model Updated Modal
	export const showUpdated = (title:string, body:string, reload: boolean) => {
		const modal: ModalSettings = {
			type: 'alert',
			title: title,
			body: body,
			buttonTextCancel: 'Cool!',
			//response: (e) => { invalidateAll() }
		};

		if (reload) {
			modal.response = (e) => { reloadPage()}
		}
		modalStore.trigger(modal);
	};

	//Save Model
	let saveDisabled = true;
	function needsSave() {
		saveDisabled = false;
		//console.log(model)
	}

	const reloadPage = async() => {
		await invalidateAll()
		model = data.model;
	}

	const saveModel = async () => {
		//console.log(model)
		const response = await fetch(_apiUrl(`/v1/model/${model._id}`), {
			method: 'PUT',
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(model)
		}).then(handleError) // skips to .catch if error is thrown
			.then((response) => {
				model._rev = response.rev;
				showUpdated('Complete', 'Model has been successfully Updated', false);
				saveDisabled = true
		}).catch((error) => {
			let errorMessage = 'Oops!  There was an error updating the model.<br/>Response was: ' + error;
			showUpdated(error, errorMessage, true);
		})
	}

	const watchableArray = target => new Proxy(target, {
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
					class="btn-icon variant-filled"
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
<h1 class="h1">{model.displayName}</h1>
<hr class="!border-t-2 my-4" />
<div class="grid grid-flow-col gap-8">
	<div class="max-w-2xl">
		<!-- Carousel -->
		<div class="card p-4 grid grid-cols-[auto_1fr_auto] gap-4 items-center">
			<!-- Button: Left -->
			<button type="button" class="btn-icon variant-filled" on:click={carouselLeft}>
				<i class="fa-solid fa-arrow-left" />
			</button>
			<!-- Full Images -->
			<div
				bind:this={elemCarousel}
				class="snap-x snap-mandatory scroll-smooth flex overflow-x-auto"
			>
				{#if (model.images.length >0)}
					{#each model.images as image}
						<img
							class="snap-center object-cover h-[480px] w-full rounded-md"
							src={ _apiUrl("/v1/model/image?path=").concat(model.basePath,'/',image.path)}
							alt={image.path}
							loading="lazy"
						/>
					{/each}
					{:else}
					<img
						class="snap-center object-cover h-[480px] w-full rounded-md"
						src="/3d-model-icon.png"
						alt="model placeholder"
						loading="lazy"
					/>
				{/if}

			</div>
			<!-- Button: Right -->
			<button type="button" class="btn-icon variant-filled" on:click={carouselRight}>
				<i class="fa-solid fa-arrow-right" />
			</button>
		</div>
		<!-- Thumbnails -->
		<div class="card p-4 grid grid-cols-6 gap-4">
			{#each model.images as image, i}
				<button type="button" on:click={() => carouselThumbnail(i)}>
					<img
						class="rounded-container-token"
						src={ _apiUrl("/v1/model/image?path=").concat(model.basePath,'/',image.path)}
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
		<hr class="!border-t-2 my-4" />
		<div class="h4">Model Summary:<span class="h6 text-xs float-right">*click to edit</span></div>
		<div contenteditable="true" bind:textContent={model.summary} on:input={needsSave} class="editable mt-2">{model.summary}</div>

		<hr class="!border-t-2 my-4" />
		<div class="h4 mb-4">Model Meta Data:</div>
		{#if data.metaData }
		<div class="grid w-fit gap-2 grid-cols-7">
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
				{:else }
			<h5 class="h6 border rounded variant-ghost-warning text-center py-1 my-6">
				No Model MetaData Available.<br/>
				Need a G-Code file uploaded to the model.
			</h5>
			{/if}
		<hr class="!border-t-2 my-4" />
		<div class="float-right">
			<button disabled={saveDisabled} type="button" class="btn w-48 my-1 variant-filled-success" on:click={saveModel}>
				<span><i class="fa-regular fa-floppy-disk"></i></span>
				<span>Save Changes</span>
			</button>
			<a href={ _apiUrl("/v1/model/export?path=").concat(model.basePath, "&filename=",model.displayName)}>
				<button type="button" class="btn w-48 my-1 variant-filled-warning">
					<span><i class="fa-solid fa-file-export" /></span>
					<span>Download</span>
				</button></a>
			<button type="button" class="btn w-48 my-1 variant-filled-error" on:click={() => showNYI()}>
				<span><i class="fa-solid fa-print" /></span>
				<span>Print</span>
			</button>
		</div>
	</div>
</div>

<hr class="!border-t-2 my-6" />
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
				<div class=" w-fit mx-auto h2">Description<span class="h6 text-xs ml-4 float-right">*click to edit</span></div>
				<div contenteditable="true" bind:textContent={model.description} on:input={needsSave}
						 class="mt-6 w-full mx-auto cursor-text editable">
					{model.description}
				</div>
			{:else if tabSet === 1}
				<div class="w-fit mx-auto">
				<Files
					on:uploadError={showError}
					modelBasePath={model.basePath}
					bind:modelFiles={modelFiles}
					bind:otherFiles={otherFiles}
					bind:printFiles={printFiles}
				/>
				</div>
				<br />&nbsp;
			{:else if tabSet === 2}
				<Notes bind:notes={notes} id={model._id} rev={model._rev} />
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
	[contenteditable="true"]:active,
	[contenteditable="true"]:focus{
		background: rgb(211, 211, 211);
		border: 1px solid rgb(133,133,133);
		border-radius: 4px;
		outline:none;
		padding: 8px;
	}
</style>

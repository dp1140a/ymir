<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { getModalStore, type ModalSettings, type ModalComponent } from '@skeletonlabs/skeleton';
	import STLModal from '$lib/stl/STLModal.svelte';
	import { STLLoader } from 'three/examples/jsm/loaders/STLLoader.js';
	import FilePond, { registerPlugin } from 'svelte-filepond'; //https://pqina.nl/filepond/docs/
	import FilePondPluginFileMetadata from 'filepond-plugin-file-metadata';
	import { _apiUrl } from '$lib/Utils';
	import type { GCodeMetaData, ModelFileType } from '$lib/Model';
	import { CheckFileType, FileUploadError } from '$lib/Files';
	import type { FilePondFile } from 'filepond';
	//import FilePondPluginImagePreview from "filepond-plugin-image-preview";
	import { _getSTLThumbnail } from './+page';
	import PrinterModal from '$lib/PrinterModal.svelte';

	const modalStore = getModalStore();
	registerPlugin(FilePondPluginFileMetadata);
	const dispatch = createEventDispatcher();
	export let modelBasePath = '';
	export let modelFiles = [];
	export let otherFiles = [];
	export let printFiles = [];
	let gCodeMetaData: GCodeMetaData;

	// Hooks to FilePond Instances
	let modelFilesPond;
	let otherFilesPond;
	let printFilesPond;

	//Modal
	// Provide the modal settings
	const modal: ModalSettings = {
		type: 'alert',
		title: 'Example Alert',
		body: 'This is an example modal.'
	};

	// Trigger the modal:
	export const showNYI = () => {
		modal.title = 'TODO:';
		modal.body = 'Not yet implemented';
		modal.buttonTextCancel = 'Got It!';
		modalStore.trigger(modal);
	};

	function checkUploadFileType(error: Error, fileItem: FilePondFile) {
		let err: FileUploadError;
		err = CheckFileType(fileItem);
		if (err !== null) {
			fileItem.abortLoad();
			dispatch('uploadError', err);
		}
	}

	async function completeUpload(err: Error, fileItem: FilePondFile) {
		if (err == null) {
			let pondName = fileItem.getMetadata().pondName;
			switch (pondName) {
				case 'Image_Files':
					break;
				case 'Model_Files':
					// eslint-disable-next-line no-case-declarations
					const modelFile = { path: fileItem.serverId } as ModelFileType;
					modelFile.thumbnail = await _getSTLThumbnail(modelFile, modelBasePath);
					modelFiles.push(modelFile);
					modelFiles = modelFiles;
					modelFilesPond.removeFiles();
					break;
				case 'Other_Files':
					otherFiles.push({ path: fileItem.serverId });
					otherFiles = otherFiles;
					otherFilesPond.removeFiles();
					break;
				case 'Print_Files':
					printFiles.push({ path: fileItem.serverId });
					printFiles = otherFiles;
					printFilesPond.removeFiles();
					break;
			}
		}
	}

	const getGCodeMetaData = async (printFiles) => {
		for (let i = 0; i < printFiles.length; i++) {
			//console.log(printFiles[i])
			const url = _apiUrl('/v1/model/gcode?path=').concat(modelBasePath, '/', printFiles[i].path);
			let res = await fetch(url);
			if (!res.ok) {
				throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
			}
			gCodeMetaData = await res.json();
			printFiles[i]['metadata'] = gCodeMetaData;
		}
	};

	const showModelSTL = (file: string) => {
		const url = _apiUrl('/v1/model/stl?path=').concat(modelBasePath, '/', file);
		let loader = new STLLoader();
		loader
			.loadAsync(url, function (geometry) {
				return geometry;
			})
			.then((geometry) => {
				const modalComponent: ModalComponent = {
					ref: STLModal,
					props: { geometry: geometry },
					slot: ''
				};

				const modal: ModalSettings = {
					type: 'component',
					backdropClasses: '--color-surface-50',
					component: modalComponent
				};
				modalStore.trigger(modal);
			});
	};

	const deleteFile = (index: number, files: string) => {
		switch (files) {
			case 'model':
				modelFiles.splice(index, 1);
				modelFiles = modelFiles;
				break;
			case 'other':
				otherFiles.splice(index, 1);
				otherFiles = otherFiles;
				break;
			case 'print':
				printFiles.splice(index, 1);
				printFiles = printFiles;
				break;
		}
	};

	const printFile = async (filePath: string) => {
		const url = _apiUrl('/v1/printer');
		let res = await fetch(url);
		if (!res.ok) {
			throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
		}
		await res.json().then((printers) => {
			const modalComponent: ModalComponent = {
				ref: PrinterModal,
				props: {
					printers: printers,
					modelBasePath: modelBasePath,
					filePath: filePath
				},
				slot: ''
			};

			new Promise<boolean>((resolve) => {
				const modal: ModalSettings = {
					type: 'component',
					backdropClasses: '--color-surface-50',
					component: modalComponent,
					response: (r: boolean) => {
						resolve(r);
					}
				};
				modalStore.trigger(modal);
			}).then(async () => {});
		});
		return;
	};
</script>

<div id="filesContainer" class="">
	<!-- Print Files -->
	<div id="printFiles" class="mt-10 rounded-md border border-black">
		<div
			class="section-header border-bottom variant-soft-surface h3 border-black px-4 pb-1 pt-4 text-center"
		>
			Print Files
			<div class="pt-6">
				<FilePond
					name="Print_Files"
					class="pond"
					allowMultiple={true}
					credits="false"
					server={{
						url: _apiUrl(`/v1/model/file?basePath=${modelBasePath}`)
					}}
					fileMetadataObject={{
						pondName: 'Print_Files'
					}}
					bind:this={printFilesPond}
					onaddfile={checkUploadFileType}
					onprocessfile={completeUpload}
					labelIdle="Drag & Drop your files or <span class='filepond--label-action'> Browse </span><br/>"
				/>
			</div>
		</div>
		{#await getGCodeMetaData(printFiles)}
			Getting Print File Details...
		<!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
		{:then _}
			{#if printFiles.length > 0}
				{#each printFiles as file, i}
					<div class="file-item-container flex justify-start px-2 py-6">
						{#if file.metadata.thumbnail}
							<div class="pr-4">
								<img
									src={file.metadata.thumbnail}
									height="90"
									alt="model thumbnail"
									class="thumbnail"
								/>
							</div>
						{:else}
							<div class="file-has-model">
								<i class="fa fa-regular fa-cube" />
							</div>
						{/if}
						<div class="info">
							<div>
								{file.path.split('/').at(-1)}
							</div>
							<div class="attributes grid grid-cols-7 gap-2">
								<div>
									<i class="icon iconfont-material-spool" />
									<div>
										{file.metadata.material ? file.metadata.material : 'unknown'}
									</div>
								</div>
								<div>
									<i class="icon iconfont-nozzle" />
									<div>
										{file.metadata.nozzleDiameter ? file.metadata.nozzleDiameter + 'mm' : 'unknown'}
									</div>
								</div>
								<div>
									<i class="icon iconfont-layer-height" /><!---->
									<div>
										{file.metadata.layerHeight ? file.metadata.layerHeight + 'mm' : 'unknown'}
									</div>
								</div>
								<div>
									<i class="icon iconfont-3d-printer" />
									<div>
										{file.metadata.printerType ? file.metadata.printerType : 'unknown'}
									</div>
								</div>
								<div>
									<i class="icon fa-regular fa-clock" />
									<div>{file.metadata.totalTime}</div>
								</div>
								<div>
									<i class="icon fa-regular fa-balance-scale" />
									<div>
										{file.metadata.filamentUsedG ? file.metadata.filamentUsedG + 'g' : 'unknown'}
									</div>
								</div>
							</div>
						</div>
						<div class="">
							<div class="float-right">
								<button type="button" on:click={() => deleteFile(i, 'print')}
									><i class="fa-regular fa-circle-xmark icon-orange float-right" /></button
								>
							</div>
							<div class="">
								<button
									type="button"
									class="variant-filled-error btn my-1 mt-6"
									on:click={() => {
										printFile(file.path);
									}}
								>
									<span><i class="fa-solid fa-print" /></span>
									<span>Print</span>
								</button>
							</div>
						</div>
					</div>
				{/each}
			{:else}
				<h5 class="variant-ghost-warning h6 my-6 rounded border py-1 text-center">
					-- There are no Print Files for this model. How about slicing a stl file? --
				</h5>
			{/if}
		{/await}
	</div>

	<!-- Model Files -->
	<div id="modelFiles" class="mt-10 rounded-md border border-black">
		<div
			class="section-header border-bottom variant-soft-surface h3 border-black px-4 pb-1 pt-4 text-center"
		>
			Model Files
			<div class="pt-6">
				<FilePond
					name="Model_Files"
					class="pond"
					allowMultiple={true}
					credits="false"
					server={{
						url: _apiUrl(`/v1/model/file?basePath=${modelBasePath}`)
					}}
					fileMetadataObject={{
						pondName: 'Model_Files'
					}}
					bind:this={modelFilesPond}
					onaddfile={checkUploadFileType}
					onprocessfile={completeUpload}
					labelIdle="Drag & Drop your files or <span class='filepond--label-action'> Browse </span><br/>"
				/>
			</div>
		</div>
		{#if modelFiles.length > 0}
			{#each modelFiles as file, i}
				<div class="file-item-container flex justify-start px-2 py-6">
					{#if file.thumbnail}
						<div class="pr-4">
							<a href="{'#'}" on:click={() => showModelSTL(file.path)}>
								<img
									src={file.thumbnail}
									height="90"
									alt="model thumbnail"
									class="thumbnail border border-neutral-400"
								/>
							</a>
						</div>
					{:else}
						<div class="mx-4">
							<i class="icon-orange fa-regular fa-cube" />
						</div>
					{/if}
					<div class="w-full">
						{file.path.split('/').at(-1)}
					</div>
					<div class="">
						<button type="button" on:click={() => deleteFile(i, 'model')}
							><i class="fa-regular fa-circle-xmark icon-orange float-right" /></button
						>
					</div>
				</div>
			{/each}
		{:else}
			<h5 class="variant-ghost-warning h6 my-6 rounded border py-1 text-center">
				-- There are no Model Files for this model. --
			</h5>
		{/if}
	</div>

	<!-- Other Files -->
	<div id="otherFiles" class="mt-10 rounded-md border border-black">
		<div
			class="section-header border-bottom variant-soft-surface h3 border-black px-4 pb-1 pt-4 text-center"
		>
			Other Files
			<div class="pt-6">
				<FilePond
					name="Other_Files"
					class="pond"
					allowMultiple={true}
					credits="false"
					server={{
						url: _apiUrl(`/v1/model/file?basePath=${modelBasePath}`)
					}}
					fileMetadataObject={{
						pondName: 'Other_Files'
					}}
					bind:this={otherFilesPond}
					onaddfile={checkUploadFileType}
					onprocessfile={completeUpload}
					labelIdle="Drag & Drop your files or <span class='filepond--label-action'> Browse </span><br/>"
				/>
			</div>
		</div>
		{#if otherFiles.length > 0}
			{#each otherFiles as file, i}
				<div class="file-item-container flex justify-start px-2 py-6">
					<div class="icon-orange mx-4">
						<i class="fa fa-regular fa-cube" />
					</div>
					<div class="w-full">
						{file.path.split('/').at(-1)}
					</div>
					<div class="">
						<button type="button" on:click={() => deleteFile(i, 'other')}
							><i class="fa-regular fa-circle-xmark icon-orange float-right" /></button
						>
					</div>
				</div>
			{/each}
		{:else}
			<h5 class="variant-ghost-warning h6 my-6 rounded border py-1 text-center">
				-- There are no Other Files for this model. --
			</h5>
		{/if}
	</div>
</div>

<style>
	.attributes {
		color: #2a2a2a;
		font-size: 11px;
	}

	.border-bottom {
		border-bottom: 1px solid #858585;
	}

	.thumbnail {
		max-height: 90px !important;
		max-width: 230px;
		width: auto;
		height: auto;
	}

	.icon-orange {
		color: #a53c23;
	}

	.icon {
		margin-bottom: 0;
	}

/**
	input[type='file'] {
		display: none;
	}

 */
</style>

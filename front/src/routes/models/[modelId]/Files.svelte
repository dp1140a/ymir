<script lang="ts">
	import { modalStore, type ModalSettings, type ModalComponent } from '@skeletonlabs/skeleton';
	import STLModal from '$lib/stl/STLModal.svelte';
	import { STLLoader } from 'three/examples/jsm/loaders/STLLoader.js';
	import FilePond, { registerPlugin, supported } from 'svelte-filepond'; //https://pqina.nl/filepond/docs/
	import FilePondPluginFileMetadata from 'filepond-plugin-file-metadata';
	import { _apiUrl } from "../../+layout";

	export let modelId = '';
	export let modelBasePath = '';
	export let modelFiles = [];
	export let otherFiles = [];
	export let printFiles = [];
	let gCodeMeta: GCodeMetaData;
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

	type GCodeMetaData = {
		filePath: string;
		gCodeType: string;
		createdBy: string;
		createdDate: string;
		totalTime: string;
		layerHeight: string;
		nozzleDiameter: string;
		material: string;
		filamentUsedG: string;
		filamentUsedM: string;
		printerType: string;
	};

	//let gCodeData = [];
	const getGCodeMetaData = async (printFiles) => {
		for (let i=0; i<printFiles.length; i++) {
			console.log(printFiles[i])
			const url = _apiUrl('/v1/model/gcode?path=').concat(modelBasePath, '/', printFiles[i].path);
			let res = await fetch(url);
			if (!res.ok) {
				throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
			}
			let metaData = await res.json();
			printFiles[i]['metadata'] = metaData;
			console.log(metaData)
			//gCodeData.push(metaData);
		}
	};

	const getSTLThumbNail = async (modelFiles) => {
		for (const file of modelFiles) {
			const url = _apiUrl('/v1/model/stlThumbnail?path=').concat(modelBasePath, '/', file.path);
			let res = await fetch(url);
			if (!res.ok) {
				throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
			}
			let imgStr = await res.text();
			return 'data:image/jpeg;base64,'.concat(imgStr);
		}
	};

	const basename = (path) => {
		return path.split('/').reverse()[0];
	};

	const showModelSTL = (file) => {
		const url = _apiUrl("/v1/model/stl?path=").concat(modelBasePath,"/",file);
		let loader = new STLLoader();
		loader.loadAsync(url, function (geometry) {
			return geometry;
		}).then((geometry) => {
			const modalComponent: ModalComponent = {
				ref: STLModal,
				props: { geometry: geometry },
				slot: ''
			};

			const modal: ModalSettings = {
				type: 'component',
				backdropClasses: '--color-surface-50',
				component: modalComponent,
			};
			modalStore.trigger(modal);
		});
	}

	function deleteFile(index: number, files: string){
		switch (files) {
			case "model":
				modelFiles.splice(index, 1);
				modelFiles = modelFiles;
				break;
			case "other":
				otherFiles.splice(index, 1);
				otherFiles = otherFiles;
				break;
			case "print":
				printFiles.splice(index, 1);
				printFiles = printFiles;
				break;
		}
	}


	function newFile(e) {
		console.log(e)
	}
</script>

<div id="filesContainer" class="">
	<!-- Print Files -->
	<div id="printFiles">
		<div class="h3 section-header border-bottom pt-6 text-center variant-soft-surface border-top">Print Files</div>
		{#await getGCodeMetaData(printFiles)}
			Getting Print File Details...
		{:then _}
			{#if printFiles.length > 0}
				{#each printFiles as file, i}
					<div class="file-item-container border-bottom py-6 flex justify-start">
						{#if file.metadata.thumbnail}
							<div class="pr-4">
								<img src={file.metadata.thumbnail} height="90" alt="model thumbnail" class="thumbnail" />
							</div>
						{:else}
							<div class="file-has-model">
								<i class="fa fa-regular fa-cube" />
							</div>
						{/if}
						<div class="info">
							<div>
								{file.path}
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
									<div>{file.metadata.layerHeight? file.metadata.layerHeight + 'mm' : 'unknown'}</div>
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
								<button type="button" on:click={() => deleteFile(i, "print")}><i class="fa-regular fa-circle-xmark float-right icon-orange" /></button>
							</div>
							<div class="">
								<button type="button" class="btn my-1 variant-filled-error mt-6" on:click={showNYI}>
									<span><i class="fa-solid fa-print" /></span>
									<span>Print</span>
								</button>
							</div>
						</div>
					</div>
				{/each}
			{:else }
				<h5 class="h6 border rounded variant-ghost-warning text-center py-1 my-6"> -- There are no Print Files for this model.  How about slicing a stl file? -- </h5>
			{/if}
		{/await}
	</div>

	<!-- Model Files -->
	<div id="modelFiles">
		<div class="h3 section-header border-bottom py-4 text-center variant-soft-surface">
			Model Files
			<label for="addModelFile" class="btn btn-sm variant-filled-error float-right mb-4 mr-3 ">
				<span><i class="fa-solid fa-add" /></span>
				<span>Add File</span>
			</label>
			<input type="file" id="addModelFile" class="" on:change={newFile} >
		</div>
		{#if modelFiles.length > 0}
		{#each modelFiles as file, i}
			<div class="file-item-container border-bottom py-6 flex justify-start">
				{#if file.thumbnail}
					<div class="pr-4 ">
						<a href="#" on:click={() => showModelSTL(file.path)}>
							<img src={file.thumbnail} height="90" alt="model thumbnail" class="thumbnail border border-neutral-400"/>
						</a>
					</div>
				{:else}
					<div class="mx-4">
						<i class="icon-orange fa-regular fa-cube" />
					</div>
				{/if}
				<div class="w-full">
					{file.path}
				</div>
				<div class="">
					<button type="button" on:click={() => deleteFile(i, "model")}><i class="fa-regular fa-circle-xmark float-right icon-orange" /></button>
				</div>
			</div>
		{/each}
			{:else }
			<h5 class="h6 border rounded variant-ghost-warning text-center py-1 my-6"> -- There are no Model Files for this model. -- </h5>
			{/if}
	</div>

	<!-- Other Files -->
	<div id="otherFiles">
		<div class="h3 section-header border-bottom pt-6 text-center variant-soft-surface">
			Other Files
			<FilePond
				name="Other_Files"
				class="filepond pond"
				allowMultiple={true}
				credits="false"
				server={{
							url: _apiUrl('/v1/model/file')
						}}
				fileMetadataObject={{
							pondName: 'Other_Files'
						}}

				labelIdle="Drag & Drop your files or <span class='filepond--label-action'> Browse </span><br/>"
			/>
		</div>
		{#if otherFiles.length > 0}
			{#each otherFiles as file, i}
				<div class="file-item-container border-bottom py-6 flex justify-start">
					<div class="icon-orange mx-4">
						<i class="fa fa-regular fa-cube" />
					</div>
					<div class="w-full">
						{file.path}
					</div>
					<div class="">
						<button type="button" on:click={() => deleteFile(i, "other")} ><i class="fa-regular fa-circle-xmark float-right icon-orange" /></button>
					</div>
				</div>
			{/each}
			{:else}
			<h5 class="h6 border rounded variant-ghost-warning text-center py-1 my-6"> -- There are no Other Files for this model. -- </h5>
		{/if}
	</div>
</div>

<style>

	.attributes {
		color: #2a2a2a;
		font-size: 11px;
	}

	.border-top{
		border-top: 1px solid #858585;
	}

	.border-bottom{
		border-bottom: 1px solid #858585;
	}

	.thumbnail {
		max-height: 90px !important;
		max-width:230px;
		width: auto;
		height: auto;
	}

	.icon-orange {
		color: #a53c23;
	}

	.icon {
		margin-bottom: 0px;
	}

	input[type="file"] {
		display: none;
	}

	.showMe {
		display: block !important;
			background-color: rgb(255 255 255 / var(--tw-bg-opacity)) !important;
	}

</style>

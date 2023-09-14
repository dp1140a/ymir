<script lang="ts">
	import { enhance } from '$app/forms';
	import { goto } from '$app/navigation';
	import { InputChip } from '@skeletonlabs/skeleton';
	import { Modal, modalStore } from '@skeletonlabs/skeleton';
	import type { ModalSettings, ModalComponent } from '@skeletonlabs/skeleton';
	import FilePond, { registerPlugin, supported } from 'svelte-filepond'; //https://pqina.nl/filepond/docs/
	import FilePondPluginImagePreview from 'filepond-plugin-image-preview';
	import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.css';
	import FilePondPluginFileMetadata from 'filepond-plugin-file-metadata';
	import CKEditor from '$lib/Ckeditor.svelte';
	import Editor from 'ckeditor5-custom-build/build/ckeditor'; // https://github.com/techlab23/ckeditor5-svelte/blob/master/src/Ckeditor.svelte
	import { _apiUrl } from "$lib/utils";

	// we need the same type
	type Data = {
		success: boolean;
		errors: Record<string, string>;
	};

	//export let data;

	// used in the template
	let form: Data;
	// Register the plugin
	registerPlugin(FilePondPluginFileMetadata, FilePondPluginImagePreview);

	let fTypes = ['image/*', 'model/stl'];
	let errorType = '';
	let errorMessage = '';
	let errorVisible: boolean = false;
	const imageTypes: string[] = ['gif', 'jpg', 'png', 'svg'];
	const modelTypes: string[] = [
		'3ds',
		'3mf',
		'amf',
		'blend',
		'dwg',
		'dxf',
		'f3d',
		'f3z',
		'factory',
		'fcstd',
		'iges',
		'ipt',
		'obj',
		'ply',
		'py',
		'rsdoc',
		'scad',
		'shape',
		'shapr',
		'skp',
		'sldasm',
		'sldprt',
		'slvs',
		'step',
		'stl',
		'stp'
	];
	const printTypes: string[] = ['gcode'];
	const otherTypes: string[] = [
		'csv',
		'doc',
		'ini',
		'json',
		'md',
		'pdf',
		'toml',
		'txt',
		'yaml',
		'yml',
		'zip'
	];
	let tags: string[] = [];

	function checkFileType(error, fileItem) {
		let pondName = fileItem.getMetadata().pondName;
		errorVisible = false;
		let extensions: string[];
		switch (pondName) {
			case 'Image_Files':
				extensions = imageTypes;
				break;
			case 'Model_Files':
				extensions = modelTypes;
				break;
			case 'Other_Files':
				extensions = otherTypes;
				break;
			case 'Print_Files':
				extensions = printTypes;
				break;
		}
		if (extensions.includes(fileItem.fileExtension)) {
		} else {
			fileItem.abortLoad();
			errorType = 'File Upload Error';
			errorMessage = `File Type ${fileItem.fileExtension} can't be uploaded to ${pondName.replace(
				'_',
				' '
			)} `;
			errorVisible = true;
		}
	}

	function hideSVG(e) {
		//console.log('Firing');
		const elms = Array.from(document.getElementsByClassName('filepond--image-preview'));
		elms.forEach((elm) => {
			//console.log(elm.style);
			elm.setAttribute("style", 'background-color: rgba(255, 255, 255, 0)');
		});
	}

	//https://joyofcode.xyz/working-with-forms-in-sveltekit
	async function handleForm(event: Event) {
		const formEl = event.target as HTMLFormElement;
		const data = new FormData(formEl);

		const response = await fetch(_apiUrl('/v1/model'), {
			method: 'POST',
			body: data
		})
			.then((response) => {
				if (!response.ok) {
					console.log(response);
					let eMsg;
					(async (response) => {
						eMsg = await response.json();
						console.log(eMsg);
						errorType = `${response.status}: ${response.statusText}`;
						errorMessage = 'Oops!  There was an error adding the model. Response was:' + eMsg;
						errorVisible = true;
					})(response);

					throw new Error(response.statusText);
				}

				// reset form
				formEl.reset();

				// rerun `load` function for the page
				//await invalidateAll()
				const modal: ModalSettings = {
					buttonTextCancel: 'OK',
					type: 'alert',
					title: 'Success!',
					body: 'The model ' + data.get('modelName') + ' has been successfully added.',
					response: () => {goto("/models")}
				};
				modalStore.trigger(modal);
				return response;
			}).then((response) => {
				console.log("Going to Models Page")
			}).catch((error) => {
				console.error(error);
				errorType = 'Model Create Error';
				errorMessage = 'There was an error on the server when creating the model';
				errorVisible = true;
			});
	}

	//Editor
	// Setting up editor prop to be sent to wrapper component
	let editor = Editor;
	//let editor = null;
	// Reference to initialised editor instance
	let editorInstance = null;
	// Setting up any initial data for the editor
	let editorData = 'write description here';

	// If needed, custom editor config can be passed through to the component
	// Uncomment the custom editor config if you need to customise the editor.
	// Note: If you don't pass toolbar object then Document editor will use default set of toolbar items.
	/**
	 * ['selectAll', 'undo', 'redo', 'alignment:left', 'alignment:right', 'alignment:center', 'alignment:justify', 'alignment', 'fontSize',
	 * 'fontFamily', 'fontColor', 'fontBackgroundColor', 'bold', 'italic', 'strikethrough', 'underline', 'blockQuote', 'link', 'ckfinder',
	 * 'uploadImage', 'imageUpload', 'heading', 'imageTextAlternative', 'toggleImageCaption', 'resizeImage:original', 'resizeImage:25',
	 * 'resizeImage:50', 'resizeImage:75', 'resizeImage', 'imageResize', 'imageStyle:inline', 'imageStyle:alignLeft', 'imageStyle:alignRight',
	 * 'imageStyle:alignCenter', 'imageStyle:alignBlockLeft', 'imageStyle:alignBlockRight', 'imageStyle:block', 'imageStyle:side',
	 * 'imageStyle:wrapText', 'imageStyle:breakText', 'indent', 'outdent', 'numberedList', 'bulletedList', 'mediaEmbed', 'insertTable',
	 * 'tableColumn', 'tableRow', 'mergeTableCells']
	 *
	 */
	let editorConfig = {};

	function onReady({ detail: editor }) {
		// Insert the toolbar before the editable area.
		//console.log(Array.from(editor.ui.componentFactory.names()))
		editorInstance = editor;

		editor.ui
			.getEditableElement()
			.parentElement.insertBefore(editor.ui.view.toolbar.element, editor.ui.getEditableElement());
	}
</script>

<div class="container mx-auto px-4">
	<form on:submit={handleForm}>
		<div>
			<div class="flex w-full">
				<div class="flex-none w-1/2"><h1 class="h1 mt-10">New Model</h1></div>
				<div class="flex-none w-1/2 ">
					<button type="submit" class="btn my-10 variant-filled-warning float-right">Submit</button></div>
			</div>
			<div>
				{#if errorVisible}
					<aside class="alert variant-filled-error py-6">
						<!-- Icon -->
						<div><i class="fa-solid fa-triangle-exclamation text" /></div>
						<!-- Message -->
						<div class="alert variant-filled-error alert-message text-sm">
							<h3 class="h3">{errorType}</h3>
							<br />
							<p>{errorMessage}</p>
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
					<br />
				{/if}
			</div>
		</div>
		<fieldset class="bg-surface-200 p-10 rounded-lg">
			<legend class="text-2xl">Basic Information</legend>
			<label class="label mb-8" for="">
				<span>Model name</span>
				<input
					class="input px-4 py-3"
					type="text"
					name="modelName"
					placeholder="model name"
					required
				/>
			</label>
			<label class="label mb-8" for="">
				<span>Summary</span>
				<textarea
					class="textarea"
					name="summary"
					rows="4"
					placeholder="Description of the model."
					required
				/>
			</label>
				<label class="label mb-8" for="">
					<span>Full Description</span>
				<div>
					<CKEditor name="description" bind:editor on:ready={onReady} bind:value={editorData} />
				</div>
			</label>
			<label class="label" for="">
				<span>Tags</span>
				<InputChip
					bind:tags
					name="tags"
					placeholder="Model tags"
					class="bg-white"
					chips="variant-ghost-warning"
					value={[]}
				/>
			</label>
		</fieldset>

		<hr class="!border-t-2 my-6" />
		<fieldset class="bg-surface-200 p-10 rounded-lg">
			<legend class="text-2xl">Files</legend>
			<!-- Print Files -->
			<label class="label" for="">
				<span>Image Files</span>
				<div class="pond">
					<FilePond
						name="Image_Files"
						class="filepond"
						allowMultiple={true}
						credits="false"
						server={{
							url: _apiUrl('/v1/model/file')
						}}
						fileMetadataObject={{
							pondName: 'Image_Files'
						}}
						imagePreviewMaxHeight="72"
						stylePanelLayout="compact"
						onprocessfilestart={hideSVG}
						onaddfile={checkFileType}
						labelIdle="Drag & Drop your files or <span class='filepond--label-action'> Browse </span><br/>{imageTypes}"
					/>
				</div>
			</label>

			<!-- Print Files -->
			<label class="label pt-4" for="">
				<span>Print Files</span>
				<div class="pond">
					<FilePond
						name="Print_Files"
						class="filepond"
						allowMultiple={true}
						credits="false"
						server={{
							url: _apiUrl('/v1/model/file')
						}}
						fileMetadataObject={{
							pondName: 'Print_Files'
						}}
						onaddfile={checkFileType}
						labelIdle="Drag & Drop your files or <span class='filepond--label-action'> Browse </span><br/>{printTypes}"
					/>
				</div>
			</label>
			<!-- Model Files -->
			<label class="label pt-4" for="">
				<span>Model Files</span>
				<div class="pond">
					<FilePond
						name="Model_Files"
						class="filepond"
						allowMultiple={true}
						credits="false"
						server={{
							url: _apiUrl('/v1/model/file')
						}}
						fileMetadataObject={{
							pondName: 'Model_Files'
						}}
						onaddfile={checkFileType}
						labelIdle="Drag & Drop your files or <span class='filepond--label-action'> Browse </span><br/>{modelTypes}"
					/>
				</div>
			</label>
			<!-- Other Files -->
			<label class="label pt-4" for="">
				<span>Other Files</span>
				<div class="pond">
					<FilePond
						name="Other_Files"
						class="filepond"
						allowMultiple={true}
						credits="false"
						server={{
							url: _apiUrl('/v1/model/file')
						}}
						fileMetadataObject={{
							pondName: 'Other_Files'
						}}
						onaddfile={checkFileType}
						labelIdle="Drag & Drop your files or <span class='filepond--label-action'> Browse </span><br/>{otherTypes}"
					/>
				</div>
			</label>
		</fieldset>
		<span style="float: right;"
			><button type="submit" class="btn my-10 variant-filled-warning">Submit</button></span
		>
	</form>
</div>

<!-- Styles -->
<style>
	.input,
	.textarea {
		background-color: rgb(255 255 255 / var(--tw-bg-opacity)) !important;
	}
</style>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { InputChip } from '@skeletonlabs/skeleton';
	import { getModalStore } from '@skeletonlabs/skeleton';
	import type { ModalSettings } from '@skeletonlabs/skeleton';
	import { popup } from '@skeletonlabs/skeleton';
	import FilePond, { registerPlugin } from 'svelte-filepond'; //https://pqina.nl/filepond/docs/
	import FilePondPluginImagePreview from 'filepond-plugin-image-preview';
	import 'filepond-plugin-image-preview/dist/filepond-plugin-image-preview.css';
	import FilePondPluginFileMetadata from 'filepond-plugin-file-metadata';
	import Markdown from 'svelte-exmarkdown';
	import { gfmPlugin } from 'svelte-exmarkdown/gfm';
	import { _apiUrl } from '$lib/Utils';

	const modalStore = getModalStore();
	// we need the same type
	// eslint-disable-next-line @typescript-eslint/no-unused-vars
	type Data = {
		success: boolean;
		errors: Record<string, string>;
	};

	// Register the plugin
	registerPlugin(FilePondPluginFileMetadata, FilePondPluginImagePreview);

	//let fTypes = ['image/*', 'model/stl'];
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
			/* empty */
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

	function hideSVG() {
		//console.log('Firing');
		const elms = Array.from(document.getElementsByClassName('filepond--image-preview'));
		elms.forEach((elm) => {
			//console.log(elm.style);
			elm.setAttribute('style', 'background-color: rgba(255, 255, 255, 0)');
		});
	}

	//https://joyofcode.xyz/working-with-forms-in-sveltekit
	async function handleForm(event: Event) {
		const formEl = event.target as HTMLFormElement;
		const data = new FormData(formEl);

		await fetch(_apiUrl('/v1/model'), {
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
					response: () => {
						goto('/models');
					}
				};
				modalStore.trigger(modal);
				return response;
			})
			.then(() => {
				console.log('Going to Models Page');
			})
			.catch((error) => {
				console.error(error);
				errorType = 'Model Create Error';
				errorMessage = 'There was an error on the server when creating the model';
				errorVisible = true;
			});
	}

	//Editor
	// Setting up any initial data for the editor
	let md = `
	# H1
## H2
### H3
#### H4
**bold**

_italic_

~~strikethrough~~

1. First ol item
2. Another ol item

* First ul Item
* Second ul item

[I'm an inline-style link](https://www.google.com)
![Image](https://github.com/adam-p/markdown-here/raw/master/src/common/images/icon48.png "Image")

\`\`\`javascript
var s = "JavaScript syntax highlighting";
alert(s);
\`\`\`

\`\`\`json
{
  "object": {
  "string": "hello
}
}
\`\`\`
sometext

\`\`\`
No language indicated, so no syntax highlighting.
But let's throw in a <b>tag</b>.
\`\`\`

| Tables        | Are           | Cool  |
| ------------- |:-------------:| -----:|
| col 3 is      | right-aligned | $1600 |
| col 2 is      | centered      |   $12 |
| zebra stripes | are neat      |    $1 |
---

> Blockquotes are very handy in email to emulate reply text.
	`;
	const plugins = [gfmPlugin()];

	let mdCodeVisible = true;
	let mdPreviewVisible = true;
	let currentMDBtn = 'both';

	const showMarkdown = (which) => {
		switch (which) {
			case 'code':
				mdCodeVisible = true;
				mdPreviewVisible = false;
				break;
			case 'both':
				mdCodeVisible = true;
				mdPreviewVisible = true;
				break;
			case 'preview':
				mdCodeVisible = false;
				mdPreviewVisible = true;
				break;
			default:
				break;
		}
		currentMDBtn = which;
	};
</script>

<div class="container mx-auto px-4">
	<form on:submit={handleForm}>
		<div>
			<div class="flex w-full">
				<div class="w-1/2 flex-none"><h1 class="h1 mt-10">New Model</h1></div>
				<div class="w-1/2 flex-none">
					<button type="submit" class="variant-filled-warning btn float-right my-10">Submit</button>
				</div>
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
								class="variant-filled btn-icon"
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
		<fieldset class="rounded-lg bg-surface-200 p-10">
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
					<div class="">
						<div class="card variant-filled-surface p-1" data-popup="editor">
							<p>Markdown Editor</p>
							<div class="variant-filled-surface arrow" />
						</div>
						<button
							type="button"
							class="btn btn-sm {currentMDBtn === 'code'
								? 'variant-ghost-primary'
								: 'variant-filled-primary'} [&>*]:pointer-events-none"
							use:popup={{ event: 'hover', target: 'editor', placement: 'top' }}
							style="padding: .3em;"
							on:click={() => showMarkdown('code')}
						>
							<span class="fa-stack">
								<i class="fa fa-code fa-stack-1x invert"></i>
							</span>
						</button>
						<div class="card variant-filled-surface p-1" data-popup="both">
							<p>Editor and Preview</p>
							<div class="variant-filled-surface arrow" />
						</div>
						<button
							type="button"
							class="btn btn-sm {currentMDBtn === 'both'
								? 'variant-ghost-secondary'
								: 'variant-filled-secondary'} [&>*]:pointer-events-none"
							use:popup={{ event: 'hover', target: 'both', placement: 'top' }}
							style="padding: .3em;"
							on:click={() => showMarkdown('both')}
						>
							<span class="fa-stack">
								<i class="fa-light fa-window-maximize fa-stack-2x invert"></i>
								<i class="fa fa-code fa-stack-1x invert"></i>
							</span>
						</button>
						<div class="card variant-filled-surface p-1" data-popup="preview">
							<p>Preview</p>
							<div class="variant-filled-surface arrow" />
						</div>
						<button
							type="button"
							class="btn btn-sm {currentMDBtn === 'preview'
								? 'variant-ghost-tertiary'
								: 'variant-filled-tertiary'} [&>*]:pointer-events-none"
							use:popup={{ event: 'hover', target: 'preview', placement: 'top' }}
							style="padding: .3em;"
							on:click={() => showMarkdown('preview')}
						>
							<span class="fa-stack">
								<i class="fa-light fa-window-maximize fa-stack-2x invert"></i>
							</span>
						</button>
					</div>
					{#if mdCodeVisible}
						<textarea class="h-64 w-full" bind:value={md} id="mdCode" />
					{/if}
					{#if mdPreviewVisible}
						<div class="ignore-css" id="md-preview">
							<Markdown {md} {plugins} />
						</div>
					{/if}
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

		<hr class="my-6 !border-t-2" />
		<fieldset class="rounded-lg bg-surface-200 p-10">
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
			><button type="submit" class="variant-filled-warning btn my-10">Submit</button></span
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

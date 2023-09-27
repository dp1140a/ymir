<script lang="ts">
	import { _apiUrl } from "../../+layout";

	export let id = '';
	export let rev = '';
	export let notes = [];

	let noteText = '';
	let disabled = true
	function isDisabled() {
		disabled = false;
	}

	function newNote(event: Event) {
		const formEl = event.target as HTMLFormElement;
		notes.unshift({
			"text": noteText,
			"date": new Date().toISOString()
		});
		notes = notes
		disabled = true
		formEl.reset();
	}
</script>

<div>
	<label class="label">
		<span>New Note</span>
		<form on:submit={newNote}>
			<input type="hidden" id="_id" name="_id" value={id} />
			<input type="hidden" id="_rev" name="_rev" value={rev} />
			<textarea
				name="noteText"
				bind:value={noteText}
				class="textarea"
				rows="4"
				on:input={isDisabled}
				placeholder="Note Text Here"
			/>
			<button type="submit" class="btn variant-filled-primary" disabled={disabled}>
				<span><i class="fa-regular fa-floppy-disk" /></span>
				<span>Save Note</span>
			</button>
		</form>
	</label>
	<hr class="!border-t-2 my-6" />
	<div class="table-container">
		<!-- Native Table Element -->
		<table class="table table-hover pb-4">
			<thead>
				<tr>
					<th class="w-24">Date</th>
					<th class="">Note</th>
				</tr>
			</thead>
			<tbody>
				{#each notes as note}
					<tr>
						<td class="w-1/6">{new Date(note.date).toLocaleString()}</td>
						<td>{note.text}</td>
					</tr>
				{/each}
			</tbody>
		</table>
		<br />&nbsp;
	</div>
</div>

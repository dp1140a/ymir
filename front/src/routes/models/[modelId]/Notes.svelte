<script lang="ts">
	export let id = '';
	export let rev = '';
	export let notes = [];

	//Note Sort
	const custom_sort = (a, b) => {
		return new Date(b.date).getTime() - new Date(a.date).getTime();
	};

	notes = notes.sort(custom_sort);
	let noteText = '';
	function isDisabled() {
		return noteText == '';
	}
</script>

<div>
	<label class="label">
		<span>New Note</span>
		<form method="POST" action="http://localhost:8081/v1/model/note">
			<input type="hidden" id="_id" name="_id" value={id} />
			<input type="hidden" id="_rev" name="_rev" value={rev} />
			<textarea
				name="text"
				bind:value={noteText}
				class="textarea"
				rows="4"
				placeholder="Lorem ipsum dolor sit amet consectetur adipisicing elit."
			/>
			<button type="submit" class="btn variant-filled-primary" disabled={isDisabled(noteText)}>
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
				{#each notes as note, i}
					<tr>
						<td>{note.date}</td>
						<td>{note.text}</td>
					</tr>
				{/each}
			</tbody>
		</table>
		<br />&nbsp;
	</div>
</div>

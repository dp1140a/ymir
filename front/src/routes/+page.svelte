<!-- src/routes/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
</script>

<div on:load|once={goto('/models')} />

<!-- src/routes/+page.svelte -->
<script lang="ts">
	import { goto } from '$app/navigation';
  import { onMount } from "svelte";

  onMount(() => goto('/models'))
</script>

<div />

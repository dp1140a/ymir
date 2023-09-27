<script lang="ts">
  import PrinterCard from "$lib/PrinterCard.svelte";

  export let data;

  const printers = data.printers;

  //https://svelte.dev/repl/e67e1a90ef3945ec988bf39f6a10b6b3?version=3.32.3
  let filteredPrinters = [];

  // For Search Input
  let nameSearch = '';
  let tagSearch = '';
  export let searchTerm = '';

  const searchByTag = () => {
    filteredPrinters = printers.filter((model: any) => {
      searchTerm = tagSearch;
      return model.tags.some((tag) => tag.toLowerCase() === tagSearch.toLowerCase());
    });

    return filteredPrinters;
  };

  const searchByName = () => {
    return (filteredPrinters = printers.filter((printer) => {
      let printerDisplayName = printer.displayName.toLowerCase();
      searchTerm = nameSearch;
      return printerDisplayName.includes(nameSearch.toLowerCase());
    }));
  };
</script>

<h1 class="h1 mt-4">Printers</h1>
<span>{printers.length} Printers</span>
<div class="flex flex-row w-[1200px]">
  <div class="w-1/3 px-4 border border-black">
    <h1 class="font-semibold">Model Filter</h1>
    <form>
      <label class="label mb-2" for="">
        <span>Name:</span>
        <input
          class="input px-4 py-1"
          title="by name"
          type="search"
          name="nameSearch"
          placeholder="name"
          bind:value={nameSearch}
          on:input={searchByName}
        />
      </label>
      <label class="label mb-8" for="">
        <span>Tag:</span>
        <input
          class="input px-4 py-1"
          title="by tag"
          type="search"
          name="tagSearch"
          placeholder="tag"
          bind:value={tagSearch}
          on:input={searchByTag}
        />
      </label>
    </form>
  </div>
  <div class="w-2/3 px-6 border border-black">
      {#if (searchTerm !== '' && filteredPrinters.length === 0) || printers.length === 0}
        No Models Found. Either Create a model or alter your search criteria.
      {:else if filteredPrinters.length > 0}
        {#each filteredPrinters as printer, i}
          <PrinterCard {printer} />
        {/each}

      {:else}
        {#each printers as printer, i}
          <PrinterCard {printer} />
        {/each}
      {/if}
  </div>
</div>

<style>

</style>
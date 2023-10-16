<script lang="ts">
  import { modalStore} from '@skeletonlabs/skeleton';
  import { ProgressRadial } from '@skeletonlabs/skeleton';
  import { type Printer, SelectedPrinter, UploadAndPrintFile } from "$lib/Printer";
  import { goto } from "$app/navigation";


  let selectedOption
  export let printers:Printer[]
  export let modelBasePath:string
  export let filePath:string

  let errorType = '';
  let errorMessage = '';
  let validExtensions ='';
  let errorVisible: boolean = false;
  let sending:boolean = false

  const doPrint = async() => {
    sending = true
    let res = await UploadAndPrintFile("".concat(modelBasePath,"/",filePath),$SelectedPrinter, true)
    if (res.error) {
      console.log(res.error)
      errorMessage = res.error.message
      errorVisible = true
      sending = false
    } else {
      goto(`/printers/${$SelectedPrinter._id}`)
      modalStore.close()
    }
  }

</script>
<div class="modal block overflow-y-auto bg-surface-100-800-token w-modal h-auto p-4 space-y-4 rounded-container-token shadow-xl " >
  <header class="modal-header text-2xl font-bold">Select Printer and Confirm</header>
    {#if errorVisible}
      <aside class="alert variant-filled-error mb-4">
        <!-- Icon -->
        <div><i class="fa-solid fa-triangle-exclamation text" /></div>
        <!-- Message -->
        <div class="alert variant-filled-error alert-message text-sm">
          <div>
            <h3 class="h3">Printer Upload Error</h3>
            <div class="h6">Message: {errorMessage}</div>
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
  {#if sending}
    <div class="mx-auto w-fit">
      <div class="my-4">Uploading . . .</div>
      <div class="my-4">
        <ProgressRadial width="w-18" stroke={200} meter="stroke-primary-500" track="stroke-primary-500/30"/>
      </div>
    </div>
  {:else}
  <div class="w-1/2 m-auto py-6">
    <div>This will upload and print the file.</div>
    <select class="select" bind:value={selectedOption} on:change={() => {SelectedPrinter.set(selectedOption)}}>
      <option value="" disabled selected>Select Printer:</option>
      {#each printers as printer,i}
        <option value="{printers[i]}">{printer.printerName}</option>
      {/each}
    </select>
  </div>
  {/if}
    <footer class="modal-footer flex justify-end space-x-2">
      <button type="button" class="btn variant-ghost-error" on:click={() => {modalStore.close()}}>
        <span><i class="fa-solid fa-cancel" /></span>
        <span>Cancel</span>
      </button>
      <button type="button" disabled={sending} class="btn variant-ghost-primary" on:click={() => {doPrint()}}>
        <span><i class="fa-solid fa-print" /></span>
        <span>Confirm</span>
      </button>
    </footer>
</div>

<style>
</style>

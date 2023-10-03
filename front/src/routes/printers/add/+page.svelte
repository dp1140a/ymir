<script lang="ts">

  import type { PrinterLocation } from "$lib/Printer.js";
  import type { PrinterType } from "$lib/Printer.js";
  import { _apiUrl } from "$lib/Utils";
  import { Modal, modalStore } from '@skeletonlabs/skeleton';
  import type { ModalSettings, ModalComponent } from '@skeletonlabs/skeleton';
  import { goto } from "$app/navigation";

  /**
  export interface Printer {
    "_id": string;
    "_rev": string;
    "printerName": string;
    "url": string;
    "apiType": string;
    "apiKey": string;
    "location": PrinterLocation;
    "type": PrinterType;
    "dateAdded": string;
    "tags": string[]
  }

  export interface PrinterLocation {
    "name": string
  }

  export interface PrinterType {
    "Make": string;
    "Model": string;
    "Version": string
  }
**/
  let errorType = '';
  let errorMessage = '';
  let errorVisible: boolean = false;

  async function handleForm(event: Event) {
    const formEl = event.target as HTMLFormElement;
    const data = new FormData(formEl);
    console.log(data)

    const response = await fetch(_apiUrl('/v1/printer'), {
      method: 'POST',
      body: data
    }).then((response) => {
      if (!response.ok) {
        //console.log(response);
        let eMsg;
        (async (response) => {
          eMsg = await response.json();
          //console.log(eMsg);
          errorType = `${response.status}: ${response.statusText}`;
          errorMessage = 'Oops!  There was an error adding the printer. Response was:' + eMsg;
          errorVisible = true;
        })(response);

        throw new Error(response.statusText);
      }

      // reset form
      formEl.reset();
      const modal: ModalSettings = {
        buttonTextCancel: 'OK',
        type: 'alert',
        title: 'Success!',
        body: `The printer ${data.get('printerName')} has been successfully added.`,
        response: () => {goto("/printers")}
      };
      modalStore.trigger(modal);
      return response;
    }).catch((error) => {
      console.error(error);
      errorType = 'Printer Add Error';
      errorMessage = 'There was an error on the server when adding the printer';
      errorVisible = true;
    });
  }
</script>

<div class="container mx-auto px-4">
  <form on:submit={handleForm}>
    <div>
      <div class="flex">
        <div class="flex-none w-1/2"><h1 class="h1 mt-10">Add Printer</h1></div>
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
        <span>Printer Name</span>
        <input
          class="input px-4 py-3"
          type="text"
          name="printerName"
          placeholder="printer name"
          required
        />
      </label>
      <label class="label mb-8" for="">
        <span>Url</span>
        <input
          class="input px-4 py-3"
          type="text"
          name="url"
          placeholder="printer URL"
          required
        />
      </label>
      <label class="label mb-8" for="">
        <span>API Type</span>
        <select class="select">
          <option selected value="octoprint">Octoprint</option>
        </select>
      </label>
      <label class="label mb-8" for="">
        <span>API key</span>
        <input
          class="input px-4 py-3"
          type="text"
          name="apiKey"
          placeholder="API key"
          required
        />
      </label>
      <label class="label mb-8" for="">
        <span>Printer Location</span>
        <input
          class="input px-4 py-3"
          type="text"
          name="location"
          placeholder="location"
          required
        />
      </label>
    </fieldset>
    <fieldset class="bg-surface-200 p-10 rounded-lg">
      <legend class="text-2xl">Printer Make and Model</legend>
      <label class="label mb-8" for="">
        <span>Make</span>
        <input
          class="input px-4 py-3"
          type="text"
          name="printerMake"
          placeholder="make"
          required
        />
      </label>
      <label class="label mb-8" for="">
        <span>Model</span>
        <input
          class="input px-4 py-3"
          type="text"
          name="printerModel"
          placeholder="model"
          required
        />
      </label>
    </fieldset>
  </form>
</div>
<style>
</style>

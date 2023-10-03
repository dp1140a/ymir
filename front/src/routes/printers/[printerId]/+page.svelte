<script lang="ts">
  import { Accordion, AccordionItem, type ModalSettings, modalStore } from "@skeletonlabs/skeleton";
  import RadialGauge from "$lib/RadialGauge.svelte";
  import { GetPrinterFiles, CheckPrinterStatus, type PrinterStatus, GetPrinterJob } from "$lib/Printer";
  import { _apiUrl, SecondsPrettyPrint } from "$lib/Utils.js";
  import type {JobInformation} from "$lib/Job"
  import { goto } from "$app/navigation";

  export let data;

  let printer= data.printer;
  let online = data.status.online;
  let printerStatus = data.status.printerStatus;
  let printerBusy = false

  let errorType = '';
  let errorMessage = '';
  let errorVisible: boolean = false;

  /**
   * Update Printer status every 10 seconds
   */
  const interval = setInterval(function() {
    (async ()=>{
      let status:{online: boolean, printerStatus:PrinterStatus} = await CheckPrinterStatus(printer)
      online = status.online
      printerStatus = status.printerStatus
    })()
  }, 10000);

  const deletePrinter = () => {
    const modal: ModalSettings = {
      buttonTextCancel: 'No',
      buttonTextConfirm: "Yes",
      type: 'confirm',
      title: 'You Sure?',
      body: `Are you sure you want to delete this printer`,
      response: (r) => {
        if (r) {
          (async () => {
            try {
              let url = _apiUrl(`/v1/printer/${printer._id}?rev=${printer._rev}`);
              let res = await fetch(url,{
                method: "DELETE"
              })
              if (!res.ok) {
                console.log(`error: ${res}`);
                errorType = `${res.status}: ${res.statusText}`;
                errorMessage = 'Oops!  There was an error deleting the printer. Response was:' + res;
                errorVisible = true;
              } else {
                goto("/printers")
              }
            } catch (err) {
              console.log(err);
            }
          })()
        }
      }
    };

    modalStore.trigger(modal);
  }

  export const printModel = async (fileResource:string) => {
    try {
      let res:Response = await fetch(fileResource, {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${printer.apiKey}`,
          "Content-Type": "application/json"
        },
        body: '{"command": "select", "print": true}'
      });
      if (!res.ok) {
        console.log(`error: ${res}`);
        let err = new Error()
        err.name = "Print Error"
        err.message = res.statusText
        showError(new Error())
      }
      printerBusy = true
      showJobInfo();
    } catch (err) {
      showError(err)
    }
  }

  let activeJob = false
  let open =false
  let jobInfo:JobInformation;

  const showJobInfo = () => {
    setInterval(function() {
      (async ()=>{
        jobInfo = await GetPrinterJob(printer)
      })()
    }, 5000);
    console.log(jobInfo)
    activeJob = true;
    open = open == false;
  }

  const showError = (err) => {
    console.log(err);
    errorType = err.detail.name;
    errorMessage = err.detail.message;
    errorVisible  = true;
  }
</script>
<!-- Error Div -->
<div>
  {#if errorVisible}
    <aside class="alert variant-filled-error mb-4">
      <!-- Icon -->
      <div><i class="fa-solid fa-triangle-exclamation text" /></div>
      <!-- Message -->
      <div class="alert variant-filled-error alert-message text-sm">
        <div>
          <h3 class="h3">{errorType}</h3>
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
</div>
<div class="container w-1/2 m-auto">
  <div class="flex w-full">
    <div class="w-1/2">
      <h1 class="h1 float-left">{printer.printerName}</h1>
    </div>
    <div class="w-1/2">
      <button type="button" class="btn ml-10 variant-filled-error float-right"  on:click={deletePrinter}>
        <span><i class="fa-solid fa-circle-xmark" /></span>
        <span>Delete Printer</span>
      </button>
    </div>

  </div>
<hr class="!border-t-2 my-4" />
<div class="grid grid-flow-col gap-8">
  <div class="">
    <img src="/mk3s.svg" class="">
  </div>
  <div class="flex flex-col ">
    <div class="">
      <div class="text-xl">Location: {printer.location.name}</div>
    </div>
    <hr class="!border-t-2 my-4" />
    <div class="">
    {#if online}
      <div class="bg-neutral-700 border-l-8 border-lime-600">
        <div class="ml-4 py-2 text-center text-lime-600">ONLINE</div>
      </div>
    {:else}
      <div class="bg-neutral-700 border-l-8 border-red-800">
        <div class="ml-4 py-2 text-center text-red-700">OFFLINE</div>
      </div>
    {/if}
    </div>
    <hr class="!border-t-2 my-4" />
    <div class="bg-neutral-700 border-l-8  border-yellow-600">
      <div class=" ml-4 py-2 text-center text-neutral-100">Status: {printerStatus.state.text}</div>
    </div>
    <hr class="!border-t-2 my-4" />
    <div class="grid grid-cols-3">
      <div class="m-auto">
        <div class="text-s text-center">Extruder Temp</div>
        <RadialGauge name="extruder" bind:value={printerStatus.temperature.tool0.actual} unitSymbol="&#176;" min=0 max=300>
        </RadialGauge>
      </div>
      <div class="m-auto">
        <div class="text-s text-center">Bed Temp</div>
        <RadialGauge name="bed" bind:value={printerStatus.temperature.bed.actual} unitSymbol="&#176;" min=0 max=100>
        </RadialGauge>
      </div>
      <div class="m-auto">
        <div class="text-s text-center">Ambient Temp</div>
        <RadialGauge name="ambient" bind:value={printerStatus.temperature.A.actual} unitSymbol="&#176;"  min=0 max=100>
        </RadialGauge>
      </div>
    </div>
  </div>
</div>
<br/>
{#if (activeJob===true)}
<div id="jobInformationDiv" >
  <Accordion>
    <AccordionItem id="jobInformation" bind:open={open}>
      <svelte:fragment slot="summary"><span class="font-bold h4">Job Information</span></svelte:fragment>
      <svelte:fragment slot="content">
        <div><span class="font-bold">Job:</span> {jobInfo.job.file.name}</div>
        <div><span class="font-bold">Estimated TIme:</span> {SecondsPrettyPrint(jobInfo.job.estimatedPrintTime)}</div>
        <div class="m-auto">
          <div class="text-s pl-6">% Complete</div>
          <RadialGauge name="jobCompletion" bind:value={jobInfo.progress.completion} unitSymbol="%"  min=0 max=100>
          </RadialGauge>
        </div>
      </svelte:fragment>
    </AccordionItem>
  </Accordion>
</div>
  {/if}
<hr class="!border-t-2 my-6" />
<div class=" text-center h2">Printer Files</div>
<hr class="!border-t-2 my-6" />
<div class="w-2/3 m-auto">
  {#await GetPrinterFiles(printer)}
    Getting Printer Files
  {:then data}
    {#if data.printerFiles.children.length > 0}
      {#each data.printerFiles.children as file, i}
        <div class="flex flex-row ">
          <div class="basis-1/12">
            <i class="icon-orange fa-regular fa-cube" />
          </div>
          <div class="basis-9/12">
            <div class="">
              {file.name}
            </div>
            {#if ("gcodeAnalysis" in file)}
              <div class="attributes flex flex-row">
                <div class="basis-1/4">
                  <i class="icon iconfont-layer-height" /><!---->
                  <div>W:{file.gcodeAnalysis.dimensions.width} x L:{file.gcodeAnalysis.dimensions.height} x D:{file.gcodeAnalysis.dimensions.depth}</div>
                </div>
                <div class="basis-1/6">
                  <i class="icon fa-regular fa-clock" />
                  <div>{new Date(file.gcodeAnalysis.estimatedPrintTime * 1000).toISOString().slice(11, 19)}</div>
                </div>
                <div class="basis-1/6">
                  <i class="icon fa-regular fa-balance-scale" />
                  <div>
                    {file.gcodeAnalysis.filament.tool0.volume.toFixed(2)}g
                  </div>
                </div>
              </div>
            {/if}
          </div>
            <div class="basis-1/6">
                <button type="button" class="btn ml-10 variant-filled-error float-right" disabled={printerBusy} on:click={() => {printModel(file.refs.resource)}}>
                  <span><i class="fa-solid fa-print" /></span>
                  <span>Print</span>
                </button>
            </div>
          </div>
        <hr class="!border-t-2 my-2" />
      {/each}
    {:else }
      <h5 class="h6 border rounded variant-ghost-warning text-center py-1 my-6"> -- There are no Print Files on this printer -- </h5>
    {/if}
  {/await}
</div>
</div>
<style>
    .attributes {
        color: #2a2a2a;
        font-size: 11px;
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
</style>
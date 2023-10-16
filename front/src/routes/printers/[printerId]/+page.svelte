<script lang="ts">
  import {type ModalSettings, modalStore } from "@skeletonlabs/skeleton";
  import RadialGauge from "$lib/RadialGauge.svelte";
  import { GetPrinterFiles, CheckPrinterStatus, type PrinterStatus} from "$lib/Printer";
  import {GetPrinterJob, type JobInformation} from "$lib/Job"
  import { _apiUrl, handleError, SecondsPrettyPrint } from "$lib/Utils.js";
  import { goto, invalidateAll } from "$app/navigation";

  export let data;

  let printer= data.printer;
  let online = data.status.online;
  let printerStatus:PrinterStatus = data.status.printerStatus;

  /**
   * The printer is available if:
   *  1. It is ONLINE
   *  AND
   *  2. activeJob == false
   *  https://docs.octoprint.org/en/master/api/datamodel.html#printer-state
   */
  let printerAvailable = true //Is the printer Available
  let activeJob = false //Is there an active job

  let editable = true //will be false if activeJob==true
  let executed = false // Have we executed this code snippet
  let cancelling = false // Was the cancel job button pressed

  let jobInfo:JobInformation = {
    job: {
      averagePrintTime: 0,
      estimatedPrintTime: 0,
      file: {
        date: 0,
        display: "",
        name: "test",
        origin: "",
        path: "",
        size: 0
      }
    },
    progress: {
      completion: 0,
      filepos: 0,
      printTime: 0,
      printTimeLeft: 0,
      printTimeLeftOrigin: ""
    },
    state: ""
  }
  let jobInterval

  let errorType = '';
  let errorMessage = '';
  let errorVisible: boolean = false;


  /**
   * Update Printer status every 10 seconds
   */
  let printerStatusInterval = setInterval(function() {
    (async ()=>{
      let status:{online: string, printerStatus:PrinterStatus, err} = await CheckPrinterStatus(printer)
      if (status.err) {
        console.log(status.err)
        clearInterval(printerStatusInterval)
        errorType = `${status.err.type}: ${status.err.statusText}`;
        errorMessage = 'Oops!  There was an error connecting to the printer. Will try to reconnect on page refresh.'
        errorVisible = true;
        printerAvailable = false
      } else {
        printerStatus = status.printerStatus
        if (status.online == "OFFLINE") {
          printerAvailable = false
        } else if (!printerStatus.state.flags.ready) {
          printerAvailable = false
        }
      }
    })()
  }, 10000);

  /**
   *  Confirm and Delete Printer
   */

  /**
   *  confirmDeletePrinter
   */
  const confirmDeletePrinter = () => {
    const modal: ModalSettings = {
      buttonTextCancel: 'No',
      buttonTextConfirm: "Yes",
      type: 'confirm',
      title: 'You Sure?',
      body: `Are you sure you want to delete this printer`,
      response: (r) => {
        if (r) {
          deletePrinter()
        }
      }
    }
    modalStore.trigger(modal);
  }

  /**
   * deletePrinter
   */
  const deletePrinter = async () => {
    try {
      editable = false
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
  }

  /*
    Confirm and Print Model
   */

  /**
   * confirmPrintModel
   * @param fileResource
   */
  export const confirmPrintModel = (fileResource:string) => {
    const modal: ModalSettings = {
      buttonTextCancel: 'No',
      buttonTextConfirm: "Yes",
      type: 'confirm',
      title: 'You Sure?',
      body: `Are you sure you want to print this model`,
      response: (r) => {
        if (r) {
          printModel(fileResource)
        }
      }
    }
    modalStore.trigger(modal);
  }

  /**
   * printModel
   * @param fileResource
   */
  export const printModel = async (fileResource:string) => {
    try {
      let res:Response = await fetch(fileResource, {
        method: "POST",
        headers: {
          "X-Api-Key": printer.apiKey,
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
        printerAvailable = true
        editable = true
      } else {
        printerAvailable = false
        editable = false
        showJobInfo();
      }
    } catch (err) {
      showError(err)
    }
  }

  /**
   * Cancels a Print Job.  Triggered from the cancel job button
   */
  const cancelPrintJob = async () => {
    cancelling = true
    try {
      let res:Response = await fetch(`${printer.url}/api/job`, {
        method: "POST",
          headers: {
            "X-Api-Key": printer.apiKey,
            "Content-Type": "application/json"
        },
        body: '{"command": "cancel"}'
      })

      if (!res.ok) {
        showError(res)
      } else if (res.status == 204){ // Successful Job cancel
        console.log("Cancelling Print Job")
        // Reset Temps just in case
        let res:Response = await fetch(`${printer.url}/api/printer/command`, {
          method: "POST",
          headers: {
            "X-Api-Key": printer.apiKey,
            "Content-Type": "application/json"
          },
          body: '{"commands": ["M104 S0", "M140 S0"]}'
        })
        if (res.status == 204) { //Successful Temp reset
          console.log("temps reset")
          clearInterval(jobInterval)
          printerAvailable = true
          editable = true
          activeJob = false
          executed = false
          cancelling = false
        }
      }
    } catch (error) {
      console.log(error)
      showError(error)
    }
  }

  /**
   * showJobInfo
   * Get the current job info form the printer and displays it
   * Then repeats every 10 seconds
   */
  clearInterval(jobInterval)
  const showJobInfo = async () => {
    jobInfo = await GetPrinterJob(printer)
    if (executed == false){
      console.log(jobInfo.state)
      executed = true
      activeJob = true
    }
    console.log(jobInfo)
    if(jobInfo.state != "Printing"){
      console.log("Cancelling Job Interval")
      clearTimeout(jobInterval)
      printerAvailable = true
      activeJob = false
      executed = false
      return
    }
    jobInterval = setTimeout(showJobInfo, 10000)
  }

  const showError = (err) => {
    console.log(err);
    errorType = err.detail.name;
    errorMessage = err.detail.message;
    errorVisible  = true;
  }

  const reloadPage = async() => {
    await invalidateAll()
    printer = data.printer;
  }

  //Save Printer
  let saveDisabled = true;
  function needsSave() {
    saveDisabled = false;
  }

  /**
   * showUpdated
   * Modal which shows that the printer was updated
   * @param title
   * @param body
   * @param reload
   */
  export const showUpdated = (title:string, body:string, reload: boolean) => {
    const modal: ModalSettings = {
      type: 'alert',
      title: title,
      body: body,
      buttonTextCancel: 'Cool!',
      //response: (e) => { invalidateAll() }
    };

    if (reload) {
      modal.response = () => { reloadPage()}
    }
    modalStore.trigger(modal);
  };

  /**
   * updatePrinter
   */
  const updatePrinter = async () => {
    //console.log(model)
    const response = await fetch(_apiUrl(`/v1/printer/${printer._id}`), {
      method: 'PUT',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(printer)
    }).then(handleError) // skips to .catch if error is thrown
      .then((response) => {
        printer._rev = response.rev;
        showUpdated('Complete', 'Printer has been successfully Updated', false);
        saveDisabled = true
      }).catch((error) => {
        let errorMessage = 'Oops!  There was an error updating the printer.<br/>Response was: ' + error;
        showUpdated(error, errorMessage, true);
      })
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
      <button type="button" class="btn ml-10 variant-filled-error float-right"  disabled={!printerAvailable} on:click={confirmDeletePrinter}>
        <span><i class="fa-solid fa-circle-xmark" /></span>
        <span>Delete Printer</span>
      </button>
    </div>
  </div>
<hr class="!border-t-2 my-4" />
<div class="grid grid-flow-col gap-8">
  <div class="">
    <img src="/mk3s.svg" alt="Printer Image" class="">
  </div>
  <div class="flex flex-col ">
    <div>
      <span class="h6 text-xs float-right">*click to edit</span>
    </div>
    <div class="">
      <span class="h4 mr-2">URL:</span>
      <span contenteditable=true bind:textContent={printer.url} on:input={needsSave} class="editable p-1 pr-10">
        {printer.url}
      </span>
    </div>
    <div class="">
      <span class="h4 mr-2">Api Key:</span>
      <span contenteditable=true bind:textContent={printer.apiKey} on:input={needsSave} class="editable p-1 pr-10">
        {printer.apiKey}
      </span>
    </div>
    <div class="">
      <span class="h4 mr-2">API Type:</span>
      <span>{printer.apiType}</span>
    </div>
    <div class="">
      <span class="h4 mr-2">Location:</span>
      <span contenteditable="true" bind:textContent={printer.location.name} on:input={needsSave} class="editable p-1 pr-10">
        {printer.location.name}
      </span>
    </div>
    <div class="">
      <span class="h4 mr-2">AutoConnect:</span>
      <input class="input checkbox" type="checkbox" name="autoConnect" bind:checked={printer.autoConnect} on:input={needsSave} bind:value={printer.autoConnect}/>
    </div>
    <div class="">
      <button disabled={saveDisabled} type="button" class="btn variant-filled-error float-right" on:click={updatePrinter}>
        <span><i class="fa-regular fa-floppy-disk"></i></span>
        <span>Save Changes</span>
      </button>
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
        <RadialGauge name="extruder" bind:value={printerStatus.temperature.tool0.actual} unitSymbol="&#176;" min={0} max={300}>
        </RadialGauge>
      </div>
      <div class="m-auto">
        <div class="text-s text-center">Bed Temp</div>
        <RadialGauge name="bed" bind:value={printerStatus.temperature.bed.actual} unitSymbol="&#176;" min={0} max={100}>
        </RadialGauge>
      </div>
      <div class="m-auto">
        <div class="text-s text-center">Ambient Temp</div>
        <RadialGauge name="ambient" bind:value={printerStatus.temperature.A.actual} unitSymbol="&#176;"  min={0} max={100}>
        </RadialGauge>
      </div>
    </div>
  </div>
</div>
<br/>
  {#if activeJob}
    <hr class="!border-t-2 my-6" />
    <div class="variant-ghost-secondary p-2 mb-2 w-full rounded flex">
      <div class="font-bold h4 ml-2 w-1/2">Job Information</div>
      <div class="w-1/2">
        <button type="button" class="btn variant-filled-error float-right" disabled={cancelling} on:click={cancelPrintJob}>
          <span><i class="fa-solid fa-circle-xmark" /></span>
          <span>Cancel Print</span>
        </button>
      </div>
    </div>
    {#if cancelling}
      <div class=" py-3 mb-3 variant-ghost-warning text-center"><span class="text-xl font-bold text-red-800">Cancelling Print Job</span></div>
      {/if}
    <div class="grid grid-cols-2 mx-auto w-fit" >
      <div class="px-4">
        <div>
          <span class="font-bold h5">Job: </span>{jobInfo.job.file.name}
        </div>
        <div>
          <span class="font-bold">Estimated Time:</span> {SecondsPrettyPrint(jobInfo.job.estimatedPrintTime)}
        </div>
      </div>
      <div class="m-auto">
        <div class="text-s pl-6">% Complete</div>
        <RadialGauge name="jobCompletion" bind:value={jobInfo.progress.completion} unitSymbol="%"  min={0} max={100}>
        </RadialGauge>
      </div>
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
                <button type="button" class="btn ml-10 variant-filled-error float-right" disabled={(!printerAvailable || activeJob)} on:click={() =>
                {confirmPrintModel(file.refs.resource)}}>
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
        margin-bottom: 0;
    }

    input[type="file"] {
        display: none;
    }

    .editable:hover,
    [contenteditable="true"]:active,
    [contenteditable="true"]:focus{
        background: rgb(211, 211, 211);
        border: 1px solid rgb(133,133,133);
        border-radius: 4px;
        outline:none;
        padding: 8px;
    }
</style>
<script lang="ts">
  export let printer = {
    _id: '',
    _rev: '',
    printerName: '',
    url: "",
    apiKey: "",
    location: {"name": ""},
    tags: []
  };

  let online = false;
  let printerStatus: PrinterStatus
  type PrinterStatus ={
    "sd": {
      "ready": boolean
    },
    "state": {
      "error": string,
      "flags": {
        "cancelling": boolean,
        "closedOrError": boolean,
        "error": boolean,
        "finishing": boolean,
        "operational": boolean,
        "paused": boolean,
        "pausing": boolean,
        "printing": boolean,
        "ready": boolean,
        "resuming": boolean,
        "sdReady": boolean
      },
      "text": string
    },
    "temperature" :{
      "A": {
        "actual": number,
        "offset": number,
        "target": number
      },
      "P": {
        "actual": number,
        "offset": number,
        "target": number
      },
      "bed": {
        "actual": number,
        "offset": number,
        "target": number
      },
      "tool0": {
        "actual": number,
        "offset": number,
        "target": number
      }
    }
  }

  const checkPrinterStatus = async () => {
    try {
      let res = await fetch(`${printer.url}/api/printer`, {
        headers: { Authorization: `Bearer ${printer.apiKey}` }
      });
      if (!res.ok) {
        console.log(`error: ${res}`)
      }
      online = true;
      printerStatus = await res.json()
      return
    } catch (err) {
      if (err.message === "Failed to fetch") {
        online = false;
      }
    }
  }

</script>

{#await checkPrinterStatus()}

{:then _}
  <div class="h-[64px] my-8 bg-surface-200 flex">
    <div class="">
      <a href="">
        <img src="/mk3s.svg" class=" printer-img">
      </a>
    </div>
    <div class="my-auto">
      <a href="/printers/{printer._id}">
        <span title="ymir" class="font-semibold text-xl">fgjwskpgjirgjj {printer.printerName}</span>
      </a>
    </div>
      {#if online}
        <div class="bg-neutral-700 m-auto border-l-4 border-lime-600 w-32">
          <div class="text-center text-lime-600">ONLINE</div>
        </div>
      {:else}
        <div class="bg-neutral-700 border-l-4 border-red-800 w-32">
          <div class="text-center mt-[4px] text-red-700">OFFLINE</div>
        </div>
      {/if}
    <div class="bg-neutral-700 border-l-4 m-auto border-black w-32">
      <div class="text-center text-neutral-100">Status:<br/>{printerStatus.state.text}</div>
    </div>
    <div class="m-auto">
      <div class="text-xs">Bed Temp</div>
      <div class="">{printerStatus.temperature.bed.actual}<span>&#176;</span></div>
    </div>
    <div class="m-auto">
      <div class="text-xs">Location</div>
      <div class="">{printer.location.name}</div>
    </div>
  </div>
{/await}

<style>
  .printer-img{
      height: 60px;
      width: 60px !important;
  }

</style>
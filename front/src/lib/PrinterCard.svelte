<script lang="ts">
  import type {PrinterStatus, Printer} from '$lib/Printer'
  import {CheckPrinterStatus} from "$lib/Printer"

  export let printer:Printer

</script>

{#await CheckPrinterStatus(printer)}

{:then status}
  {#if status.printerStatus}
    <div class="h-[64px] my-8 bg-surface-200 grid grid-cols-8">
    <div class="px-4">
      <a href="">
        <img src="/mk3s.svg" class=" printer-img">
      </a>
    </div>
    <div class="my-auto">
      <a href="/printers/{printer._id}">
        <span title="ymir" class="font-semibold text-xl">{printer.printerName}</span>
      </a>
    </div>
    <div class="bg-neutral-700 m-auto border-l-4 border-lime-600 myDiv">
      <span class="text-lime-600 middle px-10">{printer.url}<br/>{status.online}</span>
    </div>
    <div class="bg-neutral-700 m-auto border-l-4 myDiv border-yellow-600">
      <span class="text-neutral-100 middle px-10">Status: {status.printerStatus.state.text}</span>
    </div>
    <div class="m-auto">
      <div class="text-s">Bed Temp</div>
      <div class="">{status.printerStatus.temperature.bed.actual}<span>&#176;</span></div>
    </div>
    <div class="m-auto">
      <div class="text-s">Extruder Temp</div>
      <div class="">{status.printerStatus.temperature.tool0.actual}<span>&#176;</span></div>
    </div>
    <div class="m-auto">
      <div class="text-s">Ambient Temp</div>
      <div class="">{status.printerStatus.temperature.A.actual}<span>&#176;</span></div>
    </div>
    <div class="m-auto">
      <div class="text-s">Location</div>
      <div class="">{printer.location.name}</div>
    </div>
  </div>
  {:else }
    <div class="h-[64px] my-8 bg-surface-200 grid grid-cols-8">
      <div class="px-4">
        <a href="">
          <img src="/mk3s.svg" class=" printer-img">
        </a>
      </div>
      <div class="my-auto">
        <a href="/printers/{printer._id}">
          <span title="ymir" class="font-semibold text-xl">{printer.printerName}</span>
        </a>
      </div>
      <div class="bg-neutral-700 m-auto border-l-4 border-red-800 myDiv">
        <span class="text-red-700 middle px-10">{printer.url}<br/>{status.online}</span>
      </div>
      <div class="bg-neutral-700 m-auto border-l-4 myDiv border-yellow-600">
        <span class="text-neutral-100 middle px-10">Status: UNKNOWN</span>
      </div>
      <div class="m-auto">
        <div class="text-s">Bed Temp</div>
        <div class="">Unknown</div>
      </div>
      <div class="m-auto">
        <div class="text-s">Extruder Temp</div>
        <div class="">Unknown</div>
      </div>
      <div class="m-auto">
        <div class="text-s">Ambient Temp</div>
        <div class="">Unknown</div>
      </div>
      <div class="m-auto">
        <div class="text-s">Location</div>
        <div class="">{printer.location.name}</div>
      </div>
    </div>
  {/if}
{/await}

<style>
  .printer-img{
      height: 60px;
      width: 60px !important;
  }

  .myDiv {
      height: 64px;
      line-height: 64px;
      text-align: center;
  }

  .middle {
      display: inline-block;
      vertical-align: middle;
      line-height: normal;
  }

</style>
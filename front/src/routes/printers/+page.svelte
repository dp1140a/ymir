<script>
  import PrinterCard from "$lib/PrinterCard.svelte"
  import { Table } from '@skeletonlabs/skeleton';
</script>

<div class="h1">Printers</div>
<a class="btn btn-sm variant-filled-error my-4" href="/printers/add"> + New Printer </a>


<div>
<PrinterCard></PrinterCard>

  <div class="">
    <table class="table-auto w-1/2">
      <tbody>
      <tr class="">
        <td class="border border-slate-800 w-90px">
          <a class="" href=".">
            <img src="/mk3s.svg" class=" printer-img" alt="printer">
          </a>
        </td>
        <td class="border border-slate-800">
          <a class="m-0 p-0" href=".">
            <div title="ymir" class="font-semibold text-lg">ymir</div>
          </a>
        </td>
        <td class="border border-slate-800">
          <div class="bg-surface-400  max-h-7 border-l-4 border-gray-700 ">
            <div class="text-center mt-[4px] text-white">offline</div>
          </div>
        </td>
        <td class="border border-slate-800">
          <div class="pl-6">
            <div class="text-xs">Location</div>
            <div class="">Valhalla</div>
          </div>
        </td>
        <td class="border border-slate-800">
          <a class="btn btn-sm variant-filled-error" href=".">Connect</a>
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</div>
<style>
    .printer-img{
        height: 60px;
        width: 60px
    }
</style>

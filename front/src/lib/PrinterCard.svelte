<script>
import { expoIn } from "svelte/easing";
import { base } from "$app/paths";

const printer ={"displayName": "Test", "id": "123456789"}


</script>
<div class="flex mb-8 my-4 gap-6 bg-surface-200 w-1/2">
  <a class="" href=".">
    <img src="/mk3s.svg" class=" printer-img" alt="printer icon">
  </a>
    <div class="flex flex-row gap-8 mt-1">
      <div class="mt-3">
        <a class="" href=".">
          <div title="ymir" class="font-semibold text-lg">ymir</div>
        </a>
      </div>
      <div class="bg-surface-400 max-h-7 mt-3 border-l-4 border-gray-700 w-48">
          <div class="text-center mt-[4px] text-white">offline</div>
      </div>
      <div class="mt-2">
        <div class="text-xs">Location</div>
        <div class="">Valhalla</div>
      </div>
  </div>
</div>


<style>
  .printer-img{
      height: 60px;
      width: 60px
  }


</style>
<script lang="ts">
	import type { Printer, PrinterStatus } from '$lib/Printer';
	import { CheckPrinterStatus } from '$lib/Printer';

	export let printer: Printer;

	let status: { online: string; printerStatus: PrinterStatus; err: Error } = {
		online: '',
		printerStatus: null,
		err: null
	};
	status.online = 'Offline';
	status.printerStatus = {
		state: { text: 'Unknown' },
		temperature: { bed: { actual: 0 }, tool0: { actual: 0 }, A: { actual: 0 } }
	};

	$: CheckPrinterStatus(printer).then((data) => {
		status = data;
	});
</script>

<div class="my-8 grid h-[64px] grid-cols-8 bg-surface-200">
	<div class="px-4">
		<a href={'#'}>
			<img src="/mk3s.svg" class=" printer-img" alt="mk3s.svg" />
		</a>
	</div>
	<div class="my-auto">
		<a href="/printers/{printer._id}">
			<span title="ymir" class="text-xl font-semibold">{printer.printerName}</span>
		</a>
	</div>
	<div class="myDiv m-auto border-l-4 border-lime-600 bg-neutral-700">
		<span class="middle px-10 text-lime-600">{printer.url}<br />{status.online}</span>
	</div>
	<div class="myDiv m-auto border-l-4 border-yellow-600 bg-neutral-700">
		<span class="middle px-10 text-neutral-100">Status: {status.printerStatus.state.text}</span>
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

<style>
	.printer-img {
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

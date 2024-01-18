/** @type {import('./$types').PageLoad} */
import { _apiUrl } from '$lib/Utils';
import type { Printer, PrinterStatus } from '$lib/Printer';
import { CheckPrinterStatus } from '$lib/Printer';

export const load = async ({ fetch, params }) => {
	const url: string = _apiUrl(`/v1/printer/${params.printerId}`);
	const res = await fetch(url);
	if (!res.ok) {
		throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
	}
	const printer: Printer = await res.json();
	const status: { online: string; printerStatus: PrinterStatus } = await CheckPrinterStatus(printer);

	return { printer, status };
};

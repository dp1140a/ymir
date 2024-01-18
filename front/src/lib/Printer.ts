import { writable } from 'svelte/store';
import { _apiUrl } from '$lib/Utils';

export interface Printer {
	_id: string;
	_rev: string;
	printerName: string;
	url: string;
	apiType: string;
	apiKey: string;
	location: PrinterLocation;
	type: PrinterType;
	dateAdded: string;
	tags: string[];
	autoConnect: boolean;
}

export interface PrinterLocation {
	name: string;
}

export interface PrinterType {
	Make: string;
	Model: string;
	Version: string;
}

export interface PrinterStatus {
	sd: {
		ready: boolean;
	};
	state: {
		error: string;
		flags: {
			cancelling: boolean;
			closedOrError: boolean;
			error: boolean;
			finishing: boolean;
			operational: boolean;
			paused: boolean;
			pausing: boolean;
			printing: boolean;
			ready: boolean;
			resuming: boolean;
			sdReady: boolean;
		};
		text: string;
	};
	temperature: {
		A: {
			actual: number;
			offset: number;
			target: number;
		};
		P: {
			actual: number;
			offset: number;
			target: number;
		};
		bed: {
			actual: number;
			offset: number;
			target: number;
		};
		tool0: {
			actual: number;
			offset: number;
			target: number;
		};
	};
}

export const SelectedPrinter = writable<Printer>();

export const Connect = async (printer: Printer): Promise<boolean> => {
	try {
		const res: Response = await fetch(`${printer.url}/api/printer`, {
			method: 'POST',
			body: "{'command': 'connect', 'autoconnect': true}",
			headers: {
				'X-Api-Key': printer.apiKey,
				'content-type': 'application/json'
			}
		});
		if (res.status == 400) {
			console.log(`400: You sent a request that this server could not understand.`);
			return false;
		} else if (res.status == 204) {
			return true;
		}
	} catch (error) {
		console.log(error);
		return false;
	}
};

export const CheckPrinterStatus = async function (
	printer: Printer
): Promise<{ online: string; printerStatus: PrinterStatus; err: Response }> {
	let printerStatus: PrinterStatus;
	let online: string;
	let err: Response;
	try {
		const res: Response = await fetch(`${printer.url}/api/printer`, {
			headers: {
				'X-Api-Key': printer.apiKey
			}
		});
		if (!res.ok) {
			console.log(`error: ${res.status}`);
			console.log(res);
			err = res;
			if (res.status == 403) {
				online = 'FORBIDDEN';
			} else {
				online = 'OFFLINE';
			}
			if (res.status == 409 && printer.autoConnect) {
				let connected: boolean;
				let connectTimeout: number;
				while (!connected) {
					console.log('attempting to reconnect');
					connectTimeout = setTimeout(Connect, 5000);
				}
				clearTimeout(connectTimeout);
			}
		} else {
			online = 'ONLINE';
			printerStatus = await res.json();
		}
	} catch (error) {
		console.log(error);
		if (error.message === 'Failed to fetch') {
			online = 'OFFLINE';
			err = error;
		}
	}
	return { online, printerStatus, err };
};

export const GetPrinterFiles = async (printer: Printer) => {
	let printerFiles;
	try {
		const res: Response = await fetch(`${printer.url}/api/files/local/ymir`, {
			headers: {
				'X-Api-Key': printer.apiKey
			}
		});
		if (!res.ok) {
			console.log(`error: ${res}`);
		}
		printerFiles = await res.json();
		//console.log(printerFiles)
	} catch (err) {
		console.log(err);
	}
	return { printerFiles };
};

export const UploadAndPrintFile = async (filePath: string, printer: Printer, print: boolean) => {
	let resBody;
	let error: Error;
	try {
		const res: Response = await fetch(
			_apiUrl('/v1/model/file/printer?file=').concat(filePath, '&print=', print.toString()),
			{
				method: 'POST',
				body: JSON.stringify(printer),
				headers: { 'content-type': 'application/json' }
			}
		);
		if (!res.ok) {
			error = await res.json();
			error = new Error(`error printing file. Response is: ${error}`);
		} else {
			resBody = await res.json();
		}
	} catch (err) {
		error = new Error(`error printing file.  response is: ${err}`);
	}
	return { resBody, error };
};

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

export const CheckPrinterStatus = async function (printer: Printer) {
	let printerStatus: PrinterStatus;
	let online: string;
	try {
		let res: Response = await fetch(`${printer.url}/api/printer`, {
			headers: { Authorization: `Bearer ${printer.apiKey}` }
		});
		if (!res.ok) {
			console.log(`error: ${res.status}`);
			if (res.status == 403) {
				online = "FORBIDDEN"
			} else {
				online = "OFFLINE"
			}
		} else {
			online = "ONLINE";
			printerStatus = await res.json();
		}

	} catch (err) {
		if (err.message === 'Failed to fetch') {
			online = "OFFLINE";
		}
	}
	return { online, printerStatus };
};

export const GetPrinterFiles = async (printer: Printer) => {
	let printerFiles;
	try {
		let res: Response = await fetch(`${printer.url}/api/files/local/ymir`, {
			headers: { Authorization: `Bearer ${printer.apiKey}` }
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

export const GetPrinterJob = async (printer: Printer) => {
	let jobInfo;
	try {
		let res: Response = await fetch(`${printer.url}/api/job`, {
			headers: { Authorization: `Bearer ${printer.apiKey}` }
		});
		if (!res.ok) {
			console.log(`error: ${res}`);
		}
		jobInfo = await res.json();
		//console.log(printerFiles)
	} catch (err) {
		console.log(err);
	}
	return { jobInfo };
};


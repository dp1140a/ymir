import type { Printer } from '$lib/Printer';

export interface JobInformation {
	job: {
		averagePrintTime: number;
		estimatedPrintTime: number;
		filament: {
			tool0: {
				length: number;
				volume: number;
			};
		};
		file: {
			date: number;
			display: string;
			name: string;
			origin: string;
			path: string;
			size: number;
		};
		lastPrintTime: number;
		user: string;
	};
	progress: {
		completion: number;
		filepos: number;
		printTime: number;
		printTimeLeft: number;
		printTimeLeftOrigin: string;
	};
	state: string;
}

export const GetPrinterJob = async (printer: Printer) => {
	let jobInfo: JobInformation;
	try {
		const res: Response = await fetch(`${printer.url}/api/job`, {
			headers: {
				'X-Api-Key': printer.apiKey
			}
		});
		if (!res.ok) {
			console.log(`error: ${res}`);
		}
		jobInfo = await res.json();
	} catch (err) {
		console.log(err);
	}
	return jobInfo;
};

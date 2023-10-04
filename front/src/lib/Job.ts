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

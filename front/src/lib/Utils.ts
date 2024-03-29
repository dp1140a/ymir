/**
 * Common Functions for use across pages
 */

/**
 * Checks for errors on a fetch response.  Throws an error if response not ok
 * @param response
 */
export const handleError = (response: Response): Promise<unknown> => {
	if (!response.ok) {
		throw Error(response.statusText);
	}
	return response.json();
};

/**
 * Modifies path to include server host and port if in dev mode
 * @param path
 */

export const _apiUrl = (path: string): string => {
	//console.log(`${import.meta.env.VITE_API_URL}`);
	let base = '';
	if (import.meta.env.DEV) {
		base = import.meta.env.VITE_API_URL;
	}
	return `${base}${path}`;
};

/**
 * Returns the last element in a path string
 * @param path
 * @constructor
 */
export const Basename = (path: string): string => {
	return path.split('/').reverse()[0];
};

export const SecondsPrettyPrint = (totalSeconds: number): string => {
	const hours: number = Math.floor(totalSeconds / 3600);
	totalSeconds %= 3600;
	const minutes: number = Math.floor(totalSeconds / 60);
	const seconds: string = (totalSeconds % 60).toFixed(0);

	return `${hours}H ${minutes}M ${seconds}S`;
};

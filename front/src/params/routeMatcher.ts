// 1b8981b5b9de3e3a59a12176300020a4
export function match(param: string) {
	//[0-9a-f]{32}
	//return /[0-9a-f]{32}/.test(param);
	return /^[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$/.test(param);
}

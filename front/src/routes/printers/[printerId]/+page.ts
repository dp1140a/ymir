/** @type {import('./$types').PageLoad} */
import { _apiUrl } from "../../+layout";


export const load = async ({ fetch, params }) => {
  let url = _apiUrl(`/v1/printer/${params.printerId}`);
  let res = await fetch(url);
  if (!res.ok) {
    throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
  }
  const printer = await res.json();
  console.log(printer)

  return {printer}
}
import {_apiUrl} from "../+layout";

export const load = async ({ fetch, params }) => {
  const url = _apiUrl('/v1/printer');
  let res = await fetch(url);
  if (!res.ok) {
    throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
  }
  const printers = await res.json();
  console.log(printers)
  return { url, printers };
}
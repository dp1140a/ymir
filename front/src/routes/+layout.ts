// this is needed to give us force prerendering of all pages
export const prerender = false;
export const ssr = false;

export const _apiUrl = (path: string) => {
  //console.log(`${import.meta.env.VITE_API_URL}`);
  let base = ""
  if (import.meta.env.DEV) {
    base = import.meta.env.VITE_API_URL
  }
  return `${base}${path}`;
};
<!-- src/routes/+page.svelte -->
<script lang="ts">
	import { _apiUrl } from '$lib/Utils';
	import type { Version } from '$lib/Version';

	const getVersion = async () => {
		const url = _apiUrl('/v1/version');
		let version: Version = {
			AppName: '',
			Version: 'dev',
			Branch: '',
			Commit: '',
			BuildTime: ''
		};
		try {
			const res = await fetch(url);
			if (!res.ok) {
				console.log('bad shit');
				throw `Error while fetching data from ${url} (${res.status} ${res.statusText}).`;
			} else {
				version = await res.json();
			}
		} catch (err) {
			console.log(err);
		}
		console.log(version);
		return version.Version;
	};
</script>

{#await getVersion()}
	loading...
{:then Version}
	<div class="container mx-auto space-y-8 p-8">
		<span class="justify-right flex">
			<a href="about/ymir" class="logo h1">ᛃᛗᛁᚱ</a>
			<h1 class="h1">- Ymir <span class="italic">(ˈē-ˌmir')</span></h1>
		</span>
		<h3 class="h3">Version from Server: {Version}</h3>
		<h3 class="h3">Description:</h3>
		<p>
			Ymir is a 3D model manager. In a nutshell it is a light and local version of the
			printables.com website.
		</p>
		<h3 class="h3">Why:</h3>
		<p>
			While I love what the folks at Prusa (and Free3D) are doing, I created Ymir for a few reasons.
			First I am selfish and dont want to share my models with the whole world. Second I wanted a
			model manager that can tie into my printer and print directly to it (yes I know about
			Octoprint. More on that later). I also did not want it to be tied to a specific printer brand
			like printables.com is. Right now Ymir will only connect to printers using OctoPrint.
		</p>
		<p>
			I get why Prusa does that and Im cool with it. In fact I own a Prusa printer myself. But a lot
			of folks dont. Plus Im pretty sure my next printer is gonna be a Voron. Lastly I wanted
			something I could run locally and not have to create an account on.
		</p>
		<p class="mb-4">
			Now about Octoprint. I have it installed on my Prusa Mk3S, and its pretty cool. I can print
			over wifi and do other all kinds of stuff without getting my lazy ass out of my chair. And yes
			I know Octoprint has a model manager built in. But it just didnt do it for me. Octoprint is
			not a model manager first and I needed some management capabilities that Octoprint couldn't
			provide. And lets not forget that the manager plugin is crammed into the left side of the
			screen. I want to see my models in their HD fullscreen glory. So at this point you are
			probably gonna ask why I just didnt create an Octoprint plugin. Thats an easy one to answer: I
			dont Like Python and I know enough Python to hurt myself. I use Python for data manipulation
			and data science and that is where it shall stay in my world. My language of choice for
			anything backend is currently GoLang (im eyeballing Rust but that's long term).
		</p>
		<h3 class="h3">Design:</h3>
		<p>In its current version, Ymir is built with:</p>
		<ul class="myUL list-disc pl-8">
			<li>Golang</li>
			<li><a href="https://svelte.dev/">Svelte</a></li>
			<li><a href="https://kit.svelte.dev/">Sveltekit</a></li>
			<li><a href="https://www.skeleton.dev/">SkeletonUI</a></li>
			<li><a href="https://tailwindcss.com/">TailwindCSS</a></li>
		</ul>
		<h3 class="h3">License:</h3>
		<a href="https://opensource.org/licenses/gpl-3-0"
			><img src="https://img.shields.io/badge/License-gplv3-red.svg" alt="gpl3 license logo" /></a
		>
	</div>
{:catch err}
	{err}
{/await}

<style>
	.myUL {
		margin-top: 4px !important;
	}
</style>

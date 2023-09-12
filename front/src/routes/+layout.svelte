<script lang="ts">
	import { base } from '$app/paths';
	import { env } from '$env/dynamic/public';

	// The ordering of these imports is critical to your app working properly
	import '../theme.postcss';
	import '@skeletonlabs/skeleton/styles/skeleton.css';
	import { TabGroup, Tab, TabAnchor } from '@skeletonlabs/skeleton';
	import '../app.postcss';
	import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
	import { Modal, storePopup } from '@skeletonlabs/skeleton';
	import { AppShell, AppBar } from '@skeletonlabs/skeleton';
	import { popup } from '@skeletonlabs/skeleton';
	import type { PopupSettings } from '@skeletonlabs/skeleton';

	storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });
	let comboboxValue: string;

	function handleClick() {
		console.log(comboboxValue);
	}

	const popupCombobox: PopupSettings = {
		event: 'focus-click',
		target: 'popupCombobox',
		placement: 'bottom',
		closeQuery: '.listbox-item'
	};

	//export let data

	let tabsBottomNav = '/models';

	console.log("DEV: " + import.meta.env.DEV)
	console.log("PROD: " + import.meta.env.PROD)
	console.log("BASE URL: " + import.meta.env.BASE_URL)
	console.log("MODE: " + import.meta.env.MODE)

</script>

<Modal />
<!-- App Shell -->
<AppShell>
	<svelte:fragment slot="header">
		<!-- App Bar -->
		<AppBar
			gridColumns="grid-cols-3"
			slotDefault="place-self-center"
			slotTrail="place-content-end"
			padding="px-4 pt-4"
		>
			<svelte:fragment slot="lead">
				<div class="flex justify-right">
					<h1 class="h1 gradient">ᛃᛗᛁᚱ</h1>
				</div>
			</svelte:fragment>
			<div>
				<TabGroup
					justify="justify-center"
					active="variant-filled-error"
					hover="hover:variant-soft-error"
					flex="flex-1 lg:flex-none"
					rounded="rounded-tl-container-token rounded-tr-container-token"
					border="border-b border-error-400-500-token"
					class="bg-surface-100-800-token w-full"
				>
					<TabAnchor
						href="/models"
						selected={tabsBottomNav === '/models'}
						on:click={() => (tabsBottomNav = '/models')}
					>
						<svelte:fragment slot="lead"><i class="fa-solid fa-cubes" /></svelte:fragment>
						Models
					</TabAnchor>
					<TabAnchor
						href="/printers"
						selected={tabsBottomNav === '/printers'}
						on:click={() => (tabsBottomNav = '/printers')}
					>
						<svelte:fragment slot="lead"><i class="fa-solid fa-print" /></svelte:fragment>
						Printer
					</TabAnchor>
					<TabAnchor
						href="/docs"
						selected={tabsBottomNav === '/docs'}
						on:click={() => (tabsBottomNav = '/docs')}
					>
						<svelte:fragment slot="lead"><i class="fa-solid fa-book" /></svelte:fragment>
						Docs
					</TabAnchor>
				</TabGroup>
			</div>
			<svelte:fragment slot="trail">
				<a class="btn btn-sm variant-filled-error" href="/models/create"> + New Model </a>
				<button
					class="btn hover:variant-soft-error"
					use:popup={{ event: 'click', target: 'features' }}
				>
					<span>Explore</span>
					<i class="fa-solid fa-caret-down opacity-50" />
				</button>
				<!-- popup -->
				<div class="card p-4 w-60 shadow-xl" data-popup="features">
					<nav class="list-nav">
						<ul>
							<li class="hover:variant-ghost-error">
								<a href="{base}/docs">
									<span>Docs</span>
								</a>
							</li>
							<li class="hover:variant-ghost-error">
								<a href="{base}/about">
									<span>About</span>
								</a>
							</li>
							<li class="hover:variant-ghost-error">
								<a href="http://github.com/dp1140a" target="_blank">
									<span><i class="fa-brands fa-github" /> Github</span></a
								>
							</li>
						</ul>
					</nav>
					<div class="arrow bg-surface-100-800-token" />
				</div>
			</svelte:fragment>
		</AppBar>
	</svelte:fragment>
	<!-- Page Route Content -->
	<div id="content" class="w-fit mx-auto px-8 pt-4">
		<slot />
	</div>
</AppShell>

<!--
https://codepen.io/marcobiedermann/pen/GRmxWwK
Also look at :
https://codepen.io/atnyman/pen/nmEyjK
https://codepen.io/mandymichael/pen/wpYQKx
-->
<style lang="scss">
	h1.gradient {
		background-clip: text;
		background-image: linear-gradient(
			to top,
			#511d11,
			#632415,
			#7c2d1a,
			#953620,
			#a53c23,
			#c07765,
			#dbb1a7
		);
		color: #511d11;
		font-weight: 400;
		letter-spacing: calc(1em / 160);
		/*padding: calc(calc(1em / 16));*/
		-webkit-text-stroke-color: transparent;
		-webkit-text-stroke-width: calc(1em / 12);
		font-size: 64px;
	}
</style>

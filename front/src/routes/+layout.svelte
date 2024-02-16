<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { initializeStores } from '@skeletonlabs/skeleton';
	import { TabGroup, TabAnchor } from '@skeletonlabs/skeleton';
	import '../app.postcss';
	import { computePosition, autoUpdate, offset, shift, flip, arrow } from '@floating-ui/dom';
	import { Modal, storePopup } from '@skeletonlabs/skeleton';
	import { AppShell, AppBar } from '@skeletonlabs/skeleton';
	import { popup } from '@skeletonlabs/skeleton';

	initializeStores();

	storePopup.set({ computePosition, autoUpdate, offset, shift, flip, arrow });

	//let tabsBottomNav = '/models';

	console.log('DEV: ' + import.meta.env.DEV);
	console.log('PROD: ' + import.meta.env.PROD);
	console.log('BASE URL: ' + import.meta.env.BASE_URL);
	console.log('MODE: ' + import.meta.env.MODE);
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
				<div class="justify-right flex">
					<a href="/" class="logo h1">ᛃᛗᛁᚱ</a>
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
					<!--selected={tabsBottomNav === '/models'}-->
					<TabAnchor href="/models" selected={$page.url.pathname.includes('/models')}>
						<svelte:fragment slot="lead"><i class="fa-solid fa-cubes" /></svelte:fragment>
						Models
					</TabAnchor>
					<TabAnchor href="/printers" selected={$page.url.pathname.includes('/printers')}>
						<svelte:fragment slot="lead"><i class="fa-solid fa-print" /></svelte:fragment>
						Printer
					</TabAnchor>
					<TabAnchor href="/docs" selected={$page.url.pathname.includes('/docs')}>
						<svelte:fragment slot="lead"><i class="fa-solid fa-book" /></svelte:fragment>
						Docs
					</TabAnchor>
				</TabGroup>
			</div>
			<svelte:fragment slot="trail">
				<a class="variant-filled-error btn btn-sm" href="/models/create"> + New Model </a>
				<button
					class="btn hover:variant-soft-error"
					use:popup={{ event: 'click', target: 'features' }}
				>
					<span>Explore</span>
					<i class="fa-solid fa-caret-down opacity-50" />
				</button>
				<!-- popup -->
				<div class="card w-60 p-4 shadow-xl" data-popup="features">
					<nav class="list-nav">
						<ul>
							<li class="hover:variant-ghost-error">
								<a href="{base}/admin">
									<span>Admin</span>
								</a>
							</li>
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
					<div class="bg-surface-100-800-token arrow" />
				</div>
			</svelte:fragment>
		</AppBar>
	</svelte:fragment>
	<!-- Page Route Content -->
	<div id="content" class="mx-auto w-full h-full px-8 pt-4">
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
</style>

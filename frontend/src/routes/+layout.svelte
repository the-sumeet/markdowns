<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { onMount } from 'svelte';
	import { GetCurrentFilesState } from '$lib/wailsjs/go/main/App';
	import { appState } from '../store.svelte';
	import type { main } from '$lib/wailsjs/go/models';
	import { Button } from '$lib/components/ui/button/index.js';

	let { children } = $props();

	onMount(() => {
		GetCurrentFilesState().then((res: main.CurrentFilesState) => {
			appState.currentDir = res.currentDir;
			appState.currentFile = res.currentFile;
			appState.contentHash = res.contentHash;
		});
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Sidebar.Provider class="h-screen">
	<AppSidebar />
	<Sidebar.Inset class="flex max-w-full flex-col">
		<header class="flex h-16 max-w-full shrink-0 items-center gap-2">
			<div class="flex-1 flex justify-around gap-2 px-4">
				<div class="flex flex-1">
					<Sidebar.Trigger class="-ml-1" />
					<Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />
				</div>
				<div class="flex">
					{#if appState.staleContent}
					<Button>Save</Button>
					{/if}
				</div>
				<!-- <Breadcrumb.Root>
					<Breadcrumb.List>
						<Breadcrumb.Item class="hidden md:block">
							<Breadcrumb.Link href="#">Building Your Application</Breadcrumb.Link>
						</Breadcrumb.Item>
						<Breadcrumb.Separator class="hidden md:block" />
						<Breadcrumb.Item>
							<Breadcrumb.Page>Data Fetching</Breadcrumb.Page>
						</Breadcrumb.Item>
					</Breadcrumb.List>
				</Breadcrumb.Root> -->
			</div>
		</header>
		<div class="max-h-full flex-1 overflow-y-hidden p-2">
			{@render children?.()}
		</div>
	</Sidebar.Inset>
</Sidebar.Provider>

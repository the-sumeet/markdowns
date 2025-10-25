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

	let { children } = $props();

	onMount(() => {
		GetCurrentFilesState().then((res: main.CurrentFilesState) => {
			appState.currentDir = res.currentDir;
			appState.currentFile = res.currentFile;
		});
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Sidebar.Provider class="h-screen">
	<AppSidebar />
	<Sidebar.Inset class="flex flex-col max-w-full">
		<header class="flex h-16 shrink-0 items-center gap-2 max-w-full">
			<div class="flex items-center gap-2 px-4">
				<Sidebar.Trigger class="-ml-1" />
				<Separator orientation="vertical" class="mr-2 data-[orientation=vertical]:h-4" />
				
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

<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { onMount } from 'svelte';
	import { GetCurrentFilesState, ListFiles } from '$lib/wailsjs/go/main/App';
	import { appState } from '../store.svelte';
	import type { main } from '$lib/wailsjs/go/models';
	import { Button } from '$lib/components/ui/button/index.js';
	import Sidebar from './Sidebar.svelte';

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

<div class="flex h-screen flex-col overflow-y-hidden bg-background">
	
	<div class="flex flex-1 overflow-hidden">
		<Sidebar />

		<div class="flex-1">
			{@render children()}
		</div>
	</div>
</div>



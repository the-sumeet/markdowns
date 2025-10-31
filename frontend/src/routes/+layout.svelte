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
	import { goto } from '$app/navigation';

	let { children } = $props();

	onMount(() => {
		GetCurrentFilesState().then((res: main.CurrentFilesState) => {
			appState.currentDir = res.currentDir;
			appState.currentFile = res.currentFile;
		});
	});

	$effect(() => {
		if (appState.currentFile) {
			goto('/note')
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<style>
		.gradient-bg {
			background: radial-gradient(circle at 30% 70%, rgba(173, 216, 230, 0.35), transparent 60%),
						radial-gradient(circle at 70% 30%, rgba(255, 182, 193, 0.4), transparent 60%);
		}
		.dark .gradient-bg {
			background: radial-gradient(ellipse at 20% 30%, rgba(56, 189, 248, 0.4) 0%, transparent 60%),
        radial-gradient(ellipse at 80% 70%, rgba(139, 92, 246, 0.3) 0%, transparent 70%),
        radial-gradient(ellipse at 60% 20%, rgba(236, 72, 153, 0.25) 0%, transparent 50%),
        radial-gradient(ellipse at 40% 80%, rgba(34, 197, 94, 0.2) 0%, transparent 65%)
		}
	</style>
</svelte:head>




<div class="gradient-bg flex h-screen flex-col overflow-y-hidden ">
	<div class="flex flex-1 overflow-hidden">
		<Sidebar />

		<div class="flex-1">
			{@render children()}
		</div>
	</div>
</div>


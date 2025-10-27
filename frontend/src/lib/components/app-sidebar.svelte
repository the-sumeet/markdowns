<script lang="ts">
	import BookOpenIcon from '@lucide/svelte/icons/book-open';
	import BotIcon from '@lucide/svelte/icons/bot';
	import ChartPieIcon from '@lucide/svelte/icons/chart-pie';
	import FrameIcon from '@lucide/svelte/icons/frame';
	import LifeBuoyIcon from '@lucide/svelte/icons/life-buoy';
	import MapIcon from '@lucide/svelte/icons/map';
	import SendIcon from '@lucide/svelte/icons/send';
	import Settings2Icon from '@lucide/svelte/icons/settings-2';
	import SquareTerminalIcon from '@lucide/svelte/icons/square-terminal';
	import NavFiles from './nav-files.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { type ComponentProps } from 'svelte';
	import FilePlus from '@lucide/svelte/icons/file-plus';
	import FolderPlus from '@lucide/svelte/icons/folder-plus';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { CreateFile, CreateDir, ListFiles } from '$lib/wailsjs/go/main/App';
	import { appState } from '../../store.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	
	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	// Dialog state
	let showFileDialog = $state(false);
	let showDirDialog = $state(false);
	let fileName = $state('');
	let dirName = $state('');
	let errorMessage = $state('');

	// Handle file creation
	async function handleCreateFile() {
		errorMessage = '';
		try {
			await CreateFile(fileName);
			showFileDialog = false;
			fileName = '';
			// Trigger refresh by incrementing the refreshTrigger
			appState.refreshTrigger = (appState.refreshTrigger || 0) + 1;
		} catch (err) {
			errorMessage = String(err);
		}
	}

	// Handle directory creation
	async function handleCreateDir() {
		errorMessage = '';
		try {
			await CreateDir(dirName);
			showDirDialog = false;
			dirName = '';
			// Trigger refresh by incrementing the refreshTrigger
			appState.refreshTrigger = (appState.refreshTrigger || 0) + 1;
		} catch (err) {
			errorMessage = String(err);
		}
	}

	const data = {
		user: {
			name: 'shadcn',
			email: 'm@example.com',
			avatar: '/avatars/shadcn.jpg'
		},
		navMain: [
			{
				title: 'Playground',
				url: '#',
				icon: SquareTerminalIcon,
				isActive: true,
				items: [
					{
						title: 'History',
						url: '#'
					},
					{
						title: 'Starred',
						url: '#'
					},
					{
						title: 'Settings',
						url: '#'
					}
				]
			},
			{
				title: 'Models',
				url: '#',
				icon: BotIcon,
				items: [
					{
						title: 'Genesis',
						url: '#'
					},
					{
						title: 'Explorer',
						url: '#'
					},
					{
						title: 'Quantum',
						url: '#'
					}
				]
			},
			{
				title: 'Documentation',
				url: '#',
				icon: BookOpenIcon,
				items: [
					{
						title: 'Introduction',
						url: '#'
					},
					{
						title: 'Get Started',
						url: '#'
					},
					{
						title: 'Tutorials',
						url: '#'
					},
					{
						title: 'Changelog',
						url: '#'
					}
				]
			},
			{
				title: 'Settings',
				url: '#',
				icon: Settings2Icon,
				items: [
					{
						title: 'General',
						url: '#'
					},
					{
						title: 'Team',
						url: '#'
					},
					{
						title: 'Billing',
						url: '#'
					},
					{
						title: 'Limits',
						url: '#'
					}
				]
			}
		],
		navSecondary: [
			{
				title: 'Support',
				url: '#',
				icon: LifeBuoyIcon
			},
			{
				title: 'Feedback',
				url: '#',
				icon: SendIcon
			}
		],
		projects: [
			{
				name: 'Design Engineering',
				url: '#',
				icon: FrameIcon
			},
			{
				name: 'Sales & Marketing',
				url: '#',
				icon: ChartPieIcon
			},
			{
				name: 'Travel',
				url: '#',
				icon: MapIcon
			}
		]
	};
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
	<!-- <Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="lg">
					{#snippet child({ props })}
						<a href="##" {...props}>
							<div
								class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg"
							>
								<CommandIcon class="size-4" />
							</div>
							<div class="grid flex-1 text-left text-sm leading-tight">
								<span class="truncate font-medium">Acme Inc</span>
								<span class="truncate text-xs">Enterprise</span>
							</div>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header> -->
	<Sidebar.Content>
		<!-- <NavMain items={data.navMain} /> -->
		<NavFiles />
		<!-- <NavSecondary items={data.navSecondary} class="mt-auto" /> -->
	</Sidebar.Content>
	<Sidebar.Footer class="border-t">
		<div class="flex gap-2 w-full justify-end">
			<Button variant="outline" size="icon" class="size-8" onclick={() => (showFileDialog = true)}>
				<FilePlus />
			</Button>
			<Button variant="outline" size="icon" class="size-8" onclick={() => (showDirDialog = true)}>
				<FolderPlus />
			</Button>
		</div>
		<!-- <NavUser user={data.user} /> -->
	</Sidebar.Footer>
</Sidebar.Root>

<!-- New File Dialog -->
<Dialog.Root bind:open={showFileDialog}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Create New File</Dialog.Title>
			<Dialog.Description>
				Enter a name for the new file. The file will be created in the current directory.
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="fileName">File name</Label>
				<Input
					id="fileName"
					bind:value={fileName}
					placeholder="example.md"
					onkeydown={(e) => {
						if (e.key === 'Enter') {
							handleCreateFile();
						}
					}}
				/>
			</div>
			{#if errorMessage}
				<p class="text-sm text-red-500">{errorMessage}</p>
			{/if}
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showFileDialog = false)}>Cancel</Button>
			<Button onclick={handleCreateFile} disabled={!fileName.trim()}>Create File</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<!-- New Directory Dialog -->
<Dialog.Root bind:open={showDirDialog}>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>Create New Directory</Dialog.Title>
			<Dialog.Description>
				Enter a name for the new directory. It will be created in the current directory.
			</Dialog.Description>
		</Dialog.Header>
		<div class="grid gap-4 py-4">
			<div class="grid gap-2">
				<Label for="dirName">Directory name</Label>
				<Input
					id="dirName"
					bind:value={dirName}
					placeholder="my-folder"
					onkeydown={(e) => {
						if (e.key === 'Enter') {
							handleCreateDir();
						}
					}}
				/>
			</div>
			{#if errorMessage}
				<p class="text-sm text-red-500">{errorMessage}</p>
			{/if}
		</div>
		<Dialog.Footer>
			<Button variant="outline" onclick={() => (showDirDialog = false)}>Cancel</Button>
			<Button onclick={handleCreateDir} disabled={!dirName.trim()}>Create Directory</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { useSidebar } from '$lib/components/ui/sidebar/index.js';
	import EllipsisIcon from '@lucide/svelte/icons/ellipsis';
	import FolderIcon from '@lucide/svelte/icons/folder';
	import ShareIcon from '@lucide/svelte/icons/share';
	import Trash2Icon from '@lucide/svelte/icons/trash-2';
	import Folder from '@lucide/svelte/icons/folder';
	import File from '@lucide/svelte/icons/file';
	import SquareArrowDown from '@lucide/svelte/icons/square-arrow-down';
	import { OpenFile, ListFiles } from '$lib/wailsjs/go/main/App';
	import { onMount } from 'svelte';
	import type { main } from '$lib/wailsjs/go/models';
	import { appState } from '../../store.svelte';

	let files: main.FileEntry[] = $state([]);

	$effect(() => {
		if (appState.currentDir) {
			ListFiles(appState.currentDir).then((res: main.FileEntry[]) => {
				files = res;
			});
		}
	});

	const sidebar = useSidebar();

	onMount(() => {});
</script>

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<!-- <Sidebar.GroupLabel>Projects</Sidebar.GroupLabel> -->
	<Sidebar.Menu>
		{#if appState.currentDir !== '/'}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton
					onclick={() => {
						if (appState.currentDir) {
							const parentDir = appState.currentDir.split('/').slice(0, -1).join('/') || '/';
							ListFiles(parentDir).then((res: main.FileEntry[]) => {
								files = res;
								appState.currentDir = parentDir;
							});
						}
					}}
				>
					<span>..</span>
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		{/if}
		{#each files as file (file.path)}
			<Sidebar.MenuItem>
				<Sidebar.MenuButton
					disabled={!file.isDirectory &&
						!file.name.endsWith('.md') &&
						!file.name.endsWith('.markdown')}
					onclick={() => {
						OpenFile(file.path).then((res) => {
							appState.currentDir = res.currentDir;
							appState.currentFile = res.currentFile;
						});
					}}
				>
					<!-- <item.icon /> -->
					{#if file.isDirectory}
						<Folder />
					{:else if file.name.endsWith('.md') || file.name.endsWith('.markdown')}
						<SquareArrowDown />
					{:else}
						<File />
					{/if}
					<span>{file.name}</span>
				</Sidebar.MenuButton>
				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Sidebar.MenuAction showOnHover {...props}>
								<EllipsisIcon />
								<span class="sr-only">More</span>
							</Sidebar.MenuAction>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content
						class="w-48"
						side={sidebar.isMobile ? 'bottom' : 'right'}
						align={sidebar.isMobile ? 'end' : 'start'}
					>
						<DropdownMenu.Item>
							<FolderIcon class="text-muted-foreground" />
							<span>View Project</span>
						</DropdownMenu.Item>
						<DropdownMenu.Item>
							<ShareIcon class="text-muted-foreground" />
							<span>Share Project</span>
						</DropdownMenu.Item>
						<DropdownMenu.Separator />
						<DropdownMenu.Item>
							<Trash2Icon class="text-muted-foreground" />
							<span>Delete Project</span>
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</Sidebar.MenuItem>
		{/each}
		<Sidebar.MenuItem>
			<Sidebar.MenuButton>
				<EllipsisIcon />
				<span>More</span>
			</Sidebar.MenuButton>
		</Sidebar.MenuItem>
	</Sidebar.Menu>
</Sidebar.Group>

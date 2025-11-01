<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import CircleX from '@lucide/svelte/icons/circle-x';
	import * as Card from '$lib/components/ui/card/index.js';
</script>

<div class="flex h-full w-full flex-col overflow-hidden p-4">
	<div
		class="mb-4 flex shrink-0 justify-between rounded-xl bg-white/30 p-4 backdrop-blur-md dark:bg-black/30"
	>
		<div class="flex-1">
			<h1 class=" text-3xl font-extrabold">Preferences</h1>
		</div>
		<div class="flex items-center gap-1">
			<Button
				variant="outline"
				size="icon"
				class="size-8 rounded-full"
				onclick={() => {
					goto('/');
				}}
			>
				<CircleX />
			</Button>
		</div>
	</div>

	<div class="flex flex-col overflow-y-auto">
		<Card.Root>
			<Card.Header>
				<Card.Title>Files</Card.Title>
			</Card.Header>
			<Card.Content>
				<p>Card Content</p>
			</Card.Content>
			<Card.Footer>
				<p>Card Footer</p>
			</Card.Footer>
		</Card.Root>
	</div>
</div>

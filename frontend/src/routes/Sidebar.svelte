<script>
	import Button from '$lib/components/ui/button/button.svelte';
	import Label from '$lib/components/ui/label/label.svelte';
	import Plus from '@lucide/svelte/icons/plus';
	import Settings from '@lucide/svelte/icons/settings';
</script>

<div
	style="--wails-draggable:drag"
	class="flex h-full flex-col items-center gap-8 overflow-y-auto border-r p-6 pt-20"
>
	<!-- <p class="font-yellowtail">Markdowns</p> -->

	<div class="flex flex-1 flex-col">
		<Button size="icon" class="size-8 rounded-full">
			<Plus />
		</Button>
	</div>

	<div class="flex flex-col">
		<Button href="/settings" size="icon" variant='outline' class="size-8 rounded-full">
			<Settings />
		</Button>
	</div>
</div>

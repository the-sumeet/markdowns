<script>
	import { Crepe } from '@milkdown/crepe';
	import '@milkdown/crepe/theme/common/style.css';
	import '@milkdown/crepe/theme/frame-dark.css';
	import { onMount } from 'svelte';

	// Choose your preferred theme

	// Initialize the editor

	// Clean up when done
	// crepe.destroy();

	onMount(() => {
		// Create editor instance
		const crepe = new Crepe({
			root: document.getElementById('editor'),
			defaultValue: '# Hello, Crepe!\n\nStart writing your markdown...'
		});

		crepe.create();
	});
</script>

<div class="h-full ">
  <div id="editor" class="h-full "></div>
</div>

<style>
	:global(.milkdown) {
    overflow-y: hidden;
		height: 100% !important;
		margin: 0px !important;
		width: 100% !important;
    max-width: 100%!important;
    max-height: 100%!important;
    border: 0px!important;
    overflow-y: auto !important;
	}
</style>

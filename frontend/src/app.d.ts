// See https://svelte.dev/docs/kit/types#app.d.ts

import type { main } from "$lib/wailsjs/go/models";

// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}

		interface AppState {
			currentFile?: string
			currentDir?: string
		}
	}
}

export {};

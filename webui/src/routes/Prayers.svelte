<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { Card, Button, Heading, Textarea, Label, Checkbox, Badge, Alert } from 'flowbite-svelte';
	import { TrashBinOutline, HeartOutline, GlobeOutline, ShieldCheckOutline } from 'flowbite-svelte-icons';
	import { getAllPrayers, addPrayer, deletePrayer } from '../oo';

	const oo = getContext('oo');
	
	let prayers = $state([]);
	let newConcern = $state('');
	let isAnonymous = $state(false);
	let loading = $state(true);

	async function refresh() {
		loading = true;
		try {
			prayers = await getAllPrayers();
		} catch (e) {
			// Silently fail or toast if the 404 is still active
		} finally {
			loading = false;
		}
	}

	async function submit() {
		if (!newConcern.trim()) return;
		try {
			await addPrayer(newConcern, isAnonymous);
			newConcern = '';
			await refresh();
		} catch (e: any) {
			console.error(e);
		}
	}

	async function remove(pid: number) {
		if (confirm("Remove this request from the Office?")) {
			await deletePrayer(pid);
			await refresh();
		}
	}

	onMount(refresh);
</script>

<div class="max-w-4xl mx-auto p-4 space-y-6">
	<header class="text-center py-6">
		<Heading tag="h1" class="font-serif text-4xl text-primary-900">Prayer Wall</Heading>
		<p class="text-slate-500 mt-2">Requests shared here are included in the Web Amplified Daily Office.</p>
	</header>

	<Card size="none" class="p-6 shadow-lg border-t-4 border-t-primary-600">
		<div class="flex items-center gap-2 mb-4 text-primary-800">
			<GlobeOutline size="sm" />
			<span class="text-xs font-bold uppercase tracking-widest">New WADO Request</span>
		</div>
		
		<Textarea bind:value={newConcern} placeholder="What shall we pray for?" rows="4" class="mb-4" />
		
		<div class="flex flex-col sm:flex-row justify-between items-center gap-4">
			<Checkbox bind:checked={isAnonymous}>
				<span class="text-sm">Post as "A Sibling of the Order"</span>
			</Checkbox>
			<Button color="primary" onclick={submit} disabled={!newConcern}>
				Submit to the Office
			</Button>
		</div>
	</Card>

	<div class="space-y-4">
		{#if loading}
			<p class="text-center text-slate-400 animate-pulse">Gathering concerns...</p>
		{:else}
			{#each prayers as p}
				<Card size="none" class="p-5 hover:bg-slate-50 transition-colors relative group">
					<div class="flex justify-between">
						<div class="flex items-center gap-2 mb-2">
							<HeartOutline size="xs" class="text-red-500" />
							<span class="font-bold text-sm">
								{p.anonymous ? 'A Sibling of the Order' : p.MemberName}
							</span>
							<span class="text-[10px] text-slate-400">
								{new Date(p.date).toLocaleDateString()}
							</span>
						</div>
						
						{#if p.member == oo.me.id || oo.me.level > 2}
							<button onclick={() => remove(p.id)} class="text-slate-300 hover:text-red-600 transition-colors">
								<TrashBinOutline size="sm" />
							</button>
						{/if}
					</div>
					<p class="text-slate-700 whitespace-pre-wrap leading-relaxed">{p.content}</p>
				</Card>
			{:else}
				<div class="text-center py-20 border-2 border-dashed border-slate-200 rounded-xl">
					<p class="text-slate-400 italic">The wall is currently empty.</p>
				</div>
			{/each}
		{/if}
	</div>
</div>

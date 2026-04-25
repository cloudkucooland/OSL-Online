<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { Card, Button, Heading, Textarea, Label, Checkbox, Badge, Alert } from 'flowbite-svelte';
	import {
		TrashBinOutline,
		HeartOutline,
		GlobeOutline,
		ShieldCheckOutline
	} from 'flowbite-svelte-icons';
	import { getAllPrayers, addPrayer, deletePrayer } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

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
			toast.push(e.message);
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
			toast.push(e.message);
			console.error(e);
		}
	}

	async function remove(pid: number) {
		if (confirm('Remove this request from the Office?')) {
			try {
				await deletePrayer(pid);
			} catch (e: any) {
				toast.push(e.message);
			} finally {
				await refresh();
			}
		}
	}

	onMount(refresh);
</script>

<div class="mx-auto max-w-4xl space-y-6 p-4">
	<header class="py-6 text-center">
		<Heading tag="h1" class="text-primary-900 font-serif text-4xl">Prayer Requests</Heading>
		<p class="mt-2 text-slate-500">
			Requests shared here are (will be, when Br. Scot finishes it) included in WADO.
		</p>
	</header>

	<Card size="none" class="border-t-primary-600 border-t-4 p-6 shadow-lg">
		<div class="text-primary-800 mb-4 flex items-center gap-2">
			<GlobeOutline size="sm" />
			<span class="text-xs font-bold tracking-widest uppercase">New Prayer Request</span>
		</div>

		<Textarea
			bind:value={newConcern}
			placeholder="For whom or what shall we pray?"
			rows="4"
			class="mb-4"
		/>

		<div class="flex flex-col items-center justify-between gap-4 sm:flex-row">
			<Checkbox bind:checked={isAnonymous}>
				<span class="text-sm">Post anonymously</span>
			</Checkbox>
			<Button color="primary" onclick={submit} disabled={!newConcern}>Add</Button>
		</div>
	</Card>

	<div class="space-y-4">
		{#if loading}
			<p class="animate-pulse text-center text-slate-400">Gathering concerns...</p>
		{:else}
			{#each prayers as p}
				<Card size="none" class="group relative p-5 transition-colors hover:bg-slate-50">
					<div class="flex justify-between">
						<div class="mb-2 flex items-center gap-2">
							<HeartOutline size="xs" class="text-red-500" />
							<span class="text-sm font-bold">
								{p.Anonymous ? 'A Sibling of the Order' : p.OSLName}
							</span>
							<span class="text-[10px] text-slate-400">
								{new Date(p.Date).toLocaleDateString()}
							</span>
						</div>

						{#if p.MemberID == oo.me.id || oo.me.level > 2}
							<button
								onclick={() => remove(p.PrayerID)}
								class="text-slate-300 transition-colors hover:text-red-600"
							>
								<TrashBinOutline size="sm" />
							</button>
						{/if}
					</div>
					<p class="leading-relaxed whitespace-pre-wrap text-slate-700">{p.Content}</p>
				</Card>
			{:else}
				<div class="text-center py-20 border-2 border-dashed border-slate-200 rounded-xl">
					<p class="text-slate-400 italic">The wall is currently empty.</p>
				</div>
			{/each}
		{/if}
	</div>
</div>

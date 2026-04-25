<script lang="ts">
	import { onMount, getContext } from 'svelte';
	import { Select, Card, Badge, Label, Heading, Alert } from 'flowbite-svelte';
	import { InfoCircleOutline, MapPinAltOutline } from 'flowbite-svelte-icons';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { oslname, getLocalities, getLocalityMembers } from '../oo';

	const oo = getContext('oo');

	// Guard: Ensure user is logged in
	if (!oo.me) {
		push('/Login');
	}

	let { params = {} } = $props();
	let locs = $state([]);
	let members = $state([]);
	let selected = $state();

	async function updateMembers(locality: string) {
		if (!locality) return;
		try {
			members = await getLocalityMembers(locality);
		} catch (e: any) {
			toast.push(e.message);
		}
	}

	async function handleSelectChange() {
		if (selected) {
			push(`#/localitybrowser/${selected}`);
			await updateMembers(selected);
		}
	}

	async function load() {
		locs = await getLocalities();
		return locs;
	}

	onMount(async () => {
		await load();
		if (params.loc) {
			selected = params.loc;
			await updateMembers(selected);
		}
	});
</script>

<div class="mx-auto max-w-6xl space-y-6 p-4">
	{#await load()}
		<div class="flex justify-center p-12">
			<h3 class="animate-pulse text-slate-400">Loading Localities...</h3>
		</div>
	{:then}
		<Card size="none" class="border-slate-200 bg-white p-6 shadow-sm">
			<div class="text-primary-900 mb-4 flex items-center gap-2">
				<MapPinAltOutline size="lg" />
				<Heading tag="h2" class="text-2xl font-bold">Locality Browser</Heading>
			</div>

			<Alert color="blue" class="mb-6">
				<div class="flex items-center gap-2 font-medium">
					<InfoCircleOutline size="sm" />
					<span>Note on International Search</span>
				</div>
				<p class="mt-1 text-sm">
					While "search by state" is a common request, please remember we are an International
					order. Results are grouped by country and administrative region.
				</p>
			</Alert>

			<div class="space-y-2">
				<Label for="locality-select" class="text-sm font-semibold text-slate-700"
					>Select Region</Label
				>
				<Select
					id="locality-select"
					items={locs}
					bind:value={selected}
					onchange={handleSelectChange}
					placeholder="Choose a country or state..."
				/>
			</div>
		</Card>

		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each members as m}
				<Card
					size="none"
					class="group hover:border-primary-200 border-slate-100 p-5 shadow-sm transition-all hover:bg-slate-50"
				>
					<div class="flex h-full flex-col">
						<div class="mb-3 flex items-start justify-between">
							<a
								href="#/member/{m.ID}"
								class="group-hover:text-primary-700 text-lg font-bold text-slate-900 transition-colors"
							>
								{oslname(m)}
							</a>
							<Badge color={m.MemberStatus.includes('Vows') ? 'green' : 'indigo'} class="text-xs">
								{m.MemberStatus}
							</Badge>
						</div>

						<div class="mt-auto space-y-2 border-t border-slate-100 pt-3 text-sm">
							{#if m.PrimaryEmail}
								<div class="flex items-center gap-2 text-slate-600">
									<a
										href="mailto:{m.PrimaryEmail}"
										class="hover:text-primary-600 truncate hover:underline"
									>
										{m.PrimaryEmail}
									</a>
								</div>
							{/if}
							{#if m.PrimaryPhone}
								<div class="flex items-center gap-2 text-slate-600">
									<a href="tel:{m.PrimaryPhone}" class="hover:text-primary-600">
										{m.PrimaryPhone}
									</a>
								</div>
							{/if}
						</div>
					</div>
				</Card>
			{:else}
				{#if selected}
					<div class="col-span-full py-16 text-center">
						<div class="text-slate-300 mb-2"><MapPinAltOutline size="xl" class="mx-auto" /></div>
						<p class="text-slate-500 italic">No members currently listed in this locality.</p>
					</div>
				{/if}
			{/each}
		</div>
	{:catch error}
		<Alert color="red" class="mt-4">
			<span class="font-bold">Error:</span>
			{error.message}
		</Alert>
	{/await}
</div>

<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { push } from 'svelte-spa-router';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		Card,
		Spinner,
		Badge
	} from 'flowbite-svelte';
	import { getNecrology, oslname } from '../oo';
	import { CloudArrowUpOutline } from 'flowbite-svelte-icons';

	const oo = getContext('oo');

	let loading = $state(true);
	let rawData = $state([]);

	let items = $derived(
		rawData.map((i) => {
			const d = new Date(i.DateDeceased);
			return {
				name: oslname(i),
				date: d.toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' }),
				locality: i.State ? `${i.State}, ${i.Country}` : i.Country
			};
		})
	);

	onMount(async () => {
		// Guard: No dollar sign here!
		if (!oo.me) return push('/Login');

		try {
			rawData = await getNecrology();
		} catch (err) {
			console.error(err);
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>Necrology | OSL Directory</title>
</svelte:head>

<div class="mx-auto max-w-7xl space-y-8 px-6 py-10">
	<header class="flex items-center justify-between border-b-2 border-slate-100 pb-8">
		<div class="flex items-center gap-4">
			<div class="rounded-lg bg-slate-800 p-3 text-white">
				<CloudArrowUpOutline class="h-8 w-8" />
			</div>
			<div>
				<Heading tag="h1" class="text-4xl font-black tracking-tight text-slate-900"
					>Necrology</Heading
				>
				<p class="font-serif text-lg italic text-slate-500">In Memoriam: The Church Triumphant</p>
			</div>
		</div>
		<Badge color="dark" class="px-4 py-1 text-sm uppercase tracking-widest">Rest in Peace</Badge>
	</header>

	{#if loading}
		<div class="flex justify-center py-40">
			<Spinner color="dark" size="16" />
		</div>
	{:else}
		<Card size="none" class="overflow-hidden border-none shadow-xl ring-1 ring-slate-200">
			<Table striped={true} hoverable={true} class="text-base">
				<TableHead class="bg-slate-900 py-4 text-white">
					<TableHeadCell class="py-5 text-lg">Sibling Name</TableHeadCell>
					<TableHeadCell class="text-lg">Date of Decease</TableHeadCell>
					<TableHeadCell class="text-lg">Locality</TableHeadCell>
				</TableHead>
				<TableBody>
					{#each items as item}
						<TableBodyRow class="border-b border-slate-100 text-lg">
							<TableBodyCell class="py-6 font-bold text-slate-900">{item.name}</TableBodyCell>
							<TableBodyCell class="font-serif text-slate-600">{item.date}</TableBodyCell>
							<TableBodyCell class="text-slate-500">{item.locality || '—'}</TableBodyCell>
						</TableBodyRow>
					{:else}
						<TableBodyRow>
							<TableBodyCell colspan={3} class="text-center py-20 text-slate-400 text-xl italic">
								No records found in the historical log.
							</TableBodyCell>
						</TableBodyRow>
					{/each}
				</TableBody>
			</Table>
		</Card>
	{/if}
</div>

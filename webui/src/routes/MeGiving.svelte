<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Card,
		Heading,
		Spinner,
		Badge,
		Button
	} from 'flowbite-svelte';
	import { CashOutline, ArrowLeftOutline, GlobeOutline } from 'flowbite-svelte-icons';
	import { getMeFromServer, getMeGiving, oslname } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';

	const oo = getContext('oo');

	// Guard: Redirect if not logged in
	if (!oo.me) {
		push('/Login');
	}

	let loading = $state(true);
	let member = $state(null);
	let history = $state([]);

	onMount(async () => {
		try {
			const m = await getMeFromServer();
			member = m;
			history = await getMeGiving(); // oo.js shows getMeGiving() takes no args for 'me'
		} catch (err: any) {
			toast.push(err.message);
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>My Giving History | OSL</title>
</svelte:head>

<div class="mx-auto max-w-5xl space-y-6 px-4 py-8">
	<header class="flex flex-col justify-between gap-4 md:flex-row md:items-center">
		<div>
			<Button size="xs" color="alternative" onclick={() => push('/me')} class="mb-2">
				<ArrowLeftOutline class="mr-1 h-3 w-3" /> Back to My Profile
			</Button>
			<Heading tag="h2" class="text-3xl font-bold text-slate-900">Giving History</Heading>
		</div>

		{#if member}
			<div class="text-right">
				<p class="text-sm italic text-slate-500">Records for</p>
				<p class="font-semibold text-primary-700">{oslname(member)}</p>
			</div>
		{/if}
	</header>

	{#if loading}
		<div class="flex justify-center py-20">
			<Spinner color="purple" size="12" />
		</div>
	{:else}
		<Card size="none" class="overflow-hidden border-slate-200 shadow-sm">
			<Table hoverable={true} striped={true}>
				<TableHead class="border-b border-slate-200 bg-slate-50">
					<TableHeadCell>Date</TableHeadCell>
					<TableHeadCell>Description</TableHeadCell>
					<TableHeadCell>Method/Ref</TableHeadCell>
					<TableHeadCell class="text-right">Amount</TableHeadCell>
				</TableHead>
				<TableBody>
					{#each history as row}
						<TableBodyRow>
							<TableBodyCell class="font-medium">{row.Date}</TableBodyCell>
							<TableBodyCell>
								<span class="text-slate-700">{row.Description || 'General Offering'}</span>
							</TableBodyCell>
							<TableBodyCell>
								{#if row.Transaction && row.Transaction !== '0'}
									<a
										href="https://www.paypal.com/unifiedtransactions/details/payment/{row.Transaction}"
										target="_blank"
										class="flex items-center gap-1 text-xs text-primary-600 hover:underline"
									>
										<GlobeOutline class="h-3 w-3" /> PayPal: {row.Transaction.slice(0, 8)}...
									</a>
								{:else if row.Check && row.Check !== '0'}
									<Badge color="indigo" class="font-mono">Check #{row.Check}</Badge>
								{:else}
									<span class="text-xs italic text-slate-400">Other</span>
								{/if}
							</TableBodyCell>
							<TableBodyCell class="text-right font-bold text-slate-900">
								${Number(row.Amount).toLocaleString(undefined, { minimumFractionDigits: 2 })}
							</TableBodyCell>
						</TableBodyRow>
					{:else}
						<TableBodyRow>
							<TableBodyCell colspan="4" class="text-center py-12 text-slate-400 italic">
								No giving records found for the current period.
							</TableBodyCell>
						</TableBodyRow>
					{/each}
				</TableBody>
			</Table>
		</Card>

		<div
			class="flex items-center gap-2 rounded-lg border border-blue-100 bg-blue-50 p-4 text-sm text-blue-800"
		>
			<CashOutline class="h-5 w-5" />
			<p>
				Thank you for your faithfulness. These records reflect contributions received during the
				current fiscal year.
			</p>
		</div>
	{/if}
</div>

<script lang="ts">
	import { onMount, getContext } from 'svelte';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Label,
		Input,
		Button,
		Card,
		Heading,
		Badge,
		Spinner
	} from 'flowbite-svelte';
	import { PlusOutline, ArrowLeftOutline, GlobeOutline, CashOutline } from 'flowbite-svelte-icons';
	import { getMember, getGiving, postGiving, cleanDateFormat, oslname } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';

	let { params = {} } = $props();

	const oo = getContext('oo');
	// Guard: Ensure officer is logged in
	if (!oo.me) push('/Login');

	let loading = $state(true);
	let member = $state(null);
	let history = $state([]);

	// Form State
	let postdate = $state(cleanDateFormat(new Date().toISOString()));
	let amount = $state<number | undefined>();
	let description = $state('Annual Reaffirmation');
	let check = $state('0');
	let transaction = $state('0');

	// Svelte 5 Derived State
	const totalGiving = $derived(
		history.reduce((acc, curr) => acc + parseFloat(curr.Amount || 0), 0)
	);

	async function loadData() {
		try {
			const [m, g] = await Promise.all([getMember(params.id), getGiving(params.id)]);
			member = m;
			// Sort by date descending (newest first)
			history = g.sort((a, b) => new Date(b.Date).getTime() - new Date(a.Date).getTime());
		} catch (err: any) {
			toast.push(err.message);
		} finally {
			loading = false;
		}
	}

	onMount(loadData);

	async function handleAdd(event: Event) {
		event.preventDefault();
		if (!amount || amount <= 0) {
			toast.push('Please enter a valid amount');
			return;
		}

		try {
			await postGiving(
				params.id,
				cleanDateFormat(postdate),
				amount,
				description,
				check,
				transaction
			);
			toast.push(`Record Posted`);

			// Reset fields for next entry
			amount = undefined;
			check = '0';
			transaction = '0';

			await loadData(); // Refresh list
		} catch (err: any) {
			toast.push('Failed to post: ' + err.message);
		}
	}
</script>

<svelte:head>
	<title>Giving Ledger: {member ? member.LastName : 'Loading'}</title>
</svelte:head>

<div class="w-full space-y-6 px-4 py-6 sm:px-10">
	{#if loading}
		<div class="flex justify-center py-20"><Spinner color="purple" size="12" /></div>
	{:else if member}
		<div class="flex flex-col justify-between gap-6 border-b pb-6 md:flex-row md:items-center">
			<div>
				<Button
					size="xs"
					color="alternative"
					onclick={() => push(`/member/${member.ID}`)}
					class="mb-4"
				>
					<ArrowLeftOutline class="mr-1 h-3 w-3" /> Back to Record
				</Button>
				<Heading tag="h2" class="text-4xl font-black text-slate-900">
					Giving History: {oslname(member)}
				</Heading>
				<div class="mt-2 flex items-center gap-3">
					<Badge color={member.MemberStatus.includes('Vows') ? 'green' : 'dark'} class="text-xs">
						{member.MemberStatus}
					</Badge>
					<span class="font-mono text-slate-400">ID: {member.ID}</span>
				</div>
			</div>

			<div
				class="min-w-[240px] rounded-xl border-2 border-green-100 bg-white p-6 text-right shadow-sm"
			>
				<p class="mb-1 text-xs font-bold uppercase tracking-widest text-slate-500">
					Total Contribution
				</p>
				<p class="text-5xl font-black text-green-600">
					${totalGiving.toLocaleString(undefined, { minimumFractionDigits: 2 })}
				</p>
			</div>
		</div>

		<Card size="none" class="border-slate-200 bg-slate-50/50 p-8 shadow-md">
			<Heading
				tag="h5"
				class="mb-6 flex items-center gap-2 text-sm font-bold uppercase text-slate-500"
			>
				<PlusOutline class="h-4 w-4" /> Add Transaction to Ledger
			</Heading>
			<form onsubmit={handleAdd} class="grid grid-cols-1 items-end gap-6 md:grid-cols-12">
				<div class="md:col-span-2">
					<Label for="Date" class="mb-2 text-xs font-bold text-slate-600">DATE</Label>
					<Input id="Date" bind:value={postdate} placeholder="YYYY-MM-DD" />
				</div>
				<div class="md:col-span-2">
					<Label for="Amount" class="mb-2 text-xs font-bold text-slate-600">AMOUNT</Label>
					<Input
						id="Amount"
						type="number"
						step="0.01"
						bind:value={amount}
						placeholder="0.00"
						required
					/>
				</div>
				<div class="md:col-span-3">
					<Label for="Description" class="mb-2 text-xs font-bold text-slate-600">DESCRIPTION</Label>
					<Input id="Description" bind:value={description} />
				</div>
				<div class="md:col-span-2">
					<Label for="Check" class="mb-2 text-xs font-bold text-slate-600">CHECK #</Label>
					<Input id="Check" bind:value={check} />
				</div>
				<div class="md:col-span-2">
					<Label for="Transaction" class="mb-2 text-xs font-bold text-slate-600">PAYPAL ID</Label>
					<Input id="Transaction" bind:value={transaction} />
				</div>
				<div class="md:col-span-1">
					<Button type="submit" color="green" class="w-full py-2.5">POST</Button>
				</div>
			</form>
		</Card>

		<div class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
			<Table hoverable={true}>
				<TableHead class="bg-slate-100">
					<TableHeadCell>Date</TableHeadCell>
					<TableHeadCell>Description</TableHeadCell>
					<TableHeadCell>Ref / Check</TableHeadCell>
					<TableHeadCell class="text-right">Amount</TableHeadCell>
				</TableHead>
				<TableBody>
					{#each history as row}
						<TableBodyRow>
							<TableBodyCell class="font-medium">{row.Date}</TableBodyCell>
							<TableBodyCell>{row.Description}</TableBodyCell>
							<TableBodyCell>
								{#if row.Transaction && row.Transaction !== '0'}
									<a
										href="https://www.paypal.com/unifiedtransactions/details/payment/{row.Transaction}"
										target="_blank"
										class="flex items-center gap-1 font-mono text-xs text-primary-600 hover:underline"
									>
										<GlobeOutline class="h-3 w-3" />
										{row.Transaction}
									</a>
								{:else if row.Check && row.Check !== '0'}
									<Badge color="indigo" class="font-mono">Check: {row.Check}</Badge>
								{:else}
									<span class="text-slate-300">—</span>
								{/if}
							</TableBodyCell>
							<TableBodyCell class="text-right font-bold text-slate-900">
								${parseFloat(row.Amount).toLocaleString(undefined, { minimumFractionDigits: 2 })}
							</TableBodyCell>
						</TableBodyRow>
					{:else}
						<TableBodyRow>
							<TableBodyCell colspan="4" class="text-center py-12 text-slate-400">
								No transaction history found for this record.
							</TableBodyCell>
						</TableBodyRow>
					{/each}
				</TableBody>
			</Table>
		</div>
	{/if}
</div>

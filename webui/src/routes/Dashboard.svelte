<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import {
		Card,
		Heading,
		Spinner,
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow
	} from 'flowbite-svelte';
	import { ChartPieOutline, UsersGroupOutline, CashOutline } from 'flowbite-svelte-icons';
	import { getDashboard } from '../oo';
	import { push } from 'svelte-spa-router';

	const oo = getContext('oo');

	// Guard: Redirect if not logged in
	if (!oo.me) {
		push('/Login');
	}

	let data = $state(null);
	let loading = $state(true);

	onMount(async () => {
		try {
			data = await getDashboard();
		} catch (err) {
			console.error(err);
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>System Dashboard | OSL Directory</title>
</svelte:head>

<div class="space-y-8">
	<header>
		<Heading tag="h2" class="text-3xl font-extrabold text-slate-900">System Dashboard</Heading>
		<p class="italic text-slate-500">Membership and financial snapshots since July 1st.</p>
	</header>

	{#if loading}
		<div class="flex justify-center py-20">
			<Spinner color="purple" size="12" />
		</div>
	{:else if data}
		<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
			<Card padding="xl" class="border-primary-200 bg-gradient-to-br from-primary-50 to-white">
				<div class="flex items-center gap-4">
					<UsersGroupOutline class="h-10 w-10 text-primary-600" />
					<div>
						<p class="text-sm font-medium uppercase tracking-wider text-slate-500">Total Vowed</p>
						<p class="text-3xl font-bold text-slate-900">{data.LifevowCount + data.AnnualCount}</p>
					</div>
				</div>
			</Card>

			<Card padding="xl" class="border-green-200 bg-gradient-to-br from-green-50 to-white">
				<div class="flex items-center gap-4">
					<CashOutline class="h-10 w-10 text-green-600" />
					<div>
						<p class="text-sm font-medium uppercase tracking-wider text-slate-500">
							This Year Giving
						</p>
						<p class="text-3xl font-bold text-slate-900">${data.ThisYearGiving}</p>
					</div>
				</div>
			</Card>

			<Card padding="xl" class="border-purple-200 bg-gradient-to-br from-purple-50 to-white">
				<div class="flex items-center gap-4">
					<ChartPieOutline class="h-10 w-10 text-purple-600" />
					<div>
						<p class="text-sm font-medium uppercase tracking-wider text-slate-500">Subscribers</p>
						<p class="text-3xl font-bold text-slate-900">{data.SubscriberCount}</p>
					</div>
				</div>
			</Card>
		</div>

		<div class="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
			<Table hoverable={true} striped={true}>
				<TableBody>
					<TableBodyRow>
						<TableBodyCell class="font-semibold">Life Vow Members</TableBodyCell>
						<TableBodyCell>{data.LifevowCount}</TableBodyCell>
					</TableBodyRow>
					<TableBodyRow>
						<TableBodyCell>Life Vows Checked-in (since July 1)</TableBodyCell>
						<TableBodyCell>{data.LifeVowsCheckin}</TableBodyCell>
					</TableBodyRow>
					<TableBodyRow>
						<TableBodyCell>Life Vows who Gave (since July 1)</TableBodyCell>
						<TableBodyCell>{data.LifeVowsWhoGave}</TableBodyCell>
					</TableBodyRow>
					<TableBodyRow class="border-t-2 border-slate-100">
						<TableBodyCell class="font-semibold text-primary-700">Annual Vow Members</TableBodyCell>
						<TableBodyCell class="font-semibold">{data.AnnualCount}</TableBodyCell>
					</TableBodyRow>
					<TableBodyRow>
						<TableBodyCell>Annual Vows Reaffirmed (since July 1)</TableBodyCell>
						<TableBodyCell>{data.AnnualVowsReaffirmed}</TableBodyCell>
					</TableBodyRow>
					<TableBodyRow>
						<TableBodyCell>Annual Vows who Gave (since July 1)</TableBodyCell>
						<TableBodyCell>{data.AnnualVowsWhoGave}</TableBodyCell>
					</TableBodyRow>
					<TableBodyRow class="border-t-2 border-slate-100">
						<TableBodyCell class="font-semibold">Friends (non-vowed)</TableBodyCell>
						<TableBodyCell>{data.FriendCount}</TableBodyCell>
					</TableBodyRow>
					<TableBodyRow>
						<TableBodyCell>Last year member giving (Full Cycle)</TableBodyCell>
						<TableBodyCell class="text-slate-500">${data.LastYearGiving}</TableBodyCell>
					</TableBodyRow>
				</TableBody>
			</Table>
		</div>
	{:else}
		<div class="py-10 text-center text-red-500">Failed to load dashboard data.</div>
	{/if}
</div>

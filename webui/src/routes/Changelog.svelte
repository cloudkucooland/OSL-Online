<script lang="ts">
	import { onMount, getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Heading,
		Badge,
		Breadcrumb,
		BreadcrumbItem,
		Spinner,
		Card
	} from 'flowbite-svelte';
	import { ClockOutline, UserOutline, ArrowRightOutline } from 'flowbite-svelte-icons';
	import { getMember, getChangelog, oslname } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	let { params = {} } = $props();
	const oo = getContext('oo');

	// Guard: Ensure authorized officer access
	if (!oo.me) push('/Login');

	let loading = $state(true);
	let member = $state(null);
	let logs = $state([]);

	onMount(async () => {
		try {
			const [m, c] = await Promise.all([getMember(params.id), getChangelog(params.id)]);
			member = m;
			// Ensure newest changes are at the top
			logs = c.sort((a, b) => new Date(b.Date).getTime() - new Date(a.Date).getTime());
		} catch (err: any) {
			toast.push(err.message);
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>Audit Log: {member ? member.LastName : 'Loading'}</title>
</svelte:head>

<div class="mx-auto max-w-7xl space-y-6 px-4 py-8 sm:px-10">
	{#if loading}
		<div class="flex justify-center py-20"><Spinner color="purple" size="12" /></div>
	{:else if member}
		<nav class="mb-8 flex flex-col gap-6">
			<Breadcrumb
				aria-label="Member navigation"
				class="rounded-lg border border-slate-100 bg-slate-50 px-4 py-2"
			>
				<BreadcrumbItem href="#/" home>Home</BreadcrumbItem>
				<BreadcrumbItem href="#/member/{member.ID}">{oslname(member)}</BreadcrumbItem>
				<BreadcrumbItem>Audit Trail</BreadcrumbItem>
			</Breadcrumb>

			<div
				class="flex items-center justify-between rounded-xl border border-slate-200 bg-white p-6 shadow-sm"
			>
				<div class="flex items-center gap-4">
					<div class="rounded-full bg-purple-100 p-3 text-purple-600">
						<ClockArrowOutline class="h-8 w-8" />
					</div>
					<div>
						<Heading tag="h2" class="text-3xl font-black tracking-tight text-slate-900">
							Change History
						</Heading>
						<p class="font-medium text-slate-500">
							System records for {oslname(member)}
							<Badge color="purple" class="ml-2 text-[10px] tracking-widest uppercase"
								>{member.MemberStatus}</Badge
							>
						</p>
					</div>
				</div>
				<div class="hidden text-right md:block">
					<p class="text-xs font-bold tracking-widest text-slate-400 uppercase">Total Logs</p>
					<p class="font-mono text-3xl font-black text-slate-700">{logs.length}</p>
				</div>
			</div>
		</nav>

		<Card size="none" class="overflow-hidden border-slate-200 bg-white shadow-md">
			<Table hoverable={true} striped={true}>
				<TableHead class="border-b bg-slate-50 text-xs text-slate-500 uppercase">
					<TableHeadCell class="py-4"
						><div class="flex items-center gap-2">
							<ClockOutline class="h-4 w-4" /> Timestamp
						</div></TableHeadCell
					>
					<TableHeadCell
						><div class="flex items-center gap-2">
							<UserOutline class="h-4 w-4" /> Officer
						</div></TableHeadCell
					>
					<TableHeadCell>Attribute</TableHeadCell>
					<TableHeadCell>Value Assigned</TableHeadCell>
				</TableHead>
				<TableBody>
					{#each logs as entry}
						<TableBodyRow class="transition-colors hover:bg-purple-50/30">
							<TableBodyCell class="font-mono text-xs text-slate-400">
								{entry.Date}
							</TableBodyCell>
							<TableBodyCell class="text-sm font-bold text-slate-700">
								{entry.Changer}
							</TableBodyCell>
							<TableBodyCell>
								<div class="flex items-center gap-2">
									<Badge color="indigo" border class="px-2.5 py-0.5 font-bold">{entry.Field}</Badge>
									<ArrowRightOutline class="h-3 w-3 text-slate-300" />
								</div>
							</TableBodyCell>
							<TableBodyCell class="font-mono text-sm leading-relaxed break-all text-slate-600">
								{#if entry.Value}
									<span class="rounded border border-slate-200 bg-slate-100 px-2 py-1"
										>{entry.Value}</span
									>
								{:else}
									<span class="text-slate-300 italic">-- cleared --</span>
								{/if}
							</TableBodyCell>
						</TableBodyRow>
					{:else}
						<TableBodyRow>
							<TableBodyCell colspan="4" class="text-center py-20 text-slate-400">
								<ClockOutline class="w-12 h-12 mx-auto mb-4 opacity-20" />
								<p class="italic">No history records found for this member.</p>
							</TableBodyCell>
						</TableBodyRow>
					{/each}
				</TableBody>
			</Table>
		</Card>
	{/if}
</div>

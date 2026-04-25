<script lang="ts">
	import { getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { List, Li, Heading, Card, Button } from 'flowbite-svelte';
	import {
		FileChartBarOutline,
		FileCsvOutline,
		DownloadOutline,
		MailBoxOutline,
		HeartOutline
	} from 'flowbite-svelte-icons';
	import { report } from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	const oo = getContext('oo');

	// Guard: Redirect if not logged in
	if (!oo.me) {
		push('/Login');
	}

	let downloading = $state(false);

	async function runReport(name: string) {
		downloading = true;
		try {
			await report(name);
			toast.push(`Generating ${name} report...`);
		} catch (err: any) {
			toast.push(`Failed: ${err.message}`);
		} finally {
			downloading = false;
		}
	}
</script>

<svelte:head>
	<title>Administrative Reports | OSL</title>
</svelte:head>

<div class="mx-auto max-w-4xl space-y-6 px-4 py-8">
	<header class="flex items-center gap-4 border-b pb-6">
		<div class="bg-primary-100 text-primary-600 rounded-lg p-3">
			<FileChartBarOutline class="h-8 w-8" />
		</div>
		<div>
			<Heading tag="h2" class="text-3xl font-bold text-slate-900">Reports & Exports</Heading>
			<p class="text-slate-500">
				Generate CSV or PDF data for mail merges and administrative tracking.
			</p>
		</div>
	</header>

	<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
		<Card size="none" class="space-y-4 p-6">
			<Heading tag="h4" class="flex items-center gap-2 text-lg font-semibold">
				<FileCsvOutline class="h-5 w-5 text-blue-500" /> Maintenance & Cleanup
			</Heading>
			<List class="w-full">
				<Li
					onclick={() => runReport('expired')}
					class="flex cursor-pointer items-center justify-between"
				>
					<div>
						<span class="font-medium text-red-600">Expired Members</span>
						<p class="text-xs text-slate-400">Active status but no renewal for 2+ years</p>
					</div>
					<DownloadOutline class="h-4 w-4" />
				</Li>
				<Li
					onclick={() => runReport('email')}
					class="flex cursor-pointer items-center justify-between"
				>
					All Member Email List
					<DownloadOutline class="h-4 w-4" />
				</Li>
			</List>
		</Card>

		<Card size="none" class="space-y-4 p-6">
			<Heading tag="h4" class="flex items-center gap-2 text-lg font-semibold">
				<MailBoxOutline class="h-5 w-5 text-indigo-500" /> Mail Merges & Labels
			</Heading>
			<List class="w-full">
				<Li onclick={() => runReport('avery')} class="cursor-pointer">Avery Labels (All)</Li>
				<Li onclick={() => runReport('annual')} class="cursor-pointer">Annual Vows Merge</Li>
				<Li onclick={() => runReport('life')} class="cursor-pointer">Life Vows Merge</Li>
				<Li onclick={() => runReport('reaffirmation')} class="cursor-pointer"
					>Reaffirmation Form Merge</Li
				>
				<Li onclick={() => runReport('lifecheckin')} class="cursor-pointer">Life Vows Check-in</Li>
			</List>
		</Card>

		<Card size="none" class="space-y-4 p-6">
			<Heading tag="h4" class="flex items-center gap-2 text-lg font-semibold text-slate-800">
				Doxology & Subscriptions
			</Heading>
			<div class="flex flex-col gap-2">
				<Button
					color="alternative"
					outline
					onclick={() => runReport('doxprint')}
					class="justify-between text-left"
				>
					Doxology Mailing List (Full) <DownloadOutline class="ml-2 h-4 w-4" />
				</Button>
				<Button
					color="alternative"
					outline
					onclick={() => runReport('allsubscribers')}
					class="justify-between text-left"
				>
					All Subscribers Only <DownloadOutline class="ml-2 h-4 w-4" />
				</Button>
			</div>
		</Card>

		<Card size="none" class="space-y-4 border-l-4 border-l-pink-400 bg-pink-50/30 p-6">
			<Heading tag="h4" class="flex items-center gap-2 text-lg font-semibold text-pink-700">
				<HeartOutline class="h-5 w-5" /> Prayer Office
			</Heading>
			<p class="text-sm text-slate-600">
				Export prepared specifically for the Monthly Prayer List.
			</p>
			<Button color="purple" class="w-full shadow-sm" onclick={() => runReport('barb')}>
				<DownloadOutline class="mr-2 h-4 w-4" /> Sr. Barb's Monthly Prayer List
			</Button>
		</Card>
	</div>
</div>

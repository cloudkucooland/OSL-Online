<script lang="ts">
	import { onMount, getContext } from 'svelte';
	import { Card, Heading, Spinner, Badge } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { getSubscriber, updateSubscriber, cleanDateFormat } from '../oo';
	import { push } from 'svelte-spa-router';

	// Shared Components
	import AddressSection from '$lib/AddressSection.svelte';
	import ContactSection from '$lib/ContactSection.svelte';
	import FulfillmentSection from '$lib/FulfillmentSection.svelte';

	const oo = getContext('oo');

	// Guard: Ensure user is logged in
	if (!oo.me) {
		push('/Login');
	}

	let { params } = $props();

	// State Runes
	let loading = $state(true);
	let subscriber = $state(null);

	onMount(async () => {
		try {
			subscriber = await getSubscriber(params.id);
		} catch (err: any) {
			toast.push(`Failed to load: ${err.message}`);
		} finally {
			loading = false;
		}
	});

	async function handleUpdate(id: string, value: any) {
		let finalValue = value;

		// Clean date if it's the payment field
		if (id === 'DatePaid') {
			finalValue = cleanDateFormat(value);
		}

		try {
			await updateSubscriber(params.id, id, finalValue);
			toast.push(`Updated ${id}`);
		} catch (err: any) {
			toast.push(`Failed to update: ${err.message}`);
			console.error(err);
		}
	}
</script>

<svelte:head>
	<title>Subscriber: {subscriber ? subscriber.Name : 'Loading'}</title>
</svelte:head>

<div class="mx-auto w-full space-y-8 px-4 py-8">
	{#if loading}
		<div class="flex justify-center py-20">
			<Spinner color="purple" size="12" />
		</div>
	{:else if subscriber}
		<header
			class="flex flex-col items-center justify-between gap-6 rounded-xl border border-slate-200 bg-white p-6 shadow-sm md:flex-row"
		>
			<div class="flex items-center gap-4">
				<div>
					<Heading tag="h2" class="text-3xl font-bold text-slate-900"
						>{subscriber.Name || 'Institutional Subscriber'}</Heading
					>
					<div class="mt-1 flex items-center gap-2">
						<Badge color="indigo" class="text-[10px] uppercase">Institutional Subscriber</Badge>
						<span class="font-mono text-xs text-slate-400">ID: {subscriber.ID}</span>
					</div>
				</div>
			</div>
		</header>

		<div class="grid grid-cols-1 items-start gap-8 lg:grid-cols-3">
			<Card size="none" class="h-full border-slate-200 p-6 shadow-sm">
				<AddressSection
					data={subscriber}
					onUpdate={handleUpdate}
					showPrivacy={false}
					title="Mailing Address"
				/>
			</Card>

			<Card size="none" class="h-full border-slate-200 p-6 shadow-sm">
				<ContactSection
					data={subscriber}
					onUpdate={handleUpdate}
					showPrivacy={false}
					title="Contact Information"
				/>
			</Card>

			<FulfillmentSection data={subscriber} onUpdate={handleUpdate} showFinancials={true}>
				<div class="space-y-4">
					<div>
						<label for="DatePaid" class="mb-1 block text-xs font-bold text-slate-700"
							>Last Payment Date</label
						>
						<input
							id="DatePaid"
							class="focus:border-primary-500 focus:ring-primary-500 w-full rounded-lg border border-slate-300 bg-white p-2.5 text-sm"
							value={subscriber.DatePaid}
							onchange={(e) => handleUpdate('DatePaid', e.target.value)}
							placeholder="YYYY-MM-DD"
						/>
						<p class="mt-1 text-[10px] text-slate-400 italic">
							Subscription renewals are tracked by this date.
						</p>
					</div>
				</div>
			</FulfillmentSection>
		</div>
	{/if}
</div>

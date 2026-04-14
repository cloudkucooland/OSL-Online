<script lang="ts">
	import { onMount } from 'svelte';
	import { Label, Input, Select, Card, Heading, Spinner, Hr } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { getSubscriber, updateSubscriber, cleanDateFormat } from '../oo';

	let { params } = $props();

	// State Runes
	let loading = $state(true);
	let subscriber = $state(null);

	const commItems = [
		{ value: 'none', name: 'None' },
		{ value: 'mailed', name: 'Mailed' },
		{ value: 'electronic', name: 'Electronic' }
	];

	onMount(async () => {
		try {
			subscriber = await getSubscriber(params.id);
		} catch (err: any) {
			toast.push(`Failed to load: ${err.message}`);
		} finally {
			loading = false;
		}
	});

	async function handleUpdate(event: any) {
		const { id, value } = event.target;
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

	function preventDefault(event: Event) {
		event.preventDefault();
	}
</script>

<svelte:head>
	<title>Subscriber: {subscriber ? subscriber.Name : 'Loading'}</title>
</svelte:head>

{#if loading}
	<div class="flex justify-center py-20">
		<Spinner color="purple" size="12" />
	</div>
{:else if subscriber}
	<div class="mx-auto max-w-4xl space-y-6 pb-12">
		<header class="flex items-center justify-between">
			<Heading tag="h2" class="text-2xl font-bold text-slate-900">
				{subscriber.Name || 'New Subscriber'}
			</Heading>
		</header>

		<Card size="none" class="border-slate-200 p-6 shadow-sm">
			<form onsubmit={preventDefault} class="space-y-8">
				<section>
					<Heading
						tag="h5"
						class="mb-4 text-sm font-semibold uppercase tracking-wider text-slate-500"
						>Institutional Identity</Heading
					>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						<div>
							<Label for="Name" class="mb-2">Institution Name</Label>
							<Input id="Name" value={subscriber.Name} onchange={handleUpdate} />
						</div>
						<div>
							<Label for="Attn" class="mb-2">Attention To</Label>
							<Input id="Attn" value={subscriber.Attn} onchange={handleUpdate} />
						</div>
					</div>
				</section>

				<Hr />

				<section>
					<Heading
						tag="h5"
						class="mb-4 text-sm font-semibold uppercase tracking-wider text-slate-500"
						>Mailing Address</Heading
					>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-4">
						<div class="md:col-span-4">
							<Label for="Address" class="mb-2">Street Address</Label>
							<Input id="Address" value={subscriber.Address} onchange={handleUpdate} class="mb-2" />
							<Input
								id="AddressLine2"
								value={subscriber.AddressLine2}
								onchange={handleUpdate}
								placeholder="Unit, Suite, or Dept."
							/>
						</div>
						<div class="md:col-span-1">
							<Label for="City" class="mb-2">City</Label>
							<Input id="City" value={subscriber.City} onchange={handleUpdate} />
						</div>
						<div class="md:col-span-1">
							<Label for="State" class="mb-2">State/Locality</Label>
							<Input id="State" value={subscriber.State} onchange={handleUpdate} />
						</div>
						<div class="md:col-span-1">
							<Label for="Country" class="mb-2">Country</Label>
							<Input id="Country" value={subscriber.Country} onchange={handleUpdate} />
						</div>
						<div class="md:col-span-1">
							<Label for="PostalCode" class="mb-2">Postal Code</Label>
							<Input id="PostalCode" value={subscriber.PostalCode} onchange={handleUpdate} />
						</div>
					</div>
				</section>

				<Hr />

				<section class="grid grid-cols-1 gap-8 md:grid-cols-2">
					<div class="space-y-4">
						<Heading tag="h5" class="text-sm font-semibold uppercase tracking-wider text-slate-500"
							>Contact Details</Heading
						>
						<div>
							<Label for="PrimaryEmail" class="mb-1 text-xs">Primary Email</Label>
							<Input id="PrimaryEmail" value={subscriber.PrimaryEmail} onchange={handleUpdate} />
						</div>
						<div>
							<Label for="PrimaryPhone" class="mb-1 text-xs">Primary Phone</Label>
							<Input id="PrimaryPhone" value={subscriber.PrimaryPhone} onchange={handleUpdate} />
						</div>
					</div>
					<div class="space-y-4 rounded-lg border border-slate-100 bg-slate-50 p-4">
						<Heading tag="h5" class="text-sm font-semibold uppercase tracking-wider text-slate-500"
							>Financials</Heading
						>
						<div>
							<Label for="DatePaid" class="mb-2">Last Payment Date</Label>
							<Input id="DatePaid" value={subscriber.DatePaid} onchange={handleUpdate} />
							<p class="mt-2 text-xs italic text-slate-400">Format: YYYY-MM-DD</p>
						</div>
					</div>
				</section>

				<Hr />

				<section>
					<Heading
						tag="h5"
						class="mb-4 text-sm font-semibold uppercase tracking-wider text-slate-500"
						>Fulfillment Settings</Heading
					>
					<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
						<div>
							<Label for="Newsletter" class="mb-2">Newsletter</Label>
							<Select
								id="Newsletter"
								items={commItems}
								value={subscriber.Newsletter}
								onchange={handleUpdate}
							/>
						</div>
						<div>
							<Label for="Doxology" class="mb-2 font-bold text-primary-700 underline"
								>Doxology Journal</Label
							>
							<Select
								id="Doxology"
								items={commItems}
								value={subscriber.Doxology}
								onchange={handleUpdate}
							/>
						</div>
						<div>
							<Label for="Communication" class="mb-2">Communication</Label>
							<Select
								id="Communication"
								items={commItems}
								value={subscriber.Communication}
								onchange={handleUpdate}
							/>
						</div>
					</div>
				</section>
			</form>
		</Card>
	</div>
{/if}

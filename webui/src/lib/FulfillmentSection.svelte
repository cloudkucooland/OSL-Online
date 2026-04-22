<script lang="ts">
	import { Label, Select, Card, Heading } from 'flowbite-svelte';

	let { 
		data, 
		onUpdate, 
		disabled = false, 
		title = 'Fulfillment',
		showFinancials = false,
		financialsTitle = 'Financials',
		children
	} = $props();

	const newsletterItems = [
		{ value: 'none', name: 'None' },
		{ value: 'electronic', name: 'Electronic' }
	];

	const commItems = [
		{ value: 'none', name: 'None' },
		{ value: 'mailed', name: 'Mailed' },
		{ value: 'electronic', name: 'Electronic' }
	];

	function handleLocalUpdate(event: any) {
		const { id, value } = event.target;
		onUpdate(id, value);
	}
</script>

<div class="space-y-8">
	<Card size="none" class="border-slate-200 p-6 shadow-sm">
		<Heading tag="h4" class="mb-6 border-b pb-2 text-lg font-bold uppercase text-slate-800">{title}</Heading>
		<div class="space-y-4">
			<div>
				<Label for="Newsletter" class="mb-1 font-bold text-red-700">Newsletter</Label>
				<Select
					id="Newsletter"
					items={newsletterItems}
					value={data.Newsletter}
					onchange={handleLocalUpdate}
					{disabled}
				/>
			</div>
			<div>
				<Label for="Doxology" class="mb-1 font-bold text-red-700">Doxology</Label>
				<Select
					id="Doxology"
					items={commItems}
					value={data.Doxology}
					onchange={handleLocalUpdate}
					{disabled}
				/>
			</div>
			<div>
				<Label for="Communication" class="mb-1 font-bold text-red-700">General Communication</Label>
				<Select
					id="Communication"
					items={commItems}
					value={data.Communication}
					onchange={handleLocalUpdate}
					{disabled}
				/>
			</div>
		</div>
	</Card>

	{#if showFinancials}
		<Card size="none" class="border-slate-200 bg-slate-50/50 p-6 shadow-sm">
			<Heading tag="h4" class="mb-6 border-b pb-2 text-lg font-bold uppercase text-slate-800">{financialsTitle}</Heading>
			{@render children?.()}
		</Card>
	{/if}
</div>

<script lang="ts">
	import { Label, Input, Toggle, Heading } from 'flowbite-svelte';

	let { 
		data, 
		onUpdate, 
		disabled = false, 
		showPrivacy = false,
		privacyField = 'ListAddress',
		title = 'Mailing Address'
	} = $props();

	function handleLocalUpdate(event: any) {
		const { id, value, type, checked } = event.target;
		const finalValue = type === 'checkbox' ? checked : value;
		onUpdate(id, finalValue);
	}
</script>

<div class="mb-6 flex items-center justify-between border-b pb-2">
	<Heading tag="h4" class="text-lg font-bold uppercase text-slate-800">{title}</Heading>
	{#if showPrivacy}
		<Toggle
			id={privacyField}
			checked={data[privacyField]}
			onchange={handleLocalUpdate}
			color="red"
			size="small"
		/>
	{/if}
</div>

<div class="space-y-4">
	<div class="space-y-2">
		<Input
			id="Address"
			value={data.Address}
			onchange={handleLocalUpdate}
			placeholder="Address Line 1"
			{disabled}
		/>
		<Input
			id="AddressLine2"
			value={data.AddressLine2}
			onchange={handleLocalUpdate}
			placeholder="Address Line 2"
			{disabled}
		/>
		<div class="grid grid-cols-2 gap-2">
			<Input id="City" value={data.City} onchange={handleLocalUpdate} placeholder="City" {disabled} />
			<Input id="State" value={data.State} onchange={handleLocalUpdate} placeholder="State" {disabled} />
		</div>
		<div class="grid grid-cols-2 gap-2">
			<Input id="Country" value={data.Country} onchange={handleLocalUpdate} placeholder="Country" {disabled} />
			<Input
				id="PostalCode"
				value={data.PostalCode}
				onchange={handleLocalUpdate}
				placeholder="Postal Code"
				{disabled}
			/>
		</div>
	</div>

	{#if data.FormattedAddr}
		<div class="mt-6 rounded-lg border border-slate-100 bg-slate-50 p-4">
			<Label class="mb-2 text-xs font-bold uppercase text-slate-400">Formatted Label</Label>
			<p class="whitespace-pre-line font-mono text-xs leading-relaxed text-slate-600">
				{data.FormattedAddr}
			</p>
		</div>
	{/if}
</div>

<script lang="ts">
	import { Label, Input, Toggle, Heading, Hr } from 'flowbite-svelte';

	let {
		data,
		onUpdate,
		disabled = false,
		showPrivacy = false,
		title = 'Contact Information'
	} = $props();

	function handleLocalUpdate(event: any) {
		const { id, value, type, checked } = event.target;
		const finalValue = type === 'checkbox' ? checked : value;
		onUpdate(id, finalValue);
	}
</script>

<div class="mb-6 flex items-center justify-between border-b pb-2">
	<Heading tag="h4" class="text-lg font-bold text-slate-800 uppercase">{title}</Heading>
</div>

<div class="space-y-6">
	<!-- Primary Contact -->
	<div class="space-y-4">
		<div class="flex items-center gap-2">
			<div class="flex-grow">
				<Label class="text-xs italic">Primary Email</Label>
				<Input
					id="PrimaryEmail"
					value={data.PrimaryEmail}
					onchange={handleLocalUpdate}
					{disabled}
				/>
			</div>
			{#if showPrivacy}
				<Toggle
					id="ListPrimaryEmail"
					checked={data.ListPrimaryEmail}
					onchange={handleLocalUpdate}
					color="red"
					class="mt-5"
				/>
			{/if}
		</div>
		<div class="flex items-center gap-2">
			<div class="flex-grow">
				<Label class="text-xs italic">Primary Phone</Label>
				<Input
					id="PrimaryPhone"
					value={data.PrimaryPhone}
					onchange={handleLocalUpdate}
					{disabled}
				/>
			</div>
			{#if showPrivacy}
				<Toggle
					id="ListPrimaryPhone"
					checked={data.ListPrimaryPhone}
					onchange={handleLocalUpdate}
					color="red"
					class="mt-5"
				/>
			{/if}
		</div>
	</div>

	<Hr class="my-6" />

	<!-- Secondary Contact -->
	<div class="space-y-4">
		<div class="flex items-center gap-2">
			<div class="flex-grow">
				<Label class="text-xs text-slate-400 italic">Secondary Email</Label>
				<Input
					id="SecondaryEmail"
					value={data.SecondaryEmail}
					onchange={handleLocalUpdate}
					{disabled}
				/>
			</div>
			{#if showPrivacy}
				<Toggle
					id="ListSecondaryEmail"
					checked={data.ListSecondaryEmail}
					onchange={handleLocalUpdate}
					color="red"
					class="mt-5"
				/>
			{/if}
		</div>
		<div class="flex items-center gap-2">
			<div class="flex-grow">
				<Label class="text-xs text-slate-400 italic">Secondary Phone</Label>
				<Input
					id="SecondaryPhone"
					value={data.SecondaryPhone}
					onchange={handleLocalUpdate}
					{disabled}
				/>
			</div>
			{#if showPrivacy}
				<Toggle
					id="ListSecondaryPhone"
					checked={data.ListSecondaryPhone}
					onchange={handleLocalUpdate}
					color="red"
					class="mt-5"
				/>
			{/if}
		</div>
	</div>
</div>

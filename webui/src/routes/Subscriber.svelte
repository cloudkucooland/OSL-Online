<script lang="ts">
	import { getSubscriber, updateSubscriber, cleanDateFormat } from '../oo';
	import { Label, Input, Select } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';

	let { params } = $props();

	const commitems = [
		{ value: 'none', name: 'None' },
		{ value: 'mailed', name: 'Mailed' },
		{ value: 'electronic', name: 'Electronic' }
	];

	async function change(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			await updateSubscriber(params.id, event.target.id, event.target.value);
			toast.push(`Changed ${event.target.id}`);
		} catch (err) {
			toast.push('failed to change: ' + err.message);
			console.log(err);
		}
		return true;
	}

	async function changeDate(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			await updateSubscriber(params.id, event.target.id, cleanDateFormat(event.target.value));
			toast.push(`Changed ${event.target.id}`);
		} catch (err) {
			toast.push('failed to change: ' + err.message);
			console.log(err);
		}
		return true;
	}

	async function reload(event) {
		event.preventDefault();
		event.stopPropagation();
		// do nothing
		return true;
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Subscriber</title>
</svelte:head>

{#await getSubscriber(params.id)}
	<h3>... loading ...</h3>
{:then r}
	<form onsubmit={reload}>
		<section>
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-4">
					<Label for="tName" class="block">Name</Label>
					<Input id="Name" value={r.Name} onchange={change} />
				</div>
				<div class="col-span-4">
					<Label for="Attn" class="block">Attn</Label>
					<Input id="Attn" value={r.Attn} onchange={change} />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-8">
					<Label for="Address" class="block">Address</Label>
					<Input id="Address" value={r.Address} onchange={change} />
				</div>
				<div class="col-span-8">
					<Input id="AddressLine2" value={r.AddressLine2} onchange={change} />
				</div>
				<div class="col-span-2">
					<Label for="City" class="block">City</Label>
					<Input id="City" value={r.City} onchange={change} />
				</div>
				<div class="col-span-2">
					<Label for="State" class="block">State/Locality</Label>
					<Input id="State" value={r.State} onchange={change} />
				</div>
				<div class="col-span-2">
					<Label for="Country" class="block">Country</Label>
					<Input id="Country" value={r.Country} onchange={change} />
				</div>
				<div class="col-span-2">
					<Label for="PostalCode" class="block">Postal Code</Label>
					<Input id="PostalCode" value={r.PostalCode} onchange={change} />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-4">
					<Label for="PrimaryPhone" class="block">Primary Phone</Label>
					<Input id="PrimaryPhone" value={r.PrimaryPhone} onchange={change} />
				</div>
				<div class="col-span-4">
					<Label for="SecondaryPhone" class="block">Secondary Phone</Label>
					<Input id="SecondaryPhone" value={r.SecondaryPhone} onchange={change} />
				</div>
				<div class="col-span-4">
					<Label for="PrimaryEmail" class="block">Primary Email</Label>
					<Input id="PrimaryEmail" value={r.PrimaryEmail} onchange={change} />
				</div>
				<div class="col-span-4">
					<Label for="SecondaryEmail" class="block">Secondary Email</Label>
					<Input id="SecondaryEmail" value={r.SecondaryEmail} onchange={change} />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-2">
					<Label for="DatePaid" class="block">Paid</Label>
					<Input id="DatePaid" value={r.DatePaid} onchange={changeDate} />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-1">
					<Label for="Newsletter" class="block">Newsletter</Label>
					<Select id="Newsletter" items={commitems} value={r.Newsletter} onchange={change} />
				</div>
				<div class="col-span-1">
					<Label for="Doxology" class="block">Doxology</Label>
					<Select id="Doxology" items={commitems} value={r.Doxology} onchange={change} />
				</div>
				<div class="col-span-1">
					<Label for="Communication" class="block">Communication</Label>
					<Select id="Communication" items={commitems} value={r.Communication} onchange={change} />
				</div>
			</div>
		</section>
	</form>
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}

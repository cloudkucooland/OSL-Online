<script lang="ts">
	import { getContext } from 'svelte';
	import {
		getMeFromServer,
		updateMe,
		getChapters,
		updateMeChapters,
		getMeChapters,
		oslname
	} from '../oo';
	import { Label, Input, Toggle, Select, MultiSelect } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}
	let chaps = $state([]);
	let selectedchapters = $state([]);

	const cannotedit = true; // placeholder for fields we are considering enabling

	const commitems = [
		{ value: 'none', name: 'None' },
		{ value: 'mailed', name: 'Mailed' },
		{ value: 'electronic', name: 'Electronic' }
	];

	const memberstatus = [
		{ value: 'Annual Vows', name: 'Annual Vows' },
		{ value: 'Life Vows', name: 'Life Vows' },
		{ value: 'Removed', name: 'Removed' },
		{ value: 'Friend', name: 'Friend' },
		{ value: 'Benefactor', name: 'Benefactor' }
	];

	const titles = [
		{ value: ' ', name: '' },
		{ value: 'Sr.', name: 'Sr.' },
		{ value: 'Br.', name: 'Br.' },
		{ value: 'Sibling', name: 'Sibling' }
	];

	const removereasons = [
		{ value: ' ', name: '' },
		{ value: 'Bad Address', name: 'Bad Address' },
		{ value: 'Death', name: 'Death' },
		{ value: 'No Reaffirmation', name: 'No Reaffirmation' },
		{ value: 'Request', name: 'Request' }
	];

	const stati = [
		{ value: 'laity', name: 'Laity' },
		{ value: 'clergy', name: 'Clergy' },
		{ value: 'student', name: 'Student' },
		{ value: 'retired laity', name: 'Retired (laity)' },
		{ value: 'retired clergy', name: 'Retired (clergy)' }
	];

	const leadership = [
		{ value: 'member', name: 'Member' },
		{ value: 'prior', name: 'Prior' },
		{ value: 'council', name: 'General Council Officer' },
		{ value: 'canon', name: 'Canon' },
		{ value: 'elected', name: 'Elected Officer' }
	];

	async function load() {
		const m = await getMeFromServer();
		chaps = await getChapters();
		selectedchapters = await getMeChapters();

		return m;
	}

	// useless
	async function reload(event) {
		event.preventDefault();
		event.stopPropagation();
		return await load();
	}

	async function change(event) {
		try {
			await updateMe(event.target.id, event.target.value);
			toast.push(`Changed ${event.target.id}`);
		} catch (err) {
			toast.push('failed to change: ' + err.message);
			console.log(err);
		}
		return true;
	}

	async function changeCheck(event) {
		try {
			await updateMe(event.target.id, event.target.checked);
			toast.push(`Changed ${event.target.id}`);
		} catch (err) {
			toast.push('failed to change: ' + err.message);
			console.log(err);
		}
		return true;
	}

	async function setchapters() {
		try {
			await updateMeChapters(selectedchapters);
			toast.push(`Updated Chapters`);
		} catch (err) {
			toast.push('failed to set chapter: ' + err.message);
			console.log(err);
		}
		return true;
	}
</script>

<svelte:head>
	<title>OSL Member Manager: My Record</title>
</svelte:head>

{#await load()}
	<h3>... loading ...</h3>
{:then r}
	<div>
		{oslname(r)}
		<Toggle id="ListInDirectory" checked={r.ListInDirectory} onchange={changeCheck} color="red"
			><span style="color: red">List in Directory</span></Toggle
		>
	</div>

	<form onsubmit={reload}>
		<section>
			<h3>
				If any data is out-of date please use the <a href="https://saint-luke.net/reaffirmation"
					>reaffirmation form</a
				> to have the Chancellor-General update it.
			</h3>
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-3">
					<Label for="FirstName" class="block">First Name</Label>
					<Input id="FirstName" value={r.FirstName} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-1">
					<Label for="MiddleName" class="block">Middle Name</Label>
					<Input id="MiddleName" value={r.MiddleName} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-3">
					<Label for="LastName" class="block">Last Name</Label>
					<Input id="LastName" value={r.LastName} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-1">
					<Label for="Suffix" class="block">Suffix</Label>
					<Input id="Suffix" value={r.Suffix} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="LifeVowName" class="block">Life Vow Name</Label>
					<Input id="LifeVowName" value={r.LifeVowName} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="PreferedName" class="block" style="color: red"
						>Preferred Name (If different from First Name)</Label
					>
					<Input id="PreferredName" value={r.PreferredName} onchange={change} />
				</div>
				<div class="col-span-2">
					<Label for="DateReaffirmation" class="block">Last Reffirmation</Label>
					<Input id="DateReaffirmation" value={r.DateReaffirmation} disabled={true} />
				</div>
				<div class="col-span-1">
					<Label for="MemberStatus" class="block">Member Status</Label>
					<Select id="MemberStatus" items={memberstatus} value={r.MemberStatus} disabled={true} />
				</div>
				<div class="col-span-1">
					<Label for="Title" class="block">Title</Label>
					<Select
						id="Title"
						items={titles}
						value={r.Title}
						onchange={change}
						disabled={r.MemberStatus == 'Friend'}
					/>
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-7">
					<Label for="Address" class="block">Address</Label>
					<Input id="Address" value={r.Address} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-1">
					<Toggle id="ListAddress" checked={r.ListAddress} onchange={changeCheck} color="red"
						><span style="color: red">Listed</span></Toggle
					>
				</div>
				<div class="col-span-8">
					<Input id="AddressLine2" value={r.AddressLine2} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="City" class="block">City</Label>
					<Input id="City" value={r.City} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="State" class="block">State/Locality</Label>
					<Input id="State" value={r.State} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="Country" class="block">Country</Label>
					<Input id="Country" value={r.Country} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="PostalCode" class="block">Postal Code</Label>
					<Input id="PostalCode" value={r.PostalCode} onchange={change} disabled={cannotedit} />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-3">
					<Label for="PrimaryPhone" class="block">Primary Phone</Label>
					<Input id="PrimaryPhone" value={r.PrimaryPhone} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-1">
					<Toggle
						id="ListPrimaryPhone"
						checked={r.ListPrimaryPhone}
						onchange={changeCheck}
						color="red"><span style="color: red">Listed</span></Toggle
					>
				</div>
				<div class="col-span-3">
					<Label for="SecondaryPhone" class="block">Secondary Phone</Label>
					<Input
						id="SecondaryPhone"
						value={r.SecondaryPhone}
						onchange={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-1">
					<Toggle
						id="ListSecondaryPhone"
						checked={r.ListSecondaryPhone}
						onchange={changeCheck}
						color="red"><span style="color: red">Listed</span></Toggle
					>
				</div>
				<div class="col-span-3">
					<Label for="PrimaryEmail" class="block">Primary Email</Label>
					<Input id="PrimaryEmail" value={r.PrimaryEmail} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-1">
					<Toggle
						id="ListPrimaryEmail"
						checked={r.ListPrimaryEmail}
						onchange={changeCheck}
						color="red"><span style="color: red">Listed</span></Toggle
					>
				</div>
				<div class="col-span-3">
					<Label for="SecondaryEmail" class="block">Secondary Email</Label>
					<Input
						id="SecondaryEmail"
						value={r.SecondaryEmail}
						onchange={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-1">
					<Toggle
						id="ListSecondaryEmail"
						checked={r.ListSecondaryEmail}
						onchange={changeCheck}
						color="red"><span style="color: red">Listed</span></Toggle
					>
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-2">
					<Label for="DateFirstVows" class="block">First Vows</Label>
					<Input id="DateFirstVows" value={r.DateFirstVows} disabled={true} />
				</div>
				<div class="col-span-2">
					<Label for="DateNovitiate" class="block">Novitiate</Label>
					<Input id="DateNovitiate" value={r.DateNovitiate} disabled={true} />
				</div>
				<div class="col-span-2">
					<Label for="BirthDate" class="block">Birth Day</Label>
					<Input id="BirthDate" value={r.BirthDate} onchange={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="DateDeceased" class="block">Deceased</Label>
					<Input id="DateDeceased" value={r.DateDeceased} disabled={true} />
				</div>
				<div class="col-span-2">
					<Label for="DateLifeVows" class="block">Life Vows</Label>
					<Input id="DateLifeVows" value={r.DateLifeVows} disabled={true} />
				</div>
				<div class="col-span-2">
					<Label for="DateRecordCreated" class="block">Record Created</Label>
					<Input id="DateRecordCreated" value={r.DateRecordCreated} disabled={true} />
				</div>
				<div class="col-span-2">
					<Label for="DateRemoved" class="block">Removed</Label>
					<Input id="DateRemoved" value={r.DateRemoved} disabled={true} />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-2">
					<Label for="HowJoined" class="block">How Joined</Label>
					<Input id="HowJoined" value={r.HowJoined} disabled={true} />
				</div>
				<div class="col-span-2">
					<Label for="HowRemoved" class="block">How Removed</Label>
					<Select id="HowRemoved" items={removereasons} value={r.HowRemoved} disabled={true} />
				</div>
				<div class="col-span-2">
					<Label for="Status" class="block">Status</Label>
					<Select
						id="Status"
						items={stati}
						value={r.Status}
						onchange={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-2">
					<Label for="Leadership" class="block">Leadership</Label>
					<Select id="Leadership" items={leadership} value={r.Leadership} disabled={true} />
				</div>
				<div class="col-span-2">
					<Label for="Chapters" class="block" style="color: red">Chapters</Label>
					<MultiSelect
						id="Chapters"
						items={chaps}
						bind:value={selectedchapters}
						onchange={setchapters}
					/>
				</div>
				<div class="col-span-2">
					<Label for="Occupation" class="block" style="color: red">Occupation</Label>
					<Input id="Occupation" value={r.Occupation} onchange={change} />
				</div>
				<div class="col-span-2">
					<Label for="Employer" class="block" style="color: red">Employer</Label>
					<Input id="Employer" value={r.Employer} onchange={change} />
				</div>
				<div class="col-span-2">
					<Label for="Denomination" class="block" style="color: red">Denomination</Label>
					<Input id="Denomination" value={r.Denomination} onchange={change} />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<h3>
				You can only choose "Mailed" if you have donated in the past 12 months. These revert to
				"Electronic" if you have not donated in the past year.
			</h3>
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-1">
					<Label for="Newsletter" class="block" style="color: red">Newsletter</Label>
					<Select id="Newsletter" items={commitems} value={r.Newsletter} onchange={change} />
				</div>
				<div class="col-span-1">
					<Label for="Doxology" class="block" style="color: red">Doxology</Label>
					<Select id="Doxology" items={commitems} value={r.Doxology} onchange={change} />
				</div>
				<div class="col-span-1">
					<Label for="Communication" class="block" style="color: red">Communication</Label>
					<Select id="Communication" items={commitems} value={r.Communication} onchange={change} />
				</div>
			</div>
		</section>
	</form>
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}

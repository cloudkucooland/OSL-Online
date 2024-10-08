<script lang="ts">
	import { getContext } from 'svelte';
	import { getMe, getMeFromServer, updateMe } from '../oo';
	import { Label, Input, Checkbox, Select } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

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

	async function change(e) {
		try {
			await updateMe(e.target.id, e.target.value);
			toast.push(`Changed ${e.target.id}`);
			return true;
		} catch (err) {
			toast.push('failed to change: ' + err);
			console.log(err);
		}
	}

	async function changeCheck(e) {
		try {
			await updateMe(e.target.id, e.target.checked);
			toast.push(`Changed ${e.target.id}`);
			return true;
		} catch (err) {
			toast.push('failed to change: ' + err);
			console.log(err);
		}
	}
</script>

<svelte:head>
	<title>OSL Member Manager: My Record</title>
</svelte:head>

{#await getMeFromServer()}
	<h3>... loading ...</h3>
{:then r}
	<div>
		{r.Title}
		{r.FirstName}
		{r.LastName} ( {r.MemberStatus} )
		<Checkbox id="ListInDirectory" checked={r.ListInDirectory} on:change={changeCheck}
			>List in Directory</Checkbox
		>
	</div>

	<form>
		<section>
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-3">
					<Label for="FirstName" class="block">First Name</Label>
					<Input id="FirstName" value={r.FirstName} on:change={change} disabled={cannotedit} />
				</div>
				<div class="col-span-1">
					<Label for="MiddleName" class="block">Middle Name</Label>
					<Input id="MiddleName" value={r.MiddleName} on:change={change} disabled={cannotedit} />
				</div>
				<div class="col-span-3">
					<Label for="LastName" class="block">Last Name</Label>
					<Input id="LastName" value={r.LastName} on:change={change} disabled={cannotedit} />
				</div>
				<div class="col-span-1">
					<Label for="Suffix" class="block">Suffix</Label>
					<Input id="Suffix" value={r.Suffix} on:change={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="LifeVowName" class="block">Life Vow Name</Label>
					<Input id="LifeVowName" value={r.LifeVowName} on:change={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="PreferedName" class="block">Preferred Name</Label>
					<Input id="PreferredName" value={r.PreferredName} on:change={change} />
				</div>
				<div class="col-span-2">
					<Label for="DateReaffirmation" class="block">Last Reffirmation</Label>
					<Input id="DateReaffirmation" value={r.DateReaffirmation} disabled="true" />
				</div>
				<div class="col-span-1">
					<Label for="MemberStatus" class="block">Member Status</Label>
					<Select id="MemberStatus" items={memberstatus} value={r.MemberStatus} disabled="true" />
				</div>
				<div class="col-span-1">
					<Label for="Title" class="block">Title</Label>
					<Select id="Title" items={titles} value={r.Title} on:change={change} />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-7">
					<Label for="Address" class="block">Address</Label>
					<Input id="Address" value={r.Address} on:change={change} disabled={cannotedit} />
				</div>
				<div class="col-span-1">
					<Checkbox id="ListAddress" checked={r.ListAddress} on:change={changeCheck}
						>Listed</Checkbox
					>
				</div>
				<div class="col-span-8">
					<Input
						id="AddressLine2"
						value={r.AddressLine2}
						on:change={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-2">
					<Label for="City" class="block">City</Label>
					<Input id="City" value={r.City} on:change={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="State" class="block">State/Locality</Label>
					<Input id="State" value={r.State} on:change={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="Country" class="block">Country</Label>
					<Input id="Country" value={r.Country} on:change={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="PostalCode" class="block">Postal Code</Label>
					<Input id="PostalCode" value={r.PostalCode} on:change={change} disabled={cannotedit} />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-3">
					<Label for="PrimaryPhone" class="block">Primary Phone</Label>
					<Input
						id="PrimaryPhone"
						value={r.PrimaryPhone}
						on:change={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-1">
					<Checkbox id="ListPrimaryPhone" checked={r.ListPrimaryPhone} on:change={changeCheck}
						>Listed</Checkbox
					>
				</div>
				<div class="col-span-3">
					<Label for="SecondaryPhone" class="block">Secondary Phone</Label>
					<Input
						id="SecondaryPhone"
						value={r.SecondaryPhone}
						on:change={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-1">
					<Checkbox id="ListSecondaryPhone" checked={r.ListSecondaryPhone} on:change={changeCheck}
						>Listed</Checkbox
					>
				</div>
				<div class="col-span-3">
					<Label for="PrimaryEmail" class="block">Primary Email</Label>
					<Input
						id="PrimaryEmail"
						value={r.PrimaryEmail}
						on:change={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-1">
					<Checkbox id="ListPrimaryEmail" checked={r.ListPrimaryEmail} on:change={changeCheck}
						>Listed</Checkbox
					>
				</div>
				<div class="col-span-3">
					<Label for="SecondaryEmail" class="block">Secondary Email</Label>
					<Input
						id="SecondaryEmail"
						value={r.SecondaryEmail}
						on:change={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-1">
					<Checkbox id="ListSecondaryEmail" checked={r.ListSecondaryEmail} on:change={changeCheck}
						>Listed</Checkbox
					>
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-2">
					<Label for="DateFirstVows" class="block">First Vows</Label>
					<Input id="DateFirstVows" value={r.DateFirstVows} disabled="true" />
				</div>
				<div class="col-span-2">
					<Label for="DateNovitiate" class="block">Novitiate</Label>
					<Input id="DateNovitiate" value={r.DateNovitiate} disabled="true" />
				</div>
				<div class="col-span-2">
					<Label for="BirthDate" class="block">Birth Day</Label>
					<Input id="BirthDate" value={r.BirthDate} on:change={change} disabled={cannotedit} />
				</div>
				<div class="col-span-2">
					<Label for="DateDeceased" class="block">Deceased</Label>
					<Input id="DateDeceased" value={r.DateDeceased} disabled="true" />
				</div>
				<div class="col-span-2">
					<Label for="DateLifeVows" class="block">Life Vows</Label>
					<Input id="DateLifeVows" value={r.DateLifeVows} disabled="true" />
				</div>
				<div class="col-span-2">
					<Label for="DateRecordCreated" class="block">Record Created</Label>
					<Input id="DateRecordCreated" value={r.DateRecordCreated} disabled="true" />
				</div>
				<div class="col-span-2">
					<Label for="DateRemoved" class="block">Removed</Label>
					<Input id="DateRemoved" value={r.DateRemoved} disabled="true" />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-2">
					<Label for="HowJoined" class="block">How Joined</Label>
					<Input id="HowJoined" value={r.HowJoined} disabled="true" />
				</div>
				<div class="col-span-2">
					<Label for="HowRemoved" class="block">How Removed</Label>
					<Select id="HowRemoved" items={removereasons} value={r.HowRemoved} disabled="true" />
				</div>
				<div class="col-span-2">
					<Label for="Status" class="block">Status</Label>
					<Select
						id="Status"
						items={stati}
						value={r.Status}
						on:change={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-2">
					<Label for="Leadership" class="block">Leadership</Label>
					<Select id="Leadership" items={leadership} value={r.Leadership} disabled="true" />
				</div>
				<div class="col-span-2">
					<Label for="Chapter" class="block">Chapter</Label>
					<Input id="Chapter" value={r.Chapter} on:change={change} />
				</div>
				<div class="col-span-2">
					<Label for="Occupation" class="block">Occupation</Label>
					<Input id="Occupation" value={r.Occupation} on:change={change} />
				</div>
				<div class="col-span-2">
					<Label for="Employer" class="block">Employer</Label>
					<Input id="Employer" value={r.Employer} on:change={change} />
				</div>
				<div class="col-span-2">
					<Label for="Denomination" class="block">Denomination</Label>
					<Input id="Denomination" value={r.Denomination} on:change={change} />
				</div>
			</div>
		</section>

		<section>
			<hr class="px-4 py-2" />
			<div class="grid grid-cols-8 gap-4 px-4 py-2">
				<div class="col-span-1">
					<Label for="Newsletter" class="block">Newsletter</Label>
					<Select
						id="Newsletter"
						items={commitems}
						value={r.Newsletter}
						on:change={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-1">
					<Label for="Doxology" class="block">Doxology</Label>
					<Select
						id="Doxology"
						items={commitems}
						value={r.Doxology}
						on:change={change}
						disabled={cannotedit}
					/>
				</div>
				<div class="col-span-1">
					<Label for="Communication" class="block">Communication</Label>
					<Select
						id="Communication"
						items={commitems}
						value={r.Communication}
						on:change={change}
						disabled={cannotedit}
					/>
				</div>
			</div>
		</section>
	</form>
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}

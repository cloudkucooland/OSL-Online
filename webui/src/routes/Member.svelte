<script lang="ts">
	import { getContext } from 'svelte';
	import {
		getMember,
		updateMember,
		getMemberChapters,
		getChapters,
		updateMemberChapters,
		oslname,
		vcard
	} from '../oo';
	import { Button, Label, Input, Select, MultiSelect, Toggle } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}
	let chaps = $state([]);
	let selectedchapters = $state([]);

	let { params } = $props();

	const commitems = [
		{ value: 'none', name: 'None' },
		{ value: 'mailed', name: 'Mailed' },
		{ value: 'electronic', name: 'Electronic' }
	];

	// The newsletter is now digital-only
	const newsletteritems = [
		{ value: 'none', name: 'None' },
		{ value: 'electronic', name: 'Electronic' }
	];

	const memberstatus = [
		{ value: 'Annual Vows', name: 'Annual Vows' },
		{ value: 'Life Vows', name: 'Life Vows' },
		{ value: 'Removed', name: 'Removed' },
		{ value: 'Friend', name: 'Friend' },
		{ value: 'Deceased', name: 'Deceased' }
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
		{ value: ' ', name: '' },
		{ value: 'member', name: 'Member' },
		{ value: 'prior', name: 'Prior' },
		{ value: 'council', name: 'General Council Officer' },
		{ value: 'canon', name: 'Canon' },
		{ value: 'elected', name: 'Elected Officer' }
	];

	async function load(id) {
		const m = await getMember(id);
		chaps = await getChapters();
		selectedchapters = await getMemberChapters(id);

		return m;
	}

	// useless, but so when I grep for form actions, I see something
	async function reload(event) {
		event.preventDefault();
		event.stopPropagation();
		console.log('Member.svelte: reload()', event);
		load(params.id);
	}

	async function change(event) {
		event.preventDefault();
		event.stopPropagation();
		console.log('Member.svelte: change()', event);
		try {
			await updateMember(params.id, event.target.id, event.target.value);
			toast.push(`Changed ${event.target.id}`);
		} catch (err) {
			toast.push('failed to change: ' + err.message);
			console.log(err.message);
		}
		return true;
	}

	async function changeDate(event) {
		event.preventDefault();
		event.stopPropagation();
		console.log('Member.svelte: changeDate()', event);
		try {
			await updateMember(params.id, event.target.id, event.target.value);
			toast.push(`Changed ${event.target.id}`);
		} catch (err) {
			toast.push('failed to change: ' + err.message);
			console.log(err.message);
		}
		return true;
	}

	async function changeCheck(event) {
		event.preventDefault();
		event.stopPropagation();
		console.log('Member.svelte changeCheck()', event);
		try {
			await updateMember(params.id, event.target.id, event.target.checked);
			toast.push(`Changed ${event.target.id}`);
		} catch (err) {
			toast.push('failed to change: ' + err.message);
			console.log(err.message);
		}
		return true;
	}

	async function setchapters(event) {
		event.preventDefault();
		event.stopPropagation();
		console.log('Member.svelte setchapters()', event);
		try {
			await updateMemberChapters(params.id, selectedchapters);
			toast.push('Updated Chapters');
		} catch (err) {
			toast.push('failed to set chapter: ' + err.message);
			console.log(err);
		}
		return true;
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Member</title>
</svelte:head>

{#await load(params.id)}
	<h3>... loading ...</h3>
{:then r}
	{#if $me && $me.level > 0}
		<section>
			<div class="grid grid-cols-6 gap-4 px-4 py-2">
				<div class="col-span-2">
					{oslname(r)} ({r.MemberStatus})
				</div>
				<div class="col-span-1">
					<Toggle
						id="ListInDirectory"
						checked={r.ListInDirectory}
						onchange={changeCheck}
						color="red">List in Directory</Toggle
					>
				</div>
				<div class="col-span-1">
					<a href="#/giving/{r.ID}">Giving</a>
				</div>
				<div class="col-span-1">
					<a href="#/notes/{r.ID}">Notes</a>
				</div>
				<div class="col-span-1">
					<a href="#/changelog/{r.ID}">Changelog</a>
				</div>
			</div>
		</section>

		<form onsubmit={reload}>
			<section>
				<div class="grid grid-cols-8 gap-4 px-4 py-2">
					<div class="col-span-3">
						<Label for="FirstName" class="block">First Name</Label>
						<Input id="FirstName" value={r.FirstName} onchange={change} />
					</div>
					<div class="col-span-1">
						<Label for="MiddleName" class="block">Middle Name</Label>
						<Input id="MiddleName" value={r.MiddleName} onchange={change} />
					</div>
					<div class="col-span-3">
						<Label for="LastName" class="block">Last Name</Label>
						<Input id="LastName" value={r.LastName} onchange={change} />
					</div>
					<div class="col-span-1">
						<Label for="Suffix" class="block">Suffix</Label>
						<Input id="Suffix" value={r.Suffix} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="LifevowName" class="block">Life Vow Name</Label>
						<Input id="LifevowName" value={r.LifevowName} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="PreferedName" class="block">Preferred Name</Label>
						<Input id="PreferredName" value={r.PreferredName} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="DateReaffirmation" class="block">Last Reffirmation</Label>
						<Input id="DateReaffirmation" value={r.DateReaffirmation} onchange={changeDate} />
					</div>
					<div class="col-span-1">
						<Label for="MemberStatus" class="block">Member Status</Label>
						<Select
							id="MemberStatus"
							items={memberstatus}
							value={r.MemberStatus}
							onchange={change}
						/>
					</div>
					<div class="col-span-1">
						<Label for="Title" class="block">Title</Label>
						<Select id="Title" items={titles} value={r.Title} onchange={change} />
					</div>
				</div>
			</section>

			<section>
				<hr class="px-4 py-2" />
				<div class="grid grid-cols-8 gap-4 px-4 py-2">
					<div class="col-span-7">
						<Label for="Address" class="block">Address</Label>
						<Input id="Address" value={r.Address} onchange={change} />
					</div>
					<div class="col-span-1">
						<Toggle id="ListAddress" checked={r.ListAddress} onchange={changeCheck} color="red"
							>Listed</Toggle
						>
					</div>
					<div class="col-span-8">
						<Input id="AddressLine2" value={r.AddressLine2} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="City" class="block">City</Label>
						<Input id="City" value={r.City} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="State" class="block"
							>State/Locality <a href="https://en.wikipedia.org/wiki/ISO_3166-2">3166-2 code</a
							></Label
						>
						<Input id="State" value={r.State} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="Country" class="block"
							>Country <a href="https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2">2-letter code</a
							></Label
						>
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
					<div class="col-span-3">
						<Label for="PrimaryPhone" class="block">Primary Phone</Label>
						<Input id="PrimaryPhone" value={r.PrimaryPhone} onchange={change} />
					</div>
					<div class="col-span-1">
						<Toggle
							id="ListPrimaryPhone"
							checked={r.ListPrimaryPhone}
							onchange={changeCheck}
							color="red">Listed</Toggle
						>
					</div>
					<div class="col-span-3">
						<Label for="SecondaryPhone" class="block">Secondary Phone</Label>
						<Input id="SecondaryPhone" value={r.SecondaryPhone} onchange={change} />
					</div>
					<div class="col-span-1">
						<Toggle
							id="ListSecondaryPhone"
							checked={r.ListSecondaryPhone}
							onchange={changeCheck}
							color="red">Listed</Toggle
						>
					</div>
					<div class="col-span-3">
						<Label for="PrimaryEmail" class="block">Primary Email</Label>
						<Input id="PrimaryEmail" value={r.PrimaryEmail} onchange={change} />
					</div>
					<div class="col-span-1">
						<Toggle
							id="ListPrimaryEmail"
							checked={r.ListPrimaryEmail}
							onchange={changeCheck}
							color="red">Listed</Toggle
						>
					</div>
					<div class="col-span-3">
						<Label for="SecondaryEmail" class="block">Secondary Email</Label>
						<Input id="SecondaryEmail" value={r.SecondaryEmail} onchange={change} />
					</div>
					<div class="col-span-1">
						<Toggle
							id="ListSecondaryEmail"
							checked={r.ListSecondaryEmail}
							onchange={changeCheck}
							color="red">Listed</Toggle
						>
					</div>
				</div>
			</section>

			<section>
				<hr class="px-4 py-2" />
				<div class="grid grid-cols-8 gap-4 px-4 py-2">
					<div class="col-span-2">
						<Label for="DateFirstVows" class="block">First Vows</Label>
						<Input id="DateFirstVows" value={r.DateFirstVows} onchange={changeDate} />
					</div>
					<div class="col-span-2">
						<Label for="DateNovitiate" class="block">Novitiate</Label>
						<Input id="DateNovitiate" value={r.DateNovitiate} onchange={changeDate} />
					</div>
					<div class="col-span-2">
						<Label for="DateLifeVows" class="block">Life Vows</Label>
						<Input id="DateLifeVows" value={r.DateLifeVows} onchange={changeDate} />
					</div>
					<div class="col-span-2">
						<Label for="BirthDate" class="block">Birth Day</Label>
						<Input id="BirthDate" value={r.BirthDate} onchange={changeDate} />
					</div>
					<div class="col-span-2">
						<Label for="DateDeceased" class="block">Deceased</Label>
						<Input id="DateDeceased" value={r.DateDeceased} onchange={changeDate} />
					</div>
					<div class="col-span-2">
						<Label for="DateRemoved" class="block">Removed</Label>
						<Input id="DateRemoved" value={r.DateRemoved} onchange={changeDate} />
					</div>
				</div>
			</section>

			<section>
				<hr class="px-4 py-2" />
				<div class="grid grid-cols-8 gap-4 px-4 py-2">
					<div class="col-span-2">
						<Label for="HowJoined" class="block">How Joined</Label>
						<Input id="HowJoined" value={r.HowJoined} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="HowRemoved" class="block">How Removed</Label>
						<Select id="HowRemoved" items={removereasons} value={r.HowRemoved} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="Status" class="block">Status</Label>
						<Select id="Status" items={stati} value={r.Status} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="Leadership" class="block">Leadership</Label>
						<Select id="Leadership" items={leadership} value={r.Leadership} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="Chapters" class="block">Chapter</Label>
						<MultiSelect
							id="Chapters"
							items={chaps}
							bind:value={selectedchapters}
							onchange={setchapters}
						/>
					</div>
					<div class="col-span-2">
						<Label for="Occupation" class="block">Occupation</Label>
						<Input id="Occupation" value={r.Occupation} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="Employer" class="block">Employer</Label>
						<Input id="Employer" value={r.Employer} onchange={change} />
					</div>
					<div class="col-span-2">
						<Label for="Denomination" class="block">Denomination</Label>
						<Input id="Denomination" value={r.Denomination} onchange={change} />
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
							items={newsletteritems}
							value={r.Newsletter}
							onchange={change}
						/>
					</div>
					<div class="col-span-1">
						<Label for="Doxology" class="block">Doxology</Label>
						<Select id="Doxology" items={commitems} value={r.Doxology} onchange={change} />
					</div>
					<div class="col-span-1">
						<Label for="Communication" class="block">Communication</Label>
						<Select
							id="Communication"
							items={commitems}
							value={r.Communication}
							onchange={change}
						/>
					</div>
				</div>
			</section>

			<section>
				<Button onclick={() => vcard(r.ID)} color="purple">Add to Address Book</Button>
			</section>
		</form>
	{:else}
		<section>
			<div class="grid grid-cols-4 gap-1 px-6 py-6">
				<div class="col-span-1 justify-self-start">Name</div>
				<div class="col-span-3">{oslname(r)}</div>
				{#if r.ListAddress}
					<div class="col-span-1 justify-self-start">Address</div>
					<div class="col-span-3"><pre>{r.FormattedAddr}</pre></div>
				{/if}
				{#if r.ListPrimaryPhone && r.PrimaryPhone}<div class="col-span-1 justify-self-start">
						Primary Phone
					</div>
					<div class="col-span-3">{r.PrimaryPhone}</div>{/if}
				{#if r.ListSecondaryPhone && r.SecondaryPhone}<div class="col-span-1 justify-self-start">
						Secondary Phone
					</div>
					<div class="col-span-3">{r.SecondaryPhone}</div>{/if}
				{#if r.ListPrimaryEmail && r.PrimaryEmail}<div class="col-span-1 justify-self-start">
						Primary Email
					</div>
					<div class="col-span-3">{r.PrimaryEmail}</div>{/if}
				{#if r.ListSecondaryPhone && r.SecondaryEmail}<div class="col-span-1 justify-self-start">
						Secondary Email
					</div>
					<div class="col-span-3">{r.SecondaryEmail}</div>{/if}
				{#if r.Chapter}<div class="col-span-1 justify-self-start">Chapter</div>
					<div class="col-span-3">{r.Chapter}</div>{/if}
				<div class="col-span-1 justify-self-start">&nbsp;</div>
			</div>
		</section>
	{/if}
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}

<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { push } from 'svelte-spa-router';
	import {
		Button,
		Label,
		Input,
		Select,
		MultiSelect,
		Toggle,
		Card,
		Badge,
		Hr,
		Spinner,
		Heading,
		Helper
	} from 'flowbite-svelte';
	import {
		ChartMixedDollarOutline,
		EditOutline,
		ClockOutline,
		AddressBookOutline,
		InfoCircleOutline
	} from 'flowbite-svelte-icons';
	import { toast } from '@zerodevx/svelte-toast';
	import {
		getMember,
		updateMember,
		getMemberChapters,
		getChapters,
		updateMemberChapters,
		oslname,
		vcard,
		cleanDateFormat
	} from '../oo';

	let { params } = $props();
	const oo = getContext('oo');

	let loading = $state(true);
	let r = $state(null);
	let chaps = $state([]);
	let selectedChapters = $state([]);

	const commItems = [
		{ value: 'none', name: 'None' },
		{ value: 'mailed', name: 'Mailed' },
		{ value: 'electronic', name: 'Electronic' }
	];
	const newsletterItems = [
		{ value: 'none', name: 'None' },
		{ value: 'electronic', name: 'Electronic' }
	];
	const memberStatus = [
		{ value: 'Annual Vows', name: 'Annual Vows' },
		{ value: 'Life Vows', name: 'Life Vows' },
		{ value: 'Removed', name: 'Removed' },
		{ value: 'Friend', name: 'Friend' },
		{ value: 'Deceased', name: 'Deceased' }
	];
	const titles = [
		{ value: ' ', name: 'None' },
		{ value: 'Sr.', name: 'Sr.' },
		{ value: 'Br.', name: 'Br.' },
		{ value: 'Sibling', name: 'Sibling' }
	];
	const removeReasons = [
		{ value: ' ', name: 'N/A' },
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
		{ value: ' ', name: 'None' },
		{ value: 'member', name: 'Member' },
		{ value: 'prior', name: 'Prior' },
		{ value: 'council', name: 'General Council Officer' },
		{ value: 'canon', name: 'Canon' },
		{ value: 'elected', name: 'Elected Officer' }
	];

	onMount(async () => {
		if (!oo.me) return push('/Login');
		try {
			const [memberData, memberChaps] = await Promise.all([
				getMember(params.id),
				getMemberChapters(params.id)
			]);
			r = memberData;
			selectedChapters = memberChaps;
			if (oo.chaptercache) {
				chaps = oo.chaptercache;
			} else {
				chaps = await getChapters();
			}
		} catch (err: any) {
			toast.push(err.message);
		} finally {
			loading = false;
		}
	});

	async function handleUpdate(event: any) {
		const { id, value, type, checked } = event.target;
		let finalValue = type === 'checkbox' ? checked : value;

		if (id.startsWith('Date') || id === 'BirthDate') {
			finalValue = cleanDateFormat(finalValue);
		}

		try {
			await updateMember(params.id, id, finalValue);
			toast.push(`Updated ${id}`);
		} catch (err: any) {
			toast.push(`Error: ${err.message}`);
		}
	}

	async function updateChapters() {
		try {
			await updateMemberChapters(params.id, selectedChapters);
			toast.push('Chapters updated');
		} catch (err: any) {
			toast.push(err.message);
		}
	}
</script>

<svelte:head>
	<title>Member: {r ? r.LastName : 'Loading'}</title>
</svelte:head>

{#if loading}
	<div class="flex justify-center py-20"><Spinner color="purple" size="12" /></div>
{:else if r}
	<div class="mx-auto w-full space-y-8 px-4 py-8">
		{#if oo.me.level > 1}
			<header
				class="flex flex-col items-center justify-between gap-6 rounded-xl border border-slate-200 bg-white p-6 shadow-sm md:flex-row"
			>
				<div class="flex items-center gap-4">
					<div>
						<Heading tag="h2" class="text-3xl font-bold text-slate-900">{oslname(r)}</Heading>
						<div class="mt-1 flex items-center gap-2">
							<Badge color="purple" class="text-[10px] uppercase">{r.MemberStatus}</Badge>
							<span class="font-mono text-xs text-slate-400">ID: {r.ID}</span>
						</div>
					</div>
				</div>
				<div class="flex flex-wrap gap-2">
					<Button size="sm" color="alternative" href="#/giving/{r.ID}"
						><ChartMixedDollarOutline class="mr-2 h-4 w-4" />Giving</Button
					>
					<Button size="sm" color="alternative" href="#/notes/{r.ID}"
						><EditOutline class="mr-2 h-4 w-4" />Notes</Button
					>
					<Button size="sm" color="alternative" href="#/changelog/{r.ID}"
						><ClockOutline class="mr-2 h-4 w-4" />Logs</Button
					>
					<Button size="sm" color="purple" onclick={() => vcard(r.ID)}
						><AddressBookOutline class="mr-2 h-4 w-4" />vCard</Button
					>
				</div>
			</header>

			<div class="grid grid-cols-1 items-start gap-8 lg:grid-cols-3">
				<Card size="none" class="h-full border-slate-200 p-6 shadow-sm">
					<div class="mb-6 flex items-center justify-between border-b pb-2">
						<Heading tag="h4" class="text-lg font-bold uppercase text-slate-800"
							>Identity & Status</Heading
						>
						<Toggle
							id="ListInDirectory"
							checked={r.ListInDirectory}
							onchange={handleUpdate}
							color="red"
							size="small"
						/>
					</div>
					<div class="space-y-4">
						<div class="grid grid-cols-4 gap-2">
							<div class="col-span-1">
								<Label class="mb-1">Title</Label><Select
									id="Title"
									items={titles}
									value={r.Title}
									onchange={handleUpdate}
								/>
							</div>
							<div class="col-span-3">
								<Label class="mb-1">First Name</Label><Input
									id="FirstName"
									value={r.FirstName}
									onchange={handleUpdate}
								/>
							</div>
						</div>
						<div>
							<Label class="mb-1">Middle Name</Label><Input
								id="MiddleName"
								value={r.MiddleName}
								onchange={handleUpdate}
							/>
						</div>
						<div>
							<Label class="mb-1">Last Name</Label><Input
								id="LastName"
								value={r.LastName}
								onchange={handleUpdate}
							/>
						</div>
						<div>
							<Label class="mb-1">Suffix</Label><Input
								id="Suffix"
								value={r.Suffix}
								onchange={handleUpdate}
							/>
						</div>
						<div>
							<Label class="mb-1">Life Vow Name</Label><Input
								id="LifeVowName"
								value={r.LifeVowName}
								onchange={handleUpdate}
							/>
						</div>
						<div>
							<Label class="mb-1 font-bold">Preferred Name</Label><Input
								id="PreferredName"
								value={r.PreferredName}
								onchange={handleUpdate}
							/>
						</div>

						<Hr class="my-6" />

						<div class="space-y-4 rounded-lg border border-slate-200 bg-slate-50 p-4">
							<div>
								<Label class="mb-1">Member Status</Label><Select
									id="MemberStatus"
									items={memberStatus}
									value={r.MemberStatus}
									onchange={handleUpdate}
								/>
							</div>
							<div>
								<Label class="mb-1">How Removed</Label><Select
									id="HowRemoved"
									items={removeReasons}
									value={r.HowRemoved}
									onchange={handleUpdate}
								/>
							</div>
							<div>
								<Label class="mb-1">Ecclesial Status</Label><Select
									id="Status"
									items={stati}
									value={r.Status}
									onchange={handleUpdate}
								/>
							</div>
						</div>
					</div>
				</Card>

				<Card size="none" class="h-full border-slate-200 p-6 shadow-sm">
					<div class="mb-6 flex items-center justify-between border-b pb-2">
						<Heading tag="h4" class="text-lg font-bold uppercase text-slate-800">Contact</Heading>
						<Toggle
							id="ListAddress"
							checked={r.ListAddress}
							onchange={handleUpdate}
							color="red"
							size="small"
						/>
					</div>
					<div class="space-y-4">
						<div class="space-y-2">
							<Input
								id="Address"
								value={r.Address}
								onchange={handleUpdate}
								placeholder="Address Line 1"
							/>
							<Input
								id="AddressLine2"
								value={r.AddressLine2}
								onchange={handleUpdate}
								placeholder="Address Line 2"
							/>
							<div class="grid grid-cols-2 gap-2">
								<Input id="City" value={r.City} onchange={handleUpdate} placeholder="City" />
								<Input id="State" value={r.State} onchange={handleUpdate} placeholder="State" />
							</div>
							<div class="grid grid-cols-2 gap-2">
								<Input
									id="Country"
									value={r.Country}
									onchange={handleUpdate}
									placeholder="Country"
								/>
								<Input
									id="PostalCode"
									value={r.PostalCode}
									onchange={handleUpdate}
									placeholder="Postal Code"
								/>
							</div>
						</div>

						<Hr class="my-6" />

						<div class="space-y-6">
							<div class="flex items-center gap-2">
								<div class="flex-grow">
									<Label class="text-xs">Primary Email</Label><Input
										id="PrimaryEmail"
										value={r.PrimaryEmail}
										onchange={handleUpdate}
									/>
								</div>
								<Toggle
									id="ListPrimaryEmail"
									checked={r.ListPrimaryEmail}
									onchange={handleUpdate}
									color="red"
									class="mt-5"
								/>
							</div>
							<div class="flex items-center gap-2">
								<div class="flex-grow">
									<Label class="text-xs">Primary Phone</Label><Input
										id="PrimaryPhone"
										value={r.PrimaryPhone}
										onchange={handleUpdate}
									/>
								</div>
								<Toggle
									id="ListPrimaryPhone"
									checked={r.ListPrimaryPhone}
									onchange={handleUpdate}
									color="red"
									class="mt-5"
								/>
							</div>
							<div class="flex items-center gap-2 pt-2">
								<div class="flex-grow">
									<Label class="text-xs text-slate-400">Secondary Email</Label><Input
										id="SecondaryEmail"
										value={r.SecondaryEmail}
										onchange={handleUpdate}
									/>
								</div>
								<Toggle
									id="ListSecondaryEmail"
									checked={r.ListSecondaryEmail}
									onchange={handleUpdate}
									color="red"
									class="mt-5"
								/>
							</div>
							<div class="flex items-center gap-2">
								<div class="flex-grow">
									<Label class="text-xs text-slate-400">Secondary Phone</Label><Input
										id="SecondaryPhone"
										value={r.SecondaryPhone}
										onchange={handleUpdate}
									/>
								</div>
								<Toggle
									id="ListSecondaryPhone"
									checked={r.ListSecondaryPhone}
									onchange={handleUpdate}
									color="red"
									class="mt-5"
								/>
							</div>
						</div>
					</div>
				</Card>

				<div class="space-y-8">
					<Card size="none" class="border-slate-200 p-6 shadow-sm">
						<Heading tag="h4" class="mb-6 border-b pb-2 text-lg font-bold uppercase text-slate-800"
							>Data</Heading
						>
						<div class="space-y-4">
							<div class="grid grid-cols-2 gap-4">
								<div>
									<Label class="text-xs">Novitiate</Label><Input
										id="DateNovitiate"
										value={r.DateNovitiate}
										onchange={handleUpdate}
										placeholder="YYYY-MM-DD"
									/>
								</div>
								<div>
									<Label class="text-xs">First Vows</Label><Input
										id="DateFirstVows"
										value={r.DateFirstVows}
										onchange={handleUpdate}
										placeholder="YYYY-MM-DD"
									/>
								</div>
								<div>
									<Label class="text-xs">Life Vows</Label><Input
										id="DateLifeVows"
										value={r.DateLifeVows}
										onchange={handleUpdate}
										placeholder="YYYY-MM-DD"
									/>
								</div>
								<div>
									<Label class="text-xs font-bold">Reaffirmation</Label><Input
										id="DateReaffirmation"
										value={r.DateReaffirmation}
										onchange={handleUpdate}
										placeholder="YYYY-MM-DD"
									/>
								</div>
							</div>

							<div class="grid grid-cols-2 gap-4 pt-2">
								<div>
									<Label class="text-xs">Birth Date</Label><Input
										id="BirthDate"
										value={r.BirthDate}
										onchange={handleUpdate}
										placeholder="YYYY-MM-DD"
									/>
								</div>
								<div>
									<Label class="text-xs">Date Removed</Label><Input
										id="DateRemoved"
										value={r.DateRemoved}
										onchange={handleUpdate}
										placeholder="YYYY-MM-DD"
									/>
								</div>
							</div>

							<div class="rounded-lg border border-purple-100 bg-purple-50 p-3">
								<Label class="mb-1 text-[10px] font-bold uppercase text-purple-700">Deceased</Label>
								<Input
									id="DateDeceased"
									size="sm"
									value={r.DateDeceased}
									onchange={handleUpdate}
									placeholder="YYYY-MM-DD"
								/>
							</div>

							<div class="space-y-4 pt-4">
								<div>
									<Label class="mb-1">Chapters</Label><MultiSelect
										items={chaps}
										bind:value={selectedChapters}
										onchange={updateChapters}
									/>
								</div>
								<div>
									<Label class="mb-1">Leadership Role</Label><Select
										id="Leadership"
										items={leadership}
										value={r.Leadership}
										onchange={handleUpdate}
									/>
								</div>
								<div>
									<Label class="mb-1">Occupation</Label><Input
										id="Occupation"
										value={r.Occupation}
										onchange={handleUpdate}
									/>
								</div>
								<div>
									<Label class="mb-1">Employer</Label><Input
										id="Employer"
										value={r.Employer}
										onchange={handleUpdate}
									/>
								</div>
								<div>
									<Label class="mb-1">Denomination</Label><Input
										id="Denomination"
										value={r.Denomination}
										onchange={handleUpdate}
									/>
								</div>
							</div>
						</div>
					</Card>

					<Card size="none" class="border-slate-200 bg-slate-50/50 p-6 shadow-sm">
						<Heading tag="h4" class="mb-6 border-b pb-2 text-lg font-bold uppercase text-slate-800"
							>Communications</Heading
						>
						<div class="space-y-4">
							<div>
								<Label class="text-xs font-bold italic">Newsletter</Label><Select
									id="Newsletter"
									items={newsletterItems}
									value={r.Newsletter}
									onchange={handleUpdate}
								/>
							</div>
							<div>
								<Label class="text-xs font-bold italic">Doxology</Label><Select
									id="Doxology"
									items={commItems}
									value={r.Doxology}
									onchange={handleUpdate}
								/>
							</div>
							<div>
								<Label class="text-xs font-bold italic">Communication</Label><Select
									id="Communication"
									items={commItems}
									value={r.Communication}
									onchange={handleUpdate}
								/>
							</div>
						</div>
					</Card>
				</div>
			</div>
		{:else}
			<Card size="md" class="mx-auto mt-8 p-8">
				<div class="space-y-4 text-center">
					<Heading tag="h2" class="text-3xl font-bold">{oslname(r)}</Heading>
					<Badge color="purple" class="px-4 py-1 uppercase">{r.MemberStatus}</Badge>
				</div>
				<div class="mt-8 grid grid-cols-1 gap-8 border-t pt-8 md:grid-cols-2">
					{#if r.ListAddress}
						<div>
							<Label class="mb-2 text-xs font-bold uppercase text-slate-400">Mailing Address</Label>
							<p class="whitespace-pre-line font-medium leading-relaxed text-slate-700">
								{r.FormattedAddr}
							</p>
						</div>
					{/if}
					<div class="space-y-6">
						{#if r.ListPrimaryEmail && r.PrimaryEmail}
							<div>
								<Label class="mb-1 text-xs font-bold uppercase text-slate-400">Email</Label>
								<a
									href="mailto:{r.PrimaryEmail}"
									class="font-semibold text-primary-600 hover:underline">{r.PrimaryEmail}</a
								>
							</div>
						{/if}
						{#if r.ListPrimaryPhone && r.PrimaryPhone}
							<div>
								<Label class="mb-1 text-xs font-bold uppercase text-slate-400">Phone</Label>
								<p class="font-semibold text-slate-800">{r.PrimaryPhone}</p>
							</div>
						{/if}
					</div>
				</div>
			</Card>
		{/if}
	</div>
{/if}

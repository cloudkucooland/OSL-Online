<script lang="ts">
	import { getContext, onMount } from 'svelte';
	import { push } from 'svelte-spa-router';
	import {
		Label,
		Input,
		Toggle,
		Select,
		MultiSelect,
		Card,
		Badge,
		Helper,
		Heading,
		Hr,
		Spinner
	} from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import {
		getMeFromServer,
		updateMe,
		getChapters,
		updateMeChapters,
		getMeChapters,
		oslname
	} from '../oo';

	const oo = getContext('oo');

	// State Runes
	let loading = $state(true);
	let member = $state(null);
	let chaps = $state([]);
	let selectedChapters = $state([]);

	// Constants
	const cannotEdit = true;

	const newsletterItems = [
		{ value: 'none', name: 'None' },
		{ value: 'electronic', name: 'Electronic' }
	];

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
		{ value: ' ', name: 'None' },
		{ value: 'Sr.', name: 'Sr.' },
		{ value: 'Br.', name: 'Br.' },
		{ value: 'Sibling', name: 'Sibling' }
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

	onMount(async () => {
		if (!oo.me) {
			push('/Login');
			return;
		}
		try {
			const [m, myChaps] = await Promise.all([getMeFromServer(), getMeChapters()]);
			member = m;
			selectedChapters = myChaps;
			if (oo.chaptercache) {
				chaps = oo.chaptercache;
			} else {
				chaps = await getChapters();
			}
		} catch (err) {
			toast.push(err.message);
		} finally {
			loading = false;
		}
	});

	async function handleUpdate(id, value, isCheck = false) {
		try {
			await updateMe(id, value);
			toast.push(`Saved changes to ${id}`);
		} catch (err) {
			toast.push(`Error: ${err.message}`);
		}
	}

	async function updateChapters() {
		try {
			await updateMeChapters(selectedChapters);
			toast.push('Updated Chapter affiliations');
		} catch (err) {
			toast.push(err.message);
		}
	}
</script>

<svelte:head>
	<title>OSL Member Record</title>
</svelte:head>

<div class="mx-auto w-full space-y-8 px-4 py-8">
	{#if loading}
		<div class="flex justify-center p-20">
			<Spinner size="12" color="purple" />
		</div>
	{:else if member}
		<header
			class="flex flex-col items-center justify-between gap-6 rounded-xl border border-slate-200 bg-white p-8 shadow-sm md:flex-row"
		>
			<div>
				<div class="mb-2 flex items-center gap-3">
					<Heading tag="h2" class="text-3xl font-bold text-slate-900">{oslname(member)}</Heading>
					<Badge color="purple" class="px-3 text-sm">{member.MemberStatus}</Badge>
				</div>
			</div>

			<div class="flex flex-col items-center rounded-lg border border-red-100 bg-red-50 p-4">
				<Toggle
					id="ListInDirectory"
					checked={member.ListInDirectory}
					onchange={(e) => handleUpdate('ListInDirectory', e.target.checked)}
					color="red"
				>
					<span class="font-bold text-red-700">List in Directory</span>
				</Toggle>
			</div>
		</header>

		<div class="grid grid-cols-1 items-start gap-8 lg:grid-cols-3">
			<Card size="none" class="h-full border-slate-200 p-6 shadow-sm">
				<Heading tag="h4" class="mb-6 border-b pb-2 text-lg font-bold uppercase text-slate-800"
					>Personal Information</Heading
				>
				<div class="space-y-4">
					<div>
						<Label class="mb-1">First Name</Label><Input
							value={member.FirstName}
							disabled={cannotEdit}
						/>
					</div>
					<div>
						<Label class="mb-1">Middle Name</Label><Input
							value={member.MiddleName}
							disabled={cannotEdit}
						/>
					</div>
					<div>
						<Label class="mb-1">Last Name</Label><Input
							value={member.LastName}
							disabled={cannotEdit}
						/>
					</div>
					<div>
						<Label class="mb-1">Suffix</Label><Input value={member.Suffix} disabled={cannotEdit} />
					</div>
					<div>
						<Label class="mb-1">Life Vow Name</Label><Input
							value={member.LifeVowName}
							disabled={cannotEdit}
						/>
					</div>

					<div class="mt-4 rounded-lg border border-dashed border-red-200 bg-red-50/40 p-3">
						<Label class="mb-1 font-bold italic text-red-700">Preferred Name</Label>
						<Input
							id="PreferredName"
							value={member.PreferredName}
							onchange={(e) => handleUpdate('PreferredName', e.target.value)}
						/>
					</div>

					<div>
						<Label class="mb-1 font-bold text-red-700">Title</Label><Select
							id="Title"
							items={titles}
							value={member.Title}
							onchange={(e) => handleUpdate('Title', e.target.value)}
							disabled={member.MemberStatus == 'Friend'}
						/>
					</div>
					<div>
						<Label class="mb-1">Birth Date</Label><Input
							value={member.BirthDate}
							disabled={cannotEdit}
						/>
					</div>
				</div>
			</Card>

			<Card size="none" class="h-full border-slate-200 p-6 shadow-sm">
				<Heading tag="h4" class="mb-6 border-b pb-2 text-lg font-bold uppercase text-slate-800"
					>Contact & Mailing</Heading
				>
				<div class="space-y-4">
					<div class="mb-4 rounded-lg border border-slate-200 bg-slate-50 p-4">
						<div class="mb-4 flex items-center justify-between">
							<span class="text-xs font-bold uppercase text-slate-400">Mailing Address</span>
							<Toggle
								size="small"
								checked={member.ListAddress}
								onchange={(e) => handleUpdate('ListAddress', e.target.checked)}
								color="red"
							/>
						</div>
						<div class="space-y-2">
							<Input value={member.Address} disabled={cannotEdit} placeholder="Line 1" />
							<Input value={member.AddressLine2} disabled={cannotEdit} placeholder="Line 2" />
							<div class="grid grid-cols-2 gap-2">
								<Input value={member.City} disabled={cannotEdit} placeholder="City" />
								<Input value={member.State} disabled={cannotEdit} placeholder="State" />
							</div>
							<div class="grid grid-cols-2 gap-2">
								<Input value={member.Country} disabled={cannotEdit} placeholder="Country" />
								<Input value={member.PostalCode} disabled={cannotEdit} placeholder="Postal Code" />
							</div>
						</div>
					</div>

					<div class="space-y-4">
						<div class="flex items-center gap-2">
							<div class="flex-grow">
								<Label class="text-xs italic">Primary Email</Label><Input
									value={member.PrimaryEmail}
									disabled={cannotEdit}
								/>
							</div>
							<Toggle
								class="mt-5"
								checked={member.ListPrimaryEmail}
								onchange={(e) => handleUpdate('ListPrimaryEmail', e.target.checked)}
								color="red"
							/>
						</div>
						<div class="flex items-center gap-2 border-b pb-4">
							<div class="flex-grow">
								<Label class="text-xs italic">Primary Phone</Label><Input
									value={member.PrimaryPhone}
									disabled={cannotEdit}
								/>
							</div>
							<Toggle
								class="mt-5"
								checked={member.ListPrimaryPhone}
								onchange={(e) => handleUpdate('ListPrimaryPhone', e.target.checked)}
								color="red"
							/>
						</div>
						<div class="flex items-center gap-2">
							<div class="flex-grow">
								<Label class="text-xs italic text-slate-400">Secondary Email</Label><Input
									value={member.SecondaryEmail}
									disabled={cannotEdit}
								/>
							</div>
							<Toggle
								class="mt-5"
								checked={member.ListSecondaryEmail}
								onchange={(e) => handleUpdate('ListSecondaryEmail', e.target.checked)}
								color="red"
							/>
						</div>
						<div class="flex items-center gap-2">
							<div class="flex-grow">
								<Label class="text-xs italic text-slate-400">Secondary Phone</Label><Input
									value={member.SecondaryPhone}
									disabled={cannotEdit}
								/>
							</div>
							<Toggle
								class="mt-5"
								checked={member.ListSecondaryPhone}
								onchange={(e) => handleUpdate('ListSecondaryPhone', e.target.checked)}
								color="red"
							/>
						</div>
					</div>
				</div>
			</Card>

			<div class="space-y-8">
				<Card size="none" class="border-slate-200 p-6 shadow-sm">
					<Heading tag="h4" class="mb-6 border-b pb-2 text-lg font-bold uppercase text-slate-800"
						>Vocation & Status</Heading
					>
					<div class="space-y-4">
						<div class="grid grid-cols-2 gap-4">
							<div>
								<Label class="text-xs">Novitiate</Label><Input
									value={member.DateNovitiate}
									disabled={cannotEdit}
								/>
							</div>
							<div>
								<Label class="text-xs">First Vows</Label><Input
									value={member.DateFirstVows}
									disabled={cannotEdit}
								/>
							</div>
							<div>
								<Label class="text-xs">Life Vows</Label><Input
									value={member.DateLifeVows}
									disabled={cannotEdit}
								/>
							</div>
							<div>
								<Label class="text-xs font-bold">Reaffirmation</Label><Input
									value={member.DateReaffirmation}
									disabled={true}
								/>
							</div>
						</div>

						<div class="mt-2 rounded-lg border border-dashed border-red-200 bg-red-50/40 p-4">
							<Label class="mb-2 font-bold italic text-red-700">Active Chapters</Label>
							<MultiSelect items={chaps} bind:value={selectedChapters} onchange={updateChapters} />
						</div>

						<div class="space-y-4 pt-2">
							<div>
								<Label class="font-bold italic text-red-700">Status</Label><Select
									id="Status"
									items={stati}
									value={member.Status}
									onchange={(e) => handleUpdate('Status', e.target.value)}
								/>
							</div>
							<div>
								<Label class="text-xs">Occupation</Label><Input
									value={member.Occupation}
									onchange={(e) => handleUpdate('Occupation', e.target.value)}
								/>
							</div>
							<div>
								<Label class="text-xs">Employer</Label><Input
									value={member.Employer}
									onchange={(e) => handleUpdate('Employer', e.target.value)}
								/>
							</div>
							<div>
								<Label class="text-xs">Denomination</Label><Input
									value={member.Denomination}
									onchange={(e) => handleUpdate('Denomination', e.target.value)}
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
							<Label class="font-bold italic text-red-700">Newsletter</Label>
							<Select
								items={newsletterItems}
								value={member.Newsletter}
								onchange={(e) => handleUpdate('Newsletter', e.target.value)}
							/>
						</div>
						<div>
							<Label class="font-bold italic text-red-700">Doxology</Label>
							<Select
								items={commitems}
								value={member.Doxology}
								onchange={(e) => handleUpdate('Doxology', e.target.value)}
							/>
						</div>
						<div>
							<Label class="font-bold italic text-red-700">General Communication</Label>
							<Select
								items={commitems}
								value={member.Communication}
								onchange={(e) => handleUpdate('Communication', e.target.value)}
							/>
						</div>
					</div>
				</Card>
			</div>
		</div>
	{/if}
</div>

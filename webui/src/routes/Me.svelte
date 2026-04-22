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
		Heading,
		Hr,
		Spinner,
		Button
	} from 'flowbite-svelte';
	import { ChartMixedDollarOutline, AddressBookOutline } from 'flowbite-svelte-icons';
	import { toast } from '@zerodevx/svelte-toast';
	import {
		getMeFromServer,
		updateMe,
		getChapters,
		updateMeChapters,
		getMeChapters,
		oslname,
		vcard
	} from '../oo';

	// Shared Components
	import AddressSection from '$lib/AddressSection.svelte';
	import ContactSection from '$lib/ContactSection.svelte';
	import FulfillmentSection from '$lib/FulfillmentSection.svelte';

	const oo = getContext('oo');

	// State Runes
	let loading = $state(true);
	let member = $state(null);
	let chaps = $state([]);
	let selectedChapters = $state([]);

	// Constants
	const cannotEdit = true; // Most core fields are managed by office

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
				oo.chaptercache = chaps;
			}
		} catch (err: any) {
			toast.push(err.message);
		} finally {
			loading = false;
		}
	});

	async function handleUpdate(id: string, value: any) {
		try {
			await updateMe(id, value);
			toast.push(`Saved changes to ${id}`);
		} catch (err: any) {
			toast.push(`Error: ${err.message}`);
		}
	}

	async function updateChapters() {
		try {
			await updateMeChapters(selectedChapters);
			toast.push('Updated Chapter affiliations');
		} catch (err: any) {
			toast.push(err.message);
		}
	}
</script>

<svelte:head>
	<title>My Profile | OSL Directory</title>
</svelte:head>

<div class="mx-auto w-full space-y-8 px-4 py-8">
	{#if loading}
		<div class="flex justify-center py-20">
			<Spinner size="12" color="purple" />
		</div>
	{:else if member}
		<header
			class="flex flex-col items-center justify-between gap-6 rounded-xl border border-slate-200 bg-white p-6 shadow-sm md:flex-row"
		>
			<div class="flex items-center gap-4">
				<div>
					<Heading tag="h2" class="text-3xl font-bold text-slate-900">{oslname(member)}</Heading>
					<div class="mt-1 flex items-center gap-2">
						<Badge color="purple" class="text-[10px] uppercase">{member.MemberStatus}</Badge>
						<span class="font-mono text-xs text-slate-400">ID: {member.ID}</span>
					</div>
				</div>
			</div>
			
			<div class="flex flex-wrap items-center gap-4">
				<div class="flex flex-col items-center rounded-lg border border-red-100 bg-red-50 px-4 py-2">
					<Toggle
						id="ListInDirectory"
						checked={member.ListInDirectory}
						onchange={(e) => handleUpdate('ListInDirectory', e.target.checked)}
						color="red"
						size="small"
					>
						<span class="text-xs font-bold text-red-700">List in Directory</span>
					</Toggle>
				</div>
				<Button size="sm" color="alternative" href="#/mygiving">
					<ChartMixedDollarOutline class="mr-2 h-4 w-4" />My Giving
				</Button>
				<Button size="sm" color="purple" onclick={() => vcard(member.ID)}>
					<AddressBookOutline class="mr-2 h-4 w-4" />vCard
				</Button>
			</div>
		</header>

		<div class="grid grid-cols-1 items-start gap-8 lg:grid-cols-3">
			<!-- Column 1: Identity -->
			<Card size="none" class="h-full border-slate-200 p-6 shadow-sm">
				<Heading tag="h4" class="mb-6 border-b pb-2 text-lg font-bold uppercase text-slate-800">My Identity</Heading>
				<div class="space-y-4">
					<div>
						<Label class="mb-1 text-xs text-slate-400 italic">Legal Name (Office Use Only)</Label>
						<Input value={`${member.FirstName} ${member.MiddleName || ''} ${member.LastName}`} disabled={cannotEdit} />
					</div>
					
					<div class="mt-4 rounded-lg border border-dashed border-red-200 bg-red-50/40 p-4">
						<Label class="mb-1 font-bold text-red-700">Preferred Name</Label>
						<Input
							id="PreferredName"
							value={member.PreferredName}
							onchange={(e) => handleUpdate('PreferredName', e.target.value)}
							placeholder="Name as you'd like it to appear"
						/>
						<p class="mt-1 text-[10px] text-red-600 italic">This is how you will appear in the directory.</p>
					</div>

					<div>
						<Label class="mb-1 font-bold text-red-700">Title</Label>
						<Select
							id="Title"
							items={titles}
							value={member.Title}
							onchange={(e) => handleUpdate('Title', e.target.value)}
							disabled={member.MemberStatus === 'Friend'}
						/>
					</div>
					
					<Hr class="my-6" />
					
					<div>
						<Label class="mb-1 text-xs text-slate-400 italic">Birth Date</Label>
						<Input id="BirthDate" value={member.BirthDate} onchange={(e) => handleUpdate('BirthDate', e.target.value)} placeholder="YYYY-MM-DD" />
					</div>
				</div>
			</Card>

			<!-- Column 2: Contact & Mailing -->
			<Card size="none" class="h-full border-slate-200 p-6 shadow-sm">
				<AddressSection
					data={member}
					onUpdate={handleUpdate}
					showPrivacy={true}
					disabled={true}
					title="Mailing Address"
				/>
				
				<Hr class="my-6" />
				
				<ContactSection
					data={member}
					onUpdate={handleUpdate}
					showPrivacy={true}
					disabled={true}
					title="Contact Information"
				/>
				<p class="mt-4 text-[10px] italic text-slate-400">Contact the office to update your address or core email/phone.</p>
			</Card>

			<!-- Column 3: Vocation & Comms -->
			<div class="space-y-8">
				<Card size="none" class="border-slate-200 p-6 shadow-sm">
					<Heading tag="h4" class="mb-6 border-b pb-2 text-lg font-bold uppercase text-slate-800">Vocation</Heading>
					<div class="space-y-4">
						<div class="grid grid-cols-2 gap-4">
							<div>
								<Label class="text-[10px] uppercase text-slate-400">First Vows</Label>
								<p class="font-semibold text-slate-700">{member.DateFirstVows || '—'}</p>
							</div>
							<div>
								<Label class="text-[10px] uppercase text-slate-400 font-bold text-primary-700">Reaffirmation</Label>
								<p class="font-bold text-primary-800">{member.DateReaffirmation || '—'}</p>
							</div>
						</div>

						<div class="mt-2 rounded-lg border border-dashed border-red-200 bg-red-50/40 p-4">
							<Label class="mb-2 font-bold italic text-red-700">My Chapters</Label>
							<MultiSelect items={chaps} bind:value={selectedChapters} onchange={updateChapters} />
						</div>

						<div class="space-y-4 pt-2">
							<div>
								<Label class="font-bold text-red-700">Ecclesial Status</Label>
								<Select
									id="Status"
									items={stati}
									value={member.Status}
									onchange={(e) => handleUpdate('Status', e.target.value)}
								/>
							</div>
							<div>
								<Label class="text-xs">Occupation</Label>
								<Input id="Occupation" value={member.Occupation} onchange={(e) => handleUpdate('Occupation', e.target.value)} />
							</div>
							<div>
								<Label class="text-xs">Employer</Label>
								<Input id="Employer" value={member.Employer} onchange={(e) => handleUpdate('Employer', e.target.value)} />
							</div>
							<div>
								<Label class="text-xs">Denomination</Label>
								<Input id="Denomination" value={member.Denomination} onchange={(e) => handleUpdate('Denomination', e.target.value)} />
							</div>
						</div>
					</div>
				</Card>

				<FulfillmentSection data={member} onUpdate={handleUpdate} />
			</div>
		</div>
	{/if}
</div>
<script lang="ts">
	import { onMount, getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import {
		Label,
		Textarea,
		Button,
		Card,
		Heading,
		Badge,
		Spinner,
		Timeline,
		TimelineItem
	} from 'flowbite-svelte';
	import {
		MessageCaptionOutline,
		TrashBinOutline,
		PlusOutline,
		ArrowLeftOutline
	} from 'flowbite-svelte-icons';
	import {
		getMember,
		getMemberNotes,
		postMemberNote,
		deleteMemberNote,
		cleanDateFormat,
		oslname
	} from '../oo';
	import { toast } from '@zerodevx/svelte-toast';

	let { params = {} } = $props();
	const oo = getContext('oo');

	// Guard: Redirect if not logged in
	if (!oo.me) {
		push('/Login');
	}

	let loading = $state(true);
	let member = $state(null);
	let notes = $state([]);
	let newNoteText = $state('');

	async function loadData() {
		try {
			const [m, n] = await Promise.all([getMember(params.id), getMemberNotes(params.id)]);
			member = m;
			// Sort notes newest first
			notes = n.sort((a, b) => new Date(b.Date).getTime() - new Date(a.Date).getTime());
		} catch (err: any) {
			toast.push(err.message);
		} finally {
			loading = false;
		}
	}

	onMount(loadData);

	async function handleAdd(event: Event) {
		event.preventDefault();
		if (!newNoteText.trim()) return;

		try {
			await postMemberNote(params.id, newNoteText);
			newNoteText = '';
			toast.push('Note added');
			await loadData(); // Refresh list
		} catch (err: any) {
			toast.push(`Failed to add: ${err.message}`);
		}
	}

	async function handleDelete(noteId: string) {
		if (!confirm('Are you sure you want to delete this note?')) return;

		try {
			await deleteMemberNote(params.id, noteId);
			toast.push('Note deleted');
			await loadData();
		} catch (err: any) {
			toast.push(`Delete failed: ${err.message}`);
		}
	}
</script>

<svelte:head>
	<title>Notes: {member ? member.LastName : 'Loading'}</title>
</svelte:head>

<div class="mx-auto max-w-4xl space-y-6 px-4 py-6">
	{#if loading}
		<div class="flex justify-center py-20"><Spinner color="purple" size="12" /></div>
	{:else if member}
		<div class="flex items-center justify-between border-b pb-4">
			<div>
				<Button
					size="xs"
					color="alternative"
					onclick={() => push(`/member/${member.ID}`)}
					class="mb-2"
				>
					<ArrowLeftOutline class="mr-1 h-3 w-3" /> Back to Record
				</Button>
				<Heading tag="h2" class="text-2xl font-bold text-slate-900">
					Member Notes: {oslname(member)}
				</Heading>
				<p class="mt-1 text-sm italic text-slate-500">{member.MemberStatus} — ID: {member.ID}</p>
			</div>
			<MessageCaptionOutline class="h-10 w-10 text-slate-200" />
		</div>

		<Card size="none" class="border-2 border-dashed border-slate-200 bg-slate-50 p-4">
			<form onsubmit={handleAdd} class="space-y-3">
				<Label for="new-note" class="font-semibold text-slate-700">Add Internal Note</Label>
				<Textarea
					id="new-note"
					placeholder="Enter details about reaffirmation, special requests, or life events..."
					rows="3"
					bind:value={newNoteText}
				/>
				<div class="flex justify-end">
					<Button type="submit" color="green" disabled={!newNoteText.trim()}>
						<PlusOutline class="mr-2 h-4 w-4" /> Save Note
					</Button>
				</div>
			</form>
		</Card>

		<div class="space-y-4">
			{#each notes as n}
				<Card
					size="none"
					class="border-slate-200 p-4 shadow-sm transition-colors hover:border-primary-300"
				>
					<div class="mb-2 flex items-start justify-between">
						<Badge color="indigo" class="font-mono">{cleanDateFormat(n.Date)}</Badge>
						<Button
							size="xs"
							color="none"
							class="p-1 text-slate-400 hover:text-red-600"
							onclick={() => handleDelete(n.ID)}
						>
							<TrashBinOutline class="h-4 w-4" />
						</Button>
					</div>
					<p class="whitespace-pre-wrap leading-relaxed text-slate-700">
						{n.Note}
					</p>
				</Card>
			{:else}
				<div class="text-center py-12 bg-slate-50 rounded-lg border-2 border-dotted">
					<p class="text-slate-400 italic">No notes recorded for this member yet.</p>
				</div>
			{/each}
		</div>
	{/if}
</div>

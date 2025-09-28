<script lang="ts">
	import { getContext } from 'svelte';
	import {
		getMember,
		getMemberNotes,
		postMemberNote,
		deleteMemberNote,
		cleanDateFormat
	} from '../oo';
	import { Label, Input, Button } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push, replace } from 'svelte-spa-router';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	let { params } = $props();
	let note = $state();

	async function add(event) {
		event.preventDefault();
		event.stopPropagation();

		try {
			await postMemberNote(params.id, note);
			toast.push(`Added`);
			await replace(`#/notes/${params.id}`);
		} catch (err) {
			console.log(err);
			toast.push('failed to change: ' + err.message);
		}
		return true;
	}

	async function getBoth(id) {
		try {
			const m = await getMember(id);
			const n = await getMemberNotes(id);
			return { m, n };
		} catch (err) {
			toast.push('failed to load: ' + err.message);
			console.log(err);
			throw err;
		}
	}

	async function deleteNote(noteid) {
		try {
			await deleteMemberNote(params.id, noteid);
			toast.push(`deleted`);
			await replace(`#/notes/${params.id}`);
		} catch (err) {
			console.log(err);
			toast.push('failed to change: ' + err.message);
		}
		return true;
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Member Notes</title>
</svelte:head>

{#await getBoth(params.id)}
	<h3>... loading ...</h3>
{:then { m, n }}
	<div>
		Notes for : <a href="#/member/{m.ID}">
			{m.Title}
			{m.FirstName}
			{m.LastName}</a
		>
		( {m.MemberStatus} )
	</div>

	<section>
		<div class="grid grid-cols-5 gap-4 px-4 py-2">
			<div class="col-span-1">Date</div>
			<div class="col-span-1">Note</div>
			<div class="col-span-1">...</div>
		</div>
		{#each n as note}
			<div class="grid grid-cols-5 gap-4 px-4 py-2">
				<div class="col-span-1">{cleanDateFormat(note.Date)}</div>
				<div class="col-span-1">{note.Note}</div>
				<div class="col-span-1">
					<Button color="red" onclick={() => deleteNote(note.ID)}>Delete</Button>
				</div>
			</div>
		{/each}
	</section>

	<section>
		<form onsubmit={add}>
			<div class="grid grid-cols-5 gap-4 px-4 py-2">
				<div class="col-span-1">
					<Label for="Note" class="block">Note</Label>
					<Input id="Note" bind:value={note} />
				</div>
				<div class="col-span-4">&nbsp;</div>
				<div class="col-span-1">
					<Button color="green" type="submit">Add Note</Button>
				</div>
			</div>
		</form>
	</section>
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}

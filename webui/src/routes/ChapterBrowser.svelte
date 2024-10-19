<script>
	import { getContext } from 'svelte';
	import { Label, Select } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { oslname, getChapters, getChapterMembers } from '../oo';

	const { me } = getContext('oo');

	let chaps = [];
	let members = [];
	let selected;

	async function chooseChapter(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			members = await getChapterMembers(selected);
			// push('/');
			// window.location.href = '#/';
		} catch (e) {
			console.log(e);
			toast.push(e.message);
		}
	}

	async function load() {
		chaps = await getChapters();
		return chaps;
	}
</script>

{#await load()}
	<h3>... loading ...</h3>
{:then}
	<div class="grid grid-cols-8 gap-4 px-4 py-2">
		<div class="col-span-8">
			<Select class="mt-2" items={chaps} bind:value={selected} on:change={chooseChapter} />
		</div>
		{#each members as m}
			<div class="col-span-8">
				{#if m.Leadership == 'prior'}Prior{/if}
				<a href="#/member/{m.ID}">{oslname(m)}</a><br />
				{#if m.Address}
					{m.Address}<br />{m.City} {m.State} {m.PostalCode} {m.Country} <br />
				{/if}
				{#if m.PrimaryEmail}<a href="mailto:{m.PrimaryEmail}">{m.PrimaryEmail}</a><br />{/if}
				{#if m.PrimaryPhone}<a href="tel:{m.PrimaryPhone}">{m.PrimaryPhone}</a>{/if}
			</div>
		{/each}
	</div>
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}

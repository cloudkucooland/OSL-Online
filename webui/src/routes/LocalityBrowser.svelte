<script lang="ts">
	import { Select } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { oslname, getLocalities, getLocalityMembers } from '../oo';

	let { params = {} } = $props();
	let locs = $state([]);
	let members = $state([]);
	let selected = $state();

	if (params.loc) {
		selected = params.loc;
		const e = new Event('search', { bubbles: true, cancelable: true });
		chooseLocality(e);
	}

	async function chooseLocality(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			members = await getLocalityMembers(selected);
			push(`#/localitybrowser/${selected}`);
		} catch (e) {
			console.log(e);
			toast.push(e.message);
		}
	}

	async function load() {
		locs = await getLocalities();
		return locs;
	}
</script>

{#await load()}
	<h3>... loading ...</h3>
{:then}
	<div class="grid grid-cols-8 gap-4 px-4 py-2">
		<div class="col-span-8">
			<b>"Can we search by state"</b> is the most common request. We are an International order. It's
			not as easy as you think.
		</div>
		<div class="col-span-8">
			<Select class="mt-2" items={locs} bind:value={selected} onchange={chooseLocality} />
		</div>
		{#each members as m}
			<div class="col-span-8">
				<a href="#/member/{m.ID}">{oslname(m)}</a><br />
				{#if m.PrimaryEmail}<a href="mailto:{m.PrimaryEmail}">{m.PrimaryEmail}</a><br />{/if}
				{#if m.PrimaryPhone}<a href="tel:{m.PrimaryPhone}">{m.PrimaryPhone}</a>{/if}
			</div>
		{/each}
	</div>
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}

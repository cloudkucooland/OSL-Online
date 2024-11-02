<script>
	import { getContext } from 'svelte';
	import { Label, Select } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { oslname, getLocalities, getLocalityMembers } from '../oo';

	const { me } = getContext('oo');

	let locs = [];
	let members = [];
	let selected;

	async function chooseLocality(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			members = await getLocalityMembers(selected);
			// push('/');
			// window.location.href = '#/';
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
		<ul>
		  <li>"Can we search by state" is the most common request... no. We are an International order. It is not as easy as you think.</li>
		  <li>An empty locality is because people have set themselves to not be listed</li>
		  <li>Singapore isn't working correctly (Country/State/City are the same)</li>
		  <li>The international code for the UK is GB, but not working properly yet</li>
		</ul>
		</div>
		<div class="col-span-8">
			<Select class="mt-2" items={locs} bind:value={selected} on:change={chooseLocality} />
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

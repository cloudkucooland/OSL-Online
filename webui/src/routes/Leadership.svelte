<script>
	import { getContext } from 'svelte';
	import { Label, Select } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';
	import { oslname, getLeaders } from '../oo';

	const { me } = getContext('oo');

	export let params = {};
	let categories = [
		{ name: 'Priors', value: 'prior' },
		{ name: 'Abbatial Appointment', value: 'council' },
		{ name: 'Canon', value: 'canon' },
		{ name: 'General Officer', value: 'elected' }
	];
	let leaders = [];
	let selected;

	if (params.id) {
		selected = params.id;
		const e = new Event('search', { bubbles: true, cancelable: true });
		chooseType(e);
	}

	async function chooseType(event) {
		event.preventDefault();
		event.stopPropagation();
		try {
			leaders = await getLeaders(selected);
			push(`#/leadership/${selected}`);
		} catch (e) {
			console.log(e);
			toast.push(e.message);
		}
	}
</script>

<div class="grid grid-cols-8 gap-4 px-4 py-2">
	<div class="col-span-8">
		Leadership Browser
	</div>
	<div class="col-span-8">
		<Select class="mt-2" items={categories} bind:value={selected} on:change={chooseType} />
	</div>
	{#each leaders as m}
		<div class="col-span-8">
			<a href="#/member/{m.ID}">{oslname(m)}</a><br />
			{#if m.PrimaryEmail}<a href="mailto:{m.PrimaryEmail}">{m.PrimaryEmail}</a><br />{/if}
			{#if m.PrimaryPhone}<a href="tel:{m.PrimaryPhone}">{m.PrimaryPhone}</a>{/if}
		</div>
	{/each}
</div>

<script lang="ts">
	import { getContext } from 'svelte';
	import { getMe, getMember, getChangelog } from '../oo';
	import { Label, Input, Button, Select } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { push } from 'svelte-spa-router';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	export let params;
	const dd = new Date().toISOString().split('T');
	let postdate = dd[0];
	let amount;
	let description = 'Annual Reaffirmation';
	let check = 0;
	let transaction = 0;

	async function getBoth(id) {
		try {
			const m = await getMember(id);
			const c = await getChangelog(id);
			return { m, c };
		} catch (err) {
			toast.push('failed to load: ' + err);
			console.log(err);
			throw err;
		}
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Change Log</title>
</svelte:head>

{#await getBoth(params.id)}
	<h3>... loading ...</h3>
{:then { m, c }}
	<div>
		Giving History for : <a href="#/member/{m.ID}">
			{m.Title}
			{m.FirstName}
			{m.LastName}</a
		>
		( {m.MemberStatus} )
	</div>

	<section>
		<div class="grid grid-cols-5 gap-4 px-4 py-2">
			<div class="col-span-1">Date</div>
			<div class="col-span-1">Changer</div>
			<div class="col-span-1">Field</div>
			<div class="col-span-1">Value</div>
		</div>
		{#each c as cr}
			<div class="grid grid-cols-5 gap-4 px-4 py-2">
				<div class="col-span-1">{cr.Date}</div>
				<div class="col-span-1">{cr.Changer}</div>
				<div class="col-span-1">{cr.Field}</div>
				<div class="col-span-1">{cr.Value}</div>
			</div>
		{/each}
	</section>
{:catch error}
	<h3 style="color: red">{error.message}</h3>
{/await}

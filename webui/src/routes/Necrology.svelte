<script lang="ts">
	import { getContext } from 'svelte';
	import { Table } from '@flowbite-svelte-plugins/datatable';
	import { getNecrology, cleanDateFormat } from '../oo';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	$: items = [
		{
			Title: ' ',
			'First Name': 'Loading',
			'Last Name': 'Loading',
			'Date Deceased': 'Loading',
			Locality: 'Loading'
		}
	];

	const opts = {
		paging: true,
		perPage: 100,
		perPageSelect: [25, 50, 100],
		sortable: true
	};

	// should not need the #await block, but we do... onMount loaded properly, but the table never updated
	async function getit() {
		const w = [];
		try {
			const t = await getNecrology();
			t.forEach((i) => {
				// why not use cleanDateFormat ?
				const d = new Date(i.DateDeceased);
				const dd = d.toLocaleDateString();

				const j = {
					Title: i.Title,
					'First Name': i.FirstName,
					'Last Name': i.LastName,
					'Date Deceased': dd,
					Locality: i.Country + ': ' + i.State
				};
				w.push(j);
			});
			items = w;
		} catch (err) {
			console.log(err);
		}
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Necrology</title>
</svelte:head>

{#await getit()}
	... loading ....
{:then { }}
	<Table {items} dataTableOptions={opts} />
{/await}

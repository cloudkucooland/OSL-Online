<script lang="ts">
	import { getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { getMe, report } from '../oo';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}

	async function notrenewed() {
		report('notrenewed');
	}

	async function expired() {
		report('expired');
	}

	async function email() {
		report('email');
	}

	async function annual() {
		report('annual');
	}

	async function life() {
		report('life');
	}

	async function avery() {
		report('avery');
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Reports</title>
</svelte:head>

<ul>
	<li><a href="#/reports" on:click={notrenewed}>Not Renewed this year</a></li>
	<li><a href="#/reports" on:click={expired}>"Active" but not renewed for 2 years or longer</a></li>
	<li><a href="#/reports" on:click={email}>All Member Email List</a></li>
	<li><a href="#/reports" on:click={annual}>Annual Vows Mail Merge</a></li>
	<li><a href="#/reports" on:click={life}>Life Vows Mail Merge</a></li>
	<li><a href="#/reports" on:click={avery}>Active member Avery Labels</a> (not working yet)</li>
	<li>Doxology Mailing List (member + subscriber)</li>
	<li>Doxology Email List (member + subscriber)</li>
</ul>

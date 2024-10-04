<script lang="ts">
	import { getContext } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { getMe } from '../oo';

	const { me } = getContext('oo');
	if ($me === undefined) {
		push('/Login');
	}
</script>

<svelte:head>
	<title>OSL Member Manager: Subscriber</title>
</svelte:head>

Subscriber

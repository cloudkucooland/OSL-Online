<script>
	import { getContext } from 'svelte';
	import { Button, Input, Label, Card, Heading } from 'flowbite-svelte';
	import { push } from 'svelte-spa-router';
	import { toast } from '@zerodevx/svelte-toast';
	import { getJWT } from '../oo';

	const oo = getContext('oo');

	const jwt = localStorage.getItem('jwt');
	if (jwt) {
		localStorage.removeItem('jwt');
		oo.me = null;
		toast.push('Logged out');
	}
	push('#/Login');
</script>

<div class="mx-auto mt-8 max-w-2xl">Logging out</div>

<script>
	import { getContext } from 'svelte';
	import { Button, Input, Label, Card, Heading } from 'flowbite-svelte';
	import { push } from 'svelte-spa-router';
	import { toast } from '@zerodevx/svelte-toast';
	import { getJWT } from '../oo';

	const oo = getContext('oo');
	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

	// if already logged in, go to home page
	if (oo.me) {
		push('/');
	}

	let username = $state('');
	let password = $state('');

	async function doLogin(event) {
		event.preventDefault();
		try {
			const rawResult = await getJWT(username, password);
			if (rawResult) {
				oo.me = JSON.parse(window.atob(rawResult.split('.')[1]));
				push('/');
			}
		} catch (e) {
			console.error(e);
			toast.push('Incorrect password or connection error');
		}
	}
</script>

<div class="mx-auto mt-8 max-w-2xl">
	<Card size="md" padding="xl">
		<Heading tag="h2" class="mb-4 text-2xl font-bold text-gray-900">Login</Heading>
		<p class="mb-6 text-gray-600">
			Your username is your primary email address.
			<strong>If this is your first time,</strong> please register to receive your password.
		</p>

		<form onsubmit={doLogin} class="space-y-6">
			<div>
				<Label for="username" class="mb-2">Primary Email Address</Label>
				<Input
					type="email"
					id="username"
					bind:value={username}
					placeholder="name@example.com"
					required
				/>
			</div>
			<div>
				<Label for="password" class="mb-2">Password</Label>
				<Input type="password" id="password" bind:value={password} required />
			</div>

			<div class="flex items-center justify-between">
				<a href="#/register" class="text-primary-700 text-sm hover:underline"
					>Register/Lost Password</a
				>
				<Button type="submit" color="green">Login</Button>
			</div>
		</form>
	</Card>
</div>

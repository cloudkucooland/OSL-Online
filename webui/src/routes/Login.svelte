<script>
	import { getContext } from 'svelte';
	import { Button, Input, Label, Card, Heading } from 'flowbite-svelte';
	import { toast } from '@zerodevx/svelte-toast';
	import { getJWT } from '../oo';

	const oo = getContext('oo');
	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

	// Logout logic on mount
	const jwt = localStorage.getItem('jwt');
	if (jwt) {
		localStorage.removeItem('jwt');
		oo.me = undefined; // Updated: removed '$'
		toast.push('Logged out');
	}

	let username = $state('');
	let password = $state('');

	async function doLogin(event) {
		event.preventDefault();

		if (!username || !password) {
			toast.push('Fill in both username and password');
			return;
		}

		if (!emailRegex.test(username)) {
			toast.push('Please use your email address as your username');
			return;
		}

		try {
			await getJWT(username, password);
			// Force a reload or redirect to home to refresh the state
			window.location.hash = '#/';
			window.location.reload();
		} catch (e) {
			console.error(e);
			toast.push('Incorrect password for ' + username);
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
				<a href="#/register" class="text-sm text-primary-700 hover:underline"
					>Register/Lost Password</a
				>
				<Button type="submit" color="green">Login</Button>
			</div>
		</form>
	</Card>
</div>

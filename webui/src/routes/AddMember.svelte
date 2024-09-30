<script lang="ts">
import { Label, Input, Button } from 'flowbite-svelte';
import { push } from "svelte-spa-router";
import { getMe, createMember } from "../oo";
import { toast } from '@zerodevx/svelte-toast';

const me = getMe();
if (me === undefined) {
  push('/Login');
}

let firstname = "first";
let lastname = "last";

async function create() {
   try {
     const id = await createMember(firstname, lastname);
     toast.push(`created`);
     push(`/member/${id}`);
     return true;
   } catch (err) {
     toast.push("failed to create: " + err);
     console.log(err);
   }
}
</script>

<svelte:head>
  <title>OSL Member Manager: Create member</title>
</svelte:head>

<form>
<section>
<div class="grid grid-cols-2 gap-4 px-4 py-2">
<div class="col-span-1">
  <Label for="FirstName" class="block">First Name</Label>
  <Input id="FirstName" bind:value={firstname} />
</div>
<div class="col-span-1">
  <Label for="LastName" class="block">Last Name</Label>
  <Input id="LastName" bind:value={lastname} />
</div>
<div class="col-span-2">
  <Button type="submit" on:click="{create}">Add</Button>
</div>
</div>
</section>
</form>

<script lang="ts">
import { push } from "svelte-spa-router";
import { getMe } from "../oo";

const me = getMe();
if (me === undefined) {
  push('/Login');
}

</script>

<svelte:head>
  <title>OSL Member Manager: Downloads</title>
</svelte:head>

<ul>
<li>All Member Email List</li>
<li>Annual Vows Mail Merge</li>
<li>Life Vows Mail Merge</li>
</ul>

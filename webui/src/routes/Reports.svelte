<script lang="ts">
import { push } from "svelte-spa-router";
import { getMe } from "../oo";

const me = getMe();
if (me === undefined) {
  push('/Login');
}
</script>

<svelte:head>
  <title>OSL Member Manager: Reports</title>
</svelte:head>

<ul>
<li>Not Renewed this year</li>
<li>Not Renewed for 2 years or longer</li>
</ul>

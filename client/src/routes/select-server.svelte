<script>
  import { onMount } from "svelte";
  import {listServers} from "../lib/api";
  // import {goto} from "$app/navigation";

  let servers = [];

  onMount(async () => {
    const [list, err] = await listServers();
    if (err) {
      if (err instanceof Error) {
        alert(err.message);
      } else {
        return;
      }
    }

    servers = list;
  });

  const selectWorld = (id) => {
    console.log("TODO", id);
  };
</script>

<div>
  <h1>Servers</h1>
  {#if servers.length == 0}
    <p>No servers - talk to your admin</p>
  {:else}
    <div class="char-list">
      {#each servers as item}
        <p on:click={() => selectWorld(item.ID)} class="pointer">{item.Name}</p>
      {/each}
    </div>
  {/if}
  <a href="/characters">back</a>
</div>

<style>
  .pointer {
    cursor: pointer;
  }

  .pointer:hover {
    text-decoration: underline;
  }
</style>

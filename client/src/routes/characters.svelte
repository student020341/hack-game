<script>
  import { onMount } from "svelte";
  import {listCharacters} from "../lib/api";
  import {selectedCharacter} from "../lib/stores";
  import {goto} from "$app/navigation";

  let chars = [];

  onMount(async () => {
    const [list, err] = await listCharacters();
    if (err) {
      if (err instanceof Error) {
        alert(err.message);
      } else {
        return;
      }
    }

    chars = list;
  });

  const selectChar = (id) => {
    selectedCharacter.set(id);
    goto("/select-server");
  };
</script>

<div>
  <h1>Characters</h1>
  {#if chars.length == 0}
    <p>no characters - create one!</p>
  {:else}
    <div class="char-list">
      {#each chars as char}
        <p on:click={() => selectChar(char.ID)} class="pointer">{char.Name}</p>
      {/each}
    </div>
  {/if}
  <a href="/create-character">Create Character</a>
</div>

<style>
  .pointer {
    cursor: pointer;
  }

  .pointer:hover {
    text-decoration: underline;
  }
</style>

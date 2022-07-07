<script>
  import {createCharacter} from "../lib/api";
  import {goto} from "$app/navigation";

  const input = {
    name: ""
  };

  const doCreate = async () => {
    const [ok, err] = await createCharacter(input);
    if (err) {
      if (err instanceof Error) {
        alert(err.message);
      } else {
        return;
      }
    }

    goto("/characters");
  };
</script>

<div>
  <p>Character Creation (more details later, probably)</p>
  <label for="name">Name</label>
  <input name="name" type="text" bind:value={input.name} />
  <br />
  <button on:click={doCreate}>create</button>
  <br />
  <a href="/characters">cancel</a>
</div>

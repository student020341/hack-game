
<script>
  import {env} from "../lib/env";
  import {accountToken} from "../lib/stores";
  import {goto} from "$app/navigation";
import { onMount } from "svelte";

  let user = {
    username: "",
    password: ""
  };

  let status = "checking...";

  onMount(() => {
    if ($accountToken != "") {
      goto("/characters");
    }
  });

  fetch(`${env.server}/status`).then(res => {
    switch(res.status) {
      case 200: status = "Online"; break;
      default: status = "Error"; break;
    }
  })
    .catch(err => {
      console.error(err);
      status = "Offline";
    });

  const login = () => {
    fetch(`${env.server}/api/login`, {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify(user)
    }).then(async (res) => {
      const text = await res.text();
      if (res.status == 200) {
        accountToken.set(text);
        goto("/characters");
      } else if (res.status == 404) {
        alert("incorrect password");
      } else {
        alert("login eror: "+text);
      }
    })
  };

  const register = () => {
    fetch(`${env.server}/api/accounts`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(user)
    })
      .then(async (res) => {
        let text = await res.text();
        if (res.status == 200) {
          accountToken.set(text);
          goto("/characters");
        } else {
          alert("failed to register account: "+text);
        }
      });
  };

</script>

<div>
  <h1>Game Title</h1>
  <p>Game Server: {status}</p>
  <div>
    <label for="name">Username</label>
    <input name="username" type="text" bind:value={user.username} />
    <br />
    <label for="password">Password</label>
    <input name="password" type="password" bind:value={user.password} />
    <br />
    <button on:click={login}>Login</button>
    <button on:click={register}>Register</button>
  </div>
</div>

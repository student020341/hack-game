import { accountToken } from "./stores";
import { goto } from "$app/navigation";
import { env } from "./env";

let token = "";
accountToken.subscribe(t => {
  token = t;
});

/**
 * [res, err] pattern
 * 
 * [something, false] = success
 * [null, new Error("an error")] = error
 * [null, true] = die
 */

export const listCharacters = () => {
  return fetch(`${env.server}/api/characters`, {
    headers: {
      token
    }
  })
    .then(async (res) => {
      if (res.status == 401) {
        goto("/");
        return [null, true];
      }

      return [await res.json(), false];
    })
    .catch(err => ([null, err]));
};

export const createCharacter = (arg) => {
  return fetch(`${env.server}/api/characters`, {
    method: "POST",
    headers: {
      token,
      "Content-Type": "application/json"
    },
    body: JSON.stringify(arg)
  })
    .then(async (res) => {
      if (res.status == 401) {
        goto("/");
        return [null, true];
      } else if (res.status == 200) {
        return [true, false];
      }

      return [null, new Error(await res.text())];
    })
    .catch(err => ([null, err]));
};

export const listServers = () => {
  return fetch(`${env.server}/api/servers`, {
    headers: {token}
  })
    .then(async (res) => {
      if (res.status == 401) {
        goto("/");
        return [null, true];
      } else if (res.status == 200) {
        return [await res.json(), false];
      }

      return [null, new Error(await res.text())];
    })
    .catch(err => ([null, err]));
};

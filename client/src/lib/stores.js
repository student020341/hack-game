import {writable} from "svelte/store";

// auth token
const storedAccountToken = localStorage.getItem("account-token");
export const accountToken = writable(storedAccountToken ? storedAccountToken : "");

accountToken.subscribe(value => {
  localStorage.setItem("account-token", value);
});

// selected character
export const selectedCharacter = writable("");

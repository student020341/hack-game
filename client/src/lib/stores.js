import {writable} from "svelte/store";
import {browser} from "$app/env";

if (browser) {
  
}
const storedAccountToken = localStorage.getItem("account-token");
export const accountToken = writable(storedAccountToken ? storedAccountToken : "");

accountToken.subscribe(value => {
  localStorage.setItem("account-token", value);
});

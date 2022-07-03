---

kanban-plugin: basic

---

## idea

- [ ] player can join a server town
- [ ] player can create field/dungeon instances from server town gate
- [ ] player can form groups with other players and or npcs
- [ ] groups are required to enter / join an instance created from the gate
- [ ] NPCs walk around and do random things when players are in town
- [ ] spin down towns when no players are present?
- [ ] randomize npc activity when server towns spin up
- [ ] admin can send account reset link to players
- [ ] players can transfer to other servers [[server transfer | note]]
- [ ] a server can host multiple server towns with different layouts and functions
- [ ] Expire logins?<br>- [ ] determine how many sessions are ok<br>- [ ] on login, expire older tokens


## dev

- [ ] create instance versions of account and character models that are safe to share over network
- [ ] player can login/delete/transfer account<br>- [x] login<br>- [x] logout<br>- [ ] delete<br>- [ ] transfer
- [ ] server town websocket<br>- [x] basic socket connection<br>- [x] add socket to player<br>- [ ] send/rec in town


## test

- [ ] servers - player can...<br><br>- [ ] list<br>- [ ] join<br>- [ ] leave


## done

- [ ] player can register account
- [ ] integration test for server routes
- [ ] add web framework (echo)
- [ ] move models to models package
- [ ] character<br><br>- [x] data model<br>- [x] get<br>- [x] create<br>- [x] delete




%% kanban:settings
```
{"kanban-plugin":"basic"}
```
%%
### server towns
- server towns should be configurable via json in a directory on the server
- a default town or towns with a little of everything should be included
- server towns should show players and npcs regardless of group with a configurable number of both real players and npcs
- sessions should be created, destroyed, and merged over time based on players

#### Server town parameters
A server town should be defined by a json structure like this (will show up as "**Server Town Alpha**"):
```json
{
	"name": "Alpha",
	"vendors": [...], // list of vendor definitions for this town
	"tiles": [...] // sections of geometry and npc placements that make up the map
	"words": [...] // modifiers for instances, see dungeons & fields for info
}
```

#### vendors
A vendor is a static npc that can be traded with. They should be defined by a name and an inventory. Inventory is what the npc will stock for trading. Specific items and quantities can be set or groups of items can be set that the vendor will stock randomly over time.

For example, a vendor npc named Kite that sells random twin blades and medium armor pieces, but also has 3 apples available for sale per restock cycle:
```json
{
	"name": "Kite",
	"inventory": [
		{
			// specific things you want this npc to always have
			"items": [
				{
					"item": "apple", // references an item db definition
					"stock": 3 // this vendor sells at most 3 apples per restock
				}
			],
			// names of equipment categories the npc will randomly stock
			// gear will scale over time based on player levels
			"groups": [
				"twin blades",
				"medium armor"
			]
		}
	]
}
```

#### tiles
A tile describes part of a map. A tile can be part of the terrain, an npc/vendor on the terrain, or some other interactable. See this handle of tile examples for a stone square with some grass on the side and a vendor. There's a warp gate above the grass patch. Things like interactables and npcs will have their own dimensions so they only need an x/y coordinate on the map. Data/design still ongoing here.

```json
[ // tiles array, like servertown.tiles
	{
		// leading type controls what fields are expected and render layer
		"ref": "terrain/stone", // refers to a material / img db
		// x and y refer to position from bottom left of map
		"x": 2,
		"y": 2,
		// tiles can be any whole number, relative tile unit scaled with client (?)
		"width": 4,
		"height": 4
	},
	{
		"ref": "terrain/grass",
		"x": 6, // to the right of the stone patch
		"y": 2,
		// grass trail along the side of the stone
		"width": 1,
		"height": 4
	},
	{
		// vendor standing near the middle of the stone patch
		"ref": "vendors/Kite", // some reference to a defined npc/vendor
		"x": 4,
		"y": 4
	},
	{
		// warp gate above grass
		"ref": "interactive/warp-gate",
		"x": 6,
		"y": 6
	}
]
```

### dungeons and fields
A warp gate will enable players to create an instance of a location to explore. Words defined in the server config will be presented to the player in 3 columns. The first 2 columns will have the same pool and the third column will be tied to available terrain tilesets. Words in the first 2 columns cannot be repeated.

Here's an example based on the early hack games: "Burning", "Screaming", "Sea of sand". Burning and screaming were difficulty modifiers and sea of sand created an instance that looked like a desert. 

Here is what those word definitions might look like:
```json
[ // words arry, like servertown.words
	{
		// a modifier that makes things difficult
		"word": "Burning",
		"lootMod": 2, // +2 to loot
		"levelMod": 3, // +3 to enemy levels
		"densityMod": -2 // -2 to enemy density / number of encounters
	},
	{
		// a modifier to make things easier, but less rewarding
		"word": "Happy",
		"lootMod": -1,
		"levelMod":  -2,
		"densityMod": -1
	}
]
```

By combining words like Burning and Happy, a world is created that has +1 loot, +1 level, and -3 density. By combining a hard and easy modifier, the end result is a world with very few monsters that are slightly more rewarding to engage with. The third column of words will depend on what maps a server can create. 

Interaction of dungeons in fields with modifiers and field instance definitions TBD
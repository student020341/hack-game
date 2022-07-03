// server towns, instances, dungeons?
package models

// a player that is connected to an instance
type Player struct {
	Account   Account // TODO create version of these with omitted fields that are ok to share over the network for snooping clients
	Character Character
}

// server town
type Town struct {
	ID      string
	Name    string
	Players []Player `gorm:"-"`
	// TODO NPCs
	// TODO POIs
}

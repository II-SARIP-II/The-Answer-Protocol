package server

import (
	"strings"
)

// vérifie si une entité (item ou PNJ) présente dans le monde
// correspond à la saisie du joueur, qu'il ait tapé l'ID ou le nom affiché.
func (p *Player) matchNameOrID(worldEntity, userInput string) bool {
	if strings.EqualFold(worldEntity, userInput) {
		return true
	}
	id1 := p.resolveID(worldEntity)
	id2 := p.resolveID(userInput)
	return strings.EqualFold(id1, id2)
}

// associe un nom ou un ID à son identifiant canonique (items ou PNJ).
// TODO ajouter pour quest/quests?
func (p *Player) resolveID(input string) string {
	for id, itemData := range p.Server.world.Items {
		if strings.EqualFold(input, id) || strings.EqualFold(input, itemData.Name) {
			return id
		}
	}
	for id, npcData := range p.Server.world.NPCs {
		if strings.EqualFold(input, id) || strings.EqualFold(input, npcData.Name) {
			return id
		}
	}
	return input
}

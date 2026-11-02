package cli

const helpMessage = `=== TAP - Available Commands ===
Navigation:
  LOOK                        - Inspect the current room, items, and NPCs
  MOVE <direction>            - Move in a direction (north, south, east, west)

Items & Inventory:
  INVENTORY                   - List items in your bag
  TAKE <item_id>              - Pick up an item from the room
  DROP <item_id>              - Drop an item in the room
  EQUIP <weapon_id>           - Equip a weapon from your inventory

Combat:
  ATTACK <target_id>          - Attack an enemy in the room
  DEFEND                      - Defend to reduce incoming damage
  FLEE                        - Attempt to escape combat to a random exit

Quests & Stats:
  QUESTS                      - List active and completed quests
  STATUS                       - View your health, attack power, and stats

Communication & Groups:
  CHAT room <message>         - Chat with players in the same room
  CHAT global <message>       - Broadcast message to all players
  CHAT group <message>        - Chat with your group members
  INVITE <player>             - Invite a player to your group
  JOIN <leader>               - Join a group after receiving an invite
  LEAVE                       - Leave your current group

Session & Help:
  HELP                        - Display this help message
  QUIT                        - Disconnect from the server
================================
`

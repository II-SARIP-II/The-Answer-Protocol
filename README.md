*This project has been created as part of the 42 curriculum by pgougne, stmaire.*

<div align="center">
<br>
  <img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTQPzuYKu7n0cWUYa5Kbg0_LrlEQAIURWeo9A&s" alt="42 Logo" width="400" />
  <br>
</div>

# The Answer Protocol (TAP)
![Language](https://img.shields.io/badge/Language-Go-blue)
![Static Badge](https://img.shields.io/badge/Network-green)
![Static Badge](https://img.shields.io/badge/Gaming-green)
![Static Badge](https://img.shields.io/badge/Multiplayer-green)
![Static Badge](https://img.shields.io/badge/Network_%26_system_administration-pink)
![Static Badge](https://img.shields.io/badge/Graphics-pink)
![Static Badge](https://img.shields.io/badge/Group_%26_interpersonal-pink)
![Static Badge](https://img.shields.io/badge/Rigor-pink)

## Description

### Goal
The objective of this project is to build a retro multiplayer text adventure (MUD - Multi-User Dungeon) based on a custom TCP protocol called **The Answer Protocol (RFC 42TAP)**. The project includes an authoritative TCP game server, an interactive Command-Line Interface (CLI) client, and a 2D Graphical User Interface (GUI) client.

### Overview
Multiple players connect to a shared retro game world loaded from a JSON file. Players can navigate between interconnected rooms, talk to non-player characters (NPCs), pick up and drop unique items, manage their inventory, engage in turn-based combat, accept and complete quests, and communicate in real time through global, room-scoped, or private group chat channels.

The architecture is divided into three distinct programs:
- **Server**: Authoritative TCP daemon handling game state, concurrency, combat rules, quest resolution, and protocol validation.
- **CLI Client**: Terminal client supporting raw RFC commands or user-friendly input forwarding with asynchronous event display.
- **GUI Client**: 2D graphical frontend powered by Ebitengine.

---

## Instructions

### Prerequisites
- Go 1.25 or later.
- A standard POSIX terminal environment (Linux or macOS) with `make`.

### Installation
Clone the repository and install all dependencies using the build tool:
```bash
make install
```

---

## Building and Running

The project uses a standard **Makefile** to manage the compilation, code quality checks, and execution of all three components.

### 1. Build System (Makefile Targets)

| Target | Description |
| :--- | :--- |
| `make` / `make all` | Builds all executables (`server`, `client_cli`, `client_gui`) inside `bin/`. |
| `make install` | Runs `go mod tidy` and `go mod download` (alias: `install-dependencies`). |
| `make server` | Compiles only the server binary into `bin/server`. |
| `make client-cli` | Compiles only the CLI client binary into `bin/client_cli` (alias: `client_cli`). |
| `make client-gui` | Compiles only the GUI client binary into `bin/client_gui` (alias: `client_gui`). |
| `make run-server` | Starts the server in development mode (alias: `run_server`). |
| `make run-client` | Starts the CLI client in development mode (alias: `run_cli`). |
| `make run-client-gui` | Starts the GUI client in development mode (alias: `run_gui`). |
| `make lint` | Standard verification: checks formatting with `gofmt` and runs `go vet ./...`. |
| `make lint-strict` | Advanced verification: runs `make lint` and deep static analysis with `golangci-lint` (if installed). |
| `make fmt` | Automatically formats all Go files using `gofmt`. |
| `make clean` | Removes the compiled binary folder (`bin/`). |
| `make fclean` | Full clean target calling `clean`. |
| `make re` | Recompiles the complete project from scratch (`fclean` + `all`). |

---

### 2. How to Run Each Component

Each program runs in its own separate terminal window. You can either run them directly with `make` (development mode) or compile the binaries first with `make build`.

#### Server (Authoritative TCP Game Server)
- **Direct run with Make**:
  ```bash
  # Starts on default port (:4242)
  make run-server

  # Starts on a custom port
  make run-server ARGS=8080
  ```
- **From compiled binary**:
  ```bash
  ./bin/server        # default port :4242
  ./bin/server 8080   # custom port
  ```

#### CLI Client (Interactive Terminal Client)
- **Direct run with Make**:
  ```bash
  # Connects to default server (localhost:4242)
  make run-client

  # Connects to custom address
  make run-client ARGS=localhost:8080
  ```
- **From compiled binary**:
  ```bash
  ./bin/client_cli                 # default localhost:4242
  ./bin/client_cli localhost:8080  # custom address
  ```

#### GUI Client (2D Graphical Interface)
- **Direct run with Make**:
  ```bash
  # Connects to default server (localhost:4242)
  make run-client-gui

  # Connects to custom address
  make run-client-gui ARGS=localhost:8080
  ```
- **From compiled binary**:
  ```bash
  ./bin/client_gui                 # default localhost:4242
  ./bin/client_gui localhost:8080  # custom address
  ```

---

## Architecture

### Server Architecture
The server is structured around an authoritative model where clients never mutate game state directly:

1. **Connection Handling**:
   - `net.Listen("tcp", port)` accepts incoming client connections.
   - Each client runs in its own goroutine (`handleClientSession`).
   - Commands are read line by line using a `bufio.Scanner` limited to 1024 bytes per message (RFC §9.4).

2. **Non-blocking Message Sending**:
   - Each player has a buffered channel `sendChan chan string` (capacity 128) and a dedicated writer goroutine (`writeMessages`).
   - This ensures a slow client or network lag never blocks the server or other players.
   - Helper methods (`SendOK`, `SendErr`, `SendEvt`) send formatted messages to this channel.

3. **Concurrency & Synchronization**:
   - Shared data (connected players, rooms, items, NPCs) is protected by a central `sync.Mutex` (`s.mu`).
   - All state changes (movement, combat, inventory, disconnects) run under mutex protection to avoid data races.

4. **Command Routing**:
   - Incoming commands are received by `handleCommand` in `dispatcher.go`.
   - It prevents spam (rate limiting), checks authentication, and routes each command to its handler:
     - `handle_core_commands.go`: `LOOK`, `MOVE`, `QUIT`
     - `handle_combat_commands.go`: `ATTACK`, `DEFEND`, `FLEE`, `STATUS`
     - `handle_communication_commands.go`: `CHAT`, `WHO`, `WHO_ALL`
     - `handle_management_commands.go`: `GROUP CREATE`, `GROUP INVITE`, `GROUP JOIN`, `GROUP LEAVE`
     - `handle_quest_commands.go`: `QUEST`, `QUESTS`
     - `handle_resource_interaction_commands.go`: `TAKE`, `DROP`, `INVENTORY`, `TALK`, `EQUIP`

### CLI Client Architecture
The CLI client establishes a TCP connection, completes the `OK hello proto=1` handshake, prompts for a unique username, and negotiates `CONNECT <username>`.
- **Input loop**: reads user commands from standard input and sends them directly to the server.
- **Background receiver goroutine**: continuously reads server responses (`OK`, `ERR`) and asynchronous events (`EVT`), printing them directly to the terminal without blocking user input.

### GUI Client Architecture
> 🔵TODO🔵

---

## Protocol Implementation

The implementation strictly follows the **RFC 42TAP** specification.

### Deviations and Extensions

1. **`WHO_ALL` Command (Extension)**:
   - **RFC Requirement**: RFC §5.2.2 specifies that `WHO` returns `OK players=<count>`.
   - **Our Extension**: The core `WHO` command conforms strictly to the RFC format. To allow rich user listing in both GUI and CLI clients, an extended command `WHO_ALL` was introduced:
     - Request: `WHO_ALL`
     - Response: `OK {"room": ["alice", "bob"], "server": 2}`
   - **Justification**: Preserves strict RFC compliance for `WHO` while offering detailed presence inspection for enhanced user interfaces.

2. **Movement Lock during Combat (`NO_MOVE_DURING_COMBAT`)**:
   - **Mechanism**: If a player attempts `MOVE <direction>` while in active combat (`status = "combat"`), the server responds with:
     ```text
     ERR 400 NO_MOVE_DURING_COMBAT
     ```
   - **Justification**: Free movement during combat breaks turn-based game balance. Players wishing to escape combat must deliberately execute the tactical `FLEE` command, which factors in opportunity attacks.

3. **`EQUIP` Command (RPG Extension)**:
   - **Mechanism**: Players can type `EQUIP <weapon_name>` to equip a weapon from their inventory.
   - **Response**: `OK equipped=<item_id>` or `ERR 404 ITEM_NOT_IN_INVENTORY`.
   - **Justification**: Natural extension for an RPG combat system; unequipped attacks use default fist damage.

---

## Combat System

The combat system is turn-based and resolves state transitions authoritatively on the server.

### Mechanics & Damage Formulas
- **Player Stats**: Players start with 100 HP (`max_hp: 100`).
- **Unarmed Combat**: Base fist damage is set to 2 (`FistFightDamage = 2`).
- **Armed Combat**: If a player has equipped a weapon, the damage equals the weapon's configured attack rating (e.g., `knife` = 7, `katana` = 10).
- **Enemy Counter-Attack**:
  - Each `ATTACK <target>` turn deals the player's weapon damage to the target NPC.
  - If the NPC survives, it immediately launches a counter-attack dealing its own weapon damage (or base fist damage).
  - Combat status is updated to `"combat"` and broadcast to players in the room.

- **Hostile vs. Passive Enemies**:
  - **Hostile Enemies** (`enemy`, `target`): Immediately ambush the player upon entering the room, inflicting damage and initiating combat.
  - **Passive Enemies** (`passive enemy`, e.g., the *Kitchen Chef*): Neutral by default. They will not attack unprovoked, but will turn hostile if a player steals an item (`TAKE`) or attacks them (`ATTACK`). In a shared world, an enraged passive enemy remains hostile to anyone entering the room until defeated or until the stolen item is returned with `DROP`.

### Tactical Commands: DEFEND and FLEE
- **`DEFEND`**:
  - May only be used during an active combat engagement.
  - Useful in multiplayer to absorb hits ("tank") while teammates attack the enemy.
  - Reduces incoming damage, but weakens with each consecutive use:
    - 1st Defend: **Perfect Defense** (`status: "perfect_defense"`), 0 damage taken.
    - 2nd Defend: **Weakened Defense** (`status: "weakened_defense"`), damage halved.
    - 3rd+ Defend: **Broken Defense** (`status: "broken_defense"`), full damage taken.
- **`FLEE`**:
  - Allows an engaged player to escape combat.
  - Turning your back carries a real penalty: the enemy automatically strikes you in the back as you run (opportunity attack).
  - If you survive the blow, you escape to a random adjacent room and reset to `"healthy"` status. If the damage drops you to 0 HP, you die before escaping.

### Player Defeat & Respawn
- When a player drops to 0 HP:
  - The defeat event is broadcast to the room and player's group.
  - The player immediately respawns in the safe starting zone (`hotel_entrance`).
  - Max HP is halved as a death penalty (with a safe minimum floor of 25 HP).
  - Current HP is restored to the new Max HP, and status becomes `"defeated"` (status automatically resets to "healthy" as soon as the player leaves the safe zone with MOVE).

---

## Quest System

The quest system supports progressive storylines, objective tracking, prerequisites, and unique world rewards.

### Progression & Validation
Quests are defined in `data/world.json` with clear metadata:
- **Quest Types**:
  - `fetch`: requires collecting a specific quantity of world items (e.g., 2 apples or key items).
  - `defeat`: requires defeating a designated hostile NPC target in combat.
- **Quest Lifecycle**:
  1. **Available**: Offered by a designated NPC whose role is `"quest-giver"`.
  2. **Active**: Accepted by the player. Progress is updated dynamically upon taking items or defeating targets.
  3. **Completed**: Objectives met and verified. The player returns to the quest giver and executes `QUEST` (or `QUEST <npc>`) to turn it in.
- **Prerequisites**: Some quests require prior completion of earlier quests before becoming unlockable.

### Rewards & Victory
- **Unique Rewards**: Quest rewards (like weapons) exist only once. Once claimed, the item moves directly into the player's inventory and cannot be claimed again.
- **Server Victory**: When a player completes every quest, a congratulations message is broadcast to all players on the server:
  ```text
  EVT GLOBAL CHAT Server Player <name> has completed all the quests. Congratulations !
  ```

---

## World Design

The static world is defined in [`data/world.json`] and complies with all mandatory world constraints:

### Layout & Rooms (10 interconnected rooms)
The map forms a complete circuit with optional branches:
- `hotel_entrance`: Starting safe zone and respawn hub.
- `broom_closet`: Storage room connected south of the entrance.
- `hotel_reception`: Main counter where Agent Smith gives assignments.
- `hotel_kitchen`: Working kitchen guarded by the Kitchen Chef.
- `dining_room`: Dining room connected to the kitchen.
- `corridor_1` & `corridor_2`: Hallways leading toward the private suites.
- `vip_room`: VIP lounge connecting to restricted areas.
- `secret_corridor`: Hidden corridor protected behind wood paneling.
- `secret_room`: The Mafia Boss's office.

### NPCs & Roles
The world features multiple distinct NPC roles:
- **Quest-Giver**: *Agent Smith* in `hotel_reception` initiates quest lines.
- **Dialogue**: *The Receptionist* in `hotel_reception` provides world lore and directions via `TALK`.
- **Passive Enemy**: *Kitchen Chef* in `hotel_kitchen` stays peaceful until provoked by `TAKE` or `ATTACK`.
- **Hostile Enemies & Targets**: *Big Biceps*, *Big Guns*, and *Mafia Boss* engage in combat upon contact.

### Dynamic Item Distribution
Items exist as unique instances in the world:
- `knife`: Chef's weapon found in the kitchen (7 damage).
- `colt_1911`: Handgun available in the hotel (5 damage).
- `katana`: High-damage sword (10 damage).
- `M16`: Military rifle carried by syndicate guards (12 damage).
- `agent_badge`: Special trophy awarded upon completing secret missions.
- Items dropped by players stay in the current room for other players to pick up.

---

## Server Logging

The server uses the Go standard structured logger (`log/slog`) with structured key-value attributes:

### Logging Standards
- **Structured Fields**: Every log entry includes timestamps, log level, player identity, and network addresses.
- **Log Levels**:
  - `INFO`: Connection handshakes, authenticated commands, world transitions, combat turns, and quest milestones.
  - `WARN`: Authentication failures, invalid arguments, and detected abuse patterns.
  - `ERROR`: Fatal port opening errors, JSON parsing crashes, or unrecoverable TCP errors.

### Abuse Monitoring & Safeguards
The server actively monitors potential malicious or faulty client behavior:
1. **Connection Flooding**:
   - `trackConnection` tracks rapid connections per IP.
   - Logs a warning if more than 5 connections occur within a 2-second window.
2. **Command Flooding**:
   - Monitored inside `handleCommand`.
   - Rejects execution and logs a warning if a single client exceeds 20 commands per second.
3. **Repeated Authentication Failures**:
   - Closes the connection with `ERR 900 CONNECTION_FAILED` after 3 consecutive failed `CONNECT` attempts.

---

## Group Contributions

Responsibilities across the team:

- **`stmaire`**:
  - Authoritative TCP server (networking daemon, concurrency model, session management).
  - Protocol implementation & game engine (RFC 42TAP, turn-based combat, quest system).
  - Interactive CLI client, build system (Makefile), and structured logging (`slog`).

- **`pgougne`**:
  - Original world creation and narrative design (rooms layout, NPCs, dialogues, quests).
  - 2D Graphical User Interface (GUI) client architecture using Ebitengine.
  - UI components design, real-time event handling, and graphical assets integration.


---

## Testing

### 1. Multiplayer Testing
Run the server in one terminal:
```bash
make run-server
```
Open two separate terminal windows and launch two CLI clients:
```bash
# Terminal 2
./bin/client_cli
# Connect as Alice: CONNECT alice

# Terminal 3
./bin/client_cli
# Connect as Bob: CONNECT bob
```
- Verify presence events: When Alice moves (`MOVE east`), Bob receives `EVT ROOM PRESENCE ENTER alice`.
- Verify chat: Alice types `CHAT GLOBAL Hello Bob` and Bob immediately sees the message.
- Verify group invites: Alice creates a group (`GROUP CREATE`) and invites Bob (`GROUP INVITE bob`).

### 2. Combat System Testing
1. Move to a room with an enemy (e.g., `hotel_kitchen` to test a passive enemy or `corridor_1` for an ambush).
2. In `hotel_kitchen`, type `ATTACK chef` (or `TAKE knife`) to provoke the Chef.
3. Type `DEFEND` and confirm reduced or nullified damage.
4. Test `FLEE` to escape combat into an adjacent room.

### 3. Quest System Testing
1. Meet Agent Smith at `hotel_reception`.
2. Type `QUEST undercover_agent` to receive the first mission (`gather_weapons`).
3. Type `QUESTS` to inspect current progress.
4. Collect the weapons (`colt_1911` in `broom_closet` and `knife` in `hotel_kitchen`), then return to Agent Smith and type `QUEST undercover_agent` to claim your reward (`katana`).

---

## Resources

### References
- [RFC 42TAP Specification]: The Answer Protocol formal specification.
- [Go Documentation](https://golang.org/doc/): Official Go documentation for `net`, `sync`, `bufio`, and `log/slog`.
- 🔵TODO🔵 [Ebitengine Documentation](https://ebitengine.org/): 2D game engine for Go.
- [MUD Programming & Design](https://en.wikipedia.org/wiki/MUD): Historical text-based multiplayer architectures.

### AI Usage Disclosure
Artificial Intelligence (LLM) was utilized during this project for:
- Auditing the final codebase and verifying strict compliance with RFC 42TAP requirements.
- Assisting in debugging and identifying concurrency edge cases.
- Enriching the world content based on the initial world concept.
- Drafting and translating the project README.
All proposals were reviewed, verified, and tested by the team.

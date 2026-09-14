package data

import (
	"embed"
	"encoding/json"
	"fmt"
)

// Files embeds all JSON files in the data directory
//
//go:embed scenarios/*.json buildings/*.json events/*.json dungeon/*.json crafting/*.json
var Files embed.FS

// ScenarioVillage holds initial village configuration
type ScenarioVillage struct {
	Name      string         `json:"name"`
	Level     int            `json:"level"`
	Lumber    int            `json:"lumber"`
	Stone     int            `json:"stone"`
	Rations   int            `json:"rations"`
	Treasury  int            `json:"treasury"`
	Settlers  int            `json:"settlers"`
	Workers   WorkersData    `json:"workers"`
	Buildings map[string]int `json:"buildings"`
}

// WorkersData holds worker distribution in scenario
type WorkersData struct {
	Farmers     int `json:"farmers"`
	Lumberjacks int `json:"lumberjacks"`
	Miners      int `json:"miners"`
	Blacksmiths int `json:"blacksmiths"`
	Militia     int `json:"militia"`
}

// StatsData holds character attributes
type StatsData struct {
	Might     int `json:"might"`
	Agility   int `json:"agility"`
	Resolve   int `json:"resolve"`
	Ingenuity int `json:"ingenuity"`
}

// WeaponData holds starter weapon data
type WeaponData struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	WeaponType   string  `json:"weapon_type"`
	BaseDamage   [2]int  `json:"base_damage"`
	CritRate     float64 `json:"crit_rate"`
	Initiative   int     `json:"initiative"`
	Durability   int     `json:"durability"`
	MaxDura      int     `json:"max_dura"`
	SpecialAffix string  `json:"special_affix"`
}

// ScenarioPlayer holds starter player configuration
type ScenarioPlayer struct {
	Name           string     `json:"name"`
	HP             int        `json:"hp"`
	MaxHP          int        `json:"max_hp"`
	Sanity         int        `json:"sanity"`
	MaxSanity      int        `json:"max_sanity"`
	TorchMeter     int        `json:"torch_meter"`
	MaxBackpack    int        `json:"max_backpack"`
	Stats          StatsData  `json:"stats"`
	EquippedWeapon WeaponData `json:"equipped_weapon"`
}

// ScenarioData holds complete game start settings
type ScenarioData struct {
	ScenarioID  string          `json:"scenario_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Village     ScenarioVillage `json:"village"`
	Player      ScenarioPlayer  `json:"player"`
	InitialLogs []string        `json:"initial_logs"`
}

// BuildingDef holds building definition from JSON
type BuildingDef struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	MaxLevel      int    `json:"max_level"`
	BaseWoodCost  int    `json:"base_wood_cost"`
	BaseStoneCost int    `json:"base_stone_cost"`
	BaseGoldCost  int    `json:"base_gold_cost"`
}

// LoadDefaultScenario reads scenarios/default.json
func LoadDefaultScenario() (*ScenarioData, error) {
	bytes, err := Files.ReadFile("scenarios/default.json")
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file default scenario: %w", err)
	}

	var sc ScenarioData
	if err := json.Unmarshal(bytes, &sc); err != nil {
		return nil, fmt.Errorf("gagal parsing default scenario: %w", err)
	}
	return &sc, nil
}

// LoadBuildingDefs reads buildings/buildings.json
func LoadBuildingDefs() ([]BuildingDef, error) {
	bytes, err := Files.ReadFile("buildings/buildings.json")
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file buildings.json: %w", err)
	}

	var defs []BuildingDef
	if err := json.Unmarshal(bytes, &defs); err != nil {
		return nil, fmt.Errorf("gagal parsing buildings.json: %w", err)
	}
	return defs, nil
}

// LoadAtmosphereEvents reads events/atmosphere.json
func LoadAtmosphereEvents() ([]string, error) {
	bytes, err := Files.ReadFile("events/atmosphere.json")
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file atmosphere.json: %w", err)
	}

	var events []string
	if err := json.Unmarshal(bytes, &events); err != nil {
		return nil, fmt.Errorf("gagal parsing atmosphere.json: %w", err)
	}
	return events, nil
}

// EnemyDef holds enemy stats and reward definitions
type EnemyDef struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxHP       int    `json:"max_hp"`
	MinDamage   int    `json:"min_damage"`
	MaxDamage   int    `json:"max_damage"`
	Initiative  int    `json:"initiative"`
	Defense     int    `json:"defense"`
	GoldReward  int    `json:"gold_reward"`
	ExpReward   int    `json:"exp_reward"`
}

// RoomMaterial holds material reward in treasure rooms
type RoomMaterial struct {
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

// RoomDef holds dungeon room event definitions
type RoomDef struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	EnemyIDs      []string       `json:"enemy_ids,omitempty"`
	MinGold       int            `json:"min_gold,omitempty"`
	MaxGold       int            `json:"max_gold,omitempty"`
	MinRations    int            `json:"min_rations,omitempty"`
	MaxRations    int            `json:"max_rations,omitempty"`
	Materials     []RoomMaterial `json:"materials,omitempty"`
	HealPercent   int            `json:"heal_percent,omitempty"`
	TorchBonus    int            `json:"torch_bonus,omitempty"`
	StatReq       string         `json:"stat_req,omitempty"`
	ReqValue      int            `json:"req_value,omitempty"`
	RewardGold    int            `json:"reward_gold,omitempty"`
	PenaltyDamage int            `json:"penalty_damage,omitempty"`
}

// ItemDef holds dungeon consumable, material, or treasure item definitions
type ItemDef struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Value       int    `json:"value"`
	EffectValue int    `json:"effect_value"`
}

// LoadEnemyDefs reads dungeon/enemies.json
func LoadEnemyDefs() ([]EnemyDef, error) {
	bytes, err := Files.ReadFile("dungeon/enemies.json")
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file dungeon/enemies.json: %w", err)
	}

	var defs []EnemyDef
	if err := json.Unmarshal(bytes, &defs); err != nil {
		return nil, fmt.Errorf("gagal parsing dungeon/enemies.json: %w", err)
	}
	return defs, nil
}

// LoadRoomDefs reads dungeon/rooms.json
func LoadRoomDefs() ([]RoomDef, error) {
	bytes, err := Files.ReadFile("dungeon/rooms.json")
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file dungeon/rooms.json: %w", err)
	}

	var defs []RoomDef
	if err := json.Unmarshal(bytes, &defs); err != nil {
		return nil, fmt.Errorf("gagal parsing dungeon/rooms.json: %w", err)
	}
	return defs, nil
}

// LoadItemDefs reads dungeon/items.json
func LoadItemDefs() ([]ItemDef, error) {
	bytes, err := Files.ReadFile("dungeon/items.json")
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file dungeon/items.json: %w", err)
	}

	var defs []ItemDef
	if err := json.Unmarshal(bytes, &defs); err != nil {
		return nil, fmt.Errorf("gagal parsing dungeon/items.json: %w", err)
	}
	return defs, nil
}

// RecipeDef holds blacksmith weapon crafting recipe definitions
type RecipeDef struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Description     string  `json:"description"`
	MinDamage       int     `json:"min_damage"`
	MaxDamage       int     `json:"max_damage"`
	CritRate        float64 `json:"crit_rate"`
	Initiative      int     `json:"initiative"`
	Durability      int     `json:"durability"`
	WoodCost        int     `json:"wood_cost"`
	StoneCost       int     `json:"stone_cost"`
	GoldCost        int     `json:"gold_cost"`
	BlacksmithLevel int     `json:"blacksmith_level"`
	SpecialAffix    string  `json:"special_affix"`
}

// LoadRecipeDefs reads crafting/recipes.json
func LoadRecipeDefs() ([]RecipeDef, error) {
	bytes, err := Files.ReadFile("crafting/recipes.json")
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file crafting/recipes.json: %w", err)
	}

	var defs []RecipeDef
	if err := json.Unmarshal(bytes, &defs); err != nil {
		return nil, fmt.Errorf("gagal parsing crafting/recipes.json: %w", err)
	}
	return defs, nil
}



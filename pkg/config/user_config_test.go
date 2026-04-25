package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestKeybindingGroups_UnmarshalYAML_Nested(t *testing.T) {
	yamlInput := `
keybindingGroups:
  files:
    i: { name: "Ignore options" }
  branches:
    i: { name: "Git flow options" }
    u: { name: "Branch upstream options" }
`
	cfg := &UserConfig{}
	assert.NoError(t, yaml.Unmarshal([]byte(yamlInput), cfg))
	assert.Equal(t, "Ignore options", cfg.KeybindingGroups["files"]["i"].Name)
	assert.Equal(t, "Git flow options", cfg.KeybindingGroups["branches"]["i"].Name)
	assert.Equal(t, "Branch upstream options", cfg.KeybindingGroups["branches"]["u"].Name)
}

func TestKeybindingGroups_UnmarshalYAML_LegacyFlat(t *testing.T) {
	yamlInput := `
keybindingGroups:
  i: { name: "Ignore options" }
  u: { name: "Branch upstream options" }
`
	cfg := &UserConfig{}
	assert.NoError(t, yaml.Unmarshal([]byte(yamlInput), cfg))
	assert.Equal(t, "Ignore options", cfg.KeybindingGroups["global"]["i"].Name)
	assert.Equal(t, "Branch upstream options", cfg.KeybindingGroups["global"]["u"].Name)
}

func TestKeybindingGroups_UnmarshalYAML_Empty(t *testing.T) {
	cfg := &UserConfig{}
	assert.NoError(t, yaml.Unmarshal([]byte(""), cfg))
	assert.Empty(t, cfg.KeybindingGroups)
}

func TestKeybindingGroups_UnmarshalYAML_NestedWithEmptyContext(t *testing.T) {
	// An empty inner context must NOT be misinterpreted as flat.
	yamlInput := `
keybindingGroups:
  files: {}
`
	cfg := &UserConfig{}
	assert.NoError(t, yaml.Unmarshal([]byte(yamlInput), cfg))
	files, ok := cfg.KeybindingGroups["files"]
	assert.True(t, ok, "files context should be present")
	assert.Empty(t, files)
	_, leakedToGlobal := cfg.KeybindingGroups["global"]
	assert.False(t, leakedToGlobal, "empty nested context must not be reinterpreted as flat")
}

func TestKeybindingGroups_UnmarshalYAML_RejectsScalar(t *testing.T) {
	yamlInput := `keybindingGroups: "this is wrong"`
	cfg := &UserConfig{}
	err := yaml.Unmarshal([]byte(yamlInput), cfg)
	assert.Error(t, err)
}

func TestKeybindingGroups_UnmarshalYAML_RejectsList(t *testing.T) {
	yamlInput := `
keybindingGroups:
  - a
  - b
`
	cfg := &UserConfig{}
	err := yaml.Unmarshal([]byte(yamlInput), cfg)
	assert.Error(t, err)
}

func TestKeybindingGroups_UnmarshalYAML_RejectsBadInnerType(t *testing.T) {
	// Nested-shaped but with a bad innermost value. Error message
	// should reflect the nested path.
	yamlInput := `
keybindingGroups:
  files:
    i: [not, a, mapping]
`
	cfg := &UserConfig{}
	err := yaml.Unmarshal([]byte(yamlInput), cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nested shape", "error should identify the path that was attempted")
}

// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package doc

import (
	"strings"

	"github.com/larksuite/cli/shortcuts/common"
)

type docsLegacyFlag struct {
	Name        string
	Replacement string
}

func docsAPIVersionCompatFlag() common.Flag {
	return common.Flag{
		Name:    "api-version",
		Desc:    "API version; defaults to v2",
		Default: "v2",
	}
}

func docsCreateLegacyFlags() []docsLegacyFlag {
	return []docsLegacyFlag{
		{Name: "title", Replacement: "put the title in --content, for example <title>Title</title>"},
		{Name: "markdown", Replacement: "use --content with --doc-format markdown"},
		{Name: "folder-token", Replacement: "use --parent-token"},
		{Name: "wiki-node", Replacement: "use --parent-token"},
		{Name: "wiki-space", Replacement: "use --parent-position my_library or a concrete parent position"},
	}
}

func docsFetchLegacyFlags() []docsLegacyFlag {
	return []docsLegacyFlag{
		{Name: "offset", Replacement: "use --scope outline/range/keyword/section for partial reads"},
		{Name: "limit", Replacement: "use --scope outline/range/keyword/section for partial reads"},
	}
}

func docsUpdateLegacyFlags() []docsLegacyFlag {
	return []docsLegacyFlag{
		{Name: "mode", Replacement: "use --command"},
		{Name: "markdown", Replacement: "use --content with --doc-format markdown"},
		{Name: "selection-with-ellipsis", Replacement: "use --command str_replace with --pattern"},
		{Name: "selection-by-title", Replacement: "fetch block ids first, then use --command block_replace/block_insert_after with --block-id"},
		{Name: "new-title", Replacement: "update the title through XML content in --content"},
	}
}

func docsLegacyFlagDefinitions(flags []docsLegacyFlag) []common.Flag {
	out := make([]common.Flag, 0, len(flags))
	for _, flag := range flags {
		out = append(out, common.Flag{
			Name:   flag.Name,
			Desc:   "deprecated v1 compatibility flag; run `lark-cli skills read lark-doc` for the v2 CLI skill",
			Hidden: true,
		})
	}
	return out
}

func validateDocsV2Only(runtime *common.RuntimeContext, shortcut string, legacyFlags []docsLegacyFlag) error {
	switch apiVersion := strings.TrimSpace(runtime.Str("api-version")); apiVersion {
	case "", "v2":
	case "v1":
		return docsV2OnlyError(shortcut, "--api-version v1 is no longer supported; remove --api-version or use --api-version v2 for compatibility")
	default:
		return docsV2OnlyError(shortcut, "--api-version only accepts v2; remove the flag for the default v2 behavior")
	}

	var used []string
	var replacements []string
	for _, flag := range legacyFlags {
		if !runtime.Changed(flag.Name) {
			continue
		}
		used = append(used, "--"+flag.Name)
		if flag.Replacement != "" {
			replacements = append(replacements, "--"+flag.Name+" -> "+flag.Replacement)
		}
	}
	if len(used) == 0 {
		return nil
	}

	detail := "legacy v1 flag(s) " + strings.Join(used, ", ") + " are no longer supported"
	if len(replacements) > 0 {
		detail += "; " + strings.Join(replacements, "; ")
	}
	return docsV2OnlyError(shortcut, detail)
}

func docsV2OnlyError(shortcut, detail string) error {
	return common.FlagErrorf(
		"docs %s is v2-only; %s. Run `%s` for the current schema and examples. AI agents MUST read `%s` (XML) or `%s` (Markdown) and follow the latest format rules there. MUST NOT grep/open local SKILL.md files to discover this guidance; use `lark-cli skills read ...` so content stays version-matched with this CLI. Run `%s` for the latest command flags",
		shortcut,
		detail,
		docsSkillReadCommandForShortcut(shortcut),
		docsXMLSkillReadCommand,
		docsMDSkillReadCommand,
		docsHelpCommandForShortcut(shortcut),
	)
}

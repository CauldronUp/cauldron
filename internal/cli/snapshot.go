package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/CauldronUp/cauldron/internal/client"
)

// snapshotDir is where named snapshots live, relative to the project.
//
// Inside the project rather than a home directory, because the whole point is
// that a snapshot can be committed or attached to a bug report. A capture that
// only exists on one laptop is the problem, not the fix.
const snapshotDir = ".cauldron/snapshots"

func snapshotPath(dir, name string) string {
	return filepath.Join(dir, snapshotDir, name+".json")
}

func runSnapshot(ctx *context, args []string) int {
	if len(args) == 0 {
		fmt.Fprint(ctx.stderr, "cauldron: snapshot needs a subcommand\n\n  cauldron snapshot save <name>\n  cauldron snapshot load <name>\n  cauldron snapshot list\n")
		return 1
	}

	switch args[0] {
	case "save":
		return snapshotSave(ctx, args[1:])
	case "load", "restore":
		return snapshotLoad(ctx, args[1:])
	case "list", "ls":
		return snapshotList(ctx, args[1:])
	default:
		fmt.Fprintf(ctx.stderr, "cauldron: unknown snapshot subcommand %q\n", args[0])
		return 1
	}
}

func snapshotSave(ctx *context, args []string) int {
	var base, dir string

	fs := flag.NewFlagSet("snapshot save", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)
	fs.StringVar(&dir, "dir", ".", "project directory to write the snapshot into")

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	name := arg(positional, 0)
	if name == "" {
		fmt.Fprint(ctx.stderr, "cauldron: name it, e.g. 'cauldron snapshot save before-refund-bug'\n")
		return 1
	}

	data, err := client.New(base).Snapshot()
	if err != nil {
		return fail(ctx, err)
	}

	path := snapshotPath(dir, name)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
		return 1
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)
		return 1
	}

	fmt.Fprintf(ctx.stdout, "Saved %s (%s).\n", name, humanSize(len(data)))
	fmt.Fprintf(ctx.stdout, "Restore it with: cauldron snapshot load %s\n", name)

	return 0
}

func snapshotLoad(ctx *context, args []string) int {
	var base, dir string

	fs := flag.NewFlagSet("snapshot load", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)
	fs.StringVar(&dir, "dir", ".", "project directory to read the snapshot from")

	positional, err := parseFlags(fs, args)
	if err != nil {
		return 1
	}

	name := arg(positional, 0)
	if name == "" {
		fmt.Fprint(ctx.stderr, "cauldron: which snapshot? Try 'cauldron snapshot list'\n")
		return 1
	}

	path := snapshotPath(dir, name)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(ctx.stderr, "cauldron: no snapshot named %q\n", name)

			if available := savedSnapshots(dir); len(available) > 0 {
				fmt.Fprintf(ctx.stderr, "available: %s\n", strings.Join(available, ", "))
			}

			return 1
		}

		fmt.Fprintf(ctx.stderr, "cauldron: %v\n", err)

		return 1
	}

	restored, err := client.New(base).Restore(data)
	if err != nil {
		return fail(ctx, err)
	}

	fmt.Fprintf(ctx.stdout, "Restored %s into %s.\n", name, strings.Join(restored, ", "))

	return 0
}

func snapshotList(ctx *context, args []string) int {
	var base, dir string

	fs := flag.NewFlagSet("snapshot list", flag.ContinueOnError)
	fs.SetOutput(ctx.stderr)
	urlFlag(fs, &base)
	fs.StringVar(&dir, "dir", ".", "project directory to list snapshots from")

	if _, err := parseFlags(fs, args); err != nil {
		return 1
	}

	names := savedSnapshots(dir)

	if len(names) == 0 {
		fmt.Fprintf(ctx.stdout, "No snapshots in %s.\n", filepath.Join(dir, snapshotDir))
		return 0
	}

	for _, name := range names {
		info, err := os.Stat(snapshotPath(dir, name))
		if err != nil {
			continue
		}

		fmt.Fprintf(ctx.stdout, "  %-28s %s\n", name, humanSize(int(info.Size())))
	}

	return 0
}

// savedSnapshots returns the names available in a project.
func savedSnapshots(dir string) []string {
	entries, err := os.ReadDir(filepath.Join(dir, snapshotDir))
	if err != nil {
		return nil
	}

	var names []string

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		names = append(names, strings.TrimSuffix(entry.Name(), ".json"))
	}

	sort.Strings(names)

	return names
}

func humanSize(bytes int) string {
	switch {
	case bytes >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.1f kB", float64(bytes)/(1<<10))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

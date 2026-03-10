package main

import (
	"bufio"
	"embed"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type Artwork struct {
	Name   string
	Title  string
	Artist string
	Year   string
	URL    string
	Art    string
}


//go:embed artworks/*.txt
var embeddedArtworks embed.FS

var artworks []Artwork

func init() {
artworks = loadEmbeddedArtworks()
artworks = append(artworks, loadCustomArtworks()...)
}

func customArtworksDir() string {
home, err := os.UserHomeDir()
if err != nil {
return ""
}
return filepath.Join(home, ".config", "gh-art", "artworks")
}

func parseArtworkFile(data []byte) (Artwork, error) {
content := string(data)
if !strings.HasPrefix(content, "---\n") {
return Artwork{}, fmt.Errorf("missing frontmatter")
}
rest := content[4:]
idx := strings.Index(rest, "\n---\n")
if idx < 0 {
return Artwork{}, fmt.Errorf("missing closing frontmatter delimiter")
}
header := rest[:idx]
art := rest[idx+5:]

a := Artwork{Art: strings.Trim(art, "\n")}
for _, line := range strings.Split(header, "\n") {
key, val, ok := strings.Cut(line, ":")
if !ok {
continue
}
val = strings.TrimSpace(val)
switch strings.TrimSpace(key) {
case "name":
a.Name = val
case "title":
a.Title = val
case "artist":
a.Artist = val
case "year":
a.Year = val
case "url":
a.URL = val
}
}
if a.Name == "" {
return Artwork{}, fmt.Errorf("artwork missing required 'name' field")
}
return a, nil
}

func loadEmbeddedArtworks() []Artwork {
var out []Artwork
entries, err := fs.ReadDir(embeddedArtworks, "artworks")
if err != nil {
return nil
}
for _, e := range entries {
if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
continue
}
data, err := fs.ReadFile(embeddedArtworks, "artworks/"+e.Name())
if err != nil {
continue
}
a, err := parseArtworkFile(data)
if err != nil {
continue
}
out = append(out, a)
}
return out
}

func loadCustomArtworks() []Artwork {
dir := customArtworksDir()
if dir == "" {
return nil
}
entries, err := os.ReadDir(dir)
if err != nil {
return nil
}
var out []Artwork
for _, e := range entries {
if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
continue
}
data, err := os.ReadFile(filepath.Join(dir, e.Name()))
if err != nil {
continue
}
a, err := parseArtworkFile(data)
if err != nil {
continue
}
out = append(out, a)
}
return out
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "gh-art",
		Short: "Terminal art screensaver for GitHub CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			interval, _ := cmd.Flags().GetDuration("interval")
			reveal, _ := cmd.Flags().GetBool("reveal")
			revealStyle, _ := cmd.Flags().GetString("reveal-style")
			if cmd.Flags().Changed("reveal-style") {
				reveal = true
			}
			return runScreensaver(interval, reveal, revealStyle)
		},
	}

	rootCmd.Flags().Duration("interval", 8*time.Second, "rotation interval (e.g., 10s, 1m)")
	rootCmd.Flags().Bool("reveal", false, "progressively reveal artwork instead of showing it instantly")
	rootCmd.Flags().String("reveal-style", "typewriter", "reveal animation style: typewriter, random, fade, flip")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available artworks",
		Run: func(cmd *cobra.Command, args []string) {
			for _, a := range artworks {
				fmt.Printf("%-15s  %s — %s (%s)\n", a.Name, a.Artist, a.Title, a.Year)
				fmt.Printf("%-15s  %s\n", "", a.URL)
			}
		},
	}

	showCmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show a specific artwork",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			reveal, _ := cmd.Flags().GetBool("reveal")
			revealStyle, _ := cmd.Flags().GetString("reveal-style")
			if cmd.Flags().Changed("reveal-style") {
				reveal = true
			}
			for _, a := range artworks {
				if a.Name == name {
					if reveal {
						width, height, err := term.GetSize(int(os.Stdout.Fd()))
						if err != nil || width == 0 || height == 0 {
							width, height = 80, 24
						}
						stopCh := make(chan struct{})
						revealArtwork(a, revealStyle, width, height, stopCh, nil)
					} else {
						drawArtwork(a)
					}
					return nil
				}
			}
			return fmt.Errorf("unknown artwork: %s", name)
		},
	}

	showCmd.Flags().Bool("reveal", false, "progressively reveal artwork")
	showCmd.Flags().String("reveal-style", "typewriter", "reveal animation style: typewriter, random, fade, flip")

	importCmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import a custom artwork file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			a, err := parseArtworkFile(data)
			if err != nil {
				return fmt.Errorf("invalid artwork file: %w", err)
			}
			dir := customArtworksDir()
			if dir == "" {
				return fmt.Errorf("could not determine home directory")
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return err
			}
			dest := filepath.Join(dir, a.Name+".txt")
			if err := os.WriteFile(dest, data, 0o644); err != nil {
				return err
			}
			fmt.Printf("Imported %q to %s\n", a.Name, dest)
			return nil
		},
	}

	rootCmd.AddCommand(listCmd, showCmd, importCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runScreensaver(interval time.Duration, reveal bool, revealStyle string) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	// Hide cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	// Shuffle artwork order for variety each session
	order := make([]Artwork, len(artworks))
	copy(order, artworks)
	rand.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})

	clearScreen()
	idx := 0
	var prevArt *Artwork // tracks previous artwork for flip transitions

	// stopCh controls the current reveal animation; closing it cancels the reveal
	var stopCh chan struct{}

	showArtwork := func() {
		if reveal {
			width, height, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil || width == 0 || height == 0 {
				width, height = 80, 24
			}
			stopCh = make(chan struct{})
			go func(a Artwork, style string, w, h int, ch chan struct{}, prev *Artwork) {
				revealArtwork(a, style, w, h, ch, prev)
			}(order[idx], revealStyle, width, height, stopCh, prevArt)
		} else {
			drawArtwork(order[idx])
		}
	}

	cancelReveal := func() {
		if stopCh != nil {
			select {
			case <-stopCh:
				// already closed
			default:
				close(stopCh)
			}
		}
	}

	showArtwork()

	inputCh := make(chan string, 1)
	go readInput(inputCh)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	isFlip := reveal && revealStyle == "flip"

	for {
		select {
		case <-ticker.C:
			cancelReveal()
			prevArt = &order[idx]
			idx = (idx + 1) % len(order)
			if !isFlip {
				clearScreen()
			}
			showArtwork()
			ticker.Reset(interval)
		case key := <-inputCh:
			switch key {
			case "q", "ctrl-c":
				cancelReveal()
				clearScreen()
				return nil
			case "n", "tab":
				cancelReveal()
				prevArt = &order[idx]
				idx = (idx + 1) % len(order)
				if !isFlip {
					clearScreen()
				}
				showArtwork()
				ticker.Reset(interval)
			case "p", "shift-tab":
				cancelReveal()
				prevArt = &order[idx]
				idx = (idx - 1 + len(order)) % len(order)
				if !isFlip {
					clearScreen()
				}
				showArtwork()
				ticker.Reset(interval)
			}
		case <-sigCh:
			cancelReveal()
			clearScreen()
			return nil
		}
	}
}

func readInput(ch chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return
		}
		switch b {
		case 3: // Ctrl+C
			ch <- "ctrl-c"
		case 9: // Tab
			ch <- "tab"
		case 27: // ESC - start of escape sequence
			// Peek for [ Z (shift-tab)
			next, err := reader.ReadByte()
			if err != nil {
				continue
			}
			if next == '[' {
				code, err := reader.ReadByte()
				if err != nil {
					continue
				}
				if code == 'Z' {
					ch <- "shift-tab"
				}
			}
		default:
			ch <- string(b)
		}
	}
}

func drawArtwork(a Artwork) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width == 0 || height == 0 {
		width = 80
		height = 24
	}

	lines := strings.Split(strings.Trim(a.Art, "\n"), "\n")
	info := []string{
		fmt.Sprintf("%s — %s (%s)", a.Artist, a.Title, a.Year),
		a.URL,
	}

	// 2 lines padding top/bottom minimum, plus info lines + 1 blank separator
	padding := 2
	contentLines := len(lines) + len(info) + 1
	topPad := (height - contentLines) / 2
	if topPad < padding {
		topPad = padding
	}

	for i := 0; i < topPad; i++ {
		fmt.Print("\r\n")
	}

	for _, line := range lines {
		printCentered(line, width)
	}

	fmt.Print("\r\n")
	for _, line := range info {
		printCentered(line, width)
	}
}

func printCentered(line string, width int) {
	line = strings.TrimRight(line, "\n")
	if len(line) >= width {
		fmt.Printf("%s\r\n", line)
		return
	}
	pad := (width - len(line)) / 2
	fmt.Printf("%s%s\r\n", strings.Repeat(" ", pad), line)
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

// revealArtwork progressively renders an artwork using the given style.
// It stops early if stopCh is closed. prevArt is used by the "flip" style
// to transition from the previous artwork (Vanna White tile-flip effect).
func revealArtwork(a Artwork, style string, width, height int, stopCh <-chan struct{}, prevArt *Artwork) {
	lines := strings.Split(strings.Trim(a.Art, "\n"), "\n")
	info := []string{
		fmt.Sprintf("%s — %s (%s)", a.Artist, a.Title, a.Year),
		a.URL,
	}

	padding := 2
	contentLines := len(lines) + len(info) + 1
	topPad := (height - contentLines) / 2
	if topPad < padding {
		topPad = padding
	}

	switch style {
	case "random":
		revealRandom(lines, info, width, topPad, stopCh)
	case "fade":
		revealFade(lines, info, width, topPad, stopCh)
	case "flip":
		revealFlip(lines, info, width, height, topPad, stopCh, prevArt)
	default:
		revealTypewriter(lines, info, width, topPad, stopCh)
	}
}

func stopped(stopCh <-chan struct{}) bool {
	select {
	case <-stopCh:
		return true
	default:
		return false
	}
}

// moveCursor positions the cursor at (row, col) using 1-based ANSI coordinates.
func moveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func revealTypewriter(lines, info []string, width, topPad int, stopCh <-chan struct{}) {
	// Count total characters to calibrate speed (~6 seconds total)
	total := 0
	for _, line := range lines {
		total += len(line)
	}
	delay := time.Duration(6000/max(total, 1)) * time.Millisecond
	if delay < time.Millisecond {
		delay = time.Millisecond
	}
	if delay > 5*time.Millisecond {
		delay = 5 * time.Millisecond
	}

	for i, line := range lines {
		if stopped(stopCh) {
			return
		}
		row := topPad + i + 1
		pad := 0
		if len(line) < width {
			pad = (width - len(line)) / 2
		}
		for j, ch := range line {
			if stopped(stopCh) {
				return
			}
			moveCursor(row, pad+j+1)
			fmt.Printf("%c", ch)
			time.Sleep(delay)
		}
	}

	// Show info instantly after reveal
	infoStart := topPad + len(lines) + 2
	for i, line := range info {
		moveCursor(infoStart+i, 1)
		printCentered(line, width)
	}
}

func revealRandom(lines, info []string, width, topPad int, stopCh <-chan struct{}) {
	type pos struct {
		row, col int
		ch       byte
	}

	var positions []pos
	for i, line := range lines {
		pad := 0
		if len(line) < width {
			pad = (width - len(line)) / 2
		}
		for j := 0; j < len(line); j++ {
			if line[j] != ' ' {
				positions = append(positions, pos{topPad + i + 1, pad + j + 1, line[j]})
			}
		}
	}

	rand.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	// Calibrate to ~6 seconds total
	delay := time.Duration(6000/max(len(positions), 1)) * time.Millisecond
	if delay < time.Millisecond {
		delay = time.Millisecond
	}

	for _, p := range positions {
		if stopped(stopCh) {
			return
		}
		moveCursor(p.row, p.col)
		fmt.Printf("%c", p.ch)
		time.Sleep(delay)
	}

	infoStart := topPad + len(lines) + 2
	for i, line := range info {
		moveCursor(infoStart+i, 1)
		printCentered(line, width)
	}
}

func revealFade(lines, info []string, width, topPad int, stopCh <-chan struct{}) {
	// Density tiers: heaviest characters appear first
	tiers := []string{
		"#@",
		"&%$SW",
		"?*+;:,.",
	}

	tierSet := func(tier string) map[byte]bool {
		s := map[byte]bool{}
		for i := 0; i < len(tier); i++ {
			s[tier[i]] = true
		}
		return s
	}

	// Build cumulative sets — each tier adds more characters
	shown := map[byte]bool{}

	renderTier := func(charSet map[byte]bool) {
		for i, line := range lines {
			pad := 0
			if len(line) < width {
				pad = (width - len(line)) / 2
			}
			for j := 0; j < len(line); j++ {
				if charSet[line[j]] && !shown[line[j]] {
					moveCursor(topPad+i+1, pad+j+1)
					fmt.Printf("%c", line[j])
				}
			}
		}
		for k := range charSet {
			shown[k] = true
		}
	}

	// Render each density tier with pauses between
	for _, tier := range tiers {
		if stopped(stopCh) {
			return
		}
		renderTier(tierSet(tier))
		time.Sleep(800 * time.Millisecond)
	}

	// Final tier: everything remaining
	if stopped(stopCh) {
		return
	}
	for i, line := range lines {
		pad := 0
		if len(line) < width {
			pad = (width - len(line)) / 2
		}
		for j := 0; j < len(line); j++ {
			if line[j] != ' ' && !shown[line[j]] {
				moveCursor(topPad+i+1, pad+j+1)
				fmt.Printf("%c", line[j])
			}
		}
	}

	infoStart := topPad + len(lines) + 2
	for i, line := range info {
		moveCursor(infoStart+i, 1)
		printCentered(line, width)
	}
}

func revealFlip(lines, info []string, width, height, topPad int, stopCh <-chan struct{}, prevArt *Artwork) {
	// Build a screen-sized grid for both old and new artwork so we can
	// compare every cell and flip only the ones that differ.

	// Rasterize an artwork's lines into a full-width grid row.
	rasterize := func(artLines []string, pad int) [][]byte {
		rows := make([][]byte, len(artLines))
		for i, line := range artLines {
			row := make([]byte, width)
			for j := range row {
				row[j] = ' '
			}
			offset := 0
			if len(line) < width {
				offset = (width - len(line)) / 2
			}
			for j := 0; j < len(line) && offset+j < width; j++ {
				row[offset+j] = line[j]
			}
			rows[i] = row
		}
		return rows
	}

	newGrid := rasterize(lines, topPad)

	// Build the old grid from the previous artwork
	var oldGrid [][]byte
	var oldTopPad int
	if prevArt != nil {
		prevLines := strings.Split(strings.Trim(prevArt.Art, "\n"), "\n")
		prevInfo := []string{
			fmt.Sprintf("%s — %s (%s)", prevArt.Artist, prevArt.Title, prevArt.Year),
			prevArt.URL,
		}
		prevContentLines := len(prevLines) + len(prevInfo) + 1
		oldTopPad = (height - prevContentLines) / 2
		if oldTopPad < 2 {
			oldTopPad = 2
		}
		oldGrid = rasterize(prevLines, oldTopPad)

		// Clear old info lines immediately so they don't linger during the flip
		oldInfoStart := oldTopPad + len(prevLines) + 2
		for i := 0; i < len(prevInfo)+1; i++ {
			moveCursor(oldInfoStart+i, 1)
			fmt.Print(strings.Repeat(" ", width))
		}
	}

	// Collect every screen position that needs to change
	type tile struct {
		row, col int
		ch       byte
	}

	totalRows := max(len(newGrid), len(oldGrid))
	var tiles []tile

	for i := 0; i < totalRows; i++ {
		newRow := topPad + i + 1
		for j := 0; j < width; j++ {
			var oldCh byte = ' '
			var newCh byte = ' '

			// Old artwork character at this screen position
			if oldGrid != nil {
				oldScreenRow := i + topPad - oldTopPad
				if oldScreenRow >= 0 && oldScreenRow < len(oldGrid) && j < len(oldGrid[oldScreenRow]) {
					oldCh = oldGrid[oldScreenRow][j]
				}
			}

			// New artwork character
			if i < len(newGrid) && j < len(newGrid[i]) {
				newCh = newGrid[i][j]
			}

			if oldCh != newCh {
				tiles = append(tiles, tile{newRow, j + 1, newCh})
			}
		}
	}

	// Also clear any leftover old rows that extend beyond the new artwork
	if oldGrid != nil {
		for i := len(newGrid); i < len(oldGrid); i++ {
			oldScreenRow := oldTopPad + i + 1
			for j := 0; j < width; j++ {
				if oldGrid[i][j] != ' ' {
					tiles = append(tiles, tile{oldScreenRow, j + 1, ' '})
				}
			}
		}
	}

	// Shuffle for the random tile-flip effect
	rand.Shuffle(len(tiles), func(i, j int) {
		tiles[i], tiles[j] = tiles[j], tiles[i]
	})

	// Calibrate to ~6 seconds
	delay := time.Duration(6000/max(len(tiles), 1)) * time.Millisecond
	if delay < time.Millisecond {
		delay = time.Millisecond
	}

	// Flip tiles one by one — old art morphs into new art
	for _, t := range tiles {
		if stopped(stopCh) {
			return
		}
		moveCursor(t.row, t.col)
		fmt.Printf("%c", t.ch)
		time.Sleep(delay)
	}

	// Show info after flip completes
	infoStart := topPad + len(lines) + 2
	// Clear old info area first
	for i := 0; i < 3; i++ {
		moveCursor(infoStart+i, 1)
		fmt.Print(strings.Repeat(" ", width))
	}
	for i, line := range info {
		moveCursor(infoStart+i, 1)
		printCentered(line, width)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

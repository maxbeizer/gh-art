package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
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
	Art    string
}

var artworks = []Artwork{
	{
		Name:   "great-wave",
		Title:  "The Great Wave off Kanagawa",
		Artist: "Katsushika Hokusai",
		Year:   "1831",
		Art: `
                             _..-"""""-.._
                         _.-"                "-._
                      .-"    _.._      _.._     "-.
                    ."     ."    "-. ."    "-.     ".
                   /     ./  .--.  V  .--.  \.     \
                  /     /  /  _  \   /  _  \  \     \
                 ;     ;  |  ( )  | |  ( )  |  ;     ;
                 |     |   \     /   \     /   |     |
                 |     |    "-.-"  _  "-.-"    |     |
                 |     |        .-" "-.        |     |
                 ;     ;      .'  .-.  '.      ;     ;
                  \     \    /   (   )   \    /     /
                   \     \   |    \_/    |   /     /
                    '.     .  \         /  .     .'
                      "-._  \  "-.___.-"  /  _.-"
                          "-._"._       _."_.-"
                               """-._.-"""

                 _.-~~~--..__         __..--~~~-. _
             _.-"               "-._.-"               "-._
          .-"        __..--""""--..__        "-.        "-.
        ."       .-"                  "-.       ".        ".
       /      .-"   _..-"""""""-.._    "-.       \         \
      /     ."    ."                 ".    ".      \         \
     ;     /     /   _..--.   .--.._   \     \      ;         ;
     |    ;     |   /  _   | |   _  \   |     ;     |         |
     |    |     |   | (_)  | |  (_) |   |     |     |         |
     |    |     |   \   _.-' '._  _/   |     |     |         |
     |    ;     |    "-"         "-"    |     ;     |         |
     ;     \     \   .-""-.   .-""-.   /     /      ;         ;
      \     ".     "./  .  \ /  .  \."     ."      /         /
       \       "-._  \ \__/ / \__/ /  _.-"       /         /
        ".          "-._"""-.__.-"""_.-"       ."         ."
          "-.             "--..--"             .-"       .-"
             "-._                         _.-"     _.-"
                 "--..__            __..--"   __..-"
                         ""--....--""
`,
	},
	{
		Name:   "mona-lisa",
		Title:  "Mona Lisa",
		Artist: "Leonardo da Vinci",
		Year:   "1503",
		Art: `
                   .-""""""""""-.
                .-'              '-.
              .'      .-"""""-.     '.
             /      .'  _  _   '.      \
            /      /   (o)(o)    \      \
           ;      ;     .--.      ;      ;
           |      |    (____)     |      |
           |      |     (--)      |      |
           ;      ;      '--'     ;      ;
            \      \   .------.  /      /
             \      '.|        |.'     /
              '.      \  ..  /      .'
                '-._    '==''   _.-'
                    "-._    _.-"
                        """"

                .--------------------.
               /   .--------------.   \
              /   /   .------.     \   \
             /   /   /  __   \      \   \
            /   /   /  /  \   \      \   \
           /   /   /  /    \   \      \   \
          /   /   /  /      \   \      \   \
         /   /   /  /        \   \      \   \
        /   /   /  /          \   \      \   \
       /   /   /  /            \   \      \   \
      /   /   /  /              \   \      \   \
     /   /   /  /                \   \      \   \
    /   /   /  /                  \   \      \   \
   /   /   /  /                    \   \      \   \
  /   /   /  /                      \   \      \   \
 /___/___/__/________________________\___\______\___\
`,
	},
	{
		Name:   "starry-night",
		Title:  "The Starry Night",
		Artist: "Vincent van Gogh",
		Year:   "1889",
		Art: `
                         .      *        .      *
            *      .           .      *       .
                  .      *          .      .
        .   .            .   .  *        .
               .   *          .    .
     .   .          .      *       .      .

                  .-""""-.
               .-'  _   _ '-.
              /   _( )_( )_  \
             ;    /  . .  \   ;
             |    |  ___  |   |
             ;    \ (___) /   ;
              \     '---'    /
               '-.       .-'
                  '-._.-'

             ~~~~~~    ~~~~~~     ~~~~~~
        ~~~~~     ~~~~     ~~~~~     ~~~~~
    ~~~~    ~~~~    ~~~~    ~~~~   ~~~~    ~~~~

                   .-"""""-.
                .-'   _   _ '-.
              .'    _( )_( )_  '.
             /     /  . .  \     \
            ;      |  ___  |      ;
            |      \ (___) /      |
            ;       '---'        ;
             \                 /
              '.             .'
                '-._______.-'

    ~~~~~~  ~~~~~~  ~~~~~~  ~~~~~~  ~~~~~~  ~~~~~~
  ~~~~~   ~~~~~   ~~~~~   ~~~~~   ~~~~~   ~~~~~   ~~~~
 ~~~~   ~~~~   ~~~~   ~~~~   ~~~~   ~~~~   ~~~~   ~~~~

              /\                    /\
             /  \      /\          /  \
            /    \    /  \        /    \
           /      \  /    \      /      \
          /        \/      \    /        \
         /                  \  /          \
`,
	},
	{
		Name:   "the-scream",
		Title:  "The Scream",
		Artist: "Edvard Munch",
		Year:   "1893",
		Art: `
           ~~~~~~~~            ~~~~~~~~
        ~~~~~~~~~~~~~        ~~~~~~~~~~~~~
      ~~~~~   ~~~~~~~      ~~~~~~~   ~~~~~~
     ~~~~~     ~~~~~~      ~~~~~~     ~~~~~~
     ~~~~~~    ~~~~~~      ~~~~~~    ~~~~~~~
      ~~~~~~~~~~~~~        ~~~~~~~~~~~~~~~
        ~~~~~~~~~            ~~~~~~~~~~

                    .-""""-.
                  .'  .-.  '.
                 /   (   )   \
                |  .-`---`-.  |
                | /  .-.  \ | |
                | | (   ) | | |
                | \  `-'  / | |
                 \ '-.__.-' /
                  '._    _.'
                     """"

                   /\
                  /  \
                 / /\ \
                / /  \ \
               / /    \ \
              / /      \ \
             / /        \ \
            / /          \ \
           / /            \ \
          / /              \ \
         / /                \ \
        / /                  \ \
       /_/____________________\_\

        |\                     /|
        | \                   / |
        |  \                 /  |
        |   \               /   |
        |    \             /    |
        |     \           /     |
`,
	},
	{
		Name:   "girl-pearl",
		Title:  "Girl with a Pearl Earring",
		Artist: "Johannes Vermeer",
		Year:   "1665",
		Art: `
                 .-"""""""-.
               .'            '.
              /   .-""""-.     \
             /   /  .--.  \     \
            ;   ;  (    )  ;     ;
            |   |   '--'   |     |
            |   |  .----.  |     |
            |   |  |    |  |     |
            |   |  |    |  |     |
            |   |  |    |  |     |
            |   |  |    |  |     |
            |   |  |    |  |     |
            |   |  '----'  |     |
            |   |          |     |
            ;   ;     ()   ;     ;
             \   \   (  ) /     /
              \   '.  "" .'     /
               '.   "--"     .'
                 '-.______.-'

               .----------------.
              /                  \
             /                    \
            /                      \
           /                        \
          /                          \
         /                            \
        /                              \
       /                                \
      /                                  \
     /                                    \
    /                                      \
`,
	},
	{
		Name:   "creation-adam",
		Title:  "The Creation of Adam",
		Artist: "Michelangelo",
		Year:   "1512",
		Art: `
  .-"""""""""""""""""""""""""-.
 /                                \
|    .-""""""""-.        .-""""""-. |
|   /  .----.   \      /  .----.  \|
|  |  /      \   |    |  /      \  |
|  | |  .--. |   |    | |  .--. |  |
|  | | (    )|   |    | | (    )|  |
|  | |  '--' |   |    | |  '--' |  |
|  |  \      /   |    |  \      /  |
|   \  '----'   /      \  '----'  /|
|    '-.____.-'          '-.____.-' |
 \                                /
  '-._                        _.-'
      ""--..____________..--""

                 __..---""""""---..__
             _.-"                    "-._
           ."    .-""""-.   .-""""-.     ".
          /    .'  .-.  '. .'  .-.  '.     \
         ;    /   (   )   \   (   )   \     ;
         |   ;     '-'     ;   '-'     ;    |
         |   |      _      |     _      |   |
         |   ;     (_)     ;    (_)     ;   |
         ;    \           /\           /    ;
          \    '.       .'  '.       .'    /
           '.    "-._.-"      "-._.-"    .'
             "-._                      _.-"
                 ""--..__________..--""
`,
	},
	{
		Name:   "soup-can",
		Title:  "Campbell's Soup Cans",
		Artist: "Andy Warhol",
		Year:   "1962",
		Art: `
              ______________________
             /                      \
            /                        \
           /                          \
          /                            \
         /                              \
        |                                |
        |   CAMPBELL'S                   |
        |   SOUP                         |
        |                                |
        |   TOMATO                       |
        |                                |
        |   .------------------------.   |
        |   |     M'M!  M'M!        |   |
        |   |                        |   |
        |   '------------------------'   |
        |                                |
        |                                |
        |                                |
        |                                |
        |                                |
         \                              /
          \                            /
           \                          /
            \________________________/

               ____________________
              /                    \
             /                      \
            /                        \
           /                          \
          /                            \
         /                              \
        /________________________________\
`,
	},
	{
		Name:   "mondrian",
		Title:  "Composition with Red, Blue, and Yellow",
		Artist: "Piet Mondrian",
		Year:   "1930",
		Art: `
+----------------------------------------------------+
|            |                 |                     |
|            |                 |                     |
|            |                 |         #####       |
|            |                 |         #####       |
|            |                 |         #####       |
|------------+-----------------+---------#####-------|
|            |                 |                     |
|            |                 |                     |
|            |                 |                     |
|            |                 |                     |
|            |                 |                     |
|------------+-----------------+---------------------|
|            |     #####       |                     |
|            |     #####       |                     |
|            |     #####       |                     |
|            |                 |                     |
|            |                 |                     |
|            |                 |                     |
|------------+-----------------+---------------------|
|            |                 |                     |
|            |                 |         #####       |
|            |                 |         #####       |
|            |                 |         #####       |
|            |                 |                     |
+----------------------------------------------------+
`,
	},
	{
		Name:   "son-of-man",
		Title:  "The Son of Man",
		Artist: "René Magritte",
		Year:   "1964",
		Art: `
                .-""""""""""-.
              .'                '.
             /    .-""""""-.      \
            /    /  .--.   \      \
           ;    |  (    )  |      ;
           |    |   '--'   |      |
           |    |   .--.   |      |
           |    |  (    )  |      |
           |    |   '--'   |      |
           |    |          |      |
           ;    |   .--.   |      ;
            \   |  (    )  |     /
             \  |   '--'   |    /
              '.|          |  .'
                '-.______. -'

                 .-"""""-.
               .'  .--.  '.
              /   (    )   \
             |     '--'     |
             |     .--.     |
             |    (____)    |
             |              |
             |      .-.     |
             |     (   )    |
             |      `-'     |
              \            /
               '.        .'
                 '-.__.-'
`,
	},
	{
		Name:   "don-quixote",
		Title:  "Don Quixote",
		Artist: "Pablo Picasso",
		Year:   "1955",
		Art: `
                     .-.
                    /   \
                   /  .-.
                  /  /   \
                 /  /     \
                /  /       \
               /  /         \
              /  /           \
             /  /             \
            /  /               \
           /  /                 \
          /  /                   \
         /  /                     \
        /  /                       \
       /  /                         \
      /  /                           \
     /  /                             \
    /  /                               \
   /  /                                 \
  /  /                                   \
 /__/_____________________________________\

        o                 o
       /|\               /|\
       / \               / \

               *
              /|\
             / | \
               |
               |
              / \
             /   \

        ~~~~                 ~~~~
      ~~~~~~~~             ~~~~~~~~
`,
	},
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "gh-art",
		Short: "Terminal art screensaver for GitHub CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			interval, _ := cmd.Flags().GetDuration("interval")
			return runScreensaver(interval)
		},
	}

	rootCmd.Flags().Duration("interval", 8*time.Second, "rotation interval (e.g., 10s, 1m)")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available artworks",
		Run: func(cmd *cobra.Command, args []string) {
			for _, a := range artworks {
				fmt.Printf("%s\n", a.Name)
			}
		},
	}

	showCmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show a specific artwork",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			for _, a := range artworks {
				if a.Name == name {
					drawArtwork(a)
					return nil
				}
			}
			return fmt.Errorf("unknown artwork: %s", name)
		},
	}

	rootCmd.AddCommand(listCmd, showCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runScreensaver(interval time.Duration) error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	clearScreen()
	idx := 0
	drawArtwork(artworks[idx])

	inputCh := make(chan byte, 1)
	go readInput(inputCh)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			idx = (idx + 1) % len(artworks)
			clearScreen()
			drawArtwork(artworks[idx])
		case b := <-inputCh:
			switch b {
			case 'q':
				clearScreen()
				return nil
			case 'n':
				idx = (idx + 1) % len(artworks)
				clearScreen()
				drawArtwork(artworks[idx])
			case 'p':
				idx = (idx - 1 + len(artworks)) % len(artworks)
				clearScreen()
				drawArtwork(artworks[idx])
			}
		case <-sigCh:
			clearScreen()
			return nil
		}
	}
}

func readInput(ch chan<- byte) {
	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return
		}
		ch <- b
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
		fmt.Sprintf("%s — %s", a.Artist, a.Title),
		fmt.Sprintf("%s", a.Year),
	}

	contentLines := len(lines) + len(info) + 1
	topPad := (height - contentLines) / 2
	if topPad < 0 {
		topPad = 0
	}

	for i := 0; i < topPad; i++ {
		fmt.Println()
	}

	for _, line := range lines {
		printCentered(line, width)
	}

	fmt.Println()
	for _, line := range info {
		printCentered(line, width)
	}
}

func printCentered(line string, width int) {
	line = strings.TrimRight(line, "\n")
	if len(line) >= width {
		fmt.Println(line)
		return
	}
	pad := (width - len(line)) / 2
	fmt.Printf("%s%s\n", strings.Repeat(" ", pad), line)
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

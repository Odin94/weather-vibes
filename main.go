package main

import (
	"log"
	"os"

	"github.com/Odin94/weather-vibes/weather-vibes"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := weathervibes.NewModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

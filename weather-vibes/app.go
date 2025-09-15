package weathervibes

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	spinner    Spinner
	loading    bool
	location   *GeoResponse
	weather    *WeatherResponse
	err        error
	quitting   bool
}

type LocationMsg GeoResponse
type WeatherMsg WeatherResponse
type ErrMsg error

func NewModel() Model {
	return Model{
		spinner: NewSpinner(),
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick(),
		fetchLocation(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			m.loading = true
			return m, fetchLocation()
		default:
			return m, nil
		}

	case LocationMsg:
		m.location = (*GeoResponse)(&msg)
		return m, fetchWeather(m.location.Lat, m.location.Lon)

	case WeatherMsg:
		m.weather = (*WeatherResponse)(&msg)
		m.loading = false
		return m, nil

	case ErrMsg:
		m.err = error(msg)
		m.loading = false
		return m, nil

	default:
		if m.loading {
			return m, m.spinner.Tick()
		}
		return m, nil
	}
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit. Press r to reload.", m.err)
	}

	if m.loading {
		return fmt.Sprintf("\n\n   %s Loading weather data...\n\nPress q to quit.", m.spinner.View())
	}

	if m.location != nil && m.weather != nil {
		temperature := m.weather.CurrentWeather.Temperature2m
		unit := m.weather.CurrentUnits.Temperature2m
		
		style := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

		content := fmt.Sprintf("🌍 %s, %s\n", m.location.City, m.location.Country)
		content += fmt.Sprintf("🌡️ Temperature: %.1f %s", temperature, unit)

		return style.Render(content) + "\n\nPress q to quit. Press r to reload."
	}

	return "Loading..."
}

func fetchLocation() tea.Cmd {
	return func() tea.Msg {
		location, err := GetGeolocation()
		if err != nil {
			return ErrMsg(err)
		}
		return LocationMsg(location)
	}
}

func fetchWeather(lat, lon float64) tea.Cmd {
	return func() tea.Msg {
		weather, err := GetWeather(lat, lon)
		if err != nil {
			return ErrMsg(err)
		}
		return WeatherMsg(weather)
	}
}

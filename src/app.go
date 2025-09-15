package weathervibes

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	spinner       Spinner
	loading       bool
	location      *GeoResponse
	weather       *WeatherResponse
	weeklyWeather *WeeklyWeatherResponse
	selectedDay   int
	err           error
	quitting      bool
}

type LocationMsg GeoResponse
type WeatherMsg WeatherResponse
type WeeklyWeatherMsg WeeklyWeatherResponse
type ErrMsg error

func NewModel() Model {
	return Model{
		spinner:     NewSpinner(),
		loading:     true,
		selectedDay: 0, // Start with today (index 0)
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
		case "left", "h":
			// Navigate to previous day
			if m.weeklyWeather != nil && m.selectedDay > 0 {
				m.selectedDay--
			}
			return m, nil
		case "right", "l":
			// Navigate to next day
			if m.weeklyWeather != nil && m.selectedDay < len(m.weeklyWeather.Daily.Time)-1 {
				m.selectedDay++
			}
			return m, nil
		default:
			return m, nil
		}

	case LocationMsg:
		m.location = (*GeoResponse)(&msg)
		return m, tea.Batch(
			fetchWeather(m.location.Lat, m.location.Lon),
			fetchWeeklyWeather(m.location.Lat, m.location.Lon),
		)

	case WeatherMsg:
		m.weather = (*WeatherResponse)(&msg)
		// Only stop loading if we have both current and weekly weather
		if m.weeklyWeather != nil {
			m.loading = false
		}
		return m, nil

	case WeeklyWeatherMsg:
		m.weeklyWeather = (*WeeklyWeatherResponse)(&msg)
		// Only stop loading if we have both current and weekly weather
		if m.weather != nil {
			m.loading = false
		}
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

func fetchWeeklyWeather(lat, lon float64) tea.Cmd {
	return func() tea.Msg {
		weeklyWeather, err := GetWeeklyWeather(lat, lon)
		if err != nil {
			return ErrMsg(err)
		}
		return WeeklyWeatherMsg(weeklyWeather)
	}
}

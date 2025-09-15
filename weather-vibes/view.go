package weathervibes

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// TODO list:
// - Change colors based on how bright it is outside / sunrise/sunset

const (
	freshBreezeWindSpeed = 29  // according to Royal Meteorological Society
	moderateUvIndex = 3
)


func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit. Press r to reload.", m.err)
	}

	if m.loading {
		return fmt.Sprintf("\n\n   %s Loading weather data...\n\nPress q to quit.", m.spinner.View())
	}

	if m.location != nil && m.weather != nil && m.weeklyWeather != nil {
		temperature := m.weather.CurrentWeather.Temperature2m
		temperatureUnit := m.weather.CurrentUnits.Temperature2m
		rain := m.weather.CurrentWeather.Rain
		rainUnit := m.weather.CurrentUnits.Rain
		showers := m.weather.CurrentWeather.Showers
		snowfall := m.weather.CurrentWeather.Snowfall
		snowfallUnit := m.weather.CurrentUnits.Snowfall
		showersUnit := m.weather.CurrentUnits.Showers
		uvIndex := m.weather.CurrentWeather.UvIndex
		uvIndexUnit := m.weather.CurrentUnits.UvIndex
		windSpeed := m.weather.CurrentWeather.WindSpeed10m
		windSpeedUnit := m.weather.CurrentUnits.WindSpeed10m

		// current weather box (left)
		currentStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Width(30)

		currentContent := lipgloss.NewStyle().Underline(true).Render("Current") + "\n\n"
		currentContent += fmt.Sprintf("🌍 %s, %s\n", m.location.City, m.location.Country)
		currentContent += fmt.Sprintf("🌡️ Temperature: %.1f %s", temperature, temperatureUnit)
		
		if rain > 0 {
			currentContent += fmt.Sprintf("\n💧 Rain: %.1f %s", rain, rainUnit)
		}
		if showers > 0 {
			currentContent += fmt.Sprintf("\n🌧️ Showers: %.1f %s", showers, showersUnit)
		}
		if snowfall > 0 {
			currentContent += fmt.Sprintf("\n🌨️ Snowfall: %.1f %s", snowfall, snowfallUnit)
		}
		if (uvIndex > moderateUvIndex) {
			currentContent += fmt.Sprintf("\n🌞 UV Index: %.1f %s", uvIndex, uvIndexUnit)
		}
		if (windSpeed > freshBreezeWindSpeed) {
			currentContent += fmt.Sprintf("\n🌬️ Wind Speed: %.1f %s", windSpeed, windSpeedUnit)
		}

		// daily weather box (right)
		dailyStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(1, 2).
			Width(30)

		var dateDisplay string
		if m.selectedDay < len(m.weeklyWeather.Daily.Time) {
			date := m.weeklyWeather.Daily.Time[m.selectedDay]
			switch m.selectedDay {
			case 0:
				dateDisplay = "Today"
			case 1:
				dateDisplay = "Tomorrow"
			default:
				if len(date) >= 10 {
					if parsedDate, err := time.Parse("2006-01-02", date[:10]); err == nil {
						weekday := parsedDate.Weekday().String()[:3]
						day := date[8:10]   // DD
						month := date[5:7]  // MM
						dateDisplay = fmt.Sprintf("%s %s.%s", weekday, day, month)
					} else {
						day := date[8:10]   // DD
						month := date[5:7]  // MM
						dateDisplay = fmt.Sprintf("%s.%s", day, month)
					}
				} else {
					dateDisplay = date
				}
			}
		}
		
		dailyContent := lipgloss.NewStyle().Underline(true).Render(fmt.Sprintf("Daily - %s", dateDisplay)) + "\n\n"
		if m.selectedDay < len(m.weeklyWeather.Daily.Time) {
			maxTemp := m.weeklyWeather.Daily.Temperature2mMax[m.selectedDay]
			minTemp := m.weeklyWeather.Daily.Temperature2mMin[m.selectedDay]
			weatherCode := m.weeklyWeather.Daily.WeatherCode[m.selectedDay]
			tempUnit := m.weeklyWeather.DailyUnits.Temperature2mMax
			
			weatherDesc := GetWeatherDescription(weatherCode)
			weatherEmoji := GetWeatherEmoji(weatherDesc)
			
			dailyContent += fmt.Sprintf("%s %s", weatherEmoji, weatherDesc)
			
			dailyContent += fmt.Sprintf("\n🌡️ Temperature: %.0f-%.0f%s", minTemp, maxTemp, tempUnit)

			if m.selectedDay < len(m.weeklyWeather.Daily.RainSum) && m.weeklyWeather.Daily.RainSum[m.selectedDay] > 0 {
				dailyContent += fmt.Sprintf("\n💧 Rain: %.1f %s", m.weeklyWeather.Daily.RainSum[m.selectedDay], m.weeklyWeather.DailyUnits.RainSum)
			}
			if m.selectedDay < len(m.weeklyWeather.Daily.SnowfallSum) && m.weeklyWeather.Daily.SnowfallSum[m.selectedDay] > 0 {
				dailyContent += fmt.Sprintf("\n🌨️ Snow: %.1f %s", m.weeklyWeather.Daily.SnowfallSum[m.selectedDay], m.weeklyWeather.DailyUnits.SnowfallSum)
			}
			if m.selectedDay < len(m.weeklyWeather.Daily.UvIndexMax) && m.weeklyWeather.Daily.UvIndexMax[m.selectedDay] > moderateUvIndex {
				dailyContent += fmt.Sprintf("\n🌞 UV Index: %.1f %s", m.weeklyWeather.Daily.UvIndexMax[m.selectedDay], m.weeklyWeather.DailyUnits.UvIndexMax)
			}
			if m.selectedDay < len(m.weeklyWeather.Daily.WindSpeed10mMax) && m.weeklyWeather.Daily.WindSpeed10mMax[m.selectedDay] > freshBreezeWindSpeed {
				dailyContent += fmt.Sprintf("\n🌬️ Wind: %.1f %s", m.weeklyWeather.Daily.WindSpeed10mMax[m.selectedDay], m.weeklyWeather.DailyUnits.WindSpeed10mMax)
			}
		}

		dailyContent += "\n\n← → Navigate days"

		currentBox := currentStyle.Render(currentContent)
		dailyBox := dailyStyle.Render(dailyContent)
		
		combined := lipgloss.JoinHorizontal(lipgloss.Top, currentBox, dailyBox)
		
		return combined + "\n\nq to quit | r to reload"
	}

	return "Loading..."
}
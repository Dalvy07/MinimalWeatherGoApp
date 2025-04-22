package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Встраиваем статические файлы прямо в бинарник
//go:embed static/*
var staticFiles embed.FS

// Структуры для хранения данных о погоде
type CurrentWeather struct {
	Condition  struct{ Text string `json:"text"` } `json:"condition"`
	TempC      float64 `json:"temp_c"`
	Humidity   int     `json:"humidity"`
	WindKph    float64 `json:"wind_kph"`
	PressureMb float64 `json:"pressure_mb"`
	FeelslikeC float64 `json:"feelslike_c"`
	LastUpdate string  `json:"last_updated"`
}

type Location struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type WeatherResponse struct {
	Location Location      `json:"location"`
	Current  CurrentWeather `json:"current"`
}

type WeatherData struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Condition   string  `json:"condition"`
	Temperature float64 `json:"temperature"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"windSpeed"`
	Pressure    float64 `json:"pressure"`
	Feelslike   float64 `json:"feelslike"`
	LastUpdated string  `json:"last_updated"`
}

// Словарь с локациями
var locations = map[string][]string{
	"Poland":        {"Warsaw", "Krakow", "Gdansk", "Wroclaw", "Poznan", "Lublin"},
	"Germany":       {"Berlin", "Munich", "Hamburg", "Frankfurt", "Cologne"},
	"France":        {"Paris", "Marseille", "Lyon", "Toulouse", "Nice"},
	"Great Britain": {"London", "Manchester", "Liverpool", "Birmingham", "Glasgow"},
	"Italy":         {"Rome", "Milan", "Naples", "Florence", "Venice"},
}

// API ключ получаем из переменной окружения
var apiKey string

func init() {
    // Проверка, что apiKey был установлен при компиляции
    if apiKey == "" {
        log.Fatal("API ключ не установлен. Неправильная сборка приложения.")
    }
}

func main() {
	// Получаем порт из переменной окружения или используем 3000 по умолчанию
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Настраиваем маршруты
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/api/countries", getCountries)
	http.HandleFunc("/api/cities/", getCities)
	http.HandleFunc("/api/weather", getWeather)

	// Запускаем сервер
	fmt.Printf("Сервер насłuchuje на порте TCP: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Обработчик для корневого маршрута - отдает HTML страницу
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	content, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "Не удалось прочитать файл index.html", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// Обработчик для статических CSS и JS файлов
func serveStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:] // Убираем начальный слеш
	content, err := staticFiles.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	// Определяем MIME-тип на основе расширения файла
	switch {
	case strings.HasSuffix(path, ".css"):
		w.Header().Set("Content-Type", "text/css")
	case strings.HasSuffix(path, ".js"):
		w.Header().Set("Content-Type", "application/javascript")
	}
	
	w.Write(content)
}

// Обработчик для API запроса к списку стран
func getCountries(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	
	countries := make([]string, 0, len(locations))
	for country := range locations {
		countries = append(countries, country)
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

// Обработчик для API запроса к списку городов для конкретной страны
func getCities(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	
	country := r.URL.Path[len("/api/cities/"):]
	cities, exists := locations[country]
	
	w.Header().Set("Content-Type", "application/json")
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Kraj nie został znaleziony"})
		return
	}
	
	json.NewEncoder(w).Encode(cities)
}

// Обработчик для API запроса погоды
func getWeather(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	
	city := r.URL.Query().Get("city")
	country := r.URL.Query().Get("country")
	
	if city == "" || country == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Wymagane są parametry city i country"})
		return
	}
	
	weatherData, err := getWeatherFromAPI(city)
	if err != nil {
		log.Printf("Błąd podczas pobierania danych pogodowych: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Błąd podczas pobierania danych pogodowych"})
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)
}

// Функция получения данных о погоде из API
func getWeatherFromAPI(city string) (WeatherData, error) {
	var weatherData WeatherData
	
	apiURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", 
		apiKey, url.QueryEscape(city))
	
	resp, err := http.Get(apiURL)
	if err != nil {
		return weatherData, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return weatherData, fmt.Errorf("HTTP error! Status: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weatherData, err
	}
	
	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return weatherData, err
	}
	
	weatherData = WeatherData{
		City:        weatherResp.Location.Name,
		Country:     weatherResp.Location.Country,
		Condition:   weatherResp.Current.Condition.Text,
		Temperature: weatherResp.Current.TempC,
		Humidity:    weatherResp.Current.Humidity,
		WindSpeed:   weatherResp.Current.WindKph,
		Pressure:    weatherResp.Current.PressureMb,
		Feelslike:   weatherResp.Current.FeelslikeC,
		LastUpdated: weatherResp.Current.LastUpdate,
	}
	
	return weatherData, nil
}

// Функция для включения CORS
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
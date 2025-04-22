# MinimalWeatherGoApp

Repository for PAwCHO Zad1

## Setup Instructions

### Getting an API Key

1. Sign up on [Weather API](https://www.weatherapi.com/)
2. After logging in, change API Response Fields in Dashboard:
   - Unmark all responses with Imperial units (Fahrenheit, miles per hour, inches) 🦅
3. Create an `api_key.txt` file in the project root folder and paste your API key

## Docker Commands

### Building the Container

- **Default build:**
  ```bash
  docker build --secret id=api_key,src=api_key.txt -t weather-app .
  ```

- **Build with custom port:**
  ```bash
  docker build --secret id=api_key,src=api_key.txt --build-arg PORT=8080 -t weather-app .
  ```

### Running the Container

- **Run with port forwarding:**
  ```bash
  docker run -p 3000:3000 weather-app
  ```
# MinimalWeatherGoApp
Repo for PAwCHO Zad1

## Instructions
### Get API key
- First of all you have to signup on [weather api](https://www.weatherapi.com/)
- After loging you have to change API Response Fields in Dashboard. Just unmark all responses with units of freedom 🦅🦅🦅 (Faringeint, miles per hour, inches)
- You have to create a api_key.txt file and paste there your API key for 

## Commands
### Build
- Default build: docker build --secret id=api_key,src=api_key.txt -t weather-app .
- Build with different port: docker build --secret id=api_key,src=api_key.txt --build-arg PORT=8080 -t weather-app .


### Run
- Run with port forwarding: docker run -p 3000:3000 weather-app
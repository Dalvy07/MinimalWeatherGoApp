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

- **Build from local files (Must comment part that fetch from github and uncomment local copy. And must have source files locally):**
  ```bash
  docker build --secret id=api_key,src=api_key.txt --build-arg PORT=8080 -t weather-app .
  ```
- **Build using source code from github repo: (Must have only dockerfile, api_key.txt and configured ssh locally). BUT ALWAYS USE CASHH EVEN IF CODE CHANGED**
  ```bash
  docker build --ssh github=~/.ssh/your_private_key --secret id=api_key,src=api_key.txt -t weather-app .
  ```
- **Build using source code from github repo without cache**
  ```bash
  docker build --no-cache --ssh github=~/.ssh/gh_lab6 --secret id=api_key,src=api_key.txt -t weather-app .
  ```
  ```bash
  docker build --build-arg CACHE_BUST=$(date +%s) --ssh github=~/.ssh/gh_lab6 --secret id=api_key,src=api_key.txt -t weather-app .
  ```

- **You can send your image to dockerhub using:**
  ```bash
  docker tag weather-app your_username/weather-app:latest
  ```
  ```bash
  docker push weather-app your_username/weather-app:latest
  ```


### Running the Container

- **Simple run with:**
  ```bash
  docker run -p 3000:3000 weather-app
  ```

### Analize image
- **Basic information about image (can check layers, env, itp):**
  ```bash
  docker inspect weather-app
  ```
- **Build history:**
  ```bash
  docker history weather-app
  ```
- **Unfortunately, you can only view the logs using docker desktop or the terminal in which the coneytner was started. Since I tried to make the container as small as possible, the base image has nothing but scratch, which makes it impossible to connect to the container terminal to view logs, because there is no terminal there.**

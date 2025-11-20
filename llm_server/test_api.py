import requests
import json

def test_health():
    response = requests.get("http://localhost:5000/health")
    print("Health check:", response.json())

def test_analyze():
    payload = {
        "job_description": "Looking for a backend engineer with experience in Go, REST APIs, and microservices. Must have strong system design skills.",
        "projects": [
            {
                "name": "SimplifyJobsDaemon",
                "description": "Go daemon that monitors job postings, compares new listings, sends notifications, and caches data efficiently"
            },
            {
                "name": "WeatherApp",
                "description": "React frontend app that shows weather forecasts using OpenWeather API"
            },
            {
                "name": "E-commerce Backend",
                "description": "Python Flask REST API with PostgreSQL, handles orders, payments, and inventory management"
            }
        ]
    }

    print("\nSending request to /analyze...")
    response = requests.post(
        "http://localhost:5000/analyze",
        json=payload,
        timeout=120
    )

    print(f"Status: {response.status_code}")
    print(f"Response: {json.dumps(response.json(), indent=2)}")


if __name__ == "__main__":
    test_health()
    test_analyze()



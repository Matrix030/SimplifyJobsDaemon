from flask import Flask, request, jsonify
import requests
import json
from typing import List, Dict
from pydantic import BaseModel, ValidationError
from prompt_templates import create_project_selection_prompt

app = Flask(__name__)

#Ollama Configuration
OLLAMA_URL = "http://localhost:11434/api/generate"
MODEL_NAME = "gpt-oss:20b"


class ProjectInput(BaseModel):
    name: str
    description: str


class AnalyzeRequest(BaseModel):
    job_description: str
    projects: List[ProjectInput]

class AnalyzeResponse(BaseModel):
    selected_projects: List[str]
    reasoning: str

def query_ollama(prompt: str) -> str:
    payload = {
        "model": MODEL_NAME,
        "prompt": prompt,
        "stream": False,
        "temperature": 0.3,
    }

    try:
        response = requests.post(OLLAMA_URL, json=payload, timeout=120)
        response.raise_for_status()
        result =  response.json()
        return result["response"]
        
    except requests.exceptions.RequestException as e:
        raise Exception(f"Ollama request failed: {str(e)}")


def extract_json_from_response(text: str) -> dict:
    if "```json" in text:
        start = text.find("```json") + 7
        end = text.find("```", start)
        text = text[start:end].strip()

    elif "```" in text:
        start = text.find("```") + 3
        end = text.find("```", start)
        text = text[start:end].strip()

    try:
        return json.loads(text)
    except json.JSONDecodeError:
        start = text.find("{")
        end = text.rfind("}") + 1
        if start != -1 and end != 0:
            return json.loads(text[start:end])
        raise

@app.route('/health', methods=['GET'])

def health_check():
    try:
        response = requests.get("http://localhost:11434/api/tags", timeout=5)
        if response.status_code == 200:
            return jsonify({"status": "healthy", "ollama":"running"}), 200
    
    except:
        pass

    return jsonify({"status":"unhealthy", "ollama":"running"}), 503

@app.route('/analyze', methods=['POST'])

def analyze_job():
    try:
        data = AnalyzeRequest(**request.json)

        project_dict = [{"name": p.name, "description": p.description} for p in data.projects]

        prompt = create_project_selection_prompt(
            data.job_description,
            project_dict
        )

        print(f"Querying Ollama with model: {MODEL_NAME}")
        llm_response = query_ollama(prompt)
        print(f"Raw LLM response: {llm_response}")

        parsed = extract_json_from_response(llm_response)


        result = AnalyzeRequest(**parsed)
        return jsonify(result.dict()), 200
    
    except ValidationError as e:
        return jsonify({"error": "Invalid request format", "details": str(e)}), 400
    
    except json.JSONDecodeError as e:
        return jsonify({"error": "Failed to parse LLM response as JSON", "details": str(e)}), 500

    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    print(f"Starting LLM server with model: {MODEL_NAME}")
    print("Make sure Ollama is running: ollama sever")
    app.run(host="0.0.0.0", port=5000, debug=True)




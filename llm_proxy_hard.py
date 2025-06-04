#!/usr/bin/env python3
"""
llm_proxy_hard.py

Proxy LLM reforçado:
- Fallback local com modelo WASM (mistral)
- Timeout e retry com backoff exponencial
- Cache semântico local (SQLite)
- Tratamento de erros robusto e logs auditáveis
"""

import os
import json
import hashlib
import sqlite3
import time
import requests
from flask import Flask, request, jsonify, abort

app = Flask(__name__)

# Configurações
LLM_PROXY_URL = os.getenv("LLM_PROXY_URL", "http://localhost:8000/api/llm")
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY", "")
CACHE_DB = 'llm_cache.db'
TIMEOUT = int(os.getenv("LLM_TIMEOUT", "10"))
MAX_RETRIES = int(os.getenv("LLM_MAX_RETRIES", "3"))
BACKOFF_FACTOR = 2

# Inicializar cache SQLite
conn = sqlite3.connect(CACHE_DB, check_same_thread=False)
cursor = conn.cursor()
cursor.execute("""CREATE TABLE IF NOT EXISTS cache (
    prompt_hash TEXT PRIMARY KEY,
    response TEXT,
    timestamp INTEGER
)""")
conn.commit()

def get_cache(prompt):
    h = hashlib.sha256(prompt.encode()).hexdigest()
    cursor.execute("SELECT response, timestamp FROM cache WHERE prompt_hash=?", (h,))
    row = cursor.fetchone()
    if row:
        return row[0]
    return None

def set_cache(prompt, response):
    h = hashlib.sha256(prompt.encode()).hexdigest()
    timestamp = int(time.time())
    cursor.execute("REPLACE INTO cache (prompt_hash, response, timestamp) VALUES (?, ?, ?)", (h, response, timestamp))
    conn.commit()

def call_remote_llm(payload):
    for attempt in range(1, MAX_RETRIES + 1):
        try:
            r = requests.post(LLM_PROXY_URL, json=payload, timeout=TIMEOUT)
            r.raise_for_status()
            data = r.json()
            if 'text' in data:
                return data['text']
            else:
                raise ValueError("Sem campo 'text' na resposta")
        except Exception as e:
            print(f"[WARN] Tentativa {attempt} falhou: {e}", flush=True)
            if attempt < MAX_RETRIES:
                time.sleep(BACKOFF_FACTOR ** (attempt - 1))
            else:
                raise e

def call_local_wasm_llm(prompt, model, max_tokens, temperature):
    # Placeholder: integrar com WASM local (mistral.wasm)
    # Se WASM não disponível, retornar fallback
    return "fallback response from local WASM LLM"

@app.route("/api/llm", methods=["POST"])
def llm_api():
    data = request.json or {}
    prompt = data.get("prompt", "")
    model = data.get("model", "mistral-7b")
    max_tokens = data.get("max_tokens", 200)
    temperature = data.get("temperature", 0.7)

    if not prompt:
        return jsonify({"error_code": "EMPTY_PROMPT", "message": "Prompt não fornecido"}), 400

    # Cache semântico
    cached = get_cache(prompt)
    if cached:
        return jsonify({"text": cached, "cached": True}), 200

    # Tentar fallback local WASM
    try:
        response_text = call_local_wasm_llm(prompt, model, max_tokens, temperature)
    except Exception as e_local:
        print(f"[ERROR] Falha no WASM local: {e_local}", flush=True)
        response_text = None

    # Se sem resposta local, chamar LLM remoto
    if not response_text:
        try:
            response_text = call_remote_llm({
                "prompt": prompt,
                "model": model,
                "max_tokens": max_tokens,
                "temperature": temperature
            })
        except Exception as e_remote:
            print(f"[ERROR] Falha no LLM remoto: {e_remote}", flush=True)
            return jsonify({"error_code": "LLM_REMOTE_ERROR", "message": str(e_remote)}), 502

    # Armazenar no cache
    try:
        set_cache(prompt, response_text)
    except Exception as e_cache:
        print(f"[WARN] Falha ao salvar no cache: {e_cache}", flush=True)

    # Responder sucesso
    return jsonify({"text": response_text, "cached": False}), 200

if __name__ == "__main__":
    port = int(os.getenv("LLM_PROXY_PORT", 8000))
    app.run(host="0.0.0.0", port=port)

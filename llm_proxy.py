#!/usr/bin/env python3
"""
llm_proxy.py

Servidor mínimo que recebe requisições LLM e repassa ao OpenAI (ou outro serviço),
retornando a resposta em JSON. Pode ser usado como backend local para spans llm_request.
"""

import os
from flask import Flask, request, jsonify
import openai

# Inicialização do Flask
app = Flask(__name__)

# Carregar chave da OpenAI da variável de ambiente
openai.api_key = os.getenv("OPENAI_API_KEY")

@app.route("/api/llm", methods=["POST"])
def llm_api():
    data = request.json
    prompt = data.get("prompt", "")
    model = data.get("model", "gpt-4")
    max_tokens = data.get("max_tokens", 200)
    temperature = data.get("temperature", 0.7)

    if not prompt:
        return jsonify({"error": "Prompt não fornecido"}), 400

    try:
        # Chamada ao OpenAI
        response = openai.ChatCompletion.create(
            model=model,
            messages=[{"role": "user", "content": prompt}],
            max_tokens=max_tokens,
            temperature=temperature,
        )
        # Extrair texto da resposta
        text = response.choices[0].message.content
        return jsonify({"text": text}), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    port = int(os.getenv("LLM_PROXY_PORT", 8000))
    app.run(host="0.0.0.0", port=port)

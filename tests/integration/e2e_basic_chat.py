#!/usr/bin/env python3
# tests/integration/e2e_basic_chat.py
#
# Teste de integração de ponta a ponta para o FlipApp LogLineOS.
# Este script Python usa Playwright para automatizar um navegador headless,
# simulando a interação do usuário e verificando o comportamento do sistema
# do frontend ao backend LLM (via proxy).

import subprocess
import time
import requests
import json
import os
import signal
from playwright.sync_api import sync_playwright, expect

# --- Configurações do Teste ---
FLIPAPP_URL = "http://localhost:8080"
LLM_PROXY_URL = "http://localhost:8000"
HTTP_SERVER_PORT = 8080 # Porta do npx serve
LLM_PROXY_PORT = 8000   # Porta do llm_proxy_hard.py

# --- Funções Auxiliares para Gerenciar Processos ---
def start_process(command, name="Process"):
    """Inicia um processo em segundo plano e retorna o objeto Popen."""
    print(f"[{name}] Iniciando: {command}")
    process = subprocess.Popen(command, shell=True, preexec_fn=os.setsid)
    time.sleep(2) # Dar tempo para o processo iniciar
    return process

def stop_process(process, name="Process"):
    """Para um processo."""
    if process:
        print(f"[{name}] Encerrando processo...")
        # Envia SIGTERM para encerramento gracioso, depois SIGKILL se necessário
        os.killpg(os.getpgid(process.pid), signal.SIGTERM)
        time.sleep(1)
        if process.poll() is None: # Se o processo ainda estiver rodando
            os.killpg(os.getpgid(process.pid), signal.SIGKILL)
        print(f"[{name}] Processo encerrado.")

def check_llm_proxy_health():
    """Verifica se o proxy LLM está online."""
    try:
        response = requests.get(f"{LLM_PROXY_URL}/api/llm/health", timeout=5)
        response.raise_for_status()
        return response.json().get("status") == "ok"
    except (requests.exceptions.ConnectionError, requests.exceptions.Timeout, requests.exceptions.RequestException) as e:
        print(f"[LLM Proxy Health Check] Falha: {e}")
        return False

# --- Teste Principal ---
def run_e2e_test():
    llm_proxy_process = None
    http_server_process = None
    browser = None
    
    try:
        print("\n--- Iniciando Teste E2E do FlipApp LogLineOS ---")

        # 1. Iniciar o proxy LLM (se não estiver rodando)
        if not check_llm_proxy_health():
            print("[SETUP] Iniciando LLM Proxy...")
            llm_proxy_process = start_process(f"python3 llm_proxy_hard.py --port {LLM_PROXY_PORT}", "LLM Proxy")
            time.sleep(5) # Dar mais tempo para o Flask iniciar
            if not check_llm_proxy_health():
                raise Exception("LLM Proxy não iniciou corretamente!")
            print("[SETUP] LLM Proxy OK.")
        else:
            print("[SETUP] LLM Proxy já está rodando.")

        # 2. Iniciar o servidor HTTP mínimo para servir a UI LogLineOS
        print(f"[SETUP] Iniciando servidor HTTP em public/ na porta {HTTP_SERVER_PORT}...")
        http_server_process = start_process(f"npx serve public -p {HTTP_SERVER_PORT}", "HTTP Server")
        time.sleep(3) # Dar tempo para o servidor iniciar
        print("[SETUP] Servidor HTTP OK.")

        # 3. Iniciar o navegador e executar o teste
        with sync_playwright() as p:
            print("[PLAYWRIGHT] Iniciando navegador...")
            browser = p.chromium.launch()
            page = browser.new_page()

            print(f"[PLAYWRIGHT] Navegando para {FLIPAPP_URL}...")
            page.goto(FLIPAPP_URL)
            page.wait_for_selector("#loglineos-root", timeout=30000) # Espera a raiz da aplicação LogLineOS

            print("[PLAYWRIGHT] Página carregada. Verificando tela de carregamento...")
            expect(page.locator(".loading-container")).to_be_visible()
            # Espera até que o conteúdo real da UI apareça (ou o loading desapareça)
            page.wait_for_selector(".flipapp-root", timeout=30000) # Espera o container principal da UI
            expect(page.locator(".loading-container")).to_be_hidden()
            print("[PLAYWRIGHT] UI principal do FlipApp carregada.")

            # Teste de Interação: Enviar uma mensagem no chat
            print("[PLAYWRIGHT] Testando envio de mensagem no chat...")
            chat_input = page.locator(".chat-input input[type='text']")
            send_button = page.locator(".chat-input .btn-primary")

            expect(chat_input).to_be_visible()
            expect(send_button).to_be_visible()

            test_message = "Olá, FlipApp LogLineOS! Como você está hoje?"
            chat_input.fill(test_message)
            send_button.click()

            print(f"[PLAYWRIGHT] Mensagem enviada: '{test_message}'")
            # Esperar a mensagem do usuário aparecer
            user_message_bubble = page.locator(f".message-user .message-bubble:has-text('{test_message}')")
            expect(user_message_bubble).to_be_visible(timeout=10000)

            # Esperar a resposta do bot (do LLM) aparecer
            print("[PLAYWRIGHT] Aguardando resposta do bot...")
            bot_message_bubble = page.locator(".message-bot .message-bubble")
            expect(bot_message_bubble).to_be_visible(timeout=30000) # LLM pode demorar
            expect(bot_message_bubble).not_to_have_text("Loading...") # Garante que não é apenas placeholder
            print(f"[PLAYWRIGHT] Resposta do bot recebida: {bot_message_bubble.text_content()}")
            expect(bot_message_bubble.text_content()).not_to_be_empty()

            # Teste de Navegação: Ir para a aba Espelho
            print("[PLAYWRIGHT] Testando navegação para aba Espelho...")
            espelho_tab = page.locator(".nav-tab:has-text('Espelho')")
            espelho_tab.click()
            
            espelho_panel = page.locator(".espelho-panel")
            expect(espelho_panel).to_be_visible(timeout=10000)
            print("[PLAYWRIGHT] Aba Espelho carregada.")

            print("\n--- Teste E2E Concluído com SUCESSO! ---")

    except Exception as e:
        print(f"\n--- ERRO no Teste E2E ---")
        print(f"Detalhes do erro: {e}")
        # Captura screenshot em caso de falha
        if 'page' in locals() and page:
            page.screenshot(path="e2e_error_screenshot.png")
            print("Screenshot da falha salvo como e2e_error_screenshot.png")
        print("\n--- Teste E2E FAILED! ---")
        exit(1) # Sai com código de erro

    finally:
        # Garante que os processos são encerrados
        if llm_proxy_process:
            stop_process(llm_proxy_process, "LLM Proxy")
        if http_server_process:
            stop_process(http_server_process, "HTTP Server")
        if browser:
            browser.close()
            print("[PLAYWRIGHT] Navegador encerrado.")

if __name__ == "__main__":
    run_e2e_test()
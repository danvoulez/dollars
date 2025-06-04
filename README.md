```markdown
# FlipApp LogLineOS + WhatsApp Integration

![LogLineOS Logo](https://raw.githubusercontent.com/dally/logline-os/main/logo.svg)

## 🚀 Uma Nova Era de Software: LogLineOS com Módulo WhatsApp

Bem-vindo ao **FlipApp LogLineOS**, o projeto flagship que materializa a visão da linguagem LogLine com **integração completa do WhatsApp**. Este não é apenas um aplicativo; é um **manifesto operacional** onde toda a lógica, interface de usuário, estado, comportamento e até as regras de execução do sistema são descritas de forma **100% declarativa** como fluxos de spans em arquivos `.logline`.

**Este projeto não depende mais de Docker, Nginx, Vite, Jest ou Node.js para construir sua interface.** A única dependência de execução é o binário Go (`flip`) que interpreta spans LogLine, mais eventuais proxies LLM em Python ou WASM.

### ✨ Recursos do FlipApp LogLineOS + WhatsApp

- **Integração WhatsApp Web**: UI pixel-perfect, simulação de áudio, tradução automática, ghost‐sale detection, sistema de julgamento e treinamento contínuo.
- **Micro‐Kernel Declarativo**: Executor Go/WASM que lê `execution_rule` em LogLine.
- **Auditabilidade Total & Reversível**: Cada operação é um span persistido.
- **Offline‐First**: Design robusto sem necessidade de container ou servidor externo para UI.

---

## 🛠️ Pré-requisitos

- **Go 1.20+** para compilar o binário CLI `flip`
- **Python 3.8+** (opcional) para proxies LLM locais
- Nenhuma dependência de Node.js, Docker ou bundlers

## 🚀 Como Rodar Localmente

1. Clone o repositório  
   ```
   git clone https://github.com/danvoulez/dollars.git
   cd dollars
   ```

2. Compile o Executor  
   ```
   make build
   ```

3. Carregue variáveis de ambiente (se tiver `load_env.sh`)  
   ```
   source load_env.sh
   ```

4. Execute o FlipApp  
   ```
   ./bin/flip run --watch
   ```

5. (Opcional) Proxy LLM  
   ```
   python3 llm_proxy_hard.py --port 8001
   ```

---

## 🧪 Testes & Validações

- `make validate` – sintaxe de todos os spans `.logline` e `.jsonl`  
- `make test` – testes LogLine  
- Não há dependência de Jest, Vite ou Playwright

---

## 📚 Mais

- CHANGELOG.md  
- CONTRIBUTING.md  
- CHEATSHEET.md  
- LICENSE  
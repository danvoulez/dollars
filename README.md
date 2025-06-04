# FlipApp: O Flagship LogLineOS

![LogLineOS Logo](https://raw.githubusercontent.com/dally/logline-os/main/logo.svg) <!-- TODO: Criar um logo ou um placeholder -->

## 🚀 Uma Nova Era de Software: LogLineOS

Bem-vindo ao **FlipApp LogLineOS**, o projeto flagship que materializa a visão da linguagem LogLine. Este não é apenas um aplicativo; é um **manifesto operacional** onde toda a lógica, interface de usuário, estado, comportamento e até mesmo as regras de execução do sistema são descritas de forma **100% declarativa** como fluxos de spans em arquivos `.logline`.

**No FlipApp LogLineOS, não há JavaScript funcional tradicional, Docker, Nginx, Vite, Jest, ou HTML "oculto".** O sistema é executado por um **Micro-Kernel minimalista** (compilado para WASM no navegador, ou um binário Go para CLI) que entende o que fazer porque é instruído por **regras declarativas em LogLine**. Cada ação é um span, cada decisão é explícita, cada mutação é auditável e reversível.

### ✨ O Que É LogLineOS?

*   **Tudo é um Span**: Interface, lógica de negócio, configuração, testes, logs, e até as regras de execução do kernel são spans `.logline`.
*   **Auditabilidade Total**: Cada evento, cada decisão, cada falha é um span auditável e persistido.
*   **Reversibilidade**: O sistema mantém um log de mutações de estado que, em tese, permite reverter a qualquer ponto no tempo.
*   **Kernel Declarativo**: O Executor (a "VM") não possui lógica intrínseca para ações complexas; ele lê e interpreta `execution_rule`s em LogLine para saber como invocar suas ações nativas primitivas.
*   **Offline-First por Design**: Componentes e dados são projetados para funcionar de forma robusta mesmo sem conexão.

## 🌟 O Que o FlipApp LogLineOS Faz?

O FlipApp é uma plataforma conversacional que demonstra as capacidades do LogLineOS:

*   **Chat Inteligente**: Interaja com um LLM (simulado localmente ou via proxy) em uma experiência conversacional fluida.
*   **Painel Espelho**: Visualize um fluxo contínuo e auditável de eventos do sistema, com filtros e detalhamentos.
*   **Contratos Auditáveis**: Execute contratos de negócio (ex: pagamentos, ações de sistema) que são totalmente auditáveis e rastreáveis.
*   **Experiência de Usuário Premium**: Interface responsiva com animações, gestos e feedback háptico/sonoro.

---

## ⏳ Jornada do FlipApp: Ontem, Hoje, Amanhã

**Ontem (O Protótipo):**
O FlipApp começou como um protótipo JavaScript-centric. Tinha um chat básico, UI declarativa embrionária com `.logline`, e um runtime JS que fazia todo o trabalho pesado. Dependia de Webpack, Jest e um Docker complexo. Era funcional, mas opaco e difícil de auditar profundamente.

**Hoje (A Encarnação LogLineOS Pura):**
Chegamos ao ponto em que a maior parte do sistema está convertida para `.logline`.
*   **Kernel Declarativo**: O Executor é um binário minimalista que lê regras `.logline` para executar ações.
*   **Sem JS Funcional**: Classes JavaScript como `AnimationSystem`, `GestureSystem`, `ContractSystem` foram substituídas por `.logline`s.
*   **Infraestrutura Mínima**: Docker, Nginx, Vite, Jest são removidos. A compilação é para um único binário (`flip`) e arquivos `.logline`.
*   **Auditabilidade Intrínseca**: Cada operação é um span auditável.
*   **Desenvolvimento CLI-First**: A interação principal para o desenvolvedor é através da ferramenta `flip` (o Executor).

**Amanhã (Visão de Futuro):**
O FlipApp LogLineOS continuará a evoluir para:
*   **Executor totalmente auto-descritivo**: Onde o `Executor` não tem mais `switch/case`s internos e suas ações nativas são ainda mais primitivas e abstraídas.
*   **Ecossistema de Plugins LogLine**: Extensões de funcionalidades por meio de mais `.logline`s.
*   **Sincronização Distribuída de Spans**: Fluxos de spans compartilhados entre dispositivos ou nós.
*   **Interfaces Simbólicas Avançadas**: UIs com comportamento complexo descrito puramente em spans.
*   **Compilação para Hardware Dedicado**: LogLineOS rodando diretamente em hardware sem SO tradicional.

---

## 🛠️ Pré-requisitos (Para Desenvolver e Rodar)

Para desenvolver e rodar o FlipApp LogLineOS (o Executor `flip`), você precisará:

*   **Go (Golang)**: Para compilar o Executor (binário `flip`).
    *   [Instalar Go](https://go.dev/doc/install)
*   **Python 3**: Para rodar o `llm_proxy_hard.py` de exemplo.
    *   `pip install Flask requests`
*   **Opcional: Git**: Para clonar o repositório.

## 🚀 Como Instalar e Rodar Localmente

**1. Clone o Repositório:**

```bash
git clone https://github.com/sua-org/flipapp.git
cd flipapp
#!/bin/bash
# install_flip.sh
#
# Script de instalação rápido para o binário 'flip' do LogLineOS.
# Detecta o sistema operacional e arquitetura, baixa o binário correto
# do GitHub Releases e o instala em /usr/local/bin.

set -e # Sai imediatamente se um comando falhar

FLIP_VERSION="1.0.0-pure" # Versão do binário a ser instalada
GITHUB_REPO="sua-org/flipapp" # TODO: Substitua pelo seu repositório
INSTALL_DIR="/usr/local/bin" # Diretório padrão para executáveis

# --- Funções Auxiliares ---

# Função para exibir mensagens coloridas
color_echo() {
  local color=$1
  local text=$2
  case "$color" in
    "green") echo -e "\033[0;32m${text}\033[0m" ;;
    "red") echo -e "\033[0;31m${text}\033[0m" ;;
    "yellow") echo -e "\033[0;33m${text}\033[0m" ;;
    *) echo "${text}" ;;
  esac
}

# Verifica se um comando existe
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# --- Início da Instalação ---

color_echo "green" "🚀 Iniciando instalação do FlipApp LogLineOS (flip CLI) v${FLIP_VERSION}..."

# 1. Detectar SO e Arquitetura
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
  Linux)
    OS="linux"
    ;;
  Darwin)
    OS="darwin"
    ;;
  *)
    color_echo "red" "Sistema operacional ${OS} não suportado."
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  *)
    color_echo "red" "Arquitetura ${ARCH} não suportada."
    exit 1
    ;;
esac

FILENAME="flip-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/v${FLIP_VERSION}/${FILENAME}"

color_echo "yellow" "Detectado: OS=${OS}, ARCH=${ARCH}. Baixando de ${DOWNLOAD_URL}..."

# 2. Baixar o Binário
if ! command_exists curl && ! command_exists wget; then
  color_echo "red" "Erro: 'curl' ou 'wget' não encontrado. Por favor, instale um deles."
  exit 1
fi

DOWNLOAD_CMD=""
if command_exists curl; then
  DOWNLOAD_CMD="curl -L -o ${FILENAME}"
elif command_exists wget; then
  DOWNLOAD_CMD="wget -O ${FILENAME}"
fi

${DOWNLOAD_CMD} "${DOWNLOAD_URL}"

if [ ! -f "${FILENAME}" ]; then
  color_echo "red" "Erro: Falha ao baixar o binário. Verifique a URL ou sua conexão."
  exit 1
fi

color_echo "green" "Binário baixado com sucesso."

# 3. Mover para o Diretório de Instalação
if [ ! -w "${INSTALL_DIR}" ]; then
  color_echo "yellow" "Diretório de instalação ${INSTALL_DIR} não é gravável. Tentando com sudo..."
  sudo mv "${FILENAME}" "${INSTALL_DIR}/flip"
  sudo chmod +x "${INSTALL_DIR}/flip"
else
  mv "${FILENAME}" "${INSTALL_DIR}/flip"
  chmod +x "${INSTALL_DIR}/flip"
fi

if ! command_exists flip; then
  color_echo "red" "Erro: 'flip' não foi movido ou não está no PATH. Verifique ${INSTALL_DIR}."
  exit 1
fi

color_echo "green" "Binário 'flip' instalado em ${INSTALL_DIR}/flip."

# 4. Verificar Instalação
color_echo "green" "Verificando instalação..."
flip version

color_echo "green" "🎉 Instalação do FlipApp LogLineOS CLI concluída com sucesso!"
color_echo "green" "Agora você pode usar 'flip' em qualquer terminal."
color_echo "yellow" ""
color_echo "yellow" "Comandos úteis:"
color_echo "yellow" "  flip run          # Inicia o FlipApp"
color_echo "yellow" "  flip run --watch  # Modo desenvolvimento"
color_echo "yellow" "  flip logs --tail  # Ver logs em tempo real"
color_echo "yellow" "  flip test unit    # Executar testes"
color_echo "yellow" "  flip config view  # Ver configurações"
color_echo "yellow" ""
color_echo "green" "Para começar, execute: flip run"
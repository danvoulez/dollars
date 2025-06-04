#!/usr/bin/env bash
# flip_rebuild_pack.sh
# Script de transformação total para aplicar hardening no FlipApp LogLineOS

set -euo pipefail

# 1. Remover Docker e arquivos relacionados
echo "Removendo Dockerfile e configurações..."
rm -f Dockerfile docker-compose.yml
rm -rf docker/

# 2. Consolidar configurações em config.logline e gerar load_env.sh
echo "Consolidando configuração..."
CONFIG_FILE="runtime/config.logline"
LOAD_ENV="load_env.sh"
echo "#!/usr/bin/env bash" > $LOAD_ENV
echo "# Carrega variáveis de ambiente baseadas em config.logline" >> $LOAD_ENV
grep 'key:' $CONFIG_FILE | while IFS= read -r line; do
  KEY=$(echo $line | sed 's/.*key: "\(.*\)"/\1/')
  VALUE=$(grep -A1 "key: "$KEY"" $CONFIG_FILE | grep 'value:' | sed 's/.*value: "\(.*\)"/\1/')
  echo "export $KEY="$VALUE"" >> $LOAD_ENV
done
chmod +x $LOAD_ENV
rm -f env-config.sh

# 3. Atualizar README.md
echo "Atualizando README.md..."
cat << 'EOF' > README.md
# FlipApp LogLineOS - Hardening Total

## Como Executar
1. Carregue variáveis de ambiente:
   \`\`\`
   source load_env.sh
   \`\`\`
2. Inicie o proxy LLM reforçado:
   \`\`\`
   python3 llm_proxy_hard.py
   \`\`\`
3. Inicie o FlipApp:
   \`\`\`
   flip run
   \`\`\`
4. Para logs de auditoria:
   \`\`\`
   flip logs
   \`\`\`
5. Para rodar testes:
   \`\`\`
   flip test
   \`\`\`
EOF

# 4. Remover arquivos JS obsoletos e integrar Logger
echo "Removendo utils/logger.js e integrando audit.logline..."
rm -f public/src/utils/logger.js

# 5. Criar estrutura de arquivos de testes .logline
echo "Criando diretório de testes LogLine..."
mkdir -p test_logline
cat << 'EOF' > test_logline/expect_state.logline
# expect_state.logline
- type: test
  description: "Verifica estado pós-execução"
  expect_state:
    path: "state.test_key"
    equals: "expected_value"
EOF

cat << 'EOF' > test_logline/expect_effect.logline
# expect_effect.logline
- type: test
  description: "Verifica efeito de contrato"
  expect_effect:
    contract: "example_contract"
    applied: true
EOF

cat << 'EOF' > test_logline/expect_error.logline
# expect_error.logline
- type: test
  description: "Verifica erro em execução de span inválido"
  expect_error:
    span_type: "invalid_span"
    error_code: "SPAN_VALIDATION_FAILED"
EOF

echo "Hardening pack gerado com sucesso!"

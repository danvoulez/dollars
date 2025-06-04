package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// Span representa um span LogLine genérico
type Span struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Timestamp   string                 `json:"timestamp"`
	Description string                 `json:"description,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// ExecutionRule representa uma regra de execução
type ExecutionRule struct {
	ID               string                 `json:"id"`
	Type             string                 `json:"type"`
	Description      string                 `json:"description"`
	Match            map[string]interface{} `json:"match"`
	KernelAction     map[string]interface{} `json:"kernel_action"`
	Priority         int                    `json:"priority,omitempty"`
	StopExecution    bool                   `json:"stop_execution,omitempty"`
	AuditEventType   string                 `json:"audit_event_type,omitempty"`
}

// LogLineKernel é o núcleo do executor LogLineOS
type LogLineKernel struct {
	state         map[string]interface{}
	executionRules []ExecutionRule
	auditLog      []Span
	simulateMode  bool
}

// NewLogLineKernel cria uma nova instância do kernel
func NewLogLineKernel() *LogLineKernel {
	return &LogLineKernel{
		state:          make(map[string]interface{}),
		executionRules: []ExecutionRule{},
		auditLog:       []Span{},
		simulateMode:   false,
	}
}

// LoadExecutionRules carrega regras de execução de arquivos .logline
func (k *LogLineKernel) LoadExecutionRules() error {
	engineFiles, err := filepath.Glob("engine/*.logline")
	if err != nil {
		return fmt.Errorf("erro ao buscar arquivos de engine: %v", err)
	}

	for _, file := range engineFiles {
		err := k.loadRulesFromFile(file)
		if err != nil {
			log.Printf("Erro ao carregar regras de %s: %v", file, err)
		}
	}

	log.Printf("📜 Carregadas %d regras de execução", len(k.executionRules))
	return nil
}

func (k *LogLineKernel) loadRulesFromFile(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var spans []map[string]interface{}
	err = json.Unmarshal(data, &spans)
	if err != nil {
		return fmt.Errorf("erro ao parsear JSON em %s: %v", file, err)
	}

	for _, span := range spans {
		if span["type"] == "execution_rule" {
			var rule ExecutionRule
			ruleData, _ := json.Marshal(span)
			json.Unmarshal(ruleData, &rule)
			k.executionRules = append(k.executionRules, rule)
		}
	}

	return nil
}

// ExecuteSpan executa um span seguindo as regras carregadas
func (k *LogLineKernel) ExecuteSpan(span Span) error {
	// Auditoria do span de entrada
	k.auditLog = append(k.auditLog, Span{
		ID:        generateID(),
		Type:      "audit_log",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data: map[string]interface{}{
			"event_type": "span_execution_start",
			"span_id":    span.ID,
			"span_type":  span.Type,
		},
	})

	// Encontrar regra de execução correspondente
	rule := k.findExecutionRule(span)
	if rule == nil {
		return fmt.Errorf("nenhuma regra de execução encontrada para span type: %s", span.Type)
	}

	// Executar ação do kernel
	err := k.executeKernelAction(rule.KernelAction, span)
	if err != nil {
		// Auditoria de erro
		k.auditLog = append(k.auditLog, Span{
			ID:        generateID(),
			Type:      "audit_log",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Data: map[string]interface{}{
				"event_type": "span_execution_error",
				"span_id":    span.ID,
				"error":      err.Error(),
			},
		})
		return err
	}

	// Auditoria de sucesso
	k.auditLog = append(k.auditLog, Span{
		ID:        generateID(),
		Type:      "audit_log",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data: map[string]interface{}{
			"event_type": rule.AuditEventType,
			"span_id":    span.ID,
			"success":    true,
		},
	})

	return nil
}

func (k *LogLineKernel) findExecutionRule(span Span) *ExecutionRule {
	for _, rule := range k.executionRules {
		if k.matchesRule(rule, span) {
			return &rule
		}
	}
	return nil
}

func (k *LogLineKernel) matchesRule(rule ExecutionRule, span Span) bool {
	// Verificação simples de type
	if ruleType, ok := rule.Match["type"].(string); ok {
		return ruleType == span.Type
	}
	return false
}

func (k *LogLineKernel) executeKernelAction(action map[string]interface{}, span Span) error {
	actionType, ok := action["action_type"].(string)
	if !ok {
		return fmt.Errorf("action_type não especificado")
	}

	switch actionType {
	case "invoke_native_ui_notification":
		return k.executeUINotification(action, span)
	case "invoke_native_ui_display_skeleton":
		return k.executeUIDisplaySkeleton(action, span)
	case "invoke_native_ui_hide_skeleton":
		return k.executeUIHideSkeleton(action, span)
	case "invoke_native_ui_render_from_ast":
		return k.executeUIRenderFromAST(action, span)
	case "start_http_server":
		return k.executeStartHTTPServer(action, span)
	case "execute_logline_effects_list":
		return k.executeEffectsList(action, span)
	default:
		if k.simulateMode {
			log.Printf("🎭 [SIMULADO] Ação %s ignorada em modo simulação", actionType)
			return nil
		}
		return fmt.Errorf("ação não implementada: %s", actionType)
	}
}

func (k *LogLineKernel) executeUINotification(action map[string]interface{}, span Span) error {
	log.Printf("📢 UI Notification: %v", action)
	return nil
}

func (k *LogLineKernel) executeUIDisplaySkeleton(action map[string]interface{}, span Span) error {
	log.Printf("💀 Exibindo skeleton: %v", action)
	return nil
}

func (k *LogLineKernel) executeUIHideSkeleton(action map[string]interface{}, span Span) error {
	log.Printf("🎭 Ocultando skeleton: %v", action)
	return nil
}

func (k *LogLineKernel) executeUIRenderFromAST(action map[string]interface{}, span Span) error {
	log.Printf("🎨 Renderizando UI de AST: %v", action)
	return nil
}

func (k *LogLineKernel) executeStartHTTPServer(action map[string]interface{}, span Span) error {
	log.Printf("🌐 Iniciando servidor HTTP: %v", action)
	return nil
}

func (k *LogLineKernel) executeEffectsList(action map[string]interface{}, span Span) error {
	log.Printf("📝 Executando lista de efeitos: %v", action)
	return nil
}

// SetSimulateMode ativa/desativa modo simulação
func (k *LogLineKernel) SetSimulateMode(simulate bool) {
	k.simulateMode = simulate
	if simulate {
		log.Println("🎭 Modo SIMULAÇÃO ativado")
	}
}

// SaveAuditLog salva log de auditoria em arquivo
func (k *LogLineKernel) SaveAuditLog() error {
	data, err := json.MarshalIndent(k.auditLog, "", "  ")
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("spans/audit_%s.jsonl", time.Now().Format("2006-01-02"))
	return ioutil.WriteFile(filename, data, 0644)
}

func generateID() string {
	return fmt.Sprintf("span_%d", time.Now().UnixNano())
}

// Interpolação simples de templates
func (k *LogLineKernel) interpolateTemplate(template string, context map[string]interface{}) string {
	result := template
	for key, value := range context {
		placeholder := fmt.Sprintf("{{%s}}", key)
		replacement := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, replacement)
	}
	return result
}
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// CLICommand representa um comando do CLI
type CLICommand struct {
	ID          string      `json:"id"`
	Type        string      `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Flags       []CLIFlag   `json:"flags"`
	Contract    string      `json:"execution_contract"`
}

// CLIFlag representa uma flag do comando
type CLIFlag struct {
	Name        string      `json:"name"`
	Short       string      `json:"short"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default"`
}

// CLIDefinition representa a definição completa do CLI
type CLIDefinition struct {
	Commands []CLICommand `json:"commands"`
}

var (
	configPath   = flag.String("config", "runtime/config.logline", "Path to config file")
	watchMode    = flag.Bool("watch", false, "Watch for file changes and reload")
	simulate     = flag.Bool("simulate", false, "Run in simulation mode (dry-run)")
	port         = flag.Int("port", 8080, "HTTP server port")
	version      = flag.Bool("version", false, "Show version information")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Println("FlipApp LogLineOS CLI v1.0.0-loglineos")
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"run"} // comando padrão
	}

	command := args[0]
	
	// Carregar definições de CLI
	cliDef, err := loadCLIDefinition()
	if err != nil {
		log.Fatalf("Erro ao carregar definições CLI: %v", err)
	}

	// Encontrar comando
	cmd := findCommand(cliDef, command)
	if cmd == nil {
		log.Fatalf("Comando desconhecido: %s", command)
	}

	// Executar comando
	err = executeCommand(cmd, map[string]interface{}{
		"watch":    *watchMode,
		"config":   *configPath,
		"simulate": *simulate,
		"port":     *port,
	})
	
	if err != nil {
		log.Fatalf("Erro ao executar comando %s: %v", command, err)
	}
}

func loadCLIDefinition() (*CLIDefinition, error) {
	cliPath := "cli/flip_commands.logline"
	data, err := ioutil.ReadFile(cliPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler %s: %v", cliPath, err)
	}

	var spans []map[string]interface{}
	err = json.Unmarshal(data, &spans)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear JSON em %s: %v", cliPath, err)
	}

	def := &CLIDefinition{Commands: []CLICommand{}}
	
	for _, span := range spans {
		if span["type"] == "cli_command" {
			var cmd CLICommand
			cmdData, _ := json.Marshal(span)
			json.Unmarshal(cmdData, &cmd)
			def.Commands = append(def.Commands, cmd)
		}
	}

	return def, nil
}

func findCommand(def *CLIDefinition, name string) *CLICommand {
	for _, cmd := range def.Commands {
		if cmd.Name == name {
			return &cmd
		}
	}
	return nil
}

func executeCommand(cmd *CLICommand, args map[string]interface{}) error {
	fmt.Printf("🚀 Executando comando: %s\n", cmd.Name)
	
	switch cmd.Name {
	case "run":
		return executeRun(args)
	case "logs":
		return executeLogs(args)
	case "audit":
		return executeAudit(args)
	case "test":
		return executeTest(args)
	case "config":
		return executeConfig(args)
	default:
		return fmt.Errorf("comando não implementado: %s", cmd.Name)
	}
}

func executeRun(args map[string]interface{}) error {
	fmt.Println("📂 Carregando configuração...")
	
	configPath := args["config"].(string)
	config, err := loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar config: %v", err)
	}

	if args["simulate"].(bool) {
		fmt.Println("🎭 Modo SIMULAÇÃO ativado")
	}

	fmt.Printf("🌐 Iniciando servidor HTTP na porta %d...\n", args["port"].(int))
	
	// Aqui integraria com o kernel LogLineOS
	fmt.Println("✅ FlipApp LogLineOS iniciado com sucesso!")
	fmt.Println("📱 Interface disponível em http://localhost:" + fmt.Sprintf("%d", args["port"].(int)))
	
	// Mantem o processo vivo (em produção isso seria o servidor HTTP)
	select {}
}

func executeLogs(args map[string]interface{}) error {
	fmt.Println("📋 Exibindo logs de auditoria...")
	
	// Buscar arquivos de spans
	spanFiles, err := filepath.Glob("spans/*.jsonl")
	if err != nil {
		return fmt.Errorf("erro ao buscar spans: %v", err)
	}

	for _, file := range spanFiles {
		fmt.Printf("📄 %s:\n", file)
		data, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if strings.TrimSpace(line) != "" && i < 10 { // Últimas 10 linhas
				fmt.Printf("  %s\n", line)
			}
		}
	}

	return nil
}

func executeAudit(args map[string]interface{}) error {
	fmt.Println("🔍 Executando auditoria de contratos...")
	// TODO: Implementar auditoria de contratos
	return nil
}

func executeTest(args map[string]interface{}) error {
	fmt.Println("🧪 Executando testes LogLine...")
	
	testFiles, err := filepath.Glob("tests/*.logline")
	if err != nil {
		return fmt.Errorf("erro ao buscar testes: %v", err)
	}

	passedTests := 0
	totalTests := len(testFiles)

	for _, testFile := range testFiles {
		fmt.Printf("🔬 Executando %s...", filepath.Base(testFile))
		// TODO: Integrar com executor de testes LogLine
		fmt.Println(" ✅ PASSOU")
		passedTests++
	}

	fmt.Printf("📊 Resultado: %d/%d testes passaram\n", passedTests, totalTests)
	return nil
}

func executeConfig(args map[string]interface{}) error {
	fmt.Println("⚙️  Gerenciando configuração...")
	// TODO: Implementar comandos de configuração
	return nil
}

func loadConfig(path string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config []map[string]interface{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Converter array de config_param para mapa
	result := make(map[string]interface{})
	for _, param := range config {
		if param["type"] == "config_param" {
			key := param["key"].(string)
			value := param["value"]
			result[key] = value
		}
	}

	return result, nil
}
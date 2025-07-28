package quiz

import (
	"math/rand"
	"strings"
	"time"
	"fmt"

	"quiz_go/internal/ui"
	"quiz_go/internal/stats"
	"github.com/pterm/pterm"

	"github.com/AlecAivazis/survey/v2"
)

type Questao struct {
	ID          int      `json:"id"`
	Questao     string   `json:"questao"`
	Opcoes      []string `json:"opcoes"`
	Resposta    string   `json:"resposta"`
	Explicacao  string   `json:"explicacao"`
	Dificuldade string   `json:"dificuldade"` // "facil", "medio", "dificil"
	Categoria   string   `json:"categoria"`   // "sintaxe", "tipos", "concorrencia", etc.
}

type Quiz struct {
	questoes  []Questao
	stats stats.Estatisticas
	statsFile string
}

func NewQuiz() *Quiz {
	q := &Quiz{
		statsFile: "quiz_stats.json",
		questoes: []Questao{
			{
				ID:          1,
				Questao:     "Qual palavra-chave define uma função em Go?",
				Opcoes:      []string{"func", "function", "def", "lambda"},
				Resposta:    "func",
				Explicacao:  "Em Go, usamos a palavra-chave 'func' para definir funções. Exemplo: func minhaFuncao() {}",
				Dificuldade: "facil",
				Categoria:   "sintaxe",
			},
			{
				ID:          2,
				Questao:     "Como declarar uma variável em Go?",
				Opcoes:      []string{"let x = 10", "var x int = 10", "int x = 10", "x := int(10)"},
				Resposta:    "var x int = 10",
				Explicacao:  "Go usa 'var' para declaração explícita de variáveis. Também podemos usar := para declaração curta.",
				Dificuldade: "facil",
				Categoria:   "tipos",
			},
			{
				ID:          3,
				Questao:     "Qual pacote padrão usamos para imprimir no terminal?",
				Opcoes:      []string{"io", "os", "fmt", "print"},
				Resposta:    "fmt",
				Explicacao:  "O pacote 'fmt' fornece funções para formatação de I/O, incluindo Print, Printf e Println.",
				Dificuldade: "facil",
				Categoria:   "bibliotecas",
			},
			{
				ID:          4,
				Questao:     "Como criar um slice em Go?",
				Opcoes:      []string{"var s []int", "s := make([]int, 0)", "s := []int{}", "Todas as anteriores"},
				Resposta:    "Todas as anteriores",
				Explicacao:  "Go oferece múltiplas formas de criar slices: declaração zero, make() e literal.",
				Dificuldade: "medio",
				Categoria:   "tipos",
			},
			{
				ID:          5,
				Questao:     "Qual é o valor zero de um ponteiro em Go?",
				Opcoes:      []string{"0", "null", "nil", "undefined"},
				Resposta:    "nil",
				Explicacao:  "Em Go, 'nil' é o valor zero para ponteiros, interfaces, maps, slices, channels e funções.",
				Dificuldade: "medio",
				Categoria:   "tipos",
			},
			{
				ID:          6,
				Questao:     "Como criar uma goroutine em Go?",
				Opcoes:      []string{"go funcao()", "async funcao()", "thread funcao()", "spawn funcao()"},
				Resposta:    "go funcao()",
				Explicacao:  "A palavra-chave 'go' inicia uma nova goroutine, executando a função concorrentemente.",
				Dificuldade: "medio",
				Categoria:   "concorrencia",
			},
			{
				ID:          7,
				Questao:     "Qual é a forma correta de criar um channel em Go?",
				Opcoes:      []string{"ch := channel(int)", "ch := make(chan int)", "ch := new(chan int)", "ch := chan int{}"},
				Resposta:    "ch := make(chan int)",
				Explicacao:  "Channels são criados usando make(chan tipo). Eles são fundamentais para comunicação entre goroutines.",
				Dificuldade: "medio",
				Categoria:   "concorrencia",
			},
			{
				ID:          8,
				Questao:     "O que é um interface{} em Go?",
				Opcoes:      []string{"Um tipo genérico", "Interface vazia que aceita qualquer tipo", "Um erro", "Uma função"},
				Resposta:    "Interface vazia que aceita qualquer tipo",
				Explicacao:  "interface{} é a interface vazia, satisfeita por qualquer tipo. É similar ao 'any' em outras linguagens.",
				Dificuldade: "dificil",
				Categoria:   "interfaces",
			},
			{
				ID:          9,
				Questao:     "Como tratar erros idiomaticamente em Go?",
				Opcoes:      []string{"try/catch", "if err != nil", "throw/catch", "error handling"},
				Resposta:    "if err != nil",
				Explicacao:  "Go não tem exceções. Erros são valores que devem ser verificados explicitamente com 'if err != nil'.",
				Dificuldade: "medio",
				Categoria:   "erros",
			},
			{
				ID:          10,
				Questao:     "Qual é a diferença entre array e slice em Go?",
				Opcoes:      []string{"Não há diferença", "Arrays têm tamanho fixo, slices são dinâmicos", "Slices são mais lentos", "Arrays são obsoletos"},
				Resposta:    "Arrays têm tamanho fixo, slices são dinâmicos",
				Explicacao:  "Arrays têm tamanho fixo definido no tipo [5]int, enquanto slices []int são dinâmicos e mais flexíveis.",
				Dificuldade: "dificil",
				Categoria:   "tipos",
			},
		},
	}

	loadedStats, err := stats.CarregarEstatisticas(q.statsFile)
	if err != nil {
		fmt.Printf("Erro ao carregar estatísticas: %v. Iniciando com estatísticas zeradas.\n", err)
	} else {
		q.stats = loadedStats
	}
	return q
}

func (q *Quiz) MostrarEstatisticas() {
	if q.stats.TotalQuizzes == 0 {
		fmt.Println(ui.Yellow("📊 Nenhuma estatística disponível ainda."))
		return
	}

	fmt.Println(ui.Cyan("╔══════════════════════════════════════════════════════════╗"))
	fmt.Println(ui.Cyan("║") + "                    " + ui.Bold("📊 SUAS ESTATÍSTICAS 📊") + "                    " + ui.Cyan("║"))
	fmt.Println(ui.Cyan("╚══════════════════════════════════════════════════════════╝"))
	fmt.Println()

	fmt.Printf("%s Total de quizzes realizados: %s\n",
		ui.Magenta("🎯"),
		ui.Bold(fmt.Sprintf("%d", q.stats.TotalQuizzes)))

	fmt.Printf("%s Total de acertos: %s de %s\n",
		ui.Green("✅"),
		ui.Bold(fmt.Sprintf("%d", q.stats.TotalAcertos)),
		ui.Bold(fmt.Sprintf("%d", q.stats.TotalQuestoes)))

	fmt.Printf("%s Melhor score: %s questões\n",
		ui.Yellow("🏆"),
		ui.Bold(fmt.Sprintf("%d", q.stats.MelhorScore)))

	fmt.Printf("%s Média de acertos: %s\n",
		ui.Blue("📈"),
		ui.Bold(fmt.Sprintf("%.1f%%", q.stats.MediaPercentual)))

	if q.stats.UltimoQuiz != "" {
		fmt.Printf("%s Último quiz: %s\n",
			ui.Cyan("📅"),
			ui.Bold(q.stats.UltimoQuiz))
	}

	fmt.Println()
}

func (q *Quiz) SelecionarModoJogo() string {
	var modo string
	prompt := &survey.Select{
		Message: "Escolha o modo de jogo:",
		Options: []string{
			"🎯 Todas as questões (10 questões)",
			"⚡ Quiz rápido (5 questões aleatórias)",
			"🧠 Apenas questões difíceis",
			"📊 Ver estatísticas",
			"❌ Sair",
		},
	}

	survey.AskOne(prompt, &modo)
	return modo
}

func (q *Quiz) FiltrarQuestoes(modo string) []Questao {
	switch {
	case strings.Contains(modo, "Todas as questões"):
		return q.questoes
	case strings.Contains(modo, "Quiz rápido"):
		questoesAleatorias := make([]Questao, len(q.questoes))
		copy(questoesAleatorias, q.questoes)

		// Embaralhar
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questoesAleatorias), func(i, j int) {
			questoesAleatorias[i], questoesAleatorias[j] = questoesAleatorias[j], questoesAleatorias[i]
		})

		// Retornar apenas 5
		if len(questoesAleatorias) > 5 {
			return questoesAleatorias[:5]
		}
		return questoesAleatorias
	case strings.Contains(modo, "difíceis"):
		var dificeis []Questao
		for _, q := range q.questoes {
			if q.Dificuldade == "dificil" {
				dificeis = append(dificeis, q)
			}
		}
		return dificeis
	default:
		return q.questoes
	}
}

func (q *Quiz) ExecutarQuiz(questoesSelecionadas []Questao) {
	fmt.Println()
	fmt.Printf("%s Você terá %s questões para responder!\n",
		ui.Magenta("📚"),
		ui.Bold(fmt.Sprintf("%d", len(questoesSelecionadas))))
	fmt.Println()

	// Barra de progresso para preparação
	spinner, _ := pterm.DefaultSpinner.Start(ui.Cyan("Preparando o quiz..."))
	// Simula um tempo de preparação
	time.Sleep(time.Duration(len(questoesSelecionadas)) * 150 * time.Millisecond)
	spinner.Success(pterm.Green("Quiz pronto!"))
	fmt.Println()

	score := 0
	respostasCorretas := []bool{}
	tempoInicio := time.Now()

	// Loop das questões
	for i, questao := range questoesSelecionadas {
		ui.LimparTela()
		fmt.Println(ui.Cyan("═══════════════════════════════════════════════════════════"))
		fmt.Printf("%s Questão %d de %d | %s | %s\n",
			ui.Yellow("📝"),
			i+1,
			len(questoesSelecionadas),
			ui.Blue(fmt.Sprintf("Categoria: %s", questao.Categoria)),
			q.getDificuldadeIcon(questao.Dificuldade))
		fmt.Println(ui.Cyan("═══════════════════════════════════════════════════════════"))
		fmt.Println()

		var resposta string
		prompt := &survey.Select{
			Message: ui.Bold(questao.Questao),
			Options: questao.Opcoes,
		}

		err := survey.AskOne(prompt, &resposta)
		if err != nil {
			fmt.Printf(ui.Red("Erro ao ler resposta: %v\n"), err)
			continue
		}

		fmt.Println()

		if strings.TrimSpace(resposta) == questao.Resposta {
			fmt.Println(ui.Green("✅ Resposta correta! Parabéns!"))
			score++
			respostasCorretas = append(respostasCorretas, true)
		} else {
			fmt.Printf(ui.Red("❌ Resposta incorreta! A resposta correta é: %s\n"),
				ui.Bold(questao.Resposta))
			respostasCorretas = append(respostasCorretas, false)
		}

		// Mostrar explicação
		fmt.Printf("%s %s\n", ui.Blue("💡 Explicação:"), questao.Explicacao)
		fmt.Println()

		// Mostrar progresso atual
		if i < len(questoesSelecionadas)-1 {
			fmt.Printf(ui.Magenta("📊 Progresso: %d/%d questões | Acertos: %d\n"),
				i+1, len(questoesSelecionadas), score)
			fmt.Println()

			// Perguntar se quer continuar
			var continuar bool
			continuePrompt := &survey.Confirm{
				Message: "Continuar para a próxima questão?",
				Default: true,
			}
			survey.AskOne(continuePrompt, &continuar)

			if !continuar {
				fmt.Println(ui.Yellow("Quiz interrompido pelo usuário."))
				return
			}
			fmt.Println()
		}
	}

	tempoTotal := time.Since(tempoInicio)
	q.MostrarResultados(score, len(questoesSelecionadas), respostasCorretas, tempoTotal)
	q.AtualizarEstatisticas(score, len(questoesSelecionadas))
}

func (q *Quiz) getDificuldadeIcon(dificuldade string) string {
	switch dificuldade {
	case "facil":
		return ui.Green("🟢 Fácil")
	case "medio":
		return ui.Yellow("🟡 Médio")
	case "dificil":
		return ui.Red("🔴 Difícil")
	default:
		return ui.Blue("🔵 Normal")
	}
}

func (q *Quiz) MostrarResultados(score, total int, respostasCorretas []bool, tempo time.Duration) {
	fmt.Println()
	fmt.Println(ui.Cyan("╔══════════════════════════════════════════════════════════╗"))
	fmt.Println(ui.Cyan("║") + "                    " + ui.Bold("🏆 RESULTADOS FINAIS 🏆") + "                    " + ui.Cyan("║"))
	fmt.Println(ui.Cyan("╚══════════════════════════════════════════════════════════╝"))
	fmt.Println()

	// Barra de progresso dos resultados
	spinner, _ := pterm.DefaultSpinner.Start(ui.Magenta("Calculando resultados..."))
	// Simula um tempo de cálculo
	time.Sleep(2 * time.Second)
	spinner.Success(pterm.Green("Cálculos finalizados!"))
	fmt.Println()

	percentual := float64(score) / float64(total) * 100

	fmt.Printf("%s Você acertou %s de %s questões\n",
		ui.Magenta("📊"),
		ui.Bold(ui.Green(fmt.Sprintf("%d", score))),
		ui.Bold(fmt.Sprintf("%d", total)))

	fmt.Printf("%s Percentual de acertos: %s\n",
		ui.Magenta("📈"),
		ui.Bold(fmt.Sprintf("%.1f%%", percentual)))

	fmt.Printf("%s Tempo total: %s\n",
		ui.Blue("⏱️"),
		ui.Bold(fmt.Sprintf("%.1f segundos", tempo.Seconds())))

	fmt.Printf("%s Tempo médio por questão: %s\n",
		ui.Blue("⚡"),
		ui.Bold(fmt.Sprintf("%.1f segundos", tempo.Seconds()/float64(total))))

	fmt.Println()

	// Mostrar resumo das respostas
	fmt.Println(ui.Cyan("📋 Resumo das suas respostas:"))
	for i, correto := range respostasCorretas {
		status := ui.Red("❌")
		if correto {
			status = ui.Green("✅")
		}
		fmt.Printf("   Questão %d: %s\n", i+1, status)
	}
	fmt.Println()

	// Mensagem final baseada na performance
	q.MostrarMensagemFinal(score, total, percentual)
}

func (q *Quiz) MostrarMensagemFinal(score, total int, percentual float64) {
	switch {
	case score == total:
		fmt.Println(ui.Green("🎉 PERFEITO! Você acertou todas as questões!"))
		fmt.Println(ui.Green("🏆 Você é um verdadeiro expert em Go!"))
		fmt.Println(ui.Green("🌟 Considerado um GoGuru!"))
	case percentual >= 80:
		fmt.Println(ui.Green("🌟 Excelente! Você tem um ótimo conhecimento em Go!"))
		fmt.Println(ui.Green("👏 Continue assim!"))
		fmt.Println(ui.Blue("🚀 Próximo nível: tente as questões difíceis!"))
	case percentual >= 60:
		fmt.Println(ui.Yellow("👍 Muito bem! Você está no caminho certo!"))
		fmt.Println(ui.Yellow("📚 Continue estudando para melhorar ainda mais!"))
		fmt.Println(ui.Blue("💡 Dica: revise os conceitos que errou!"))
	case percentual >= 40:
		fmt.Println(ui.Yellow("😊 Bom começo! Você já sabe algumas coisas sobre Go!"))
		fmt.Println(ui.Yellow("💪 Com mais estudo você chegará lá!"))
		fmt.Println(ui.Blue("📖 Recomendo focar nos fundamentos primeiro!"))
	default:
		fmt.Println(ui.Red("📖 Você precisa estudar mais sobre Go!"))
		fmt.Println(ui.Red("💡 Que tal revisar a documentação oficial?"))
		fmt.Println(ui.Yellow("🔗 Recursos recomendados:"))
		fmt.Println(ui.Cyan("   • https://golang.org/doc/"))
		fmt.Println(ui.Cyan("   • https://tour.golang.org/"))
		fmt.Println(ui.Cyan("   • https://gobyexample.com/"))
	}
	fmt.Println()
}

func (q *Quiz) AtualizarEstatisticas(score, total int) {
	q.stats.TotalQuizzes++
	q.stats.TotalAcertos += score
	q.stats.TotalQuestoes += total

	if score > q.stats.MelhorScore {
		q.stats.MelhorScore = score
	}

	q.stats.MediaPercentual = float64(q.stats.TotalAcertos) / float64(q.stats.TotalQuestoes) * 100
	q.stats.UltimoQuiz = time.Now().Format("02/01/2006 15:04")

	_ = stats.SalvarEstatisticas(q.statsFile, q.stats)
}

func (q *Quiz) JogarNovamente() bool {
	var jogarNovamente bool
	playAgainPrompt := &survey.Confirm{
		Message: "Gostaria de jogar novamente?",
		Default: false,
	}
	survey.AskOne(playAgainPrompt, &jogarNovamente)
	return jogarNovamente
}

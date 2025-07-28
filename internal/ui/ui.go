package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

var (
	Cyan    = color.New(color.FgCyan).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Magenta = color.New(color.FgMagenta).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
)

func MostrarTelaInicial() {
	fmt.Println()
	fmt.Println(Cyan("╔══════════════════════════════════════════════════════════╗"))
	fmt.Println(Cyan("║") + "                                                          " + Cyan("║"))
	fmt.Println(Cyan("║") + "                " + Bold("🐹 QUIZ INTERATIVO DE GO 🐹") + "               " + Cyan("║"))
	fmt.Println(Cyan("║") + "                                                          " + Cyan("║"))
	fmt.Println(Cyan("║") + "          Teste seus conhecimentos sobre Golang!          " + Cyan("║"))
	fmt.Println(Cyan("║") + "                                                          " + Cyan("║"))
	fmt.Println(Cyan("╚══════════════════════════════════════════════════════════╝"))
	fmt.Println()
}

func MostrarDespedida() {
	fmt.Println()
	fmt.Println(Cyan("┌─────────────────── 👋 ATÉ LOGO! ────────────────────┐"))
	fmt.Println(Cyan("|                                                     |"))
	fmt.Println(Cyan("|    Obrigado por testar seus conhecimentos em Go!    |"))
	fmt.Println(Cyan("|                                                     |"))
	fmt.Println(Cyan("|    🐹 Continue aprendendo e praticando!             |"))
	fmt.Println(Cyan("|    🚀 Go é uma linguagem incrível!                  |"))
	fmt.Println(Cyan("|                                                     |"))
	fmt.Println(Cyan("|    Recursos úteis:                                  |"))
	fmt.Println(Cyan("|    • https://golang.org/doc/                        |"))
	fmt.Println(Cyan("|    • https://tour.golang.org/                       |"))
	fmt.Println(Cyan("|    • https://gobyexample.com/                       |"))
	fmt.Println(Cyan("|    • https://pkg.go.dev/                            |"))
	fmt.Println(Cyan("|                                                     |"))
	fmt.Println(Cyan("└─────────────────────────────────────────────────────┘"))
	fmt.Println()
}


// LimparTela limpa a tela do console de forma compatível com Windows, Linux e macOS.
func LimparTela() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
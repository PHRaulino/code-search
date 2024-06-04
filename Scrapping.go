package main

import (
    "context"
    "fmt"
    "time"

    "github.com/go-rod/rod"
    "github.com/go-rod/rod/lib/input"
    "github.com/go-rod/rod/lib/launcher"
    "github.com/go-rod/rod/lib/proto"
)

func main() {
    // Inicie o navegador com modo headless desativado para depuração
    url := launcher.New().Headless(false).MustLaunch()
    browser := rod.New().ControlURL(url).MustConnect()
    defer browser.MustClose()

    // Crie um novo contexto com timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Navegue até a página desejada
    page := browser.MustPage("http://example.com").MustWaitLoad()

    // Acesse o menu (exemplo de interação com um elemento)
    err := page.MustElement("css selector do menu").MustClick().WaitVisible()
    if err != nil {
        fmt.Println("Erro ao clicar no menu:", err)
        return
    }

    // Aguarde o carregamento do conteúdo do menu
    page.MustWaitLoad()

    // Tirar uma screenshot da página inteira
    screenshot, err := page.Screenshot(true, &proto.PageCaptureScreenshot{})
    if err != nil {
        fmt.Println("Erro ao tirar screenshot:", err)
        return
    }

    // Salvar a screenshot em um arquivo
    err = rod.OutputFile("screenshot.png", screenshot)
    if err != nil {
        fmt.Println("Erro ao salvar screenshot:", err)
        return
    }

    fmt.Println("Screenshot salva com sucesso")
}

package main

import (
    "bufio"
    "fmt"
    "log"
    "os/exec"
    "strings"
)

func main() {
    files, err := getAllNonIgnoredFilesOptimized()
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        fmt.Println(file)
    }
}

func getAllNonIgnoredFilesOptimized() ([]string, error) {
    // Usa git ls-files com flags para pegar todos os arquivos relevantes de uma vez
    // --cached: arquivos no índice
    // --others: arquivos não rastreados
    // --exclude-standard: aplica .gitignore, .git/info/exclude, etc.
    cmd := exec.Command("git", "ls-files", "--cached", "--others", "--exclude-standard")
    
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, err
    }

    if err := cmd.Start(); err != nil {
        return nil, err
    }

    var files []string
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        file := strings.TrimSpace(scanner.Text())
        if file != "" {
            files = append(files, file)
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return files, cmd.Wait()
}

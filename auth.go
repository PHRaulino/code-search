// cmd/auth.go
package cmd

import (
“fmt”
“os”

```
"github.com/charmbracelet/bubbles/textinput"
tea "github.com/charmbracelet/bubbletea"
"github.com/charmbracelet/lipgloss"
"github.com/spf13/cobra"
```

)

var setupCmd = &cobra.Command{
Use:   “setup”,
Short: “Setup authentication credentials interactively”,
RunE:  runSetup,
}

func init() {
authCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) error {
p := tea.NewProgram(newCredentialsModel())
if _, err := p.Run(); err != nil {
return fmt.Errorf(“failed to run TUI: %w”, err)
}
return nil
}

// tui/credentials.go
package tui

import (
“fmt”
“strings”

```
"github.com/charmbracelet/bubbles/textinput"
tea "github.com/charmbracelet/bubbletea"
"github.com/charmbracelet/lipgloss"
```

)

type credentialType int

const (
stackspotUser credentialType = iota
stackspotService
hashicorpVault
)

type step int

const (
stepSelectType step = iota
stepEnterCredentials
stepConfirm
stepComplete
)

type credentialsModel struct {
step           step
selectedType   credentialType
cursor         int
inputs         []textinput.Model
currentInput   int
savedSuccessfully bool
errorMsg       string
}

var (
titleStyle = lipgloss.NewStyle().
Bold(true).
Foreground(lipgloss.Color(”#7C3AED”)).
Padding(1, 0)

```
optionStyle = lipgloss.NewStyle().
	Padding(0, 2)

selectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#7C3AED")).
	Bold(true).
	Padding(0, 2).
	Background(lipgloss.Color("#E0E7FF"))

inputStyle = lipgloss.NewStyle().
	Padding(0, 1)

successStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#10B981")).
	Bold(true)

errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#EF4444")).
	Bold(true)
```

)

func newCredentialsModel() credentialsModel {
return credentialsModel{
step:   stepSelectType,
cursor: 0,
}
}

func (m credentialsModel) Init() tea.Cmd {
return nil
}

func (m credentialsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
switch msg := msg.(type) {
case tea.KeyMsg:
switch m.step {
case stepSelectType:
return m.updateSelectType(msg)
case stepEnterCredentials:
return m.updateEnterCredentials(msg)
case stepConfirm:
return m.updateConfirm(msg)
case stepComplete:
return m, tea.Quit
}
}
return m, nil
}

func (m credentialsModel) updateSelectType(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
switch msg.String() {
case “up”, “k”:
if m.cursor > 0 {
m.cursor–
}
case “down”, “j”:
if m.cursor < 2 {
m.cursor++
}
case “enter”:
m.selectedType = credentialType(m.cursor)
m.step = stepEnterCredentials
m.inputs = m.createInputs()
return m, m.inputs[0].Focus()
case “q”, “ctrl+c”:
return m, tea.Quit
}
return m, nil
}

func (m credentialsModel) updateEnterCredentials(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
switch msg.String() {
case “tab”, “down”:
if m.currentInput < len(m.inputs)-1 {
m.inputs[m.currentInput].Blur()
m.currentInput++
return m, m.inputs[m.currentInput].Focus()
}
case “shift+tab”, “up”:
if m.currentInput > 0 {
m.inputs[m.currentInput].Blur()
m.currentInput–
return m, m.inputs[m.currentInput].Focus()
}
case “enter”:
if m.currentInput == len(m.inputs)-1 {
// Todos os campos preenchidos
m.step = stepConfirm
} else {
m.inputs[m.currentInput].Blur()
m.currentInput++
return m, m.inputs[m.currentInput].Focus()
}
case “esc”:
m.step = stepSelectType
m.currentInput = 0
case “ctrl+c”:
return m, tea.Quit
}

```
var cmd tea.Cmd
m.inputs[m.currentInput], cmd = m.inputs[m.currentInput].Update(msg)
return m, cmd
```

}

func (m credentialsModel) updateConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
switch msg.String() {
case “y”, “Y”, “enter”:
err := m.saveCredentials()
if err != nil {
m.errorMsg = err.Error()
} else {
m.savedSuccessfully = true
}
m.step = stepComplete
case “n”, “N”, “esc”:
m.step = stepEnterCredentials
return m, m.inputs[m.currentInput].Focus()
case “ctrl+c”:
return m, tea.Quit
}
return m, nil
}

func (m credentialsModel) createInputs() []textinput.Model {
var inputs []textinput.Model

```
switch m.selectedType {
case stackspotUser:
	clientID := textinput.New()
	clientID.Placeholder = "Enter Client ID"
	clientID.CharLimit = 100

	clientSecret := textinput.New()
	clientSecret.Placeholder = "Enter Client Secret"
	clientSecret.EchoMode = textinput.EchoPassword
	clientSecret.CharLimit = 200

	inputs = []textinput.Model{clientID, clientSecret}

case stackspotService:
	clientID := textinput.New()
	clientID.Placeholder = "Enter Client ID"
	clientID.CharLimit = 100

	clientSecret := textinput.New()
	clientSecret.Placeholder = "Enter Client Secret"
	clientSecret.EchoMode = textinput.EchoPassword
	clientSecret.CharLimit = 200

	inputs = []textinput.Model{clientID, clientSecret}

case hashicorpVault:
	vaultURL := textinput.New()
	vaultURL.Placeholder = "Enter Vault URL (https://vault.example.com)"
	vaultURL.CharLimit = 200

	vaultPath := textinput.New()
	vaultPath.Placeholder = "Enter Secret Path (/secret/myapp)"
	vaultPath.CharLimit = 100

	vaultToken := textinput.New()
	vaultToken.Placeholder = "Enter Vault Token"
	vaultToken.EchoMode = textinput.EchoPassword
	vaultToken.CharLimit = 200

	inputs = []textinput.Model{vaultURL, vaultPath, vaultToken}
}

return inputs
```

}

func (m credentialsModel) saveCredentials() error {
switch m.selectedType {
case stackspotUser:
clientID := m.inputs[0].Value()
clientSecret := m.inputs[1].Value()
// TODO: Salvar credenciais de usuário no keyring
return saveUserCredentials(clientID, clientSecret)

```
case stackspotService:
	clientID := m.inputs[0].Value()
	clientSecret := m.inputs[1].Value()
	// TODO: Salvar credenciais de serviço no keyring
	return saveServiceCredentials(clientID, clientSecret)

case hashicorpVault:
	vaultURL := m.inputs[0].Value()
	vaultPath := m.inputs[1].Value()
	vaultToken := m.inputs[2].Value()
	// TODO: Salvar credenciais do Vault no keyring
	return saveVaultCredentials(vaultURL, vaultPath, vaultToken)
}
return nil
```

}

func (m credentialsModel) View() string {
switch m.step {
case stepSelectType:
return m.viewSelectType()
case stepEnterCredentials:
return m.viewEnterCredentials()
case stepConfirm:
return m.viewConfirm()
case stepComplete:
return m.viewComplete()
}
return “”
}

func (m credentialsModel) viewSelectType() string {
s := titleStyle.Render(“🔐 Stackspot CLI - Authentication Setup\n”)
s += “Select the type of credentials to configure:\n\n”

```
options := []string{
	"Stackspot User (Personal credentials)",
	"Stackspot Service (Service account)",
	"Hashicorp Vault (External secret management)",
}

for i, option := range options {
	if i == m.cursor {
		s += selectedStyle.Render("→ " + option)
	} else {
		s += optionStyle.Render("  " + option)
	}
	s += "\n"
}

s += "\n"
s += lipgloss.NewStyle().Faint(true).Render("Use ↑/↓ to navigate, Enter to select, q to quit")

return s
```

}

func (m credentialsModel) viewEnterCredentials() string {
var title string
switch m.selectedType {
case stackspotUser:
title = “📱 Stackspot User Credentials”
case stackspotService:
title = “🏢 Stackspot Service Credentials”
case hashicorpVault:
title = “🔒 Hashicorp Vault Configuration”
}

```
s := titleStyle.Render(title + "\n")

for i, input := range m.inputs {
	var label string
	switch m.selectedType {
	case stackspotUser, stackspotService:
		if i == 0 {
			label = "Client ID:"
		} else {
			label = "Client Secret:"
		}
	case hashicorpVault:
		switch i {
		case 0:
			label = "Vault URL:"
		case 1:
			label = "Secret Path:"
		case 2:
			label = "Vault Token:"
		}
	}

	s += fmt.Sprintf("%s\n", label)
	s += inputStyle.Render(input.View()) + "\n\n"
}

s += lipgloss.NewStyle().Faint(true).Render("Tab/↓: Next field, Shift+Tab/↑: Previous field, Enter: Continue, Esc: Back")

return s
```

}

func (m credentialsModel) viewConfirm() string {
var title string
switch m.selectedType {
case stackspotUser:
title = “📱 Stackspot User Credentials”
case stackspotService:
title = “🏢 Stackspot Service Credentials”
case hashicorpVault:
title = “🔒 Hashicorp Vault Configuration”
}

```
s := titleStyle.Render(title + "\n")
s += "Please confirm your credentials:\n\n"

switch m.selectedType {
case stackspotUser, stackspotService:
	s += fmt.Sprintf("Client ID: %s\n", m.inputs[0].Value())
	s += fmt.Sprintf("Client Secret: %s\n", strings.Repeat("*", len(m.inputs[1].Value())))
case hashicorpVault:
	s += fmt.Sprintf("Vault URL: %s\n", m.inputs[0].Value())
	s += fmt.Sprintf("Secret Path: %s\n", m.inputs[1].Value())
	s += fmt.Sprintf("Vault Token: %s\n", strings.Repeat("*", len(m.inputs[2].Value())))
}

s += "\n"
s += lipgloss.NewStyle().Faint(true).Render("Save these credentials? (y/N)")

return s
```

}

func (m credentialsModel) viewComplete() string {
if m.savedSuccessfully {
s := successStyle.Render(“✅ Credentials saved successfully!\n\n”)
s += “You can now use the Stackspot CLI with your configured credentials.\n”
s += lipgloss.NewStyle().Faint(true).Render(“Press any key to exit…”)
return s
} else {
s := errorStyle.Render(“❌ Failed to save credentials\n\n”)
s += fmt.Sprintf(“Error: %s\n\n”, m.errorMsg)
s += lipgloss.NewStyle().Faint(true).Render(“Press any key to exit…”)
return s
}
}

// Funções helper para salvar credenciais
func saveUserCredentials(clientID, clientSecret string) error {
// TODO: Implementar salvamento no keyring
fmt.Printf(“Saving user credentials: %s\n”, clientID)
return nil
}

func saveServiceCredentials(clientID, clientSecret string) error {
// TODO: Implementar salvamento no keyring
fmt.Printf(“Saving service credentials: %s\n”, clientID)
return nil
}

func saveVaultCredentials(vaultURL, vaultPath, vaultToken string) error {
// TODO: Implementar salvamento no keyring
fmt.Printf(“Saving vault credentials: %s, %s\n”, vaultURL, vaultPath)
return nil
}
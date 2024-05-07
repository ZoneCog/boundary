// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package search

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/boundary/api/targets"
	daemoncmd "github.com/hashicorp/boundary/internal/clientcache/cmd/daemon"
	"github.com/hashicorp/boundary/internal/gen/controller/api"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func (c *SearchCommand) tui(ctx context.Context) error {
	cl, err := c.Client()
	if err != nil {
		return err
	}
	t := cl.Token()
	if t == "" {
		return fmt.Errorf("Auth Token selected for searching is empty.")
	}
	tSlice := strings.SplitN(t, "_", 3)
	if len(tSlice) != 3 {
		return fmt.Errorf("Auth Token selected for searching is in an unexpected format.")
	}
	dotPath, err := daemoncmd.DefaultDotDirectory(ctx)
	if err != nil {
		return err
	}

	searchFn := searchFn(strings.Join(tSlice[:2], "_"), dotPath)
	m := initModel(searchFn)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}

	if m.err != nil {
		return m.err
	}
	if m.selectedTargetId != "" {
		if err := connectFn(m.selectedTargetId, m.selectedConnectSubCmd); err != nil {
			return err
		}
	}

	return nil
}

func initModel(searchFn func(input string) tea.Cmd) *model {
	return &model{
		searchInput: initSearchInput(),
		resources: []string{
			"targets",
		},
		targetTable: initTargetTable(),
		searchFn:    searchFn,
		detailsView: initDetailsView(),
	}
}

func initSearchInput() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Type in your search..."
	ta.Focus()

	ta.Prompt = "> "
	ta.CharLimit = 280
	ta.SetWidth(30)
	ta.SetHeight(1)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	return ta
}

func initTargetTable() table.Model {
	columns := []table.Column{
		{Title: "Id", Width: 12},
		{Title: "Type", Width: 6},
		{Title: "Name", Width: 50},
		{Title: "Description", Width: 50},
	}
	tarTable := table.New(
		table.WithColumns(columns),
		table.WithFocused(false),
		table.WithHeight(20),
		table.WithFocused(false),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	tarTable.SetStyles(s)
	return tarTable
}

const detailsWidth = 118

func initDetailsView() viewport.Model {

	vp := viewport.New(detailsWidth, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(detailsWidth),
	)
	if err != nil {
		return vp
	}

	str, err := renderer.Render("")
	if err != nil {
		return vp
	}

	vp.SetContent(str)
	return vp
}

const tuiSectionCount = 2

type model struct {
	selectedTargetId      string
	selectedConnectSubCmd string

	// Section 0
	searchInput textarea.Model

	// Section 1
	resources       []string
	resourcesCursor int

	// Section 2
	targetTable   table.Model
	targetResults []*targets.Target

	detailsView viewport.Model

	sectionSelected int

	searchFn func(input string) tea.Cmd
	apiErr   *api.Error
	err      error
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, m.searchFn("")) //, tea.EnterAltScreen)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		searchCmd  tea.Cmd
		tableCmd   tea.Cmd
		detailsCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		}
		if cmd := m.updateExecDetection(msg); cmd != nil {
			return m, cmd
		}
		// Process the input
		m.updateSectionFocus(msg)
		searchCmd = m.updateSearch(msg)
		tableCmd = m.updateTableNavigation(msg)
		detailsCmd = m.updateDetailsView(msg)
	case *resultsMsg:
		m.targetResults = msg.targets
		var rows []table.Row
		for _, t := range msg.targets {
			rows = append(rows, table.Row{t.Id, t.Type, t.Name, t.Description})
		}
		m.targetTable.SetRows(rows)
		m.targetTable, tableCmd = m.targetTable.Update(msg)

	case *api.Error:
		m.apiErr = msg
	case errMsg:
		// TODO: something
	}

	return m, tea.Batch(searchCmd, tableCmd, detailsCmd)
}

func (m *model) updateDetailsView(msg tea.Msg) tea.Cmd {
	if !m.targetTable.Focused() || m.targetTable.SelectedRow() == nil {
		return nil
	}
	idx := m.targetTable.Cursor()
	if len(m.targetResults) <= idx {
		return nil
	}
	tar := m.targetResults[idx]
	const content = `
	# ID %s (Type: %s)
	**Name**:        %s
	**Description**: %s
	**Address**:     %s
	**Created**:     %s
	**Updated**:     %s
	**Scope**:       %s
	`

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(detailsWidth),
	)
	if err != nil {
		return nil
	}

	str, err := renderer.Render(fmt.Sprintf(content,
		tar.Id, tar.Type, tar.Name, tar.Description, tar.Address, tar.CreatedTime, tar.UpdatedTime, tar.ScopeId))
	if err != nil {
		return nil
	}

	m.detailsView.SetContent(str)
	return nil
}

func (m *model) updateSearch(msg tea.Msg) tea.Cmd {
	if !m.searchInput.Focused() {
		return nil
	}
	var inputCmd tea.Cmd
	m.searchInput, inputCmd = m.searchInput.Update(msg)
	searchCmd := m.searchFn(m.searchInput.Value())
	return tea.Sequence(inputCmd, searchCmd)
}

func (m *model) updateExecDetection(msg tea.Msg) tea.Cmd {
	if !m.targetTable.Focused() {
		return nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEnter.String():
			targetId := m.targetTable.SelectedRow()[0]
			m.selectedTargetId = targetId
			return tea.Quit
		case "s":
			targetId := m.targetTable.SelectedRow()[0]
			m.selectedTargetId = targetId
			m.selectedConnectSubCmd = "ssh"
			return tea.Quit
		case "h":
			targetId := m.targetTable.SelectedRow()[0]
			m.selectedTargetId = targetId
			m.selectedConnectSubCmd = "http"
			return tea.Quit
		case "p":
			targetId := m.targetTable.SelectedRow()[0]
			m.selectedTargetId = targetId
			m.selectedConnectSubCmd = "postgres"
			return tea.Quit
		case "r":
			targetId := m.targetTable.SelectedRow()[0]
			m.selectedTargetId = targetId
			m.selectedConnectSubCmd = "rdp"
			return tea.Quit
		case "k":
			targetId := m.targetTable.SelectedRow()[0]
			m.selectedTargetId = targetId
			m.selectedConnectSubCmd = "kube"
			return tea.Quit
		}
	}
	return nil
}

func (m *model) updateTableNavigation(msg tea.Msg) tea.Cmd {
	if !m.targetTable.Focused() {
		return nil
	}
	var tableCmd tea.Cmd
	m.targetTable, tableCmd = m.targetTable.Update(msg)
	return tableCmd
}

func (m *model) updateSectionFocus(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type != tea.KeyTab {
			return
		}
		m.sectionSelected = (m.sectionSelected + 1) % tuiSectionCount
	default:
		return
	}

	m.targetTable.Blur()
	m.searchInput.Blur()
	switch m.sectionSelected {
	case 0:
		m.searchInput.Focus()
	case 1:
		m.targetTable.Focus()
	}
}

func (m *model) View() string {
	// The header
	s := "Boundary Target Search\n"
	s += fmt.Sprintf("%s\n", m.searchInput.View())
	if m.apiErr != nil {
		s += fmt.Sprintf("Error: %s\n", m.apiErr.Message)
	}
	s += m.targetTable.View()
	s += "\n"

	// The footer
	if m.targetTable.Focused() && m.targetTable.SelectedRow() != nil {
		targetId := m.targetTable.SelectedRow()[0]
		s += fmt.Sprintf("Press ENTER to connect to %q. 's' with ssh, 'p' with postgres, 'h' with http, 'r' with rdp, 'k' with kube\n", targetId)
	} else {
		s += "\n"
	}
	s += m.detailsView.View()
	s += "\n"
	s += "\nPress Ctrl+C or ESC to quit. Press TAB to switch between search and results.\n"

	// Send the UI for rendering
	return s
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type resultsMsg struct {
	targets []*targets.Target
}

func searchFn(at, dotPath string) func(input string) tea.Cmd {
	return func(input string) tea.Cmd {
		return func() tea.Msg {
			var query []string
			query = append(query, fmt.Sprintf("id %% '%s'", input))
			query = append(query, fmt.Sprintf("type %% '%s'", input))
			query = append(query, fmt.Sprintf("name %% '%s'", input))
			query = append(query, fmt.Sprintf("description %% '%s'", input))
			query = append(query, fmt.Sprintf("address %% '%s'", input))
			query = append(query, fmt.Sprintf("scope_id %% '%s'", input))
			tf := filterBy{
				flagQuery:   strings.Join(query, " or "),
				resource:    "targets",
				authTokenId: at,
			}
			_, res, apiErr, err := search(context.TODO(), dotPath, tf)
			if err != nil {
				return errMsg{err: err}
			}

			if apiErr != nil {
				return apiErr
			}

			return &resultsMsg{
				targets: res.Targets,
			}
		}
	}
}

func connectFn(tid, subCmd string) error {
	cmdName, err := os.Executable()
	if err != nil {
		return err
	}

	args := []string{"connect"}
	if subCmd != "" {
		args = append(args, subCmd)
	}
	args = append(args, "-target-id", tid)

	c := exec.Command(cmdName, args...)
	c.Env = os.Environ()
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	return c.Run()
}

type execFinishedMsg struct{ err error }

func openEditor() tea.Cmd {
	c := exec.Command("uname") //nolint:gosec
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return execFinishedMsg{err}
	})
}

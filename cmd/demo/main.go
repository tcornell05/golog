// cmd/tool/demo/main.go
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/tcornell05/golog"
)

// Example struct for demoing structure logging
type UserInfo struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
	Roles     []string
	Settings  map[string]interface{}
	IsActive  bool
}

// Example nested struct
type APIResponse struct {
	Status  int
	Message string
	Data    UserInfo
	Errors  []string
	Timing  struct {
		RequestTime time.Time
		Duration    time.Duration
	}
}

// fakeFilePath returns a fake file path using a random buzz word and a file name.
func fakeFilePath() string {
	// Use BS() for buzzwords and combine with other functions for filename
	return "/" + gofakeit.BS() + "/" + gofakeit.Word() + gofakeit.FileExtension()
}

// fakeDuration returns a fake duration string in milliseconds.
func fakeDuration() string {
	return fmt.Sprintf("%dms", gofakeit.Number(50, 500))
}

// generateFakeUser creates a fake user structure for demo
func generateFakeUser() UserInfo {
	return UserInfo{
		ID:        gofakeit.Number(1000, 9999),
		Username:  gofakeit.Username(),
		Email:     gofakeit.Email(),
		CreatedAt: gofakeit.Date(),
		Roles:     []string{"user", "member", gofakeit.JobTitle()},
		Settings: map[string]interface{}{
			"theme":             gofakeit.Color(),
			"notifications":     gofakeit.Bool(),
			"twoFactorAuth":     gofakeit.Bool(),
			"preferredLanguage": gofakeit.Language(),
		},
		IsActive: gofakeit.Bool(),
	}
}

// generateFakeAPIResponse creates a complex nested structure
func generateFakeAPIResponse() APIResponse {
	resp := APIResponse{
		Status:  gofakeit.Number(200, 299),
		Message: gofakeit.Sentence(5),
		Data:    generateFakeUser(),
		Errors:  []string{},
	}

	resp.Timing.RequestTime = time.Now().Add(-time.Duration(gofakeit.Number(100, 500)) * time.Millisecond)
	resp.Timing.Duration = time.Duration(gofakeit.Number(10, 200)) * time.Millisecond

	// Occasionally add errors
	if gofakeit.Bool() {
		resp.Status = gofakeit.Number(400, 599)
		resp.Errors = append(resp.Errors, gofakeit.Sentence(3))
		resp.Errors = append(resp.Errors, gofakeit.Sentence(4))
	}

	return resp
}

func init() {
	// Seed the fake data generator.
	gofakeit.Seed(time.Now().UnixNano())
}

// runDevelopmentDemo logs five realistic messages using the Development handler.
func runDevelopmentDemo() {
	logger := golog.NewDevelopment()

	for i := 0; i < 3; i++ {
		// Emit debug log.
		logger.Debug(gofakeit.Sentence(10),
			slog.String("config", fakeFilePath()),
			slog.Int("retry", gofakeit.Number(1, 5)),
		)
		// Emit info log.
		logger.Info(gofakeit.Sentence(10),
			slog.String("db", gofakeit.AppName()),
		)
		// Emit warning log.
		logger.Warn(gofakeit.Sentence(10),
			slog.String("module", gofakeit.AppName()),
		)
		// Emit error log - this will now automatically capture the code snippet
		logger.Error("Error: "+gofakeit.Sentence(10),
			// Use real source information instead of fake
			slog.Any("details", map[string]interface{}{
				"reason":  gofakeit.HackerPhrase(),
				"attempt": gofakeit.Number(1, 3),
			}),
		)
	}
}

// runDiscardDemo logs five messages using the Discard handler.
func runDiscardDemo() {
	logger := golog.NewDiscard()
	for i := 0; i < 5; i++ {
		logger.Info(gofakeit.Sentence(10),
			slog.String("test", "discard"),
		)
	}
	os.Stdout.WriteString("Discard handler: no log output should appear.\n")
}

// runProductionDemo logs five messages using the Production handler.
func runProductionDemo() {
	logger := golog.NewProduction()
	for i := 0; i < 5; i++ {
		logger.Info(gofakeit.Sentence(10),
			slog.String("user", gofakeit.Username()),
			slog.String("session", gofakeit.UUID()),
		)
		logger.Error("Critical error: "+gofakeit.Sentence(10),
			slog.String("component", gofakeit.HackerVerb()),
			slog.Int("error_code", gofakeit.Number(400, 599)),
		)
	}
}

// runTeraDemo logs five messages using the custom Tera handler.
func runTeraDemo() {
	teraOpts := golog.TeraOptions{
		HeaderTitle:      "Traceback (most recent call last)",
		BorderColor:      lipgloss.Color("#FF5F87"),
		PanelBorderColor: lipgloss.Color("#5F87FF"),
		DateFormat:       time.RFC822,
	}
	teraHandler := golog.NewTeraHandler(teraOpts)
	logger := slog.New(teraHandler)

	for i := 0; i < 3; i++ {
		// Emit error log - this will automatically capture and display the code snippet
		logger.Error("Error: "+gofakeit.Sentence(10),
			// No need to specify fake file and line - it will be captured automatically
			slog.Any("locals", map[string]interface{}{
				"input": gofakeit.Word(),
				"retry": gofakeit.Bool(),
			}),
		)

		logger.Debug("Debug: "+gofakeit.Sentence(10),
			slog.String("duration", fakeDuration()),
		)
	}
}

// runStructureDemo shows logging of various data structures
func runStructureDemo() {
	logger := golog.NewDevelopment()

	fmt.Println("\n=== Demo: Simple Structure Logging ===")
	user := generateFakeUser()
	logger.Info("User information",
		slog.Any("user", user),
	)

	fmt.Println("\n=== Demo: Complex Nested Structure Logging ===")
	apiResponse := generateFakeAPIResponse()
	logger.Info("API Response received",
		slog.Any("response", apiResponse),
	)

	fmt.Println("\n=== Demo: Error with Structure Context ===")
	if apiResponse.Status >= 400 {
		logger.Error("API request failed",
			slog.Int("statusCode", apiResponse.Status),
			slog.Any("errors", apiResponse.Errors),
			slog.Any("context", apiResponse),
		)
	}

	fmt.Println("\n=== Demo: Tera Handler with Structure ===")
	teraOpts := golog.TeraOptions{
		HeaderTitle:      "Structure Demo",
		BorderColor:      lipgloss.Color("#5F87FF"),
		PanelBorderColor: lipgloss.Color("#87FF5F"),
		DateFormat:       time.RFC822,
	}
	teraHandler := golog.NewTeraHandler(teraOpts)
	teraLogger := slog.New(teraHandler)

	teraLogger.Info("Complex structure visualization",
		slog.Any("userInfo", user),
		slog.Any("apiResponse", apiResponse),
	)
}

// runCodeSnippetDemo shows real-world examples of code snippets in errors
func runCodeSnippetDemo() {
	logger := golog.NewDevelopment()

	// Demo 1: Log an error directly - will show this exact line in the snippet
	fmt.Println("\n=== Demo 1: Direct Error Logging ===")
	logger.Error("This error occurs at this exact line in main.go")

	// Demo 2: Function with a deliberate error
	fmt.Println("\n=== Demo 2: Error with Call Stack ===")
	demonstrateRealError(logger)

	// Demo 3: Error propagation through multiple functions
	fmt.Println("\n=== Demo 3: Error Propagation ===")
	if err := nestedFunction(3); err != nil {
		// Log the error at the top level - will show both this location
		// and the deeper location where the error originated
		logger.Error("Top-level error occurred",
			slog.Any("error", err),
		)
	}

	// Demo 4: TeraHandler with real error
	fmt.Println("\n=== Demo 4: TeraHandler with Real Error ===")
	teraOpts := golog.TeraOptions{
		HeaderTitle:      "Error Traceback",
		BorderColor:      lipgloss.Color("#FF5F87"),
		PanelBorderColor: lipgloss.Color("#5F87FF"),
		DateFormat:       time.RFC822,
	}
	teraHandler := golog.NewTeraHandler(teraOpts)
	teraLogger := slog.New(teraHandler)

	// This will show the exact line in the code snippet
	teraLogger.Error("Error with TeraHandler showing this exact line")

	// Show different log levels with their respective colors and messages
	teraLogger.Debug("Debug message with source line")
	teraLogger.Info("Info message with source line")
	teraLogger.Warn("Warning message with source line")
}

// Helper function to demonstrate a real error
func demonstrateRealError(logger *slog.Logger) {
	// Try to access an index out of bounds to create a real error
	defer func() {
		if r := recover(); r != nil {
			// Log the recovered panic - will show where the panic occurred
			logger.Error("Recovered from panic",
				slog.Any("error", r),
				slog.String("operation", "array access"),
			)
		}
	}()

	// This will cause a panic - the code snippet should highlight this line
	var data []string
	_ = data[5] // Deliberate out of bounds access
}

// Helper function for nested error demo
func nestedFunction(depth int) error {
	if depth <= 0 {
		// When we reach depth 0, return an error
		// The code snippet will show this line as the error origin
		return fmt.Errorf("reached maximum depth")
	}

	// Recurse deeper
	err := nestedFunction(depth - 1)
	if err != nil {
		// Wrap and propagate the error
		return fmt.Errorf("error at depth %d: %w", depth, err)
	}

	return nil
}

func main() {
	var devFlag, discardFlag, prodFlag, teraFlag, snippetFlag, structFlag bool

	rootCmd := &cobra.Command{
		Use:   "golog-demo",
		Short: "Demo tool for golog handlers; use flags --development, --discard, --production, --tera, --snippet, --struct",
		Run: func(cmd *cobra.Command, args []string) {
			// If any flag is provided, run only those handlers.
			ran := false
			if devFlag {
				runDevelopmentDemo()
				ran = true
			}
			if discardFlag {
				runDiscardDemo()
				ran = true
			}
			if prodFlag {
				runProductionDemo()
				ran = true
			}
			if teraFlag {
				runTeraDemo()
				ran = true
			}
			if snippetFlag {
				runCodeSnippetDemo()
				ran = true
			}
			if structFlag {
				runStructureDemo()
				ran = true
			}
			if !ran {
				// Run all demos if no flag is provided.
				runDevelopmentDemo()
				runDiscardDemo()
				runProductionDemo()
				runTeraDemo()
				runCodeSnippetDemo()
				runStructureDemo()
			}
		},
	}

	// Define flags for each handler demo.
	rootCmd.Flags().BoolVar(&devFlag, "development", false, "Run the development handler demo")
	rootCmd.Flags().BoolVar(&discardFlag, "discard", false, "Run the discard handler demo")
	rootCmd.Flags().BoolVar(&prodFlag, "production", false, "Run the production handler demo")
	rootCmd.Flags().BoolVar(&teraFlag, "tera", false, "Run the tera handler demo")
	rootCmd.Flags().BoolVar(&snippetFlag, "snippet", false, "Run the code snippet demo")
	rootCmd.Flags().BoolVar(&structFlag, "struct", false, "Run the structure logging demo")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

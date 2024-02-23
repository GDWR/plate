package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	templateDir   = "$HOME/.local/share/plate"
	configDir     = "$HOME/.config/plate"
	configName    = "config"
	defaultEditor = "vim" // Fallback if $EDITOR doesn't exist, or is empty

	rootCmd = &cobra.Command{
		Use:   "plate",
		Short: "Generate templated files quickly from your terminal.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	addCmd = &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Example: "plate add <template>",
		Short:   "Add a new template",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = defaultEditor
			}

			if err := runEditor(editor, args[0]); err != nil {
				return err
			}

			return nil
		},
	}

	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List all templates",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := os.ExpandEnv(templateDir)
			files, err := os.ReadDir(dir)
			if err != nil {
				return err
			}

			for _, file := range files {
				fmt.Printf("- %s\n", file.Name())
			}

			return nil
		},
	}
	createCmd = &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Example: "plate create <template> <destination>",
		Short:   "Create a new file from a template",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			dst := args[1]
			templateFile := fmt.Sprintf("%s/%s", os.ExpandEnv(templateDir), templateName)
			if _, err := os.Stat(templateFile); os.IsNotExist(err) {
				return fmt.Errorf("template %s does not exist", templateName)
			}

			data, err := os.ReadFile(templateFile)
			if err != nil {
				return fmt.Errorf("template %s could not be read", templateName)
			}

			if err := os.WriteFile(dst, data, 0644); err != nil {
				return fmt.Errorf("%s could not be written to", dst)
			}

			return nil
		},
	}
	deleteCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d", "rm"},
		Example: "plate delete <template>",
		Short:   "Delete a template",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			templateName := args[0]
			templateFile := fmt.Sprintf("%s/%s", os.ExpandEnv(templateDir), templateName)
			if _, err := os.Stat(templateFile); os.IsNotExist(err) {
				return fmt.Errorf("template %s does not exist", templateName)
			}

			if err := os.Remove(templateFile); err != nil {
				return fmt.Errorf("template %s could not be deleted", templateName)
			}

			return nil
		},
	}
)

func runEditor(editor string, templateName string) error {
	templateFile := fmt.Sprintf("%s/%s", os.ExpandEnv(templateDir), templateName)
	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		if _, err := os.Create(templateFile); err != nil {
			return err
		}
	}

	cmd := exec.Command(editor, templateFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	// Ensure that the config and template directories exist
	if err := os.MkdirAll(os.ExpandEnv(configDir), 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := os.MkdirAll(os.ExpandEnv(templateDir), 0755); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(configDir)
	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.SafeWriteConfig()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

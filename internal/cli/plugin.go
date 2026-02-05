package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func newPluginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "プラグイン管理",
		Long:  "プラグインの一覧表示や管理を行います",
	}

	cmd.AddCommand(newPluginListCmd())

	return cmd
}

func newPluginListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "利用可能なプラグイン一覧を表示",
		Long:  "Arsenal で利用可能な全プラグイン（ツール）の一覧を表示します",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluginList()
		},
	}
}

func runPluginList() error {
	// レジストリから全プラグインを取得
	plugins := registry.All()

	// プラグイン名でソート
	names := make([]string, 0, len(plugins))
	for name := range plugins {
		names = append(names, name)
	}
	sort.Strings(names)

	// 表示
	fmt.Println("利用可能ツール:")
	fmt.Println()

	for _, name := range names {
		p := plugins[name]
		fmt.Printf("  %s\n", p.DisplayName)
		fmt.Printf("    名称: %s\n", p.Name)
		if p.Description != "" {
			fmt.Printf("    説明: %s\n", p.Description)
		}
		fmt.Println()
	}

	return nil
}

package main

import (
	"path/filepath"
	"strings"
)

// VsWorkspace is the structure that gets unmarshaled into the correct
// Structure folders to be used as .code-workspace files for vscode
type VsWorkspace struct {
	Folders    []folderPath      `json:"folders"`
	Settings   map[string]string `json:"settings"`
	Extensions extensionRecs     `json:"extensions"`
}

type folderPath struct {
	Path string `json:"path"`
}

type extensionRecs struct {
	Recommendations []string `json:"recommendations"`
}

// Vscode creates .code-workspace files. It uses the config file to
// create the folder structure and the key to make the correct
// .code-workspace file. Key represents the base path ie. "code"
func MakeVsWorkspace(c Config, key string) VsWorkspace {
	mapper := getPathMap(c)
	var newVs VsWorkspace

	for _, path := range mapper[key] {
		fp := makeFolderPath(path)
		newVs.Folders = append(newVs.Folders, fp)
	}

	for k, v := range c.VsCodeSettings {
		newVs.addSetting(k, v)
	}

	for _, ext := range c.VsCodeExtensions {
		newVs.addExtension(ext)
	}

	return newVs
}

// getPathMap uses the Structure from Config to return a map of the base
// directory as key and the full path to the folder as value
func getPathMap(c Config) map[string][]string {
	mapper := map[string][]string{}

	for _, v := range c.Structure {
		base := strings.Split(v, "/")[0]

		mapper[base] = append(mapper[base], filepath.Join(c.Path, v))
	}
	return mapper
}

// makeFolderPath is a factory function to create folder path structs
func makeFolderPath(path string) folderPath {
	newFolderPath := folderPath{
		Path: path,
	}

	return newFolderPath
}

// addSetting adds a new setting to VsWorkspace. If the settings map does not
// have any k,v pairs, create a newSetting and set it as vs.Settings. Otherwise,
// add an additional k,v pair to vs.Settings.
func (vs *VsWorkspace) addSetting(k, v string) {
	switch len(vs.Settings) {
	case 0:
		newSetting := map[string]string{k: v}
		vs.Settings = newSetting
	default:
		vs.Settings[k] = v
	}

}

// addExtension adds a new extension to the VsWorkspace Extension list.
func (vs *VsWorkspace) addExtension(ex string) {
	vs.Extensions.Recommendations = append(vs.Extensions.Recommendations, ex)
}

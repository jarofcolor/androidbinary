package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/jarofcolor/androidbinary"
	"github.com/jarofcolor/androidbinary/apk"
)

type ApkInfo struct {
	FileSize       int64
	FileSizeFormat string
	Name           string
	NameCN         string
	Pkg            string
	Icon           string
	LaunchActivty  string
	VersionCode    int32
	VersionName    string
	MinVersion     int32
	TargetVersion  int32
	Permissions    []string
	MetaInfo       map[string]string
}

func ParseApk(apkPath string, isAll bool) (*ApkInfo, error) {
	pkg, err := apk.OpenFile(apkPath)
	if err != nil {
		return nil, err
	}
	defer pkg.Close()

	versionCode, _ := pkg.Manifest().VersionCode.Int32()
	versionName, _ := pkg.Manifest().VersionName.String()
	name, _ := pkg.Manifest().App.Label.WithResTableConfig(&androidbinary.ResTableConfig{
		Language: [2]uint8{uint8('e'), uint8('n')},
	}).String()
	nameCN, _ := pkg.Manifest().App.Label.WithResTableConfig(&androidbinary.ResTableConfig{
		Language: [2]uint8{uint8('z'), uint8('h')},
		Country:  [2]uint8{uint8('C'), uint8('N')},
	}).String()
	minVersion, _ := pkg.Manifest().SDK.Min.Int32()
	targetVersion, _ := pkg.Manifest().SDK.Target.Int32()

	densities := []uint16{640, 560, 480, 440, 420, 360, 240, 260, 120}
	iconData := []byte{}
	for _, density := range densities {
		iconData, err = pkg.Icon(&androidbinary.ResTableConfig{
			Density: density,
		})
		if err == nil {
			break
		}
	}

	iconBase64Str := ""
	if err == nil {
		iconBase64Str = base64.StdEncoding.EncodeToString(iconData)
	}

	info := &ApkInfo{
		Pkg:           pkg.PackageName(),
		Name:          name,
		NameCN:        nameCN,
		Icon:          iconBase64Str,
		VersionCode:   versionCode,
		VersionName:   versionName,
		MinVersion:    minVersion,
		TargetVersion: targetVersion,
		Permissions:   []string{},
		MetaInfo:      map[string]string{},
	}

	mainActivity, _ := pkg.MainActivity()
	info.LaunchActivty = mainActivity

	f, _ := os.Stat(apkPath)
	info.FileSize = f.Size()
	info.FileSizeFormat = formatFileSize(f.Size())

	if !isAll {
		return info, nil
	}

	for _, v := range pkg.Manifest().UsesPermissions {
		info.Permissions = append(info.Permissions, v.Name.MustString())
	}

	for _, v := range pkg.Manifest().App.MetaData {
		metaName, _ := v.Name.String()
		metaValue, _ := v.Value.String()
		if metaName != "" && metaValue != "" {
			info.MetaInfo[metaName] = metaValue
		}
	}

	return info, nil
}

func formatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

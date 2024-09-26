package gores

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func NewResDate() ResDate {
	goRes := ResDate{
		ICOName:                           "icon.png",
		Name:                              "WPS Office",
		Version:                           "12.1.0.16399",
		Description:                       "WPS Office",
		MinimumOs:                         "win7",
		ExecutionLevel:                    "requireAdministrator",
		UIAccess:                          false,
		AutoElevate:                       true,
		DpiAwareness:                      "system",
		DisableTheming:                    false,
		DisableWindowFiltering:            false,
		HighResolutionScrollingAware:      false,
		UltraHighResolutionScrollingAware: false,
		LongPathAware:                     false,
		PrinterDriverIsolation:            false,
		GDIScaling:                        false,
		SegmentHeap:                       false,
		UseCommonControlsV6:               false,
		FixedFileVersion:                  "12.1.0.16399",
		FixedProductVersion:               "WPS Office",
		Comments:                          "",
		CompanyName:                       "",
		FileDescription:                   "WPS Office",
		FileVersion:                       "12.1.0.16399",
		InternalName:                      "",
		LegalCopyright:                    "Copyright©2024 Kingsoft Corporation. All rights reserved.",
		LegalTrademarks:                   "",
		OriginalFilename:                  "wps_host.exe",
		PrivateBuild:                      "",
		ProductName:                       "WPS Office",
		ProductVersion:                    "WPS Office",
		SpecialBuild:                      "",
	}

	return goRes
}

func (rt ResTmpl) ResRender() error {
	content, err := os.ReadFile(rt.ResPath)
	if err != nil {
		return fmt.Errorf("<os.ReadFile()> err: %s", err)
	}

	tmpl, err := template.New(filepath.Base(rt.ResPath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("<template.New()> err: %s", err)
	}

	resFile, err := os.Create(filepath.Join(rt.OutputDir, "winres.json"))
	if err != nil {
		return fmt.Errorf("<os.Create()> err: %s", err)
	}
	defer resFile.Close()

	data := NewResDate()
	err = tmpl.Execute(resFile, data)
	if err != nil {
		return fmt.Errorf("<render.Execute()> err: %s", err)
	}

	return nil
}

package tools

import (
	"github.com/fatih/color"
)

var Red = color.New(color.FgRed).PrintfFunc()
var Green = color.New(color.FgGreen).PrintfFunc()
var Orange = color.New(color.FgYellow).PrintfFunc() // orange uses yellow in terminal
var Blue = color.New(color.FgBlue).PrintfFunc()
var Yellow = color.New(color.FgYellow).PrintfFunc()
var Cyan = color.New(color.FgCyan).PrintfFunc()
var Magenta = color.New(color.FgMagenta).PrintfFunc()
var White = color.New(color.FgWhite).PrintfFunc()

var BoldRed = color.New(color.FgRed, color.Bold).PrintfFunc()
var BoldGreen = color.New(color.FgGreen, color.Bold).PrintfFunc()
var BoldOrange = color.New(color.FgYellow, color.Bold).PrintfFunc()

var RedString = color.New(color.FgRed).SprintfFunc()
var GreenString = color.New(color.FgGreen).SprintfFunc()
var OrangeString = color.New(color.FgYellow).SprintfFunc()

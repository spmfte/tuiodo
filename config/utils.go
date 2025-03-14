package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color format regexes
var (
	hexColorRegex  = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
	rgbColorRegex  = regexp.MustCompile(`^rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)$`)
	ansiColorRegex = regexp.MustCompile(`^ansi(\d+)$`)
)

// ParseColor parses a color string into a lipgloss.Color
// Supports hex (#RGB, #RRGGBB), RGB (rgb(r,g,b)), and ANSI (ansi0-255) formats
func ParseColor(colorStr string) (lipgloss.Color, error) {
	// Trim whitespace
	colorStr = strings.TrimSpace(colorStr)

	// Handle empty or "none" color
	if colorStr == "" || strings.ToLower(colorStr) == "none" {
		return lipgloss.Color(""), nil
	}

	// Check if it's a named color
	if namedColor, ok := namedColors[strings.ToLower(colorStr)]; ok {
		return lipgloss.Color(namedColor), nil
	}

	// Check if it's a hex color
	if hexColorRegex.MatchString(colorStr) {
		return lipgloss.Color(colorStr), nil
	}

	// Check if it's an RGB color
	if match := rgbColorRegex.FindStringSubmatch(colorStr); match != nil {
		r, _ := strconv.Atoi(match[1])
		g, _ := strconv.Atoi(match[2])
		b, _ := strconv.Atoi(match[3])

		// Convert to hex format
		hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)
		return lipgloss.Color(hex), nil
	}

	// Check if it's an ANSI color
	if match := ansiColorRegex.FindStringSubmatch(colorStr); match != nil {
		code, _ := strconv.Atoi(match[1])
		if code >= 0 && code <= 255 {
			return lipgloss.Color(strconv.Itoa(code)), nil
		}
	}

	// If it's just a number between 0-255, treat as ANSI color code
	if code, err := strconv.Atoi(colorStr); err == nil && code >= 0 && code <= 255 {
		return lipgloss.Color(colorStr), nil
	}

	return lipgloss.Color(""), fmt.Errorf("invalid color format: %s", colorStr)
}

// ColorToHex converts various color formats to hex string
func ColorToHex(c lipgloss.Color) string {
	// If it's already hex, return it
	colorStr := string(c)

	if hexColorRegex.MatchString(colorStr) {
		return colorStr
	}

	// Try to parse as ANSI color code
	if code, err := strconv.Atoi(colorStr); err == nil && code >= 0 && code <= 255 {
		// This is a simplified conversion from ANSI -> RGB
		// In a real implementation, this would be more accurate
		r, g, b := ansiToRGB(code)
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}

	return "#000000" // Default black if couldn't parse
}

// ansiToRGB converts an ANSI color code to RGB values
// This is a simplified implementation
func ansiToRGB(code int) (int, int, int) {
	// Simple ANSI color table - in a real implementation this would be more accurate
	ansiColors := []struct{ r, g, b int }{
		{0, 0, 0},       // Black
		{170, 0, 0},     // Red
		{0, 170, 0},     // Green
		{170, 85, 0},    // Yellow
		{0, 0, 170},     // Blue
		{170, 0, 170},   // Magenta
		{0, 170, 170},   // Cyan
		{170, 170, 170}, // White
		{85, 85, 85},    // Bright Black
		{255, 85, 85},   // Bright Red
		{85, 255, 85},   // Bright Green
		{255, 255, 85},  // Bright Yellow
		{85, 85, 255},   // Bright Blue
		{255, 85, 255},  // Bright Magenta
		{85, 255, 255},  // Bright Cyan
		{255, 255, 255}, // Bright White
	}

	if code < 16 {
		return ansiColors[code].r, ansiColors[code].g, ansiColors[code].b
	}

	// Handle 216 colors (16-231)
	if code >= 16 && code <= 231 {
		// 6x6x6 color cube
		code -= 16
		r := (code / 36) * 51
		g := ((code % 36) / 6) * 51
		b := (code % 6) * 51
		return r, g, b
	}

	// Handle grayscale (232-255)
	code -= 232
	v := code*10 + 8
	return v, v, v
}

// HexToRGB converts a hex color string to RGB
func HexToRGB(hex string) (r, g, b int, err error) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) == 3 {
		// Convert 3-digit hex to 6-digit
		var rInt64 int64
		rInt64, err = strconv.ParseInt(string(hex[0])+string(hex[0]), 16, 0)
		if err != nil {
			return 0, 0, 0, err
		}
		r = int(rInt64)

		var gInt64 int64
		gInt64, err = strconv.ParseInt(string(hex[1])+string(hex[1]), 16, 0)
		if err != nil {
			return 0, 0, 0, err
		}
		g = int(gInt64)

		var bInt64 int64
		bInt64, err = strconv.ParseInt(string(hex[2])+string(hex[2]), 16, 0)
		if err != nil {
			return 0, 0, 0, err
		}
		b = int(bInt64)
	} else if len(hex) == 6 {
		// Parse 6-digit hex
		var rInt64 int64
		rInt64, err = strconv.ParseInt(hex[0:2], 16, 0)
		if err != nil {
			return 0, 0, 0, err
		}
		r = int(rInt64)

		var gInt64 int64
		gInt64, err = strconv.ParseInt(hex[2:4], 16, 0)
		if err != nil {
			return 0, 0, 0, err
		}
		g = int(gInt64)

		var bInt64 int64
		bInt64, err = strconv.ParseInt(hex[4:6], 16, 0)
		if err != nil {
			return 0, 0, 0, err
		}
		b = int(bInt64)
	} else {
		return 0, 0, 0, fmt.Errorf("invalid hex color: %s", hex)
	}

	return r, g, b, nil
}

// Common named colors
var namedColors = map[string]string{
	"black":   "#000000",
	"white":   "#ffffff",
	"red":     "#ff0000",
	"green":   "#00ff00",
	"blue":    "#0000ff",
	"yellow":  "#ffff00",
	"cyan":    "#00ffff",
	"magenta": "#ff00ff",
	"orange":  "#ffa500",
	"purple":  "#800080",
	"pink":    "#ffc0cb",
	"brown":   "#a52a2a",
	"gray":    "#808080",
	"grey":    "#808080",
}

// DetermineColorMode automatically determines the best color mode
// based on terminal capabilities
func DetermineColorMode() string {
	// This is a simplified implementation
	// In a real app, you would check terminal capabilities

	// Default to "auto" which lets lipgloss decide
	return "auto"
}

// ExpandPath expands ~ to the user's home directory
func ExpandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, path[1:]), nil
}

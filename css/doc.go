// Package css provides CSS parsing, styling, and stylesheet management for
// the Textual TUI framework.
//
// This package is a Go port of the Python Textual CSS engine. It supports
// a subset of CSS adapted for terminal user interfaces (TCSS — Textual CSS).
//
// # Overview
//
// The package is organized as follows:
//
//   - [Token] and [Tokenizer] handle low-level CSS tokenization.
//   - [Tokenize] runs the full TCSS tokenization state machine.
//   - [Scalar] and [ScalarOffset] represent CSS dimension values.
//   - [Styles] holds all CSS property values for a node.
//   - [RenderStyles] merges a base [Styles] with inline overrides.
//   - [StylesBuilder] converts parsed declarations into a [Styles] object.
//   - [ParseSelectors], [ParseDeclarations], and [Parse] parse selectors,
//     inline declarations, and full stylesheets respectively.
//   - [Match] matches selector sets against DOM nodes.
//   - [Stylesheet] manages multiple CSS sources and applies them to nodes.
package css

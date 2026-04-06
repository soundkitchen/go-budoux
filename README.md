# BudouX for Go

BudouX is a standalone, small, and language-neutral phrase segmenter tool that provides beautiful and legible line breaks.
This fork targets Go 1.23+, adopting the iterator APIs and model optimisations introduced in Go 1.23.
The bundled models are currently synced against upstream BudouX `v0.8.1`.

## Supported Bundled Models

- Japanese default (ja)
- Japanese KNBC base model (ja_knbc)
- Simplified Chinese (zh-hans)
- Traditional Chinese (zh-hant)
- Thai (th)

## Installation

```bash
go get github.com/soundkitchen/go-budoux
```

## Usage

### Simple usage

```go
package main

import (
    "fmt"
    budoux "github.com/soundkitchen/go-budoux"
)

func main() {
    parser := budoux.NewDefaultJapaneseParser()
    phrases := parser.Parse("今日は良い天気ですね。")
    fmt.Println(phrases)
    // Output: [今日は 良い 天気ですね。]
}
```

### Bundled parser constructors

- Japanese: `budoux.NewDefaultJapaneseParser()`
- Japanese KNBC base model: `budoux.NewJapaneseKNBCParser()`
- Simplified Chinese: `budoux.NewDefaultSimplifiedChineseParser()`
- Traditional Chinese: `budoux.NewDefaultTraditionalChineseParser()`
- Thai: `budoux.NewDefaultThaiParser()`

## API Reference

### Parser

The main interface for text segmentation.

#### Constructor Functions

```go
// Create parsers with default models
func NewDefaultJapaneseParser() *Parser
func NewJapaneseKNBCParser() *Parser
func NewDefaultSimplifiedChineseParser() *Parser
func NewDefaultTraditionalChineseParser() *Parser
func NewDefaultThaiParser() *Parser

// Create parser with custom model
func New(model models.Model) *Parser
```

#### Methods

```go
// Parse segments input text into phrases
func (p *Parser) Parse(sentence string) []string
```

## Testing

Run the test suite:

```bash
GOCACHE=$(pwd)/.cache go test ./...
```

Run tests with verbose output:

```bash
GOCACHE=$(pwd)/.cache go test -v ./...
```

Run the parser benchmark:

```bash
GOCACHE=$(pwd)/.cache go test -bench=BenchmarkParserParse -benchmem ./...
```

## Regenerating models

The bundled language models are generated from the upstream BudouX `v0.8.1` JSON assets. To refresh them, run:

```bash
GOCACHE=$(pwd)/.cache go generate ./gen
```

## License

Copyright 2021 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

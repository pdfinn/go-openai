
# Go OpenAI Library

## Overview
This Go library provides a simple interface to interact with the OpenAI GPT-3 API.

## Installation
To include this package in your project, use `go get` to download it.

```bash
go get -u github.com/pdfinn/go-openai/openai
```

## Usage

### Importing the Library
```go
import "github.com/pdfinn/go-openai/openai"
```

### Setting up the Configuration
Create a `Config` struct and populate its fields:

```go
cfg := go-openai.Config{
	APIKey:      "your-api-key-here",
	Instruction: "Translate the following English text to French: '{}'.",
	Input:       "Hello, World!",
	Temperature: 0.7,
	Model:       "gpt-3.5-turbo",
	Debug:       false,
}
```

### Making an API Call
Call the `Callgo-openai` function with the `Config` struct:

```go
response, err := go-openai.CallOpenAI(cfg, nil)
if err != nil {
	log.Fatalf("Failed to call OpenAI: %v", err)
}
```

### Handling the Response
The `CallOpenAI` function returns the API response as a string:

```go
fmt.Println("API Response:", response)
```

## Functions

### `CallOpenAI(cfg Config, httpClient *http.Client) (string, error)`
Makes an API call to OpenAI.

**Parameters**:
- `cfg`: Configuration settings.
- `httpClient`: Optional custom HTTP client.

**Returns**:
- API response as a string.
- Error, if any.

### `ValidateModel(model string) bool`
Checks if the model is supported.

**Parameters**:
- `model`: Model name as a string.

**Returns**:
- `true` if the model is supported, `false` otherwise.

## Errors
- `ErrInvalidModel`: Unsupported or invalid model.
- `ErrNoChoices`: No choices returned from the API.
- `ErrNoAssistantMessage`: No assistant message found in the API response.

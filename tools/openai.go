package tools

import (
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

// ForOpenAI returns all Gomind tools in OpenAI Chat Completions API format.
func ForOpenAI() []openai.ChatCompletionToolUnionParam {
	defs := Definitions()
	tools := make([]openai.ChatCompletionToolUnionParam, len(defs))
	for i, def := range defs {
		tools[i] = defToOpenAI(def)
	}
	return tools
}

// ForOpenAIResponses returns all Gomind tools in OpenAI Responses API format.
func ForOpenAIResponses() []responses.ToolUnionParam {
	defs := Definitions()
	tools := make([]responses.ToolUnionParam, len(defs))
	for i, def := range defs {
		tools[i] = defToOpenAIResponses(def)
	}
	return tools
}

func defToOpenAI(def Definition) openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
		Name:        def.Name,
		Description: param.NewOpt(def.Description),
		Parameters:  paramsToOpenAI(def.Parameters),
	})
}

func defToOpenAIResponses(def Definition) responses.ToolUnionParam {
	return responses.ToolUnionParam{
		OfFunction: &responses.FunctionToolParam{
			Name:        def.Name,
			Description: param.NewOpt(def.Description),
			Parameters:  paramsToOpenAI(def.Parameters),
			Strict:      param.NewOpt(true),
		},
	}
}

func paramsToOpenAI(params []*Param) shared.FunctionParameters {
	properties := make(map[string]any)
	var required []string

	for _, p := range params {
		properties[p.Name] = paramToOpenAI(p)
		if p.Required {
			required = append(required, p.Name)
		}
	}

	result := shared.FunctionParameters{
		"type":                 "object",
		"properties":           properties,
		"additionalProperties": false,
	}
	if len(required) > 0 {
		result["required"] = required
	}
	return result
}

func paramToOpenAI(p *Param) map[string]any {
	result := map[string]any{
		"type":        string(p.Type),
		"description": p.Description,
	}

	if p.Type == TypeArray && p.Items != nil {
		result["items"] = paramToOpenAI(p.Items)
	}

	if p.Type == TypeObject && p.Properties != nil {
		props := make(map[string]any)
		var required []string
		for name, prop := range p.Properties {
			props[name] = paramToOpenAI(prop)
			if prop.Required {
				required = append(required, name)
			}
		}
		result["properties"] = props
		result["additionalProperties"] = false
		if len(required) > 0 {
			result["required"] = required
		}
	}

	return result
}

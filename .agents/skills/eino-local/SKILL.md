---
name: eino-local
description: Use for any Go/Eino work in this repo — pinned version, the model/tool/agent APIs actually validated in AI01, retry and cancel limits, and the rule to verify against pinned source before coding.
---

# Eino, pinned for this repo

This repo pins `github.com/cloudwego/eino v0.9.23` (stable release). Upstream Eino docs and maintainer skills describe `main` and have drifted from this tag; two API guesses failed to compile during AI01 because of that. Use them for orientation only, never as the final authority.

## Rules

1. Use the APIs below. They were compiled and tested against the pinned tag in [experiments/ai01](../../../experiments/ai01/).
2. Before writing framework code, confirm the signature in the pinned module source in the local module cache (`$GOMODCACHE/github.com/cloudwego/eino@v0.9.23/...`) or with `go doc`. Browse the web only if the pin cannot answer the question.
3. When the Eino version changes, update this file in the same commit. This is the only place in the repo that claims version-specific API facts.
4. This spike is evidence, not production code. No streaming, paid calls, real data or Core writes. Core owns authorization and every business write; our code owns scoped identity, output validation, limits and safe errors.

## Validated surface (v0.9.23)

**Model** — `model.BaseChatModel` is the alias for `BaseModel[*schema.Message]`; `model.BaseModel` is generic and must not be used bare.

```go
Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error)
Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error)
```

Personal requests use only `Generate` and get one complete message. `Stream` must stay unused; a fake treats any stream call as a test failure. `model.ToolCallingChatModel` adds `WithTools([]*schema.ToolInfo) (ToolCallingChatModel, error)`. ADK passes tools as a common call option, readable with `model.GetCommonOptions(&model.Options{}, opts...)` → `opts.Tools`.

**Tool** — implement `tool.InvokableTool`, or derive one with `utils.InferTool`/`utils.NewTool`.

```go
Info(ctx context.Context) (*schema.ToolInfo, error)
InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error)
```

`Info` uses `schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{...})` or struct tags via `utils.GoStruct2ParamsOneOf[T]()`. A tool is read-only here. Strict JSON decoding (`DisallowUnknownFields`, trailing-content check), required fields and actor/organization checks are **our code**, not the framework's.

**Agent and runner** — use ADK, not a hand-built ReAct graph.

```go
agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name: "...", Instruction: "...",
    Model: cm,
    ToolsConfig: adk.ToolsConfig{
        ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []tool.BaseTool{myTool}},
    },
    MaxIterations: 3,
    ModelRetryConfig: &adk.ModelRetryConfig{MaxRetries: 2, IsRetryAble: fn, BackoffFunc: fn},
})
iter := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent}).Query(ctx, prompt)
```

`MaxIterations` stops a looping tool call. `MaxRetries: 2` means three total model attempts; only errors accepted by `IsRetryAble` retry, permanent errors stop after one attempt. Event messages come from `iter.Next()`, then `event.Output.MessageOutput.GetMessage()`; use `MessageOutput.Role`, `.ToolName` and `schema.Message.ToolCallID` to correlate a tool result with its call. Ignore `adk.WillRetryError` when a final answer exists.

**Cancellation** — create a fresh option per run: `cancelOpt, cancel := adk.WithCancel()`, pass `cancelOpt` to `Query`, then `handle, ok := cancel()` and check `handle.Wait()` (it returns one value, not two). For deadlines use `context.WithTimeout`; every blocking step must select on `ctx.Done()`. On timeout return a safe error with no partial output.

**Structured drafts** — the classic ChatModel path does not enforce output against a schema. `utils.GoStruct2ParamsOneOf[Draft]()` + `ToJSONSchema()` only describe fields and required flags. Validate model text locally with strict JSON decoding and explicit bounds; reject unknown action fields.

## Deliberately out of scope

`Stream`, the `AgenticModel`/`AgenticMessage` path, hosted provider calls without approval, and the untested ADK surface (middleware, checkpoints, multi-agent, TurnLoop). The older `flow/agent/react` manual shows stale signatures such as `Info()` without context; do not code from it.

## Evidence and optional deeper reading

- Pinned decisions and hypotheses: [PLAN.md](../../../experiments/ai01/PLAN.md), [AI01_VALIDATION_RECORD.md](../../../docs/engineering/AI01_VALIDATION_RECORD.md).
- Worked code: [fakemodel.go](../../../experiments/ai01/internal/fakemodel/fakemodel.go), [lookup.go](../../../experiments/ai01/internal/lookup/lookup.go), [agentrun.go](../../../experiments/ai01/internal/agentrun/agentrun.go), [draft.go](../../../experiments/ai01/internal/draft/draft.go).
- Raw failures and fixes: [iterations.md](../../../experiments/ai01/evidence/iterations.md).
- Optional upstream orientation only: [Eino docs](https://www.cloudwego.io/docs/eino/), [eino-ext skills](https://github.com/cloudwego/eino-ext/tree/main/skills). Verify anything you take from them against the pinned tag.

package ai01

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cloudwego/eino/adk"

	"campusly.local/ai01/internal/agentrun"
	"campusly.local/ai01/internal/fakemodel"
)

type requestResult struct {
	text string
	err  error
}

// Eight synthetic requests, one shared agent, unique prompts. This is an
// isolation and race check, not a capacity benchmark.
func TestE05ConcurrentRequestsStayIsolated(t *testing.T) {
	ctx := context.Background()
	const requests = 8
	script := make(fakemodel.Script, requests)
	for i := 0; i < requests; i++ {
		prompt := fmt.Sprintf("synthetic:concurrent-%d", i)
		script[prompt] = []fakemodel.Step{{Content: fmt.Sprintf("answer-%d", i), Delay: 5 * time.Millisecond}}
	}
	model := fakemodel.New(script)
	agent, err := agentrun.NewAgent(ctx, model, nil)
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}
	runner := agentrun.NewRunner(agent)

	results := make([]requestResult, requests)
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			prompt := fmt.Sprintf("synthetic:concurrent-%d", i)
			text, runErr := agentrun.Outcome(agentrun.Collect(runner.Query(ctx, prompt)))
			results[i] = requestResult{text: text, err: runErr}
		}(i)
	}
	wg.Wait()

	for i, result := range results {
		want := fmt.Sprintf("answer-%d", i)
		if result.err != nil {
			t.Errorf("request %d: %v", i, result.err)
		}
		if result.text != want {
			t.Errorf("request %d: got %q, want %q", i, result.text, want)
		}
	}
	if model.StreamCalls() != 0 {
		t.Fatalf("stream calls = %d, want 0", model.StreamCalls())
	}
}

func TestE05DeadlineReturnsSafeError(t *testing.T) {
	const prompt = "synthetic:deadline"
	model := fakemodel.New(fakemodel.Script{
		prompt: {{Content: "too late", Delay: 200 * time.Millisecond}},
	})
	agent, err := agentrun.NewAgent(context.Background(), model, nil)
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	text, runErr := agentrun.Outcome(agentrun.Collect(agentrun.NewRunner(agent).Query(ctx, prompt)))
	if text != "" {
		t.Fatalf("deadline produced an answer: %q", text)
	}
	if !errors.Is(runErr, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context.DeadlineExceeded", runErr)
	}
	if strings.Contains(fmt.Sprint(runErr), "too late") {
		t.Fatal("error leaked the synthetic model output")
	}
}

func TestE05ExplicitCancelStopsRun(t *testing.T) {
	const prompt = "synthetic:cancel"
	model := fakemodel.New(fakemodel.Script{
		prompt: {{Content: "should not appear", Delay: 200 * time.Millisecond}},
	})
	agent, err := agentrun.NewAgent(context.Background(), model, nil)
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}
	cancelOption, cancel := adk.WithCancel()
	iter := agentrun.NewRunner(agent).Query(context.Background(), prompt, cancelOption)
	handle, contributed := cancel()
	if !contributed {
		t.Fatal("cancel did not contribute to the running request")
	}
	if err := handle.Wait(); err != nil {
		t.Fatalf("cancel handle: %v", err)
	}
	text, runErr := agentrun.Outcome(agentrun.Collect(iter))
	if text != "" {
		t.Fatalf("cancelled run produced an answer: %q", text)
	}
	var cancelError *adk.CancelError
	if !errors.As(runErr, &cancelError) && !errors.Is(runErr, context.Canceled) {
		t.Fatalf("error = %v, want a cancellation error", runErr)
	}
}

func TestE05BoundedRetryUnderConcurrency(t *testing.T) {
	const requests = 4
	script := make(fakemodel.Script, requests)
	for i := 0; i < requests; i++ {
		prompt := fmt.Sprintf("synthetic:retry-concurrent-%d", i)
		script[prompt] = []fakemodel.Step{
			{Err: fakemodel.ErrTransient},
			{Err: fakemodel.ErrTransient},
			{Content: fmt.Sprintf("done-%d", i)},
		}
	}
	model := fakemodel.New(script)
	agent, err := agentrun.NewAgent(context.Background(), model, nil)
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}
	runner := agentrun.NewRunner(agent)
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			prompt := fmt.Sprintf("synthetic:retry-concurrent-%d", i)
			text, runErr := agentrun.Outcome(agentrun.Collect(runner.Query(context.Background(), prompt)))
			if runErr != nil || text != fmt.Sprintf("done-%d", i) {
				t.Errorf("request %d: text=%q err=%v", i, text, runErr)
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < requests; i++ {
		prompt := fmt.Sprintf("synthetic:retry-concurrent-%d", i)
		if got := model.Attempts(prompt); got != agentrun.MaxRetries+1 {
			t.Errorf("request %d: attempts = %d, want %d", i, got, agentrun.MaxRetries+1)
		}
	}
}

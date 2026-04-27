// Package reporter provides formatted output for envchain validation results.
package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envchain/internal/chain"
)

// Format represents the output format for reports.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes validation results to an output stream.
type Reporter struct {
	out    io.Writer
	format Format
}

// New creates a Reporter writing to the given writer with the specified format.
func New(out io.Writer, format Format) *Reporter {
	if out == nil {
		out = os.Stdout
	}
	return &Reporter{out: out, format: format}
}

// Report writes the validation results for all stages.
func (r *Reporter) Report(results []chain.StageResult) error {
	switch r.format {
	case FormatJSON:
		return r.reportJSON(results)
	default:
		return r.reportText(results)
	}
}

func (r *Reporter) reportText(results []chain.StageResult) error {
	for _, res := range results {
		status := "✓ PASS"
		if !res.OK {
			status = "✗ FAIL"
		}
		fmt.Fprintf(r.out, "[%s] stage: %s\n", status, res.Stage)
		if len(res.Missing) > 0 {
			fmt.Fprintf(r.out, "  missing: %s\n", strings.Join(res.Missing, ", "))
		}
		if len(res.Empty) > 0 {
			fmt.Fprintf(r.out, "  empty:   %s\n", strings.Join(res.Empty, ", "))
		}
	}
	return nil
}

func (r *Reporter) reportJSON(results []chain.StageResult) error {
	fmt.Fprintln(r.out, "[")
	for i, res := range results {
		missing := jsonStringSlice(res.Missing)
		empty := jsonStringSlice(res.Empty)
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(r.out, "  {\"stage\":%q,\"ok\":%v,\"missing\":%s,\"empty\":%s}%s\n",
			res.Stage, res.OK, missing, empty, comma)
	}
	fmt.Fprintln(r.out, "]")
	return nil
}

func jsonStringSlice(ss []string) string {
	if len(ss) == 0 {
		return "[]"
	}
	quoted := make([]string, len(ss))
	for i, s := range ss {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return "[" + strings.Join(quoted, ",") + "]"
}

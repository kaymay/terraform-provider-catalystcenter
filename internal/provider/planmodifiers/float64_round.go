// Package planmodifiers contains hand-owned plan modifiers used by the
// Catalyst Center provider. Files here are NOT generated and are safe to edit.
package planmodifiers

import (
	"context"
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RoundFloat64 returns a plan modifier that rounds a float64 attribute to the
// given number of decimal places (issue #145: lat/long phantom drift).
// The read path (fromBody) MUST round to the same 'decimals' or the diff persists.
func RoundFloat64(decimals int) planmodifier.Float64 {
	return roundFloat64Modifier{decimals: decimals}
}

// Coordinate5DecimalPlanModifier rounds a float64 to 5 decimals. It is a
// zero-value struct so the generator can emit it as `Coordinate5DecimalPlanModifier{}`
// (the same convention used by other custom_modifier attributes).
type Coordinate5DecimalPlanModifier struct{}

func (m Coordinate5DecimalPlanModifier) Description(_ context.Context) string {
	return "Rounds the coordinate to 5 decimal places to match Catalyst Center's stored precision."
}

func (m Coordinate5DecimalPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m Coordinate5DecimalPlanModifier) PlanModifyFloat64(_ context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	resp.PlanValue = types.Float64Value(RoundTo(req.PlanValue.ValueFloat64(), 5))
}

type roundFloat64Modifier struct {
	decimals int
}

func (m roundFloat64Modifier) Description(_ context.Context) string {
	return fmt.Sprintf("Rounds the value to %d decimal places to match Catalyst Center's stored precision.", m.decimals)
}

func (m roundFloat64Modifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m roundFloat64Modifier) PlanModifyFloat64(_ context.Context, req planmodifier.Float64Request, resp *planmodifier.Float64Response) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	rounded := RoundTo(req.PlanValue.ValueFloat64(), m.decimals)
	resp.PlanValue = types.Float64Value(rounded)
}

// RoundTo rounds v to 'decimals' places, half-away-from-zero (math.Round
// semantics), matching the module's historical floor(x*1e5 + 0.5)/1e5.
// Used by both the plan modifier and fromBody to guarantee plan/state symmetry.
func RoundTo(v float64, decimals int) float64 {
	if decimals < 0 {
		decimals = 0
	}
	scale := math.Pow(10, float64(decimals))
	return math.Round(v*scale) / scale
}

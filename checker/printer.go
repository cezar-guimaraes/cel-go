// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package checker

import (
	"github.com/google/cel-go/common/debug"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type semanticAdorner struct {
	checks *exprpb.CheckedExpr
}

var _ debug.Adorner = &semanticAdorner{}

func (a *semanticAdorner) GetMetadata(elem interface{}) string {
	result := ""
	e, isExpr := elem.(*exprpb.Expr)
	if !isExpr {
		return result
	}
	t := a.checks.TypeMap[e.Id]
	if t != nil {
		result += "~"
		result += FormatCheckedType(t)
	}

	switch e.ExprKind.(type) {
	case *exprpb.Expr_IdentExpr,
		*exprpb.Expr_CallExpr,
		*exprpb.Expr_StructExpr,
		*exprpb.Expr_SelectExpr:
		if ref, found := a.checks.ReferenceMap[e.Id]; found {
			if len(ref.GetOverloadId()) == 0 {
				result += "^" + ref.Name
			} else {
				for i, overload := range ref.OverloadId {
					if i == 0 {
						result += "^"
					} else {
						result += "|"
					}
					result += overload
				}
			}
		}
	}

	return result
}

func print(e *exprpb.Expr, checks *exprpb.CheckedExpr) string {
	a := &semanticAdorner{checks: checks}
	return debug.ToAdornedDebugString(e, a)
}

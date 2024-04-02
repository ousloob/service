package tests

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ardanlabs/service/app/services/sales-api/apis/crud/userapi"
	"github.com/ardanlabs/service/business/api/errs"
	"github.com/ardanlabs/service/business/data/dbtest"
	"github.com/google/go-cmp/cmp"
)

func userUpdate200(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/users/%s", sd.Users[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Model: &userapi.AppUpdateUser{
				Name:            dbtest.StringPointer("Jack Kennedy"),
				Email:           dbtest.StringPointer("jack@ardanlabs.com"),
				Department:      dbtest.StringPointer("IT"),
				Password:        dbtest.StringPointer("123"),
				PasswordConfirm: dbtest.StringPointer("123"),
			},
			Resp: &userapi.AppUser{},
			ExpResp: &userapi.AppUser{
				ID:          sd.Users[0].ID.String(),
				Name:        "Jack Kennedy",
				Email:       "jack@ardanlabs.com",
				Roles:       []string{"USER"},
				Department:  "IT",
				Enabled:     true,
				DateCreated: sd.Users[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Users[0].DateUpdated.Format(time.RFC3339),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*userapi.AppUser)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*userapi.AppUser)
				gotResp.DateUpdated = expResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:       "role",
			URL:        fmt.Sprintf("/v1/users/role/%s", sd.Admins[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Model: &userapi.AppUpdateUserRole{
				Roles: []string{"USER"},
			},
			Resp: &userapi.AppUser{},
			ExpResp: &userapi.AppUser{
				ID:          sd.Admins[0].ID.String(),
				Name:        sd.Admins[0].Name,
				Email:       sd.Admins[0].Email.Address,
				Roles:       []string{"USER"},
				Department:  sd.Admins[0].Department,
				Enabled:     true,
				DateCreated: sd.Admins[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Admins[0].DateUpdated.Format(time.RFC3339),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*userapi.AppUser)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*userapi.AppUser)
				gotResp.DateUpdated = expResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func userUpdate400(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:       "bad-input",
			URL:        fmt.Sprintf("/v1/users/%s", sd.Users[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Model: &userapi.AppUpdateUser{
				Email:           dbtest.StringPointer("bill@"),
				PasswordConfirm: dbtest.StringPointer("jack"),
			},
			Resp:    &errs.Error{},
			ExpResp: toErrorPtr(errs.Newf(errs.FailedPrecondition, "validate: [{\"field\":\"email\",\"error\":\"email must be a valid email address\"},{\"field\":\"passwordConfirm\",\"error\":\"passwordConfirm must be equal to Password\"}]")),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-role",
			URL:        fmt.Sprintf("/v1/users/role/%s", sd.Admins[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Model: &userapi.AppUpdateUserRole{
				Roles: []string{"BAD ROLE"},
			},
			Resp:    &errs.Error{},
			ExpResp: toErrorPtr(errs.Newf(errs.FailedPrecondition, "parse: invalid role \"BAD ROLE\"")),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func userUpdate401(sd dbtest.SeedData) []dbtest.AppTable {
	table := []dbtest.AppTable{
		{
			Name:       "emptytoken",
			URL:        fmt.Sprintf("/v1/users/%s", sd.Users[0].ID),
			Token:      "",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			Resp:       &errs.Error{},
			ExpResp:    toErrorPtr(errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments")),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badsig",
			URL:        fmt.Sprintf("/v1/users/%s", sd.Users[0].ID),
			Token:      sd.Users[0].Token + "A",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			Resp:       &errs.Error{},
			ExpResp:    toErrorPtr(errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]")),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        fmt.Sprintf("/v1/users/%s", sd.Admins[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			Model: &userapi.AppUpdateUser{
				Name:            dbtest.StringPointer("Bill Kennedy"),
				Email:           dbtest.StringPointer("bill@ardanlabs.com"),
				Department:      dbtest.StringPointer("IT"),
				Password:        dbtest.StringPointer("123"),
				PasswordConfirm: dbtest.StringPointer("123"),
			},
			Resp:    &errs.Error{},
			ExpResp: toErrorPtr(errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]")),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "roleadminonly",
			URL:        fmt.Sprintf("/v1/users/role/%s", sd.Users[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			Model: &userapi.AppUpdateUserRole{
				Roles: []string{"ADMIN"},
			},
			Resp:    &errs.Error{},
			ExpResp: toErrorPtr(errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[[{USER}]] rule[rule_admin_only]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]")),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

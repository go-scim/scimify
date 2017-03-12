package validation

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mutabilityRulesTest struct {
	name         string
	resourcePath string
	assertion    func(bool, error)
}

func BenchmarkMutabilityRulesValidator_Validate(b *testing.B) {
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}

	ref, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	validator := &mutabilityRulesValidator{}
	opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}

	ctx := context.Background()
	ctx = context.WithValue(context.Background(), resource.CK_Schema, schema)
	ctx = context.WithValue(ctx, resource.CK_Reference, ref)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := validator.Validate(r, opt, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestMutabilityRulesValidator_Validate(t *testing.T) {
	validator := &mutabilityRulesValidator{}

	for _, test := range []mutabilityRulesTest{
		{
			"test same resource passes",
			"../test_data/single_test_user_david.json",
			func(ok bool, err error) {
				assert.True(t, ok)
				assert.Nil(t, err)
			},
		},
		{
			"test changes readOnly string attribute",
			"../test_data/changes_readonly_string.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "id", err.(*validationError).FullPath)
			},
		},
		{
			"test changes readOnly complex attribute",
			"../test_data/changes_readonly_complex.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "meta", err.(*validationError).FullPath)
			},
		},
	} {
		schema, err := loadSchema("../test_data/test_user_schema_all.json")
		if err != nil {
			t.Fatal(err)
		}

		referenceData := loadTestDataFromJson(t, "../test_data/single_test_user_david.json")
		reference := resource.NewResourceFromMap(referenceData)

		resourceData := loadTestDataFromJson(t, test.resourcePath)
		r := resource.NewResourceFromMap(resourceData)

		opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}

		ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)
		ctx = context.WithValue(ctx, resource.CK_Reference, reference)

		ok, err := validator.Validate(r, opt, ctx)
		test.assertion(ok, err)
	}
}

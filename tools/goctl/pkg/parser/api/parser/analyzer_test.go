package parser

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/assertx"
)

func Test_Parse(t *testing.T) {
	t.Run(
		"valid", func(t *testing.T) {
			apiSpec, err := Parse("./testdata/example.api", nil)
			assert.Nil(t, err)
			ast := assert.New(t)
			ast.Equal(
				spec.Info{
					Title:   "type title here",
					Desc:    "type desc here",
					Version: "type version here",
					Author:  "type author here",
					Email:   "type email here",
					Properties: map[string]string{
						"title":   "type title here",
						"desc":    "type desc here",
						"version": "type version here",
						"author":  "type author here",
						"email":   "type email here",
					},
				}, apiSpec.Info,
			)
			ast.True(
				func() bool {
					for _, group := range apiSpec.Service.Groups {
						value, ok := group.Annotation.Properties["summary"]
						if ok {
							return value == "test"
						}
					}
					return false
				}(),
			)
		},
	)

	t.Run(
		"invalid", func(t *testing.T) {
			data, err := os.ReadFile("./testdata/invalid.api")
			assert.NoError(t, err)
			splits := bytes.Split(data, []byte("-----"))
			var testFile []string
			for idx, split := range splits {
				replacer := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\r", "", "\f", "")
				r := replacer.Replace(string(split))
				if len(r) == 0 {
					continue
				}
				filename := filepath.Join(t.TempDir(), fmt.Sprintf("invalid%d.api", idx))
				err := os.WriteFile(filename, split, 0666)
				assert.NoError(t, err)
				testFile = append(testFile, filename)
			}
			for _, v := range testFile {
				_, err := Parse(v, nil)
				assertx.Error(t, err)
			}
		},
	)

	t.Run(
		"circleImport", func(t *testing.T) {
			_, err := Parse("./testdata/base.api", nil)
			assertx.Error(t, err)
		},
	)

	t.Run(
		"link_import", func(t *testing.T) {
			_, err := Parse("./testdata/link_import.api", nil)
			assert.Nil(t, err)
		},
	)

	t.Run(
		"duplicate_types", func(t *testing.T) {
			_, err := Parse("./testdata/duplicate_type.api", nil)
			assertx.Error(t, err)
		},
	)

	t.Run(
		"duplicate_path_expression", func(t *testing.T) {
			_, err := Parse("./testdata/duplicate_path_expression.api", nil)
			assertx.Error(t, err)
		},
	)
	t.Run(
		"duplicate_path_expression_different_prefix", func(t *testing.T) {
			_, err := Parse("./testdata/duplicate_path_expression_different_prefix.api", nil)

			assert.Nil(t, err)
		},
	)
	t.Run(
		"multi_level_embedding", func(t *testing.T) {
			apiSpec, err := Parse("./testdata/multi_level_embedding.api", nil)
			assert.Nil(t, err)
			assert.NotNil(t, apiSpec)

			// Verify that types are resolved correctly
			var ageStruct, fullProfile spec.DefineStruct
			for _, typ := range apiSpec.Types {
				if defineStruct, ok := typ.(spec.DefineStruct); ok {
					switch defineStruct.RawName {
					case "AgeStruct":
						ageStruct = defineStruct
					case "FullProfile":
						fullProfile = defineStruct
					}
				}
			}

			// Test 2-level embedding: AgeStruct embeds NameStruct which embeds IdStruct
			assert.NotEmpty(t, ageStruct.RawName, "AgeStruct should be found")
			assert.Equal(t, 2, len(ageStruct.Members), "AgeStruct should have 2 members: embedded NameStruct and Age field")
			
			// Verify the embedded NameStruct is fully resolved
			var embeddedNameStruct spec.Member
			for _, member := range ageStruct.Members {
				if member.IsInline {
					embeddedNameStruct = member
					break
				}
			}
			assert.True(t, embeddedNameStruct.IsInline, "NameStruct should be embedded (IsInline=true)")
			
			// The embedded type should be fully resolved NameStruct with its members
			if resolvedNameStruct, ok := embeddedNameStruct.Type.(spec.DefineStruct); ok {
				assert.Equal(t, "NameStruct", resolvedNameStruct.RawName)
				assert.Equal(t, 2, len(resolvedNameStruct.Members), "Embedded NameStruct should have 2 members")
				
				// Check that NameStruct has the embedded IdStruct resolved
				var embeddedIdStruct spec.Member
				for _, member := range resolvedNameStruct.Members {
					if member.IsInline {
						embeddedIdStruct = member
						break
					}
				}
				assert.True(t, embeddedIdStruct.IsInline, "IdStruct should be embedded in NameStruct")
				
				if resolvedIdStruct, ok := embeddedIdStruct.Type.(spec.DefineStruct); ok {
					assert.Equal(t, "IdStruct", resolvedIdStruct.RawName)
					assert.Equal(t, 1, len(resolvedIdStruct.Members), "IdStruct should have 1 member")
					assert.Equal(t, "Id", resolvedIdStruct.Members[0].Name)
				} else {
					t.Errorf("Embedded IdStruct should be resolved to DefineStruct, got %T", embeddedIdStruct.Type)
				}
			} else {
				t.Errorf("Embedded NameStruct should be resolved to DefineStruct, got %T", embeddedNameStruct.Type)
			}

			// Test 3-level embedding: FullProfile -> ExtendedProfile -> UserProfile -> BaseUser
			assert.NotEmpty(t, fullProfile.RawName, "FullProfile should be found")
			assert.Equal(t, 3, len(fullProfile.Members), "FullProfile should have 3 members: embedded ExtendedProfile, Bio, and CreatedAt")
			
			// Verify the embedded ExtendedProfile is fully resolved
			var embeddedExtendedProfile spec.Member
			for _, member := range fullProfile.Members {
				if member.IsInline {
					embeddedExtendedProfile = member
					break
				}
			}
			assert.True(t, embeddedExtendedProfile.IsInline, "ExtendedProfile should be embedded")
			
			if resolvedExtendedProfile, ok := embeddedExtendedProfile.Type.(spec.DefineStruct); ok {
				assert.Equal(t, "ExtendedProfile", resolvedExtendedProfile.RawName)
				assert.Equal(t, 3, len(resolvedExtendedProfile.Members), "ExtendedProfile should have 3 members")
				
				// Check that ExtendedProfile has UserProfile embedded and resolved
				var embeddedUserProfile spec.Member
				for _, member := range resolvedExtendedProfile.Members {
					if member.IsInline {
						embeddedUserProfile = member
						break
					}
				}
				
				if resolvedUserProfile, ok := embeddedUserProfile.Type.(spec.DefineStruct); ok {
					assert.Equal(t, "UserProfile", resolvedUserProfile.RawName)
					assert.Equal(t, 3, len(resolvedUserProfile.Members), "UserProfile should have 3 members")
					
					// Check that UserProfile has BaseUser embedded and resolved
					var embeddedBaseUser spec.Member
					for _, member := range resolvedUserProfile.Members {
						if member.IsInline {
							embeddedBaseUser = member
							break
						}
					}
					
					if resolvedBaseUser, ok := embeddedBaseUser.Type.(spec.DefineStruct); ok {
						assert.Equal(t, "BaseUser", resolvedBaseUser.RawName)
						assert.Equal(t, 2, len(resolvedBaseUser.Members), "BaseUser should have 2 members")
						// Verify the deepest fields are present
						assert.Equal(t, "UserId", resolvedBaseUser.Members[0].Name)
						assert.Equal(t, "Email", resolvedBaseUser.Members[1].Name)
					} else {
						t.Errorf("Embedded BaseUser should be resolved to DefineStruct, got %T", embeddedBaseUser.Type)
					}
				} else {
					t.Errorf("Embedded UserProfile should be resolved to DefineStruct, got %T", embeddedUserProfile.Type)
				}
			} else {
				t.Errorf("Embedded ExtendedProfile should be resolved to DefineStruct, got %T", embeddedExtendedProfile.Type)
			}
		},
	)
}

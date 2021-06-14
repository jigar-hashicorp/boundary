package node

import (
	"context"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/boundary/internal/errors"
	"github.com/hashicorp/boundary/internal/kms"
	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptFilter_filterValue(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	var nilBytePtr []byte
	testStr := "fido"
	testInt := 22

	wrapper := TestWrapper(t)
	testFilter := &EncryptFilter{
		Wrapper:  wrapper,
		HmacSalt: []byte("salt"),
		HmacInfo: []byte("info"),
	}

	// optWrapper := TestWrapper(t)
	tests := []struct {
		name            string
		ef              *EncryptFilter
		opt             []Option
		fv              reflect.Value
		classification  *tagInfo
		decryptWrapper  wrapping.Wrapper
		wantValue       string
		wantErrMatch    *errors.Template
		wantErrContains string
	}{
		{
			name:            "missing-tag",
			ef:              testFilter,
			fv:              reflect.ValueOf(&testStr).Elem(),
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "missing classification tag",
		},
		{
			name:            "missing-wrapper-encrypt",
			ef:              &EncryptFilter{},
			fv:              reflect.ValueOf(&testStr).Elem(),
			classification:  &tagInfo{Classification: SensitiveClassification, Operation: EncryptOperation},
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "missing wrapper",
		},
		{
			name:            "missing-wrapper-hmac",
			ef:              &EncryptFilter{},
			fv:              reflect.ValueOf(&testStr).Elem(),
			classification:  &tagInfo{Classification: SensitiveClassification, Operation: HmacSha256Operation},
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "missing wrapper",
		},
		{
			name:            "not-string-or-bytes",
			ef:              testFilter,
			fv:              reflect.ValueOf(&testInt).Elem(),
			classification:  &tagInfo{Classification: SensitiveClassification, Operation: EncryptOperation},
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "field value is not a string or []byte",
		},
		{
			name:           "nil",
			ef:             testFilter,
			fv:             reflect.ValueOf(nil),
			classification: &tagInfo{Classification: SensitiveClassification, Operation: EncryptOperation},
		},
		{
			name:           "nil-byte-ptr",
			ef:             testFilter,
			fv:             reflect.ValueOf(nilBytePtr),
			classification: &tagInfo{Classification: SensitiveClassification, Operation: EncryptOperation},
			decryptWrapper: wrapper,
		},
		{
			name:            "unknown-filter-operation",
			ef:              testFilter,
			fv:              reflect.ValueOf(&testStr).Elem(),
			classification:  &tagInfo{Classification: SecretClassification, Operation: NoOperation},
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "unknown filter operation",
		},
		{
			name:           "success-public",
			ef:             testFilter,
			fv:             reflect.ValueOf(&testStr).Elem(),
			classification: &tagInfo{Classification: PublicClassification},
			wantValue:      testStr,
		},
		{
			name:           "success-secret-hmac",
			ef:             testFilter,
			fv:             reflect.ValueOf(&testStr).Elem(),
			classification: &tagInfo{Classification: SecretClassification, Operation: HmacSha256Operation},
			wantValue:      testHmacSha256(t, []byte("fido"), wrapper, []byte("salt"), []byte("info")),
		},
		{
			name:           "success-secret-encrypt",
			ef:             testFilter,
			fv:             reflect.ValueOf(&testStr).Elem(),
			classification: &tagInfo{Classification: SecretClassification, Operation: EncryptOperation},
			decryptWrapper: wrapper,
			wantValue:      "fido",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStr = "fido" // reset it everytime
			assert, require := assert.New(t), require.New(t)
			err := tt.ef.filterValue(ctx, tt.fv, tt.classification, tt.opt...)
			if tt.wantErrMatch != nil {
				require.Error(err)
				assert.Truef(errors.Match(tt.wantErrMatch, err), "want err %q and got %q", tt.wantErrMatch, err.Error())
				if tt.wantErrContains != "" {
					assert.Contains(err.Error(), tt.wantErrContains)
				}
				return
			}
			require.NoError(err)

			switch tt.classification.Classification {
			case PublicClassification:
				assert.Equal(tt.wantValue, fmt.Sprintf("%s", tt.fv))
			case SecretClassification, SensitiveClassification:
				switch tt.classification.Operation {
				case EncryptOperation:
					switch {
					case tt.fv == reflect.ValueOf(nil):
						assert.Equal(tt.wantValue, "")
					case tt.fv.Type() == reflect.TypeOf([]uint8(nil)):
						assert.Equal(fmt.Sprintf("%s", TestDecryptValue(t, tt.decryptWrapper, tt.fv.Bytes())), tt.wantValue)
					case tt.fv.Type() == reflect.TypeOf(""):
						assert.Equal(fmt.Sprintf("%s", TestDecryptValue(t, tt.decryptWrapper, []byte(tt.fv.String()))), tt.wantValue)
					}
				case HmacSha256Operation:
					assert.Equal(tt.wantValue, fmt.Sprintf("%s", tt.fv))
				case RedactOperation:
					assert.Equal(RedactedData, fmt.Sprintf("%s", tt.fv))
				}
			default:
				assert.Equal(RedactedData, fmt.Sprintf("%s", tt.fv))
			}
		})
	}
}

func TestEncryptFilter_encrypt(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	wrapper := TestWrapper(t)
	testFilter := &EncryptFilter{
		Wrapper:  wrapper,
		HmacSalt: []byte("salt"),
		HmacInfo: []byte("info"),
	}

	optWrapper := TestWrapper(t)
	tests := []struct {
		name            string
		ef              *EncryptFilter
		opt             []Option
		data            []byte
		decryptWrapper  wrapping.Wrapper
		wantErrMatch    *errors.Template
		wantErrContains string
	}{
		{
			name:            "missing-data",
			ef:              testFilter,
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "missing data",
		},
		{
			name:            "missing-wrapper",
			ef:              &EncryptFilter{},
			data:            []byte("fido"),
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "missing wrapper",
		},
		{
			name:           "success",
			ef:             testFilter,
			data:           []byte("fido"),
			decryptWrapper: wrapper,
		},
		{
			name:           "success-with-wrapper",
			ef:             testFilter,
			opt:            []Option{WithWrapper(optWrapper)},
			data:           []byte("fido"),
			decryptWrapper: optWrapper,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			got, err := tt.ef.encrypt(ctx, tt.data, tt.opt...)
			if tt.wantErrMatch != nil {
				require.Error(err)
				assert.Truef(errors.Match(tt.wantErrMatch, err), "want err %q and got %q", tt.wantErrMatch, err.Error())
				if tt.wantErrContains != "" {
					assert.Contains(err.Error(), tt.wantErrContains)
				}
				return
			}
			require.NoError(err)

			assert.Equal(TestDecryptValue(t, tt.decryptWrapper, []byte(got)), tt.data)
		})
	}
}

func TestEncryptFilter_hmacSha256(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	wrapper := TestWrapper(t)
	testFilter := &EncryptFilter{
		Wrapper:  wrapper,
		HmacSalt: []byte("salt"),
		HmacInfo: []byte("info"),
	}

	optWrapper := TestWrapper(t)

	tests := []struct {
		name            string
		ef              *EncryptFilter
		opt             []Option
		data            []byte
		want            string
		wantErrMatch    *errors.Template
		wantErrContains string
	}{
		{
			name:            "missing-data",
			ef:              testFilter,
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "missing data",
		},
		{
			name:            "missing-wrapper",
			ef:              &EncryptFilter{},
			data:            []byte("fido"),
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "missing wrapper",
		},
		{
			name: "success",
			ef:   testFilter,
			data: []byte("fido"),
			want: testHmacSha256(t, []byte("fido"), wrapper, []byte("salt"), []byte("info")),
		},
		{
			name: "success-with-wrapper",
			ef:   testFilter,
			opt:  []Option{WithWrapper(optWrapper)},
			data: []byte("fido"),
			want: testHmacSha256(t, []byte("fido"), optWrapper, []byte("salt"), []byte("info")),
		},
		{
			name: "success-with-info",
			ef:   testFilter,
			data: []byte("fido"),
			opt:  []Option{WithInfo([]byte("opt-info"))},
			want: testHmacSha256(t, []byte("fido"), wrapper, []byte("salt"), []byte("opt-info")),
		},
		{
			name: "success-with-salt",
			ef:   testFilter,
			data: []byte("fido"),
			opt:  []Option{WithSalt([]byte("opt-salt"))},
			want: testHmacSha256(t, []byte("fido"), wrapper, []byte("opt-salt"), []byte("info")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			got, err := tt.ef.hmacSha256(ctx, tt.data, tt.opt...)
			if tt.wantErrMatch != nil {
				require.Error(err)
				assert.Truef(errors.Match(tt.wantErrMatch, err), "want err %q and got %q", tt.wantErrMatch, err.Error())
				if tt.wantErrContains != "" {
					assert.Contains(err.Error(), tt.wantErrContains)
				}
				return
			}
			require.NoError(err)
			assert.Equal(tt.want, got)
		})
	}
}

func Test_setValue(t *testing.T) {
	t.Parallel()
	testInt := 22
	testStr := "fido"
	tests := []struct {
		name            string
		fv              reflect.Value
		newVal          string
		wantErrMatch    *errors.Template
		wantErrContains string
	}{
		{
			name:            "not-string-or-bytes",
			fv:              reflect.ValueOf(&testInt).Elem(),
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "field value is not a string or []byte",
		},
		{
			name:            "not-settable",
			fv:              reflect.ValueOf(&testStr),
			wantErrMatch:    errors.T(errors.InvalidParameter),
			wantErrContains: "unable to set value",
		},
		{
			name:   "string-with-value",
			fv:     reflect.ValueOf(&testStr).Elem(),
			newVal: "alice",
		},
		{
			name:   "empty-string",
			fv:     reflect.ValueOf(&testStr).Elem(),
			newVal: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert, require := assert.New(t), require.New(t)
			err := setValue(tt.fv, tt.newVal)
			if tt.wantErrMatch != nil {
				require.Error(err)
				assert.Truef(errors.Match(tt.wantErrMatch, err), "want err %q and got %q", tt.wantErrMatch, err.Error())
				if tt.wantErrContains != "" {
					assert.Contains(err.Error(), tt.wantErrContains)
				}
				return
			}
			require.NoError(err)
			assert.Equal(fmt.Sprintf("%s", tt.fv), tt.newVal)
		})
	}

}

func testHmacSha256(t *testing.T, data []byte, w wrapping.Wrapper, salt, info []byte) string {
	t.Helper()
	require := require.New(t)
	reader, err := kms.NewDerivedReader(w, 32, salt, info)
	require.NoError(err)
	key, _, err := ed25519.GenerateKey(reader)
	require.NoError(err)

	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(data)
	return "hmac-sh256:" + base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

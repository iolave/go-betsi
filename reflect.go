package betsi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/iolave/go-errors"
)

// encodeAppRequest takes a request url with path params as ".../{paramName}/..." and a struct v with/wo "ar" tags.
//
//   - It detects and replaces url path params from "ar" tags ("path=*")
//   - It detects the body type from "ar" tags ("body=json") and use it's value to encode the request body.
//
// If some error occurs, an error of type [github.com/iolave/go-errors.GenericError] will be returned.
func encodeAppRequest(url string, v any) (newUrl string, r io.Reader, err error) {
	for i := range reflect.ValueOf(v).Elem().NumField() {
		f := reflect.ValueOf(v).Elem().Field(i)
		t := reflect.TypeOf(v).Elem().Field(i)
		fullTag := t.Tag.Get("ar")
		tags := strings.SplitSeq(fullTag, ",")
		for tag := range tags {
			splittedTag := strings.Split(tag, "=")
			k := splittedTag[0]
			switch k {
			case "path":
				if len(splittedTag) != 2 {
					return "", nil, errors.NewWithName(
						ERR_NAME_ENCODER,
						fmt.Sprintf(ERR_ENCDEC_PATH_NO_VAL, t.Name),
					)
				}
				if t.Type.Kind() != reflect.String {
					return "", nil, errors.NewWithName(
						ERR_NAME_ENCODER,
						fmt.Sprintf(ERR_ENCDEC_PATH_INVALID, t.Name),
					)
				}
				v := splittedTag[1]
				newUrl = strings.ReplaceAll(url, fmt.Sprintf("{%s}", v), f.String())
			case "body":
				if len(splittedTag) != 2 {
					return "", nil, errors.NewWithName(
						ERR_NAME_ENCODER,
						fmt.Sprintf(ERR_ENCDEC_BODY_NO_VAL, t.Name),
					)
				}
				switch typ := splittedTag[1]; typ {
				case "json":
					v := f.Interface()
					b, err := json.Marshal(v)
					if err != nil {
						return "", nil, errors.NewWithNameAndErr(
							ERR_NAME_ENCODER,
							fmt.Sprintf(ERR_ENCDEV_BODY_INVALID, t.Name),
							err,
						)
					}

					r = bytes.NewReader(b)
				default:
					return "", nil, errors.NewWithName(
						ERR_NAME_ENCODER,
						fmt.Sprintf(ERR_ENCDEC_TAG_VAL_INVALID, "body", typ, t.Name),
					)
				}
			default:
				return "", nil, errors.NewWithName(
					ERR_NAME_ENCODER,
					fmt.Sprintf(ERR_ENCDEC_INVALID_TAG, k, t.Name),
				)
			}
		}

	}

	return newUrl, r, nil
}

// decodeAppRequest takes a request and a struct v with/wo "ar" tags.
//
//   - It detects path ar tags ("path=*") from v, retrieves the [http.Request] path values and stores them in the corresponding v properties.
//   - It detects body ar tags ("body=json") from v, decodes the request body and stores the decoded value in the corresponding v property.
//
// If some error occurs, an error of type [github.com/iolave/go-errors.GenericError] will be returned.
func decodeAppRequest(r *http.Request, v any) error {
	// check if v is a pointer, otherwise return an error
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return errors.NewWithName(
			ERR_NAME_DECODER,
			fmt.Sprintf(ERR_ENC_DEC_EXPECT_PTR, reflect.ValueOf(v).Kind().String()),
		)
	}

	// check if v is a pointer to a struct, otherwise return an error
	if reflect.ValueOf(v).Elem().Kind() != reflect.Struct {
		return errors.NewWithName(
			ERR_NAME_DECODER,
			fmt.Sprintf(ERR_ENC_DEC_EXPECT_PTR, reflect.ValueOf(v).Kind().String()),
		)
	}

	for i := range reflect.ValueOf(v).Elem().NumField() {
		f := reflect.ValueOf(v).Elem().Field(i)
		t := reflect.TypeOf(v).Elem().Field(i)
		fullTag := t.Tag.Get("ar")
		tags := strings.SplitSeq(fullTag, ",")
		for tag := range tags {
			splittedTag := strings.Split(tag, "=")
			k := splittedTag[0]
			switch k {
			case "path":
				if len(splittedTag) != 2 {
					return errors.NewWithName(
						ERR_NAME_DECODER,
						ERR_ENCDEC_PATH_NO_VAL,
					)
				}
				if t.Type.Kind() != reflect.String {
					return errors.NewWithName(
						ERR_NAME_DECODER,
						fmt.Sprintf(ERR_ENCDEC_PATH_INVALID, t.Name),
					)
				}
				v := splittedTag[1]
				vv := r.PathValue(v)
				f.SetString(vv)
			case "body":
				if len(splittedTag) != 2 {
					return errors.NewWithName(
						ERR_NAME_DECODER,
						fmt.Sprintf(ERR_ENCDEC_BODY_NO_VAL, t.Name),
					)
				}
				switch typ := splittedTag[1]; typ {
				case "json":
					if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
						return errors.NewBadRequestError(
							"invalid content type",
							nil,
						)
					}
					b, err := io.ReadAll(r.Body)
					if err != nil {
						return errors.Wrap(err)
					}

					addr := f.Addr()
					in := addr.Interface()

					err = json.Unmarshal(b, in)
					if err != nil {
						return errors.NewBadRequestError(
							ERR_ENCDEC_PARSE,
							err,
						)
					}

					f.Set(reflect.ValueOf(in).Elem())
				default:
					return errors.NewWithName(
						ERR_NAME_DECODER,
						fmt.Sprintf(ERR_ENCDEC_TAG_VAL_INVALID, "body", typ, t.Name),
					)
				}

			default:
				return errors.NewWithName(
					ERR_NAME_DECODER,
					fmt.Sprintf(ERR_ENCDEC_INVALID_TAG, k, t.Name),
				)
			}
		}

	}

	return nil
}

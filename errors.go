package betsi

// App related error codes
const (
	ERR_NAME               = "app_error"
	ERR_NIL_LOGGER         = "logger cannot be nil"
	ERR_START_W_NIL_SERVER = "server config not provided"
	ERR_START_W_NIL_ROUTER = "router cannot be nil"
	ERR_INVALID_TYPE       = "expected type %s, got %s"
)

// Server related error codes
const (
	ERR_SRV_AR_NIL_ERR                  = "nil error passed"
	ERR_SRV_AR_GENERIC_ERR              = "Internal server error"
	ERR_SRV_AR_SEND_JSON_VALIDATION_ERR = "failed to send response, response doesn't meet validation requirements"
	ERR_SRV_AR_SEND_JSON_MARSHALL_ERR   = "failed to send response, unable to marshal response"
)

// AR encoder/decoder related error codes
const (
	ERR_NAME_ENCODER           = "app_request_encoder_error"
	ERR_NAME_DECODER           = "app_request_decoder_error"
	ERR_ENCDEC_BODY_NO_VAL     = "body tag expected a value (name:%s,oneof:json)"
	ERR_ENCDEV_BODY_INVALID    = "failed to encode body (name:%s)"
	ERR_ENCDEC_PATH_NO_VAL     = "path tag expected a value (name:%s)"
	ERR_ENCDEC_PATH_INVALID    = "path tag can only be used with string (name:%s)"
	ERR_ENCDEC_TAG_VAL_INVALID = "%s tag value %s not supported (name:%s)"
	ERR_ENCDEC_INVALID_TAG     = "tag %s is not supported (name:%s)"
	ERR_ENC_DEC_EXPECT_PTR     = "v has to be a pointer to a struct, got %s"
)

// AR errors
const (
	ERR_NAME_PARSE      = "parse_request_error"
	ERR_AR_NIL_REQ      = "ar.Req is nil"
	ERR_AR_NIL_REQ_BODY = "ar.Req.Body is nil"
)

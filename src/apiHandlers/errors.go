package apiHandlers

// error messages
const (
	ErrMalformedJWT               = "Missing or malformed JWT"
	ErrInvalidJWT                 = "Invalid or expired JWT"
	ErrInvalidBodyData            = "Invalid body data"
	ErrInvalidUsernamePass        = "Invalid username or password"
	ErrInvalidAuth                = "Invalid auth provided, please relogin"
	ErrPermissionDenied           = "Permission denied"
	ErrInvalidInputPass           = "Invalid input for password entered"
	ErrUsernameExists             = "Username already exists"
	ErrInternalServerError        = "Internal Server Error happened. Please try again later."
	ErrInvalidFileData            = "Invalid file data provided"
	ErrInvalidPhoneNumber         = "Invalid phone number provided: %s"
	ErrPhoneNumberAlreadyImported = "Phone number already imported: %s"
	ErrInvalidUsername            = "Invalid username provided"
	ErrNoPhonesDonated            = "No phones donated by this account"
	ErrAgentNotConnected          = "Agent is not connected; try again later"
	ErrInvalidPagination          = "Invalid pagination parameters provided"
	ErrMaxContactImportLimit      = "Maximum contact import limit reached"
	ErrPhoneNumberNotFound        = "Phone number not found: %s"
	ErrParameterRequired          = "Parameter required but not provided: %s"
	ErrUserBanned                 = "User is banned"
	ErrLabelInfoNotFound          = "Label info with specified id '%d' not found"
	ErrLabelAlreadyApplied        = "Label with id '%d' already applied to phone number: %s"
	ErrLabelAlreadyExistsByName   = "Label with name '%s' already exists"
	ErrTooManyChatLabelInfo       = "Too many chat label info"
	ErrLabelNameTooLong           = "Label name is too long; max length is %d characters"
	ErrLabelDescriptionTooLong    = "Label description is too long; max length is %d characters"
	ErrInvalidColor               = "Invalid color provided: %s"
	ErrLabelNotApplied            = "Label with id '%d' not applied to phone number: %s"
	ErrCannotDeleteBuiltInLabel   = "Cannot delete built-in label: '%d'"
	ErrDuplicatePhoneNumber       = "Duplicate phone number: '%s'"
	ErrPhoneNotWorking            = "None of the phone numbers are working"
	ErrInvalidPmsPass             = "Invalid PMS password provided"
	ErrInvalidAgentId             = "Invalid agent id provided: %d"
	ErrInvalidAppSettingName      = "Invalid app setting name provided: %s"
	ErrAppSettingNotFound         = "App setting with name '%s' not found"
	ErrTextEmpty                  = "Text is empty"
	ErrTextTooLong                = "The provided text is too long"
	ErrInvalidClientRId           = "Invalid client rId provided: %s"
	ErrInvalidCaptcha             = "Invalid captcha id/answer provided"
)

// error codes
const (
	ErrCodeMalformedJWT               = 2100
	ErrCodeInvalidJWT                 = 2101
	ErrCodeInvalidBodyData            = 2102
	ErrCodeInvalidUsernamePass        = 2103
	ErrCodeInvalidAuth                = 2104
	ErrCodePermissionDenied           = 2105
	ErrCodeInvalidInputPass           = 2106
	ErrCodeUsernameExists             = 2107
	ErrCodeInternalServerError        = 2108
	ErrCodeInvalidFileData            = 2109
	ErrCodeInvalidPhoneNumber         = 2110
	ErrCodePhoneNumberAlreadyImported = 2111
	ErrCodeInvalidUsername            = 2112
	ErrCodeNoPhonesDonated            = 2113
	ErrCodeAgentNotConnected          = 2114
	ErrCodeInvalidPagination          = 2115
	ErrCodeMaxContactImportLimit      = 2116
	ErrCodePhoneNumberNotFound        = 2117
	ErrCodeParameterRequired          = 2118
	ErrCodeUserBanned                 = 2119
	ErrCodeLabelInfoNotFound          = 2120
	ErrCodeLabelAlreadyApplied        = 2121
	ErrCodeLabelAlreadyExistsByName   = 2122
	ErrCodeTooManyChatLabelInfo       = 2123
	ErrCodeLabelNameTooLong           = 2124
	ErrCodeLabelDescriptionTooLong    = 2125
	ErrCodeInvalidColor               = 2126
	ErrCodeLabelNotApplied            = 2127
	ErrCodeCannotDeleteBuiltInLabel   = 2128
	ErrCodeDuplicatePhoneNumber       = 2129
	ErrCodePhoneNotWorking            = 2130
	ErrCodeInvalidPmsPass             = 2131
	ErrCodeInvalidAgentId             = 2132
	ErrCodeInvalidAppSettingName      = 2133
	ErrCodeAppSettingNotFound         = 2134
	ErrCodeTextEmpty                  = 2135
	ErrCodeTextTooLong                = 2136
	ErrCodeInvalidClientRId           = 2137
	ErrCodeInvalidCaptcha             = 2138
)

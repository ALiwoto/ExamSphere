package apiHandlers

// error messages
const (
	ErrMalformedJWT                  = "Missing or malformed JWT"
	ErrInvalidJWT                    = "Invalid or expired JWT"
	ErrInvalidBodyData               = "Invalid body data"
	ErrInvalidUsernamePass           = "Invalid username or password"
	ErrInvalidAuth                   = "Invalid auth provided, please relogin"
	ErrPermissionDenied              = "Permission denied"
	ErrInvalidInputPass              = "Invalid input for password entered"
	ErrUsernameExists                = "Username already exists"
	ErrInternalServerError           = "Internal Server Error happened. Please try again later."
	ErrInvalidFileData               = "Invalid file data provided"
	ErrInvalidPhoneNumber            = "Invalid phone number provided: %s"
	ErrPhoneNumberAlreadyImported    = "Phone number already imported: %s"
	ErrInvalidUsername               = "Invalid username/userId provided"
	ErrNoPhonesDonated               = "No phones donated by this account"
	ErrAgentNotConnected             = "Agent is not connected; try again later"
	ErrInvalidPagination             = "Invalid pagination parameters provided"
	ErrMaxContactImportLimit         = "Maximum contact import limit reached"
	ErrPhoneNumberNotFound           = "Phone number not found: %s"
	ErrParameterRequired             = "Parameter required but not provided: %s"
	ErrUserBanned                    = "User is banned"
	ErrLabelInfoNotFound             = "Label info with specified id '%d' not found"
	ErrLabelAlreadyApplied           = "Label with id '%d' already applied to phone number: %s"
	ErrLabelAlreadyExistsByName      = "Label with name '%s' already exists"
	ErrTooManyChatLabelInfo          = "Too many chat label info"
	ErrLabelNameTooLong              = "Label name is too long; max length is %d characters"
	ErrLabelDescriptionTooLong       = "Label description is too long; max length is %d characters"
	ErrInvalidColor                  = "Invalid color provided: %s"
	ErrLabelNotApplied               = "Label with id '%d' not applied to phone number: %s"
	ErrCannotDeleteBuiltInLabel      = "Cannot delete built-in label: '%d'"
	ErrDuplicatePhoneNumber          = "Duplicate phone number: '%s'"
	ErrPhoneNotWorking               = "None of the phone numbers are working"
	ErrInvalidPmsPass                = "Invalid PMS password provided"
	ErrInvalidAgentId                = "Invalid agent id provided: %d"
	ErrInvalidAppSettingName         = "Invalid app setting name provided: %s"
	ErrAppSettingNotFound            = "App setting with name '%s' not found"
	ErrTextEmpty                     = "Text is empty"
	ErrTextTooLong                   = "The provided text is too long"
	ErrInvalidClientRId              = "Invalid client rId provided: %s"
	ErrInvalidCaptcha                = "Invalid captcha id/answer provided"
	ErrQueryParameterNotProvided     = "Query parameter required but not provided: %s"
	ErrTooManyPasswordChangeAttempts = "Too many password change attempts. Please try again later"
	ErrRequestExpired                = "Request expired. Please try again later"
)

// error codes
const (
	ErrCodeMalformedJWT APIErrorCode = 2100 + iota
	ErrCodeInvalidJWT
	ErrCodeInvalidBodyData
	ErrCodeInvalidUsernamePass
	ErrCodeInvalidAuth
	ErrCodePermissionDenied
	ErrCodeInvalidInputPass
	ErrCodeUsernameExists
	ErrCodeInternalServerError
	ErrCodeInvalidFileData
	ErrCodeInvalidPhoneNumber
	ErrCodePhoneNumberAlreadyImported
	ErrCodeInvalidUsername
	ErrCodeNoPhonesDonated
	ErrCodeAgentNotConnected
	ErrCodeInvalidPagination
	ErrCodeMaxContactImportLimit
	ErrCodePhoneNumberNotFound
	ErrCodeParameterRequired
	ErrCodeUserBanned
	ErrCodeLabelInfoNotFound
	ErrCodeLabelAlreadyApplied
	ErrCodeLabelAlreadyExistsByName
	ErrCodeTooManyChatLabelInfo
	ErrCodeLabelNameTooLong
	ErrCodeLabelDescriptionTooLong
	ErrCodeInvalidColor
	ErrCodeLabelNotApplied
	ErrCodeCannotDeleteBuiltInLabel
	ErrCodeDuplicatePhoneNumber
	ErrCodePhoneNotWorking
	ErrCodeInvalidPmsPass
	ErrCodeInvalidAgentId
	ErrCodeInvalidAppSettingName
	ErrCodeAppSettingNotFound
	ErrCodeTextEmpty
	ErrCodeTextTooLong
	ErrCodeInvalidClientRId
	ErrCodeInvalidCaptcha
	ErrCodeQueryParameterNotProvided
	ErrCodeTooManyPasswordChangeAttempts
	ErrCodeRequestExpired
)

package core

const (
	ERROR_UNABLE_TO_OPEN_SRC_FILE = "unable_to_open_file"
	ERROR_EMPTY_SRC_FILE          = "empty_source_file"
	ERROR_INVALID_VERSION_FORMAT  = "invalid_version_format"
)

type ResFileError struct {
	version, tp, code string
	parent            error
}

func (e ResFileError) Error() string {
	return e.code
}

func (e ResFileError) Message() string {
	if e.code == ERROR_UNABLE_TO_OPEN_SRC_FILE {
		return "Unable to read source file " + e.version + "." + e.tp + ".sql, error: " + e.parent.Error()
	} else if e.code == ERROR_EMPTY_SRC_FILE {
		return "Source file " + e.version + "." + e.tp + ".sql is empty"
	} else if e.code == ERROR_INVALID_VERSION_FORMAT {
		return "Invalid version format for version. Error: " + e.parent.Error()
	}
	return ""
}

func NewSrcFileError(code, version, tp string, parent error) ResFileError {
	return ResFileError{version, tp, code, parent}
}

package file_parser

type NginxLogParser struct {
	fp *FileParser
}

func (np *NginxLogParser) SetFormatString(s string) error {
	panic("")
}

func (np *NginxLogParser) Parse() map[string]string {
	panic("")
}

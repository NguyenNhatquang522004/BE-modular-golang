package enum

// =============================================================================
// FILE TYPE (Extension Category)
// =============================================================================

//go:generate enumer -type=FileType -json -transform=lower -trimprefix=FileType
type FileType int

const (
	FileTypeOther FileType = iota // 'other' (Các loại lạ)
	FileTypePDF                   // 'pdf'
	FileTypeDOCX                  // 'docx' (Word)
	FileTypeXLSX                  // 'xlsx' (Excel)
	FileTypePPTX                  // 'pptx' (PowerPoint)
	FileTypeZIP                   // 'zip'  (Archive)
	FileTypeRAR                   // 'rar'
	FileTypeTXT                   // 'txt'
	FileTypeCSV                   // 'csv'
)
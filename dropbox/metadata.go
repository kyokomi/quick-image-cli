package dropbox

type metadata struct {
	Bytes    float64 `json:"bytes"`
	Contents []struct {
		Bytes       float64 `json:"bytes"`
		ClientMtime string  `json:"client_mtime"`
		Icon        string  `json:"icon"`
		IsDir       bool    `json:"is_dir"`
		MimeType    string  `json:"mime_type"`
		Modified    string  `json:"modified"`
		Path        string  `json:"path"`
		Rev         string  `json:"rev"`
		Revision    float64 `json:"revision"`
		Root        string  `json:"root"`
		Size        string  `json:"size"`
		ThumbExists bool    `json:"thumb_exists"`
	} `json:"contents"`
	Hash        string  `json:"hash"`
	Icon        string  `json:"icon"`
	IsDir       bool    `json:"is_dir"`
	Modified    string  `json:"modified"`
	Path        string  `json:"path"`
	Rev         string  `json:"rev"`
	Revision    float64 `json:"revision"`
	Root        string  `json:"root"`
	Size        string  `json:"size"`
	ThumbExists bool    `json:"thumb_exists"`
}

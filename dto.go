package nhentai

type PageData struct {
	PageCount int `json:"page_count"`
}

type ComicSimple struct {
	Id          int    `json:"id"`
	MediaId     int    `json:"media_id"`
	Title       string `json:"title"`
	TagIds      []int  `json:"tag_ids"`
	Lang        string `json:"lang"`
	Thumb       string `json:"thumb"`
	ThumbWidth  int    `json:"thumb_width"`
	ThumbHeight int    `json:"thumb_height"`
}

type ComicPageData struct {
	PageData
	Records []ComicSimple `json:"records"`
}

type TagPageData struct {
	PageData
	Records []TagPageTag `json:"records"`
}

type TagPageTag struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Count string `json:"count"` // like "7k"
}

type ComicInfo struct {
	Id           int             `json:"id"`
	MediaId      string          `json:"media_id"`
	Title        ComicInfoTitle  `json:"title"`
	Images       ComicInfoImages `json:"images"`
	Scanlator    string          `json:"scanlator"`
	UploadDate   int             `json:"upload_date"`
	Tags         []ComicInfoTag  `json:"tags"`
	NumPages     int             `json:"num_pages"`
	NumFavorites int             `json:"num_favorites"`
}

type ComicInfoTitle struct {
	English  string `json:"english"`
	Japanese string `json:"japanese"`
	Pretty   string `json:"pretty"`
}

type ComicInfoImages struct {
	Pages     []ImageInfo `json:"pages"`
	Cover     ImageInfo   `json:"cover"`
	Thumbnail ImageInfo   `json:"thumbnail"`
}

type ComicInfoTag struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
	Type  string `json:"type"`
	Url   string `json:"url"`
}

type ImageInfo struct {
	// T type "j" -> jpg
	T string `json:"t"`
	// W width
	W int `json:"w"`
	// H height
	H int `json:"h"`
}

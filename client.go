package nhentai

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Client nHentai客户端
type Client struct {
	// http.Client 继承HTTP客户端
	http.Client
	// Mirror 分流
	Mirror string
}

// GetMirror 获取当前的镜像网站
func (c *Client) GetMirror() string {
	if c.Mirror == "" {
		return MirrorOrigin
	}
	return c.Mirror
}

// Comics 列出漫画
// https://nhentai.net/?page=1
func (c *Client) Comics(page int) (*ComicPageData, error) {
	urlStr := fmt.Sprintf("https://%s/?page=%d", c.GetMirror(), page)
	return c.parsePage(urlStr)
}

// ComicsByTag 列出标签下的漫画
// https://nhentai.net/tag/group/?page=1
func (c *Client) ComicsByTag(tag string, page int) (*ComicPageData, error) {
	urlStr := fmt.Sprintf("https://%s/tag/%s/?page=%d", c.GetMirror(), tag, page)
	return c.parsePage(urlStr)
}

// parsePage 获取页面上的漫画列表
func (c *Client) parsePage(urlStr string) (*ComicPageData, error) {
	doc, err := c.parseUrlToDoc(urlStr)
	if err != nil {
		return nil, err
	}
	var divSelection *goquery.Selection
	doc.Find(".container.index-container:not(.index-popular)").Each(func(i int, selection *goquery.Selection) {
		divSelection = selection
	})
	if divSelection == nil {
		return nil, errors.New("NOT MATCH CONTAINER")
	}
	gallerySelection := divSelection.Find("div.gallery")
	galleries := make([]ComicSimple, gallerySelection.Size())
	gallerySelection.Each(func(i int, selection *goquery.Selection) {
		idStr, _ := selection.Find("a").First().Attr("href")
		idStr = strings.TrimPrefix(idStr, "/g/")
		idStr = strings.TrimSuffix(idStr, "/")
		id, _ := strconv.Atoi(idStr)
		title := selection.Find(".caption").Text()
		thumb, thumbWidth, thumbHeight, mediaId := c.parseCover(selection)
		tagIdsStr, _ := selection.Attr("data-tags")
		tsp := strings.Split(tagIdsStr, " ")
		tagIds := make([]int, len(tsp))
		for i2 := range tsp {
			tagIds[i2], _ = strconv.Atoi(tsp[i2])
		}
		lang := lang(tagIds)
		galleries[i] = ComicSimple{
			Id:          id,
			Title:       title,
			MediaId:     mediaId,
			TagIds:      tagIds,
			Lang:        lang,
			Thumb:       thumb,
			ThumbWidth:  thumbWidth,
			ThumbHeight: thumbHeight,
		}
	})
	lastPage := c.parseLastPage(doc)
	return &ComicPageData{
		PageData: PageData{
			PageCount: lastPage,
		},
		Records: galleries,
	}, nil
}

// parseCover 分析媒体信息
func (c *Client) parseCover(selection *goquery.Selection) (string, int, int, int) {
	lazyload := selection.Find(".lazyload")
	thumb, _ := lazyload.Attr("data-src")
	thumbWidthStr, _ := lazyload.Attr("width")
	thumbHeightStr, _ := lazyload.Attr("height")
	width, _ := strconv.Atoi(thumbWidthStr)
	thumbHeight, _ := strconv.Atoi(thumbHeightStr)
	mediaIdStr := thumb[strings.Index(thumb, "galleries")+10 : strings.LastIndex(thumb, "/")]
	mediaId, _ := strconv.Atoi(mediaIdStr)
	return thumb, width, thumbHeight, mediaId
}

// ComicInfo 获取漫画的信息
func (c *Client) ComicInfo(id int) (*ComicInfo, error) {
	urlStr := fmt.Sprintf("https://%s/api/gallery/%d", c.GetMirror(), id)
	rsp, err := c.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	buff, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	var comicInfo ComicInfo
	err = json.Unmarshal(buff, &comicInfo)
	if err != nil {
		return nil, err
	}
	// []Tag
	return &comicInfo, nil
}

// Tags 获取标签
// https://nhentai.net/tags/?page=1
func (c *Client) Tags(page int) (*TagPageData, error) {
	urlStr := fmt.Sprintf("https://%s/tags/?page=%d", c.GetMirror(), page)
	doc, err := c.parseUrlToDoc(urlStr)
	if err != nil {
		return nil, err
	}
	tags := c.parseTags(doc.Find("div.container#tag-container>section>a"))
	lastPage := c.parseLastPage(doc)
	return &TagPageData{
		PageData: PageData{
			PageCount: lastPage,
		},
		Records: tags,
	}, nil
}

// parseTags 解析标签数据
func (c *Client) parseTags(tagSelections *goquery.Selection) []TagPageTag {
	tags := make([]TagPageTag, tagSelections.Size())
	tagSelections.Each(func(i int, selection *goquery.Selection) {
		aClass, _ := selection.Attr("class")
		aClass = strings.TrimPrefix(aClass, "tag tag-")
		aClass = strings.TrimSpace(aClass)
		id, _ := strconv.Atoi(aClass)
		name := selection.Find(".name").Text()
		count := selection.Find(".count").Text()
		tags[i] = TagPageTag{
			Id:    id,
			Name:  name,
			Count: count,
		}
	})
	return tags
}

// parseLastPage 获取一共多少页
func (c *Client) parseLastPage(doc *goquery.Document) int {
	lastPageHref, _ := doc.Find(".pagination>.last").Attr("href")
	lastPageHref = lastPageHref[strings.Index(lastPageHref, "page=")+5:]
	lastPage, _ := strconv.Atoi(lastPageHref)
	return lastPage
}

// parseUrlToDoc 从网址读取网页并且转换成document
func (c *Client) parseUrlToDoc(str string) (*goquery.Document, error) {
	rsp, err := c.Get(str)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	return goquery.NewDocumentFromReader(rsp.Body)
}

// CoverUrl 拼接封面的URL
// "https://t.nhentai.net/galleries/{media_id}/cover{cover_ext}"
func (c *Client) CoverUrl(mediaId int, t string) string {
	return fmt.Sprintf("https://t.%s/galleries/%d/cover.%s", c.GetMirror(), mediaId, c.GetExtension(t))
}

// ThumbnailUrl 拼接缩略图的URL
// "https://t.nhentai.net/galleries/{media_id}/thumbnail{thumbnail_ext}"
func (c *Client) ThumbnailUrl(mediaId int, t string) string {
	return fmt.Sprintf("https://t.%s/galleries/%d/thumbnail.%s", c.GetMirror(), mediaId, c.GetExtension(t))
}

// PageUrl
// https://i.nhentai.net/galleries/{media_id}/{num}{extension}
// {num} is {index + 1} (begin is 1)
func (c *Client) PageUrl(mediaId int, num int, t string) string {
	return fmt.Sprintf("https://i.%s/galleries/%d/%d.%s", c.GetMirror(), mediaId, num, c.GetExtension(t))
}

// GetExtension 使用type获得拓展名
func (c *Client) GetExtension(t string) string {
	// Official only j
	if t == "j" {
		return "jpg"
	}
	// redundancy
	if t == "p" {
		return "png"
	}
	return ""
}

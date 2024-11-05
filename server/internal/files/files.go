package files

import (
	"bytes"
	"cmp"
	"html/template"
	"os"
	"path/filepath"
	"quicknotes/model"
	"regexp"
	"slices"
	"strings"
	"time"

	markdown "github.com/teekennedy/goldmark-markdown"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/text"
)

const filesDir = "/data"

type FileHandler struct{}

func (h FileHandler) GetNote(path string) (*model.Note, error) {
	renderer := markdown.NewRenderer(markdown.WithHeadingStyle(markdown.HeadingStyleATX))
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
		goldmark.WithRenderer(renderer),
	)

	fullPath := filepath.Join(filesDir, path) + ".md"
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	document := md.Parser().Parse(text.NewReader(content))
	metadata := document.OwnerDocument().Meta()

	lastModifiedStr := metadata["lastModified"].(string)
	lastModified, err := time.Parse(time.DateTime, lastModifiedStr)
	if err != nil {
		return nil, err
	}

	mTags := metadata["tags"].([]any)
	tags := []string{}
	for _, t := range mTags {
		tags = append(tags, t.(string))
	}

	var b bytes.Buffer
	err = md.Renderer().Render(&b, content, document)
	if err != nil {
		return nil, err
	}

	return &model.Note{
		Title:        metadata["title"].(string),
		Content:      b.String(),
		Tags:         tags,
		LastModified: lastModified,
	}, nil
}

func (h FileHandler) CreateNote(input model.AddNoteInput) error {
	var tmplFile = "templates/note.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		return err
	}
	cleanTitle := strings.TrimSpace(input.Title)
	re := regexp.MustCompile("[^0-9a-zA-Z._-]")
	cleanTitle = re.ReplaceAllString(cleanTitle, "")
	f, err := os.Create(filepath.Join(filesDir, cleanTitle) + ".md")
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, map[string]any{
		"Title":        input.Title,
		"Content":      input.Content,
		"Tags":         input.Tags,
		"LastModified": time.Now().Format(time.DateTime),
	})
	if err != nil {
		return err
	}
	return nil
}

func (h FileHandler) UpdateNote(path string, input model.PatchNoteInput) error {
	var tmplFile = "templates/note.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		return err
	}
	fullPath := filepath.Join(filesDir, path)
	if _, err := os.Stat(fullPath); err != nil {
		return err
	}
	f, err := os.Open(fullPath + ".md")
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, map[string]any{
		"Title":        input.Title,
		"Content":      input.Content,
		"Tags":         input.Tags,
		"LastModified": time.Now().Format(time.DateTime),
	})
	if err != nil {
		return err
	}
	return nil
}

func (h FileHandler) SearchNotes(term string, tag string, sort string, order model.OrderType, limit int) ([]model.PartialNote, error) {
	filePaths, err := h.listFiles()
	if err != nil {
		return nil, err
	}

	notes := []model.PartialNote{}
	i := 0
	for _, path := range filePaths {
		if i > limit {
			break
		}
		note, err := h.getNoteMetadata(path)
		if err != nil {
			return nil, err
		}
		if term != "" && !strings.Contains(strings.ToLower(note.Title), strings.ToLower(term)) {
			continue
		}
		if tag != "" {
			matchTag := false
			for _, noteTag := range note.Tags {
				if noteTag == tag {
					matchTag = true
				}
			}
			if !matchTag {
				continue
			}
		}

		notes = append(notes, *note)
	}

	slices.SortFunc(notes, func(a, b model.PartialNote) int {
		v := 0
		switch sort {
		case "title":
			v = cmp.Compare(a.Title, b.Title)
		case "lastModified":
			v = a.LastModified.Compare(b.LastModified)
		}
		if order == model.AscOrderType {
			return v
		}
		return v * -1
	})

	return notes, nil
}

func (h FileHandler) GetTags() ([]string, error) {
	filePaths, err := h.listFiles()
	if err != nil {
		return nil, err
	}

	tags := []string{}
	tagMap := map[string]int{}
	for _, path := range filePaths {
		note, err := h.getNoteMetadata(path)
		if err != nil {
			return nil, err
		}
		for _, tag := range note.Tags {
			if _, ok := tagMap[tag]; !ok {
				tags = append(tags, tag)
				tagMap[tag] = 1
			}
		}
	}

	slices.Sort(tags)

	return tags, nil
}

func (h FileHandler) getNoteMetadata(file string) (*model.PartialNote, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.New(
				meta.WithStoresInDocument(),
			),
		),
	)
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	document := md.Parser().Parse(text.NewReader(content))
	metadata := document.OwnerDocument().Meta()

	lastModifiedStr := metadata["lastModified"].(string)
	lastModified, err := time.Parse(time.DateTime, lastModifiedStr)
	if err != nil {
		return nil, err
	}

	mTags := metadata["tags"].([]any)
	tags := []string{}
	for _, t := range mTags {
		tags = append(tags, t.(string))
	}

	return &model.PartialNote{
		Title:        metadata["title"].(string),
		Tags:         tags,
		LastModified: lastModified,
	}, nil
}

func (h FileHandler) listFiles() ([]string, error) {
	files := []string{}
	entries, err := os.ReadDir(filesDir)
	for _, entry := range entries {
		if err != nil {
			return nil, err
		}
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			files = append(files, entry.Name())
		}
	}
	if err != nil {
		return nil, err
	}
	return files, nil
}

package cmd

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type MDFile struct {
	Name     string
	Headings [][]string
}

func (mdf *MDFile) String() string {
	return fmt.Sprintf("{Name: %q, Headings: %q}", mdf.Name, mdf.Headings)
}

type MDDir struct {
	Path    string
	Name    string
	MDFiles []*MDFile
	SubDir  []*MDDir
}

func (mdd *MDDir) String() string {
	return fmt.Sprintf("{Path: %q, Name: %q, MDFiles: %v, SubDir: %v}", mdd.Path, mdd.Name, mdd.MDFiles, mdd.SubDir)
}

func GenerateTOCFile(root string) {
	start := time.Now()
	rootNode := &MDDir{Path: root}
	skipMap := GenerateSkipMap(skipDirs)
	err := WalkMDDir(rootNode, skipMap)
	check(err)
	err = WriteReadme(rootNode)
	check(err)
	end := time.Now()
	delta := end.Sub(start)
	log.Printf("Processing Time: %v", delta)
}

func GenerateSkipMap(skipDirs []string) map[string]struct{} {
	skipMap := make(map[string]struct{})
	for _, v := range skipDirs {
		skipMap[v] = struct{}{}
	}
	return skipMap
}

func WalkMDDir(root *MDDir, skip map[string]struct{}) error {
	headingMatcher := regexp.MustCompile(`(?m)^(\s*#\s+)(.+[^\r\n])`)
	frontmatterMatcher := regexp.MustCompile(`(?s)(---)(.*?)(---)`)

	f, err := os.Open(root.Path)
	if err != nil {
		return err
	}

	dirs, err := f.ReadDir(-1)
	for _, dir := range dirs {
		dirName := dir.Name()
		if strings.HasPrefix(dir.Name(), ".") {
			continue
		}
		if _, ok := skip[dir.Name()]; ok {
			continue
		}
		if dir.Name() == "README.md" {
			continue
		}
		subPath := filepath.Join(root.Path, dirName)
		if dir.IsDir() {
			subDir := &MDDir{Path: subPath, Name: dirName}
			err := WalkMDDir(subDir, skip)
			if err != nil {
				return err
			}
			root.SubDir = append(root.SubDir, subDir)
		}
		if !dir.IsDir() && strings.HasSuffix(dirName, ".md") {
			content, err := os.ReadFile(subPath)
			if err != nil {
				return err
			}

			contentStr := string(content)
			headings := make([][]string, 6)

			// read title from frontmatter
			// frontmatter start from fist line of file
			if len(contentStr) > 0 && contentStr[0:3] == "---" {
				matchesF := frontmatterMatcher.FindAllStringSubmatch(contentStr, -1)
				if len(matchesF) > 0 && len(matchesF[0]) > 0 {
					frontmatter := matchesF[0][2]
					//log.Printf(">>> File: %q,Frontmatter: %q", subPath, frontmatter)
					m := map[string]interface{}{}
					err := yaml.Unmarshal([]byte(frontmatter), &m)
					if err != nil {
						return err
					}
					if t, ok := m[titleField]; ok {
						title := t.(string)
						headings[0] = append(headings[0], title)
					}
				}
			}

			// read markdown Lv heading
			matchesH := headingMatcher.FindAllStringSubmatch(contentStr, -1)
			for i := range matchesH {
				headings[0] = append(headings[0], matchesH[i][2])
			}
			mdFile := &MDFile{Name: dirName, Headings: headings}
			root.MDFiles = append(root.MDFiles, mdFile)
		}

	}

	return nil
}

func WriteReadme(root *MDDir) error {
	//fmt.Printf(">>> Current Path: %s, Name: %s\n", root.Path, root.Name)
	tocFilename := "README.md"
	readmePath := filepath.Join(root.Path, tocFilename)

	f, err := os.OpenFile(readmePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	var title string
	title = root.Name
	if root.Name == "" {
		title = "README"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n", title))
	err = WriteTOC(root, root, &sb, 1)
	if err != nil {
		return err
	}

	_, err = f.WriteString(sb.String())
	if err != nil {
		return err
	}
	return nil
}

func WriteTOC(root *MDDir, currentDir *MDDir, sb *strings.Builder, depth int) error {
	for _, mdir := range currentDir.SubDir {
		relativePath := strings.TrimPrefix(mdir.Path, root.Path)
		relativePath = strings.TrimPrefix(relativePath, string(os.PathSeparator))
		//fmt.Printf(">>>>>> Root: %q, Current: %q, Relative Path: %q\n", root.Path, mdir.Path, relativePath)
		sb.WriteString(fmt.Sprintf("%s- [%s](%s)\n", strings.Repeat(" ", depth*2), mdir.Name, relativePath))
		err := WriteTOC(root, mdir, sb, depth+1)
		if err != nil {
			return err
		}

		if recursive {
			err = WriteReadme(mdir)
			if err != nil {
				return err
			}
		}

	}

	for _, mf := range currentDir.MDFiles {
		if len(mf.Headings) == 0 || len(mf.Headings[0]) == 0 {
			log.Printf("Warning: No Frontmatter title or Lv1 Heading, %s\n", filepath.Join(currentDir.Path, mf.Name))
			continue
		}
		relativePath := strings.TrimPrefix(currentDir.Path, root.Path)
		relativePath = strings.TrimPrefix(relativePath, string(os.PathSeparator))
		sb.WriteString(fmt.Sprintf("%s - [%s](%s)\n", strings.Repeat(" ", depth*2), mf.Headings[0][0], filepath.Join(relativePath, mf.Name)))
	}
	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

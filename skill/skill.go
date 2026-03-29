package skill

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type meta map[string]string

type skill struct {
	name         string
	descriptions string
	meta         meta
	body         string
	path         string
}

type SkillLoader struct {
	skillsDir string
	skills    map[string]skill
}

func NewSkillLoader(skillsDir string) *SkillLoader {
	s := &SkillLoader{
		skillsDir: skillsDir,
		skills:    make(map[string]skill),
	}
	s.loadAll()
	return s
}

func (s *SkillLoader) loadAll() {
	if _, err := os.Stat(s.skillsDir); os.IsNotExist(err) {
		return
	}
	err := filepath.Walk(s.skillsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "SKILL.md" {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			metaData, body := s.parseFrontmatter(string(data))
			name := metaData["name"]
			if name == "" {
				name = filepath.Base(filepath.Dir(path))
			}
			s.skills[name] = skill{
				name:         name,
				descriptions: metaData["description"],
				meta:         metaData,
				body:         body,
				path:         path,
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("遍历技能目录出错: %v\n", err)
	}
}

func (s *SkillLoader) parseFrontmatter(text string) (meta, string) {
	// Parse YAML frontmatter between --- delimiters.
	re := regexp.MustCompile(`(?s)^---\n(.*?)\n---\n(.*)`)
	match := re.FindStringSubmatch(text)
	if len(match) < 3 {
		return meta{}, text
	}
	frontmatter := match[1]
	body := strings.TrimSpace(match[2])
	m := make(meta)

	lines := strings.Split(frontmatter, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			switch key {
			case "tags":
				m["tags"] = val
			case "name":
				m["name"] = val
			case "description":
				m["description"] = val
			}
		}
	}
	return m, body
}

func (s SkillLoader) GetDescriptions() string {
	// Layer 1: short descriptions for the system prompt.
	if len(s.skills) == 0 {
		return "(no skills available)"
	}
	lines := make([]string, 0)
	for name, skill := range s.skills {
		desc := skill.descriptions
		if desc == "" {
			desc = "No description"
		}
		tags := skill.meta["tags"]
		line := fmt.Sprintf("  - %s: %s", name, desc)
		if tags != "" {
			line += fmt.Sprintf(" [%s]", tags)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (s SkillLoader) GetContent(name string) string {
	// Layer 2: full skill body returned in tool_result.
	if skill, ok := s.skills[name]; ok {
		return fmt.Sprintf("<skill name=\"%s\">\n%s\n</skill>", name, skill.body)
	}
	skillList := ""
	for k, _ := range s.skills {
		skillList += k + ","
	}
	return fmt.Sprintf("Error: Unknown skill %s. Available: %s", name, skillList)
}

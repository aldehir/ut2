package ini

import (
	"fmt"
	"io"
	"strings"
)

type Config struct {
	Sections []*Section
}

func (c *Config) Write(w io.Writer) error {
	var err error
	for _, s := range c.Sections {
		err = s.Write(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) Values(section string, key string) ([]string, bool) {
	sec, found := c.Section(section)
	if !found {
		return nil, false
	}

	values, found := sec.Values(key)
	if !found {
		return nil, false
	}

	return values, true
}

func (c *Config) Section(name string) (*Section, bool) {
	for _, s := range c.Sections {
		if s.Name == name {
			return s, true
		}
	}

	return nil, false
}

type Section struct {
	Name  string
	Items []*Item
}

func (s *Section) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "[%s]\n", s.Name)
	if err != nil {
		return err
	}

	for _, item := range s.Items {
		err = item.Write(w)
		if err != nil {
			return err
		}
	}

	fmt.Fprint(w, "\n")
	return nil
}

func (s *Section) Values(key string) ([]string, bool) {
	for _, i := range s.Items {
		if strings.EqualFold(i.Key, key) {
			return i.Values, true
		}
	}

	return nil, false
}

type Item struct {
	Key    string
	Values []string
}

func (i *Item) Write(w io.Writer) error {
	var err error
	for _, val := range i.Values {
		_, err = fmt.Fprintf(w, "%s=%s\n", i.Key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

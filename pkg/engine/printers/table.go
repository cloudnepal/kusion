package printers

import (
	k8swatch "k8s.io/apimachinery/pkg/watch"

	"kusionstack.io/kusion/pkg/util/pretty"
)

type Table struct {
	IDs  []string
	Rows map[string]*Row
}

type Row struct {
	Type   k8swatch.EventType
	Kind   string
	Name   string
	Detail string
}

func NewTable(ids []string) *Table {
	return &Table{
		IDs:  ids,
		Rows: make(map[string]*Row),
	}
}

func NewRow(t k8swatch.EventType, kind, name, detail string) *Row {
	return &Row{
		Type:   t,
		Kind:   kind,
		Name:   name,
		Detail: detail,
	}
}

const READY k8swatch.EventType = "READY"

func (t *Table) Update(id string, row *Row) {
	t.Rows[id] = row
}

func (t *Table) AllCompleted() bool {
	if len(t.Rows) < len(t.IDs) {
		return false
	}
	for _, row := range t.Rows {
		if row.Type != READY {
			return false
		}
	}
	return true
}

func (t *Table) Print() [][]string {
	data := [][]string{{"Type", "Kind", "Name", "Detail"}}
	for _, id := range t.IDs {
		var eventType k8swatch.EventType
		row := t.Rows[id]
		if row != nil {
			eventType = row.Type
		} else {
			// In case that the row of the table hasn't updated contents.
			continue
		}

		// Colored type
		eventTypeS := pretty.Normal(string(eventType))

		data = append(data, []string{eventTypeS, row.Kind, row.Name, row.Detail})
	}
	return data
}

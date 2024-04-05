package gsdlist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
  Task string
  Status bool
  CreatedAt time.Time
}

type List []item

// Add task to list
func (t *List) Add(task string) {
  todo := item{
    Task: task,
    Status: false,
    CreatedAt: time.Now(),
  }

  *t = append(*t, todo)
}

// Mark task as completed
func (t *List) SetStatus(index int) error {  
  ls := *t

  if index <= 0 || index > len(ls) {
    return errors.New("Invalid index")
  }

  // set status as true
  ls[index-1].Status = !ls[index-1].Status

  return nil
}

// Delete method
func (t *List) Delete(index int) error {
  ls := *t

  if index <= 0 || index > len(ls) {
    return errors.New("invalid index")
  }

  *t = append(ls[:index-1], ls[index:]...)

  return nil
}

func (t *List) Load(filename string) error {
  file, err := ioutil.ReadFile(filename) 
  if err != nil {
    if errors.Is(err, os.ErrNotExist) {
      return nil
    }
    return err
  }

  if len(file) == 0 {
    return err
  }

  err = json.Unmarshal(file, t)
  if err != nil {
    return err
  }

  return nil
}

func (t *List) Store (filename string) error {
  data, err := json.Marshal(t)
  if err != nil {
    return err
  }

  return ioutil.WriteFile(filename, data, 0644)
}

func (t *List) Print() {
  table := simpletable.New()

  table.Header = &simpletable.Header{
    Cells: []*simpletable.Cell {
      {
        Align: simpletable.AlignCenter,
        Text: "#",
      },
      {
        Align: simpletable.AlignCenter,
        Text: "Task",
      },
      {
        Align: simpletable.AlignCenter,
        Text: "Status",
      },
    },
  }

  var cells [][]*simpletable.Cell

  for index, item := range *t {
    index++
    cells = append(cells, *&[]*simpletable.Cell{
      {Text: fmt.Sprintf("%d", index)},
      {Text: item.Task},
      {Text: fmt.Sprintf("%t", item.Status)}, 
    })}

  table.Body = &simpletable.Body{Cells: cells}
  table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell {
    {Align: simpletable.AlignCenter, Span: 3, Text: "Tasks"},
  }}

  table.SetStyle(simpletable.StyleUnicode)
  table.Println()
}


package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/anuragcsangal/gsdlist"
)

const (
  todoFile = "todos.json"
)

func main() {

  // Flags
  add := flag.Bool("add", false, "add a new todo")
  setStatus := flag.Int("setStatus", 0, "mark a todo complete")
  del := flag.Int("del", 0, "delete a todo")
  list := flag.Bool("list", false, "list all todos")

  flag.Parse()

  todos := &gsdlist.List{}

  if err := todos.Load(todoFile); err != nil {
    fmt.Fprintln(os.Stderr, err.Error())
    os.Exit(1)
  }

  switch {
  // Operation for flag Add
  case *add:
    task, err := getInput(os.Stdin, flag.Args()...)
    if err != nil {
      fmt.Fprintln(os.Stderr, err.Error())
      os.Exit(1)
    }
    todos.Add(task)
    err = todos.Store(todoFile)
    if err != nil {
      fmt.Fprintln(os.Stderr, err.Error())
      os.Exit(1)
    }

  // Operation for flag setStatus
  case *setStatus > 0:
    err := todos.SetStatus(*setStatus)
    if err != nil {
      fmt.Fprintln(os.Stderr, err.Error())
      os.Exit(1)
    }
    err = todos.Store(todoFile)
    if err != nil {
      fmt.Fprintln(os.Stderr, err.Error())
      os.Exit(1)
    }
    
  // Operation for flag del (delete)
  case *del > 0:
    err := todos.Delete(*del)
    if err != nil {
      fmt.Fprintln(os.Stderr, err.Error())
      os.Exit(1)
    }
    err = todos.Store(todoFile)
    if err != nil {
      fmt.Fprintln(os.Stderr, err.Error())
      os.Exit(1)
    }

  // Flag to show all todos
  case *list:
    todos.Print()

  default:
    fmt.Fprintln(os.Stdout, "Invalid command")
    os.Exit(1)
  }
}

func getInput(r io.Reader, args ...string) (string, error) {
  if len(args) > 0 {
    return strings.Join(args, " "), nil
  }

  scanner := bufio.NewScanner(r)
  scanner.Scan()
  if err := scanner.Err(); err != nil {
    return "", err
  }

  if len(scanner.Text()) == 0 {
    return "", errors.New("Empty todo is not allowed")
  }

  return scanner.Text(), nil
}

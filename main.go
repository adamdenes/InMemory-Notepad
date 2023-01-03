package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Action string

const (
	Exit   Action = "exit"
	Create Action = "create"
	List   Action = "list"
	Clear  Action = "clear"
	Update Action = "update"
	Delete Action = "delete"
)

type Notepad []string

const CommandString string = "Enter a command and data:"

func main() {
	np := make(Notepad, 0, getInput())
	process(&np)
}

func process(np *Notepad) {
	for {
		var restS string
		fmt.Printf("%s %s", CommandString, restS)

		cmd, rest := reader()
		restS = strings.Join(rest, " ")

		switch Action(cmd) {
		case Exit:
			quit()
		case Create:
			add(np, restS)
		case List:
			show(np)
		case Clear:
			deleteAll(np)
		case Update:
			update(np, restS)
		case Delete:
			remove(np, restS)
		default:
			fmt.Printf("[Error] Unknown command\n")
		}
	}
}

func reader() (string, []string) {
	r := bufio.NewReader(os.Stdin)
	s, err := r.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	s = strings.Trim(s, "\r\n ")
	slice := strings.Split(s, " ")
	return slice[0], slice[1:]
}

func getInput() int {
	var numOfNotes int
	fmt.Print("Enter the maximum number of notes: ")
	_, err := fmt.Scanf("%d", &numOfNotes)
	if err != nil {
		return 0
	}
	fmt.Println()
	return numOfNotes
}

func quit() {
	fmt.Printf("[Info] Bye!\n")
	os.Exit(0)
}

func add(np *Notepad, record string) {
	if cap(*np) == len(*np) {
		fmt.Printf("[Error] Notepad is full\n")
		return
	}
	if record == " " || len(record) == 0 {
		fmt.Printf("[Error] Missing note argument\n")
		return
	}
	*np = append(*np, record)
	fmt.Printf("[OK] The note was successfully created\n")
}

func show(np *Notepad) {
	if len(*np) == 0 {
		fmt.Printf("[Info] Notepad is empty\n")
		return
	}
	for i, note := range *np {
		if note == "" {
			continue
		}
		fmt.Printf("[Info] %d: %s\n", i+1, note)
	}
}

func deleteAll(np *Notepad) {
	*np = make(Notepad, 0, cap(*np))
	fmt.Printf("[OK] All notes were successfully deleted\n")
}

func update(np *Notepad, record string) {
	slice, idx, err := checkError(np, "update", record)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	(*np)[idx] = strings.Join(slice[1:], " ")
	fmt.Printf("[OK] The note at position %d was successfully updated\n", idx+1)
}

func remove(np *Notepad, record string) {
	_, idx, err := checkError(np, "delete", record)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	*np = append((*np)[:idx], (*np)[idx+1:]...)
	fmt.Printf("[OK] The note at position %d was successfully deleted\n", idx+1)
}

func checkError(np *Notepad, a Action, record string) ([]string, int, error) {
	str := ""
	switch a {
	case Update:
		str = string(a)
	case Delete:
		str = string(a)
	}

	if len(*np) == 0 {
		return nil, 0, fmt.Errorf("[Error] There is nothing to %s\n", str)
	}
	if len(record) == 0 {
		return nil, 0, errors.New("[Error] Missing position argument\n")
	}
	slice, idx, err := slicer(record)
	if len(slice) < 2 && a == "update" {
		return nil, 0, errors.New("[Error] Missing note argument\n")
	}
	if err != nil {
		return nil, 1, fmt.Errorf("[Error] Invalid position: %s\n", slice[0])
	}
	if idx >= len(*np) {
		return nil, 1, fmt.Errorf(
			"[Error] Position %d is out of the boundary [1, %d]\n",
			idx+1,
			cap(*np),
		)
	}
	return slice, idx, err
}

func slicer(s string) ([]string, int, error) {
	slice := strings.Split(s, " ")
	idx, err := strconv.Atoi(slice[0])
	// decrement idx due to zero based indexing in Notepad
	idx--
	return slice, idx, err
}

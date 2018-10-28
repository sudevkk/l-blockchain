package main

func main() {

	bc := MakeNewBlockchain("bc")
	defer bc.db.Close()

	exec := NewCommandExecutor(bc)
	exec.run()
}

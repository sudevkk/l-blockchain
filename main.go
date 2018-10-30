package main

func main() {

	bc := MakeNewBlockchain("dev")
	defer bc.db.Close()

	exec := NewCommandExecutor(bc)
	exec.run()
}

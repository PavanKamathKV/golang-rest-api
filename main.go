package main

func main() {
	a := App{}

	a.Initialize("root", "Kamath@123", "golang_restAPI")
	a.Run(":8080")

}

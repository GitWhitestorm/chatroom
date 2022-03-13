package tcp

import "net"

func main() {
	listener, err := net.Listen("tcp", ":2000")
	if err != nil {
		panic(err)
	}
	go broadcaster()

	for 

	

}
// 记录用户 用于广播
func broadcaster(){
	
}
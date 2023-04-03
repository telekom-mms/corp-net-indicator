package main

func main() {
	iS := NewIdentityServer(true)
	defer iS.Close()
	vS := NewVPNServer(true)
	defer vS.Close()

	select {}
}

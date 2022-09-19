package main

import "flag"

const (
	port              = "3000"
	host              = "plex.tv"
	officialAppURL    = "https://www.dropbox.com/s/f17hx2w7tvofjqr/Plex_2.014_11112020.zip?dl=1"
	officialAppChksum = "8c6b2bb25a4c2492fd5dbde885946dcb6b781ba292e5038239559fd7a20e707e"
	officialAppFile   = "Plex_2.014_11112020.zip"
	modifiedAppName   = "Plex_2.014_net"
	modifiedAppFile   = "Plex_2.014_11112020_net.zip"
)

func main() {
	disableSideloadingPtr := flag.Bool("disable-sideloading", false, "This Option will disable sideloading and prevent the use of Port 80 by the Application")
	flag.Parse()

	checkIPs()
	runProxy(*disableSideloadingPtr)
}

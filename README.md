# proxy4plex
Proxy Server for older Samsung Smart-TVs (E/F - Series) which only support TLS up to 1.1

## Installation

### Preparing the Proxy

Download the corresponding binary for your system from the [latest releases](https://github.com/kadrim/proxy4plex/releases/latest) and run it. You might want to run the binary on system start. This is different for each system but in the future the proxy will install itself if needed.

### Preparing the Smart-TV

#### Using proxy4plex as sideloader

Notice: this will only work, if the system you are running on, is on the same network as your TV and has not yet a service running on Port 80, because this is needed by the TV to download new Apps. I.e: This will not work on non-rooted Android Phones, as the Port 80 is blocked. Running the proxy on non-rooted Android phones does work. So the initial Setup has do be done via another machine.

This process has only to be done once. So after the app has been installed via this method, it stays on the TV!

1. Start the `proxy4plex` binary
2. Open `smarthub` on your TV
3. Hit the `Red` button on the remote to login to an account
4. Use `develop` for username, password should fill automatically, if it does not, use `111111`
5. Hit the `Tools` button on the remote and select `Settings` while within `smarthub`
6. Select the menu-entry `Development` and accept the agreement after reading it
7. Select `Server-IP-Settings` and enter the local IP-Adress of your machine running `proxy4plex`
8. Select `Synchronize` and the app will be installed permanentely on the TV

#### Manual sideloading

1. Download the latest version of Orca's Plex-App for Samsung TVs from here https://www.dropbox.com/s/f17hx2w7tvofjqr/Plex_2.014_11112020.zip?dl=1
2. Install it using USB-Sideloading or via SammyWidgets
3. Check this thread if you need more information: https://forums.plex.tv/t/samsung-tv-cannot-connect-to-plex/650100/8

### Putting it together

1. Run the app on your TV.
2. You will see that the app cannot connect to plex.tv and a popup will show up
3. Select the "Home" option and a new screen will appear that has a "Configure Proxy" option. Select it and press enter to show the proxy entry screen.
4. Enter the IP-Address of the server that is running the Proxy server and click "Add Server"
5. If the proxy server is reachable, you will get a message and the app will close down automatically.
6. Now restart the app, and it will use the Proxy to connect to plex.tv

Note that you can run the TV-App without the proxy after you logged in the first time. The app will then run in "Offline" mode. However you cannot switch users without the proxy though.

#### Using Android

On Android this package is currently not yet available as native app. Nevertheless it can be executed, albeit with a little tinkering. To do that, follow these steps:

1. Install [termux](https://play.google.com/store/apps/details?id=com.termux) on your device
2. Run termux, a terminal will open
3. Run the command `pkg install git golang` which will install the needed tools
4. Get the current repository by typing `git clone https://github.com/kadrim/proxy4plex`
5. Enter the directory `cd proxy4plex`
6. Compile a binary for your android device `go build`
7. Execute the proxy! `./proxy4plex`

After that you can use the IP-Adress of your android device for your Plex-App on the TV. If you close termux or reboot your device you can simply re-run the proxy by starting termux and then starting the proxy via the command `proxy4plex/proxy4plex`

Beware: Sideloading (i.e. installing the app on the TV) does not work this way, because non-rooted Android devices are not allowed to use Port 80

## Compiling

At the time of writing, this package needs at least [golang](https://golang.org/) v1.15

To compile a binary for your currently running system, simply run this command:

`go build`

To compile for all possible architectures you can run the command `go-build-all.sh` from a bash terminal (works also on windows in a git-bash terminal). The output binaries will be put into the directory `build/`

## TODOs

- detect OS and allow User to install the proxy as a boot-service
- porting to Android and iOS as native apps
- auto-update notifications 

## Thanks

- Thanks go out to @makeworld for the great makefile to build all platforms at https://gist.github.com/makeworld-the-better-one/e1bb127979ae4195f43aaa3ad46b1097
- Very specials thanks to @Orca for his work on the plex-app for smarthub devices. Without that, nothing here would exist!

## Donation
I developed this package in my free time. If you like it and want to support future updates, feel free to donate here:

[Donate via PayPal](https://www.paypal.com/donate?hosted_button_id=RDJ8ZWG3GRWE8)

Thanks in advance :-)

## Disclaimer
THIS SOFTWARE IS PROVIDED "AS IS" AND ANY EXPRESSED OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


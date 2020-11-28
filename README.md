# proxy4plex
Proxy Server for older Samsung Smart-TVs (E/F - Series) which only support TLS up to 1.1

## Installation

### Preparing the Smart-TV

- Download the latest version of Orca's Plex-App for Samsung TVs from here https://www.dropbox.com/s/f17hx2w7tvofjqr/Plex_2.014_11112020.zip?dl=1
- Install it sing USB-Sideloading or via SammyWidgets
- Check this thread if you need more information: https://forums.plex.tv/t/samsung-tv-cannot-connect-to-plex/650100/8

### Preparing the Proxy

Download the corresponding binary for your system from the [latest releases](https://github.com/kadrim/proxy4plex/releases/latest) and run it. You might want to run the binary on system start. This is different for each system but in the future the proxy will install itself if needed.

### Putting it together

1. Run the app on your TV.
2. You will see that the app cannot connect to plex.tv and a popup will show up
3. Select the "Home" option and a new screen will appear that has a "Configure Proxy" option. Select it and press enter to show the proxy entry screen.
4. Enter the IP-Address of the server that is running the Proxy server and click "Add Server"
5. If the proxy server is reachable, you will get a message and the app will close down automatically.
6. Now restart the app, and it will use the Proxy to connect to plex.tv

Note that you can run the app without the proxy after you logged in the first time. The app will then run in "Offline" mode. However you cannot switch users without the proxy though.

## TODOs

- detect OS and allow User to install the proxy as a boot-service
- serve as sideloading server so no need to manually run SammyWidgets
- porting to Android and iOS
- auto-update notifications 

## Thanks

- Thanks go out to @makeworld for the great makefile to build all platforms at https://gist.github.com/makeworld-the-better-one/e1bb127979ae4195f43aaa3ad46b1097
- Very specials thanks to @Orca for his work on the plex-app for smarthub devices. Without that, nothing here would exist!

## Disclaimer
THIS SOFTWARE IS PROVIDED "AS IS" AND ANY EXPRESSED OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


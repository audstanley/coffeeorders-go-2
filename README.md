# Coffeeorders-go

is an API for [Front-End Web Development: The Big Nerd Ranch Guide](https://a.co/d/egDG6ub)
The book is missing some endpoints, and this GoLang Project recreates those endpoints and is compiled for windows, mac, and linux.


## Running the REST Server:

```bash
# first, install golang
# navigate to the project folder in whatever terminal:
go build .;
# if you are on linux or mac:
./coffeeorders-go-2
# or in windows powershell:
.\coffeeoders-go-2.exe
# binaries are built for the OS you are currently on.
```

### Once you have the binaries running, you will want this extension:

[REST](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)

Once you install the rest extension, you can copy the .http files that are in this repositorie's rest folder,
don't forget to copy the rest/.env file, since that is where are the http calls are assigned to.

In the rest/.env
```bash
baseUrl=https://co.audstanley.com
# can be changed to
baseUrl=http://127.0.0.1:3001
# and now all the http calls will be made to your own computer - you must be running the REST Server.
```

You are welcome to:

<a href='https://ko-fi.com/A687KA8' target='_blank'><img height='36' style='border:0px;height:36px;' src='https://storage.ko-fi.com/cdn/useruploads/1e427eff-8866-46b7-bd55-03c85e75c6c1_9aa005d7-97dd-4a17-86d5-b01d4e7b47bb.png' border='0' alt='Buy Me a Coffee at ko-fi.com' /></a>

Or stop by my blog [audstanley.com](https://www.audstanley.com)

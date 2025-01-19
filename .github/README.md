
<p align="center">
<img src="https://i.ibb.co/J5r1zCP/who-dat-square.png" width="128" /><br />
<i>Free & Open Source WHOIS Lookup Service</i>
<br />
<i>No-CORS, no auth API that's publicly available or easily self-hostable</i>
<br />
<b>🌐 <a href="https://who-dat.as93.net/">who-dat.as93.net</a></b><br />
</p>

---

<details>
  <summary>Contents</summary>
  
- [API Usage](#api-usage)
  - [Base URL](#base-url)
  - [Endpoints](#endpoints)
    - [Single Domain](#single-domain-lookup-domain)
    - [Bulk Domains](#multiple-domain-lookup-multi)
- [Deployment](#deployment)
  - [Option 1: Vercel](#option-1-vercel)
  - [Option 2: Docker](#option-2-docker)
  - [Option 3: Binary](#option-3-binary)
  - [Option 4: Build from Source](#option-4-build-from-source)
- [Adding Auth](#authentication)
- [Development](#development)
- [Contributing](#contributing)
- [Web Interface](#web-interface)
- [Mirror](#mirror)
- [Credits](#credits)
- [More Like This](#more-like-this)
- [License](#license)

</details>

## API Usage

> **TL;DR** Get the WHOIS records for a site: `curl https://who-dat.as93.net/example.com`

For detailed request + response schemas, and to try the API out, you can reference the [spec](https://who-dat.as93.net/docs.html)

### Base URL

The base URL for the public API is [`who-dat.as93.net`](https://who-dat.as93.net)

If you're self-hosting (recommended) then replace this with your own base URL.

### Endpoints

<details>
  <summary><h4>Single Domain Lookup <code>/[domain]</code></h4></summary>

- **URL**: `/[domain]`
- **Method**: `GET`
- **URL Params**: None
- **Success Response**:
  - **Code**: 200
  - **Content**: WHOIS data for the specified domain in JSON format.
- **Error Response**:
  - **Code**: 400 BAD REQUEST
  - **Content**: `{ "error": "Domain not specified" }`
  - **Code**: 404 NOT FOUND
  - **Content**: `{ "error": "Domain not found" }`
- **Sample Call**:

##### Command Line

```bash
curl https://who-dat.as93.net/example.com
```

##### JavaScript

```javascript
fetch('https://who-dat.as93.net/example.com')
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

##### Python

```python
import requests

response = requests.get('https://who-dat.as93.net/example.com')
if response.status_code == 200:
    print(response.json())
else:
    print("Error:", response.status_code)
```

</details>

<details>
  <summary><h4>Multiple Domain Lookup <code>/multi</code></h4></summary>

- **URL**: `/multi`
- **Method**: `GET`
- **Query Params**: 
  - **domains**: A comma-separated list of domains.
- **Success Response**:
  - **Code**: 200
  - **Content**: Array of WHOIS data for the specified domains in JSON format.
- **Error Response**:
  - **Code**: 400 BAD REQUEST
  - **Content**: `{ "error": "No domains specified" }`
  - **Code**: 500 INTERNAL SERVER ERROR
  - **Content**: `{ "error": "[error message]" }`
- **Sample Call**:

```
curl "https://who-dat.as93.net/multi?domains=example.com,example.net"
```

</details>

[![Who-Dat Swagger Docs](https://img.shields.io/badge/Swagger-Docs-85EA2D?style=for-the-badge&logo=swagger&labelColor=1b2744&link=https%3A%2F%2Fwho-dat.as93.net%2Fdocs.html)](https://who-dat.as93.net/docs.html)


---

## Deployment

#### Option 1: Vercel

This is the quickest and easiest way to get up-and-running. Simply fork the repository, then login to Vercel (using GitHub), and after importing your fork, it will be deployed! There's no additional config or keys needed, and it should work just fine on the free plan.

Alternatively, just hit the button below for 1-click deploy 👇

[![1-Click Deploy to Vercel](https://img.shields.io/badge/Deploy-Vercel-ffffff?style=for-the-badge&logo=vercel&labelColor=1b2744&link=https%3A%2F%2Fwho-dat.as93.net%2Fdocs.html)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Flissy93%2Fwho-dat&demo-title=Who-Dat%20Demo&demo-url=https%3A%2F%2Fwho-dat.as93.net&demo-image=https%3A%2F%2Fi.ibb.co%2FJ5r1zCP%2Fwho-dat-square.png)

#### Option 2: Docker

The light-weight Docker image is published to DockerHub ([hub.docker.com/r/lissy93/who-dat](https://hub.docker.com/r/lissy93/who-dat)), as well as GHCR ([here](https://github.com/Lissy93/who-dat/pkgs/container/who-dat)).

Providing you've got Docker installed, you can get everything by running:

```shell
docker run -p 8080:8080 --dns 8.8.8.8 --dns 8.8.4.4 lissy93/who-dat
```

[![Deploy from Docker](https://img.shields.io/badge/Deploy-Docker-2496ED?style=for-the-badge&logo=docker&labelColor=1b2744&link=https%3A%2F%2Fwho-dat.as93.net%2Fdocs.html)](https://hub.docker.com/r/lissy93/who-dat)


#### Option 3: Binary

Head to the [Releases Tab](https://github.com/Lissy93/who-dat/releases), download and extract the pre-built executable for your system, then run it.

<details>

<summary>Example</summary>

If you're using the command line, you can do something like this<br>
Don't forget to update (v1.0) with the version number you want, and (linux-amd64) with your system's architecture.
  
```bash
# Download the binary for your system (from releases tab)
wget https://github.com/Lissy93/who-dat/releases/download/v0.9/who-dat-v0.9-linux-amd64.tar.gz -O ./who-dat.tar.gz

# Extract the compressed file
tar -xzvf who-dat.tar.gz

# Make it executable
chmod +x who-dat

# Run Who-Dat!
./who-dat
```

(Or, if you're a Microsoft fanboy, you can just double-click the `who-dat.exe` after extracting in Windows Explorer)

</details>



#### Option 4: Build from Source

Follow the setup instructions in the [Development](#development) section.<br>
Then run `go build -a -installsuffix cgo -o who-dat .` to generate the binary for your system.<br>
You'll then be able to execute the newly built `./who-dat` file directly to start the application.

---

## Authentication

Authentication is optional, and can be enabled by setting the `AUTH_KEY` environment variable.

When authentication is enabled, all API requests must include the key in the Authorization header, using one of the formats indicated below.

#### Raw API Key

```
curl -H "Authorization: your-secret-key" https://who-dat.yourdomain.com/example.com
```

#### Bearer Token Format

```
curl -H "Authorization: Bearer your-secret-key" https://who-dat.yourdomain.com/example.com
```

If authentication is not configured (no `AUTH_KEY` set), the API will remain publicly accessible.

---

## Development

Prerequisites: You'll need [Go](https://go.dev/) and [Node](https://nodejs.org/) installed. You will likley also want [Git](https://git-scm.com/) and/or [Docker](https://www.docker.com/).

```
git clone git@github.com:Lissy93/who-dat.git
cd who-dat
go get
npm install
npm run build
```

Then run either `npx vercel dev`, or `go run main.go`

Alternativley, build the Docker container with `docker build -t who-dat .`

[![Open in GitPod](https://img.shields.io/badge/GitPod-Try_Live-FFAE33?style=for-the-badge&logo=gitpod&labelColor=1b2744&link=https%3A%2F%2Fcodeberg.org%2Falicia%2Fwho-dat)](https://gitpod.io/#https://github.com/lissy93/who-dat)
[![Open in VS Code](https://img.shields.io/badge/CodeSpaces-Try_Live-007ACC?style=for-the-badge&logo=visualstudiocode&labelColor=1b2744&link=https%3A%2F%2Fcodeberg.org%2Falicia%2Fwho-dat)](https://codespaces.new/Lissy93/who-dat)

---

## Web Interface

There's a very simple frontend included in the app. This is built with Alpine.js, so is super light-weight, and only adds about 100kb to the total executable.
The web interface is used to view WHOIS records for a given domain, and also hosts the API documentation.

<p align="center">
<img width="600" src="https://i.ibb.co/1dYcdZC/who-dat-screenshot.png" />
</p>

---

## Contributing

Contributions of any kind are welcome (and would be much appreciated!). Be sure to follow our [Code of Conduct](https://github.com/Lissy93/who-dat/blob/main/.github/CODE_OF_CONDUCT.md).<br>
If you're new to open source, I've put together some guides in [Git-In](https://github.com/Lissy93/git-into-open-source/), but feel free to reach out if you need any support.

Not a coder? You can still help, by raising bugs you find, updating docs, or consider sponsoring me on GitHub

[![Sponsor](https://img.shields.io/badge/Sponsor-Lissy93-EA4AAA?style=for-the-badge&logo=githubsponsors&labelColor=1b2744&link=https%3A%2F%2Fgithub.com%2Fsponsors%2FLissy93)](https://github.com/sponsors/Lissy93)

---

## Mirror

We've got a (non-Microsoft) mirror of this repository hosted on CodeBerg, at [codeberg.org/alicia/who-dat](https://codeberg.org/alicia/who-dat)

[![CodeBerg Mirror](https://img.shields.io/badge/Mirror-Who_Dat-2185D0?style=for-the-badge&logo=codeberg&labelColor=1b2744&link=https%3A%2F%2Fcodeberg.org%2Falicia%2Fwho-dat)](https://codeberg.org/alicia/who-dat)


---

## Credits

##### Inspiration
This project was inspired by [someshkar/whois-api](https://github.com/someshkar/whois-api) by [Somesh Kar](https://someshkar.com/).

##### Tech Credits
- The frontend is built with Alpine.js[^alpinejs], Vite[^vite], TS[^typescript] and SCSS[^scss] (plus the usual web tech stack).
- The backend is written in Go[^golang], and was made possible thanks to [json-iterator/go](https://github.com/json-iterator/go) and [likexian/whois-parser](https://github.com/likexian/whois-parser)
- Demo deployed to Vercel[^vercel] (but also available on DockerHub[^dockerhub]), and source of course on GitHub[^github] and CodeBerg[^codeberg].

[^alpinejs]: [Alpine.js](https://alpinejs.dev/) - A rugged, minimal framework for composing JavaScript behavior in your markup.
[^vite]: [Vite](https://vitejs.dev/) - A build tool that aims to provide a faster and leaner development experience for modern web projects.
[^typescript]: [TypeScript](https://www.typescriptlang.org/) - A typed superset of JavaScript that compiles to plain JavaScript.
[^scss]: [SCSS](https://sass-lang.com/) - A preprocessor scripting language that is interpreted or compiled into Cascading Style Sheets (CSS).
[^golang]: [Go Lang](https://golang.org/) - An open source programming language that makes it easy to build simple, reliable, and efficient software.
[^github]: [GitHub](https://github.com/) - A platform for version control and collaboration. It lets you and others work together on projects from anywhere.
[^codeberg]: [Codeberg](https://codeberg.org/) - A free and open-source forge for collaborative software development.
[^vercel]: [Vercel](https://vercel.com/) - Static hosting and shit
[^dockerhub]: [DockerHub](https://hub.docker.com/) - Container registry hosting and shit

##### Contributors

<!-- readme: contributors -start -->
<table>
<tr>
    <td align="center">
        <a href="https://github.com/liss-bot">
            <img src="https://avatars.githubusercontent.com/u/87835202?v=4" width="80;" alt="liss-bot"/>
            <br />
            <sub><b>Alicia Bot</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/Lissy93">
            <img src="https://avatars.githubusercontent.com/u/1862727?v=4" width="80;" alt="Lissy93"/>
            <br />
            <sub><b>Alicia Sykes</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/SamLS98">
            <img src="https://avatars.githubusercontent.com/u/63964062?v=4" width="80;" alt="SamLS98"/>
            <br />
            <sub><b>Sammy LS</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/cedwardsmedia">
            <img src="https://avatars.githubusercontent.com/u/1514767?v=4" width="80;" alt="cedwardsmedia"/>
            <br />
            <sub><b>Corey Edwards</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/phill-holland">
            <img src="https://avatars.githubusercontent.com/u/32714574?v=4" width="80;" alt="phill-holland"/>
            <br />
            <sub><b>Phill Holland</b></sub>
        </a>
    </td></tr>
</table>
<!-- readme: contributors -end -->

##### Sponsors

<!-- readme: sponsors -start -->
<table>
<tr>
    <td align="center">
        <a href="https://github.com/vincentkoc">
            <img src="https://avatars.githubusercontent.com/u/25068?u=cbf098fc04c0473523d373b0dd2145b4ec99ef93&v=4" width="80;" alt="vincentkoc"/>
            <br />
            <sub><b>Vincent Koc</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/AnandChowdhary">
            <img src="https://avatars.githubusercontent.com/u/2841780?u=747e554b3a7f12eb20b7910e1c87d817844f714f&v=4" width="80;" alt="AnandChowdhary"/>
            <br />
            <sub><b>Anand Chowdhary</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/shrippen">
            <img src="https://avatars.githubusercontent.com/u/2873570?v=4" width="80;" alt="shrippen"/>
            <br />
            <sub><b>Shrippen</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/bile0026">
            <img src="https://avatars.githubusercontent.com/u/5022496?u=aec96ad173c0ea9baaba93807efa8a848af6595c&v=4" width="80;" alt="bile0026"/>
            <br />
            <sub><b>Zach Biles</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/UlisesGascon">
            <img src="https://avatars.githubusercontent.com/u/5110813?u=3c41facd8aa26154b9451de237c34b0f78d672a5&v=4" width="80;" alt="UlisesGascon"/>
            <br />
            <sub><b>Ulises Gascón</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/digitalarche">
            <img src="https://avatars.githubusercontent.com/u/6546135?u=564756d7f44ab2206819eb3148f6d822673f5066&v=4" width="80;" alt="digitalarche"/>
            <br />
            <sub><b>Digital Archeology</b></sub>
        </a>
    </td></tr>
<tr>
    <td align="center">
        <a href="https://github.com/InDieTasten">
            <img src="https://avatars.githubusercontent.com/u/7047377?u=8d8f8017628b38bc46dcbf3620e194b01d3fb2d1&v=4" width="80;" alt="InDieTasten"/>
            <br />
            <sub><b>InDieTasten</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/araguaci">
            <img src="https://avatars.githubusercontent.com/u/7318668?v=4" width="80;" alt="araguaci"/>
            <br />
            <sub><b>Araguaci</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/bmcgonag">
            <img src="https://avatars.githubusercontent.com/u/7346620?u=2a0f9284f3e12ac1cc15288c254d1ec68a5081e8&v=4" width="80;" alt="bmcgonag"/>
            <br />
            <sub><b>Brian McGonagill</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/arcestia">
            <img src="https://avatars.githubusercontent.com/u/7936962?u=41e34bb816ad09323e1650f3efc0bec4fb2bc5dd&v=4" width="80;" alt="arcestia"/>
            <br />
            <sub><b>Laurensius Jeffrey</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/vlad-tim">
            <img src="https://avatars.githubusercontent.com/u/11474041?u=eee43705b54d2ec9f51fc4fcce5ad18dd17c87e4&v=4" width="80;" alt="vlad-tim"/>
            <br />
            <sub><b>Vlad</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/helixzz">
            <img src="https://avatars.githubusercontent.com/u/12218889?u=d06d0c103dfbdb99450623064f7da3c5a3675fb6&v=4" width="80;" alt="helixzz"/>
            <br />
            <sub><b>HeliXZz</b></sub>
        </a>
    </td></tr>
<tr>
    <td align="center">
        <a href="https://github.com/patvdv">
            <img src="https://avatars.githubusercontent.com/u/12430107?u=e8911c2fb91af4d30432f76da8c40927b2830bd7&v=4" width="80;" alt="patvdv"/>
            <br />
            <sub><b>Patrick Van Der Veken</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/mryesiller">
            <img src="https://avatars.githubusercontent.com/u/24632172?u=0d20f2d615158f87cd60a3398d3efb026c32f291&v=4" width="80;" alt="mryesiller"/>
            <br />
            <sub><b>Göksel Yeşiller</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/sushibait">
            <img src="https://avatars.githubusercontent.com/u/26634535?v=4" width="80;" alt="sushibait"/>
            <br />
            <sub><b>Shiverme Timbers</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/forwardemail">
            <img src="https://avatars.githubusercontent.com/u/32481436?v=4" width="80;" alt="forwardemail"/>
            <br />
            <sub><b>Forward Email - Open-source & Privacy-focused Email Service (2023)</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/getumbrel">
            <img src="https://avatars.githubusercontent.com/u/59408891?v=4" width="80;" alt="getumbrel"/>
            <br />
            <sub><b>Umbrel</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/OlliVHH">
            <img src="https://avatars.githubusercontent.com/u/84959562?v=4" width="80;" alt="OlliVHH"/>
            <br />
            <sub><b>HamburgerJung</b></sub>
        </a>
    </td></tr>
<tr>
    <td align="center">
        <a href="https://github.com/frankdez93">
            <img src="https://avatars.githubusercontent.com/u/87549420?v=4" width="80;" alt="frankdez93"/>
            <br />
            <sub><b>Frankdez93</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/terminaltrove">
            <img src="https://avatars.githubusercontent.com/u/121595180?v=4" width="80;" alt="terminaltrove"/>
            <br />
            <sub><b>Terminal Trove</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/st617">
            <img src="https://avatars.githubusercontent.com/u/128325650?v=4" width="80;" alt="st617"/>
            <br />
            <sub><b>St617</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/nrvo">
            <img src="https://avatars.githubusercontent.com/u/151435968?u=e1dcb307fd0efdc45cddbe9490a7b956e4da6835&v=4" width="80;" alt="nrvo"/>
            <br />
            <sub><b>Nrvo</b></sub>
        </a>
    </td>
    <td align="center">
        <a href="https://github.com/hudsonrock-partnerships">
            <img src="https://avatars.githubusercontent.com/u/163282900?u=5f2667f7fe5d284ac7a2da6b0800ea8970b0fcbf&v=4" width="80;" alt="hudsonrock-partnerships"/>
            <br />
            <sub><b>Hudsonrock-partnerships</b></sub>
        </a>
    </td></tr>
</table>
<!-- readme: sponsors -end -->

---

## More Like This

You might be interested in [Web-Check](https://github.com/Lissy93/web-check), an all-in-one tool for fetching info on a given domain name.

If you like projects like these, consider [following me](https://github.com/Lissy93) on GitGub 😊<br>
I'm often putting out new (free & open source) utilities, relating to security, privacy, OSINT, Linux and self-hosting.

---

## License

> _**[Lissy93/Who-Dat](https://github.com/Lissy93/who-dat)** is licensed under [MIT](https://github.com/Lissy93/who-dat/blob/HEAD/LICENSE) © [Alicia Sykes](https://aliciasykes.com) 2024._<br>
> <sup align="right">For information, see <a href="https://tldrlegal.com/license/mit-license">TLDR Legal > MIT</a></sup>

<details>
<summary>Expand License</summary>

```
The MIT License (MIT)
Copyright (c) Alicia Sykes <alicia@omg.com> 

Permission is hereby granted, free of charge, to any person obtaining a copy 
of this software and associated documentation files (the "Software"), to deal 
in the Software without restriction, including without limitation the rights 
to use, copy, modify, merge, publish, distribute, sub-license, and/or sell 
copies of the Software, and to permit persons to whom the Software is furnished 
to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included install 
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANT ABILITY, FITNESS FOR A
PARTICULAR PURPOSE AND NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```

</details>


<!-- License + Copyright -->
<p  align="center">
  <i>© <a href="https://aliciasykes.com">Alicia Sykes</a> 2024</i><br>
  <i>Licensed under <a href="https://gist.github.com/Lissy93/143d2ee01ccc5c052a17">MIT</a></i><br>
  <a href="https://github.com/lissy93"><img src="https://i.ibb.co/4KtpYxb/octocat-clean-mini.png" /></a><br>
  <sup>Thanks for visiting :)</sup>
</p>

###### References

<small><sub>➧ See [Credits](#credits)</sub></small>

<!-- Dinosaurs are Awesome -->
<!-- 
                        . - ~ ~ ~ - .
      ..     _      .-~               ~-.
     //|     \ `..~                      `.
    || |      }  }              /       \  \
(\   \\ \~^..'                 |         }  \
 \`.-~  o      /       }       |        /    \
 (__          |       /        |       /      `.
  `- - ~ ~ -._|      /_ - ~ ~ ^|      /- _      `.
              |     /          |     /     ~-.     ~- _
              |_____|          |_____|         ~ - . _ _~_-_
-->

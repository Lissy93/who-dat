
<h1 align="center">Who Dat?</h1>
<p align="center">
<img src="https://i.ibb.co/J5r1zCP/who-dat-square.png" width="128" /><br />
<i>Free & Open Source WHOIS Lookup Service</i>
<br />
<b>🌐 <a href="https://who-dat.as93.net/">who-dat.as93.net</a></b><br />
</p>

---

## API Usage

> **TL;DR** Get the WHOIS records for a site: `curl https://who-dat.as93.net/example.com`

For detailed request + response schemas, and to try the API out, you can reference the [Swagger Docs](https://who-dat.as93.net/docs.html)

### Base URL

The base URL for the public API is [`who-dat.as93.net`](https://who-dat.as93.net)

If you're self-hosting (reccomended) then replace this with your own base URL.

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

---

## Deployment

#### Option 1: Vercel

Click the button below to deploy to Vercel 👇

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Flissy93%2Fwho-dat&demo-title=Who-Dat%20Demo&demo-url=https%3A%2F%2Fwho-dat.as93.net&demo-image=https%3A%2F%2Fi.ibb.co%2FJ5r1zCP%2Fwho-dat-square.png)

#### Option 2: Docker

The Docker image is published to DockerHub ([hub.docker.com/r/lissy93/who-dat](https://hub.docker.com/r/lissy93/who-dat)), as well as GHCR

```shell
docker run -p 8080:8080 lissy93/who-dat
```

#### Option 3: Binary

Head to the [Releases Tab](https://github.com/Lissy93/who-dat/releases), download the pre-built executable for your system, then run it.

#### Option 4: Build from Source

Follow the setup instructions in the [Development](#development) section.<br>
Then run `go build -a -installsuffix cgo -o who-dat .` to generate the binary for your system.<br>
You'll then be able to execute the newly built `./who-dat` file directly to start the application.

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
<!-- readme: contributors -end -->

##### Sponsors

<!-- readme: sponsors -start -->
<!-- readme: sponsors -end -->

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

###### Credits and References

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

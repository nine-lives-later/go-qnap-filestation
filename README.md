# QNAP File Station API for Go

This library provides access to the QNAP File Station API (http://www.qnap.com).

**[Click here to open the documentation.](https://pkg.go.dev/github.com/nine-lives-later/go-qnap-filestation)**

## Usage

To use the API, simply create a new session:

```go
import "github.com/nine-lives-later/go-qnap-filestation"

func main() {
    session, _ := filestation.Connect("storage:8443", "admin", "admin", nil)

    session.Logout()
}
```

## Authors

We thank all the authors who provided code to this library:

* Felix Kollmann
* marvin + konsorten GmbH (who sponsored this library in 2018)
* Scott Ware (v1.0 of this library is loosely based on the [go-bigip package](https://github.com/scottdware/go-bigip))

## License

(The MIT License)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

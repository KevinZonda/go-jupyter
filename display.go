package jupyter

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Support an interface similar - but not identical - to the IPython (canonical Jupyter kernel).
// See http://ipython.readthedocs.io/en/stable/api/generated/IPython.display.html#IPython.display.display
// for a good overview of the support types.

const (
	MIMETypeHTML       = "text/html"
	MIMETypeJavaScript = "application/javascript"
	MIMETypeJPEG       = "image/jpeg"
	MIMETypeJSON       = "application/json"
	MIMETypeLatex      = "text/latex"
	MIMETypeMarkdown   = "text/markdown"
	MIMETypePNG        = "image/png"
	MIMETypePDF        = "application/pdf"
	MIMETypeSVG        = "image/svg+xml"
	MIMETypeText       = "text/plain"
)

/**
 * general interface, allows libraries to fully specify
 * how their data is displayed by Jupyter.
 * Supports multiple MIME formats.
 *
 * Note that Data defined above is an alias:
 * libraries can implement Renderer without importing gophernotes
 */

// if vals[] contain a single non-nil value which is auto-renderable,
// convert it to Data and return it.
// otherwise return MakeData("text/plain", fmt.Sprint(vals...))
func (kernel *Kernel) autoRenderResults(vals []any) Data {
	for _, val := range vals {
		if x, ok := val.(Data); ok {
			return x
		}
	}
	return Data{}
}

func anyToString(vals ...interface{}) string {
	var buf strings.Builder
	for i, val := range vals {
		if i != 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprint(&buf, val)
	}
	return buf.String()
}

// return true if data type should be auto-rendered graphically
func (kernel *Kernel) canAutoRender(data interface{}) bool {
	return true
}

// detect and render data types that should be auto-rendered graphically
func (kernel *Kernel) autoRender(mimeType string, arg interface{}) Data {
	// try Data
	if x, ok := arg.(Data); ok {
		return x
	}

	return Data{}
}

func fillDefaults(data Data, arg interface{}, s string, b []byte, mimeType string, err error) Data {
	if err != nil {
		return makeDataErr(err)
	}
	if data.Data == nil {
		data.Data = make(MIMEMap)
	}
	// cannot autodetect the mime type of a string
	if len(s) != 0 && len(mimeType) != 0 {
		data.Data[mimeType] = s
	}
	// ensure plain text is set
	if data.Data[MIMETypeText] == "" {
		if len(s) == 0 {
			s = fmt.Sprint(arg)
		}
		data.Data[MIMETypeText] = s
	}
	// if []byte is available, use it
	if len(b) != 0 {
		if len(mimeType) == 0 {
			mimeType = http.DetectContentType(b)
		}
		if len(mimeType) != 0 && mimeType != MIMETypeText {
			data.Data[mimeType] = b
		}
	}
	return data
}

// do our best to render data graphically
func render(mimeType string, data interface{}) Data {
	var kernel *Kernel // intentionally nil
	if kernel.canAutoRender(data) {
		return kernel.autoRender(mimeType, data)
	}
	var s string
	var b []byte
	var err error
	switch data := data.(type) {
	case string:
		s = data
	case []byte:
		b = data
	case io.Reader:
		b, err = io.ReadAll(data)
	case io.WriterTo:
		var buf bytes.Buffer
		data.WriteTo(&buf)
		b = buf.Bytes()
	default:
		panic(fmt.Errorf("unsupported type, cannot render: %T", data))
	}
	return fillDefaults(Data{}, data, s, b, mimeType, err)
}

func makeDataErr(err error) Data {
	return Data{
		Data: MIMEMap{
			"ename":     "ERROR",
			"evalue":    err.Error(),
			"traceback": nil,
			"status":    "error",
		},
	}
}

func Any(mimeType string, data interface{}) Data {
	return render(mimeType, data)
}

// same as Any("", data), autodetects MIME type
func Auto(data interface{}) Data {
	return render("", data)
}

func MakeData(mimeType string, data interface{}) Data {
	d := Data{
		Data: MIMEMap{
			mimeType: data,
		},
	}
	if mimeType != MIMETypeText {
		d.Data[MIMETypeText] = fmt.Sprint(data)
	}
	return d
}

func MakeData3(mimeType string, plaintext string, data interface{}) Data {
	return Data{
		Data: MIMEMap{
			MIMETypeText: plaintext,
			mimeType:     data,
		},
	}
}

func File(mimeType string, path string) Data {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return Any(mimeType, bytes)
}

func HTML(html string) Data {
	return MakeData(MIMETypeHTML, html)
}

func JavaScript(javascript string) Data {
	return MakeData(MIMETypeJavaScript, javascript)
}

func JPEG(jpeg []byte) Data {
	return MakeData(MIMETypeJPEG, jpeg)
}

func JSON(json map[string]interface{}) Data {
	return MakeData(MIMETypeJSON, json)
}

func Latex(latex string) Data {
	return MakeData3(MIMETypeLatex, latex, "$"+strings.Trim(latex, "$")+"$")
}

func Markdown(markdown string) Data {
	return MakeData(MIMETypeMarkdown, markdown)
}

func Math(latex string) Data {
	return MakeData3(MIMETypeLatex, latex, "$$"+strings.Trim(latex, "$")+"$$")
}

func PDF(pdf []byte) Data {
	return MakeData(MIMETypePDF, pdf)
}

func PNG(png []byte) Data {
	return MakeData(MIMETypePNG, png)
}

func SVG(svg string) Data {
	return MakeData(MIMETypeSVG, svg)
}

// MIME encapsulates the data and metadata into a Data.
// The 'data' map is expected to contain at least one {key,value} pair,
// with value being a string, []byte or some other JSON serializable representation,
// and key equal to the MIME type of such value.
// The exact structure of value is determined by what the frontend expects.
// Some easier-to-use functions for common formats supported by the Jupyter frontend
// are provided by the various functions above.
func MIME(data, metadata MIMEMap) Data {
	return Data{data, metadata, nil}
}

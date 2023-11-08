package http

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alipourhabibi/exercises-journal/logging/logger"
)

var osPathSeparator = string(filepath.Separator)

type directoryListingFileData struct {
	Name  string
	Size  int64
	IsDir bool
	URL   *url.URL
}
type directoryListingData struct {
	Files []directoryListingFileData
}

var directoryListingTemplate = template.Must(template.New("").Parse(directoryListingTemplateText))

const directoryListingTemplateText = `
<html>
<body>
{{ if .Files }}
<table>
	<thead>
		<th></th>
		<th colspan=2 class=number>Size (bytes)</th>
	</thead>
	<tbody>
	{{ range .Files }}
	<tr>
		{{ if (not .IsDir) }}
		<td class=text><a href="{{ .URL.String }}">{{ .Name }}</td>
		<td class=number>{{ .Size | printf "%d" }}</td>
		{{ else }}
		<td colspan=3 class=text><a href="{{ .URL.String }}">{{ .Name }}</td>
		{{ end }}
	</tr>
	{{ end }}
	</tbody>
</table>
{{ end }}
</body>
</html>
`

type fileServer struct {
	route  string
	path   string
	logger logger.Logger
}

func NewFileServer(logger logger.Logger, path, route string) *fileServer {
	return &fileServer{
		logger: logger,
		route:  route,
		path:   path,
	}
}

type responseLogger struct {
	http.ResponseWriter
	statusCode int
	logger     logger.Logger
}

func (f *fileServer) serveStatus(w http.ResponseWriter, r *http.Request, status int) error {
	w.WriteHeader(status)
	_, err := w.Write([]byte(http.StatusText(status)))
	if err != nil {
		return err
	}
	return nil
}

func (f *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f.logger.GetLevel() >= logger.LevelDebug {
		f.logger.Debug("Start [%s] %s %s %s", f.path, r.RemoteAddr, r.Method, r.URL.String())
	}
	urlPath := r.URL.Path
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}
	urlPath = strings.TrimPrefix(urlPath, f.route)
	urlPath = strings.TrimPrefix(urlPath, "/"+f.route)

	osPath := strings.ReplaceAll(urlPath, "/", osPathSeparator)
	osPath = filepath.Clean(osPath)
	osPath = filepath.Join(f.path, osPath)
	info, err := os.Stat(osPath)
	switch {
	case os.IsNotExist(err):
		_ = f.serveStatus(w, r, http.StatusNotFound)
		if f.logger.GetLevel() >= logger.LevelInfo {
			f.logger.Info("Request %s, %s | Response %d", r.URL.String(), r.Method, http.StatusNotFound)
		}
	case os.IsPermission(err):
		_ = f.serveStatus(w, r, http.StatusForbidden)
		if f.logger.GetLevel() >= logger.LevelInfo {
			f.logger.Info("Request %s, %s | Response %d", r.URL.String(), r.Method, http.StatusForbidden)
		}
	case err != nil:
		_ = f.serveStatus(w, r, http.StatusInternalServerError)
		if f.logger.GetLevel() >= logger.LevelError {
			f.logger.Info("Request %s, %s | Response %d", r.URL.String(), r.Method, http.StatusInternalServerError)
		}
	case info.IsDir():
		err := f.serveDir(w, r, osPath)
		if err != nil {
			_ = f.serveStatus(w, r, http.StatusInternalServerError)
			if f.logger.GetLevel() >= logger.LevelError {
				f.logger.Info("Request %s, %s | Response %d", r.URL.String(), r.Method, http.StatusInternalServerError)
			}
		} else {
			if f.logger.GetLevel() >= logger.LevelInfo {
				f.logger.Info("Request %s, %s | Response %d", r.URL.String(), r.Method, http.StatusOK)
			}
		}
	default:
		http.ServeFile(w, r, osPath)
		if f.logger.GetLevel() >= logger.LevelInfo {
			f.logger.Info("Request %s, %s | Response %d", r.URL.String(), r.Method, http.StatusOK)
		}
	}
}

func (f *fileServer) serveDir(w http.ResponseWriter, r *http.Request, osPath string) error {
	d, err := os.Open(osPath)
	if err != nil {
		return err
	}
	files, err := d.Readdir(-1)
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
	return directoryListingTemplate.Execute(w, directoryListingData{
		Files: func() (out []directoryListingFileData) {
			for _, d := range files {
				name := d.Name()
				if d.IsDir() {
					name += osPathSeparator
				}
				fileData := directoryListingFileData{
					Name:  name,
					IsDir: d.IsDir(),
					Size:  d.Size(),
					URL: func() *url.URL {
						url := *r.URL
						url.Path = path.Join(url.Path, name)
						if d.IsDir() {
							url.Path += "/"
						}
						return &url
					}(),
				}
				out = append(out, fileData)
			}
			return out
		}(),
	})
}

func (f *fileServer) Run(port uint16) error {
	f.logger.Info("Starting server at %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), f)
	if err != nil {
		f.logger.Error("Starting server at %d, err: %v", port, err)
	}
	return nil
}

package devserver

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"sync"
	"time"

	"github.com/twharmon/gouix/build"
	"github.com/twharmon/gouix/utils"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"github.com/twharmon/slices"
)

type Server struct {
	upgrader  websocket.Upgrader
	listeners []*websocket.Conn
	mu        sync.Mutex
	watcher   *fsnotify.Watcher
	loaded    bool
	build     *build.Build
}

func New() (*Server, error) {
	os.Setenv("DEBUG", "true")
	s := &Server{
		build: build.New(),
	}
	if err := s.watchAll(); err != nil {
		return nil, fmt.Errorf("devserver.New: %w", err)
	}
	go s.watch()
	return s, nil
}

func (s *Server) Run() error {
	http.HandleFunc("/hot", s.ws())
	http.HandleFunc("/", s.files())
	if err := s.build.Run(); err != nil {
		go s.reportBuildError(fmt.Errorf("devserver.Server.Run: %s", err))
	}
	s.openBrowser()
	return http.ListenAndServe(utils.DevServerPort(), nil)
}

func (s *Server) Shutdown() error {
	if err := os.RemoveAll(s.build.BuildDir()); err != nil {
		return fmt.Errorf("build.Run: %w", err)
	}
	return s.watcher.Close()
}

func (s *Server) files() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(s.build.BuildDir(), r.URL.Path))
	}
}

func (s *Server) ws() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("devserver.Server.ws: %s\n", err)
			return
		}
		defer conn.Close()
		s.mu.Lock()
		s.listeners = append(s.listeners, conn)
		s.mu.Unlock()
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				break
			}
			if bytes.Equal(p, []byte("loaded")) {
				s.loaded = true
			}
		}
		s.mu.Lock()
		s.listeners = slices.Filter(s.listeners, func(listener *websocket.Conn) bool {
			return listener != conn
		})
		s.mu.Unlock()
	}
}

func (s *Server) watch() {
	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) > 0 {
				if err := s.build.Run(); err != nil {
					s.reportBuildError(fmt.Errorf("devserver.Server.watch: %s", err))
				} else {
					s.sendMessage("reload")
				}
			}
		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("devserver.Server.watch: %s\n", err)
		}
	}
}

func (s *Server) reportBuildError(err error) {
	fmt.Printf("\nError: %s\n\n", err)
	s.sendMessage(err.Error())
}

func (s *Server) watchAll() error {
	var err error
	s.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("devserver.Server.getWatcher: %w", err)
	}
	if err := s.watchDir("src"); err != nil {
		return fmt.Errorf("devserver.Server.watchAll: %w", err)
	}
	if err := s.watchDir("public"); err != nil {
		return fmt.Errorf("devserver.Server.watchAll: %w", err)
	}
	return nil
}

func (s *Server) watchDir(dir string) error {
	if err := s.watcher.Add(dir); err != nil {
		return fmt.Errorf("devserver.Server.watchDir: %w", err)
	}
	if fi, err := os.Stat(dir); err != nil {
		return fmt.Errorf("devserver.Server.watchDir: %w", err)
	} else if fi.IsDir() {
		fis, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("devserver.Server.watchDir: %w", err)
		}
		for _, fi := range fis {
			if fi.IsDir() {
				if err := s.watchDir(path.Join(dir, fi.Name())); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Server) sendMessage(msg string) {
	for {
		if s.loaded {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	s.mu.Lock()
	for i := len(s.listeners) - 1; i >= 0; i-- {
		if s.listeners[i].WriteMessage(websocket.TextMessage, []byte(msg)) != nil {
			s.listeners = slices.Splice(s.listeners, i, 1)
		}
	}
	s.mu.Unlock()
}

func (s *Server) openBrowser() {
	url := utils.DevServerUrl()
	var err error
	switch runtime.GOOS {
	case "linux":
		err = utils.Command("xdg-open", url)
	case "windows":
		err = utils.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		err = utils.Command("open", url)
	default:
		return
	}
	if err != nil {
		fmt.Printf("devserver.Server.openBrowser: %s\n", err)
	}
}

package modelcache

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const payload = "onnx-model-bytes"

func payloadHash() string { return fmt.Sprintf("%x", sha256.Sum256([]byte(payload))) }

func newServer(t *testing.T, status int) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(payload))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func setHome(t *testing.T, home string) {
	t.Helper()
	t.Setenv("HOME", home)        // unix
	t.Setenv("USERPROFILE", home) // windows
}

func TestDir(t *testing.T) {
	home := t.TempDir()
	setHome(t, home)

	dir, err := Dir("models/ocr")
	if err != nil {
		t.Fatalf("Dir() error = %v", err)
	}
	if want := filepath.Join(home, Root, "models", "ocr"); dir != want {
		t.Errorf("Dir() = %q, want %q", dir, want)
	}
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		t.Errorf("Dir() did not create %s", dir)
	}
}

func TestExists(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "f")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}

	if !Exists(file) {
		t.Error("Exists(file) = false, want true")
	}
	if Exists(dir) {
		t.Error("Exists(dir) = true, want false")
	}
	if Exists(filepath.Join(dir, "missing")) {
		t.Error("Exists(missing) = true, want false")
	}
}

func TestDownload(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		hash    string
		wantErr error
	}{
		{"ok with hash", http.StatusOK, payloadHash(), nil},
		{"ok without hash", http.StatusOK, "", nil},
		{"hash mismatch", http.StatusOK, "deadbeef", ErrHashMismatch},
		{"http error", http.StatusNotFound, "", ErrDownloadFailed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dest := filepath.Join(t.TempDir(), "sub", "model.onnx")

			err := Download(newServer(t, tt.status).URL, dest, tt.hash)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Download() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				if Exists(dest) || Exists(dest+".tmp") {
					t.Error("failed download left a file behind")
				}
				return
			}
			got, err := os.ReadFile(dest) //nolint:gosec // G304: path built from t.TempDir()
			if err != nil || string(got) != payload {
				t.Errorf("downloaded file = %q, %v; want %q", got, err, payload)
			}
		})
	}
}

func TestResolveOrDownload(t *testing.T) {
	existing := filepath.Join(t.TempDir(), "m.onnx")
	if err := os.WriteFile(existing, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("config path wins", func(t *testing.T) {
		got, err := ResolveOrDownload(existing, "", t.TempDir(), "m.onnx", "http://unused", "")
		if err != nil || got != existing {
			t.Errorf("got %q, %v; want %q", got, err, existing)
		}
	})
	t.Run("missing config path", func(t *testing.T) {
		_, err := ResolveOrDownload(existing+".nope", "", t.TempDir(), "m.onnx", "http://unused", "")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("error = %v, want ErrNotFound", err)
		}
	})
	t.Run("env var", func(t *testing.T) {
		t.Setenv("PDF2MD_TEST_MODEL", existing)
		got, err := ResolveOrDownload("", "PDF2MD_TEST_MODEL", t.TempDir(), "m.onnx", "http://unused", "")
		if err != nil || got != existing {
			t.Errorf("got %q, %v; want %q", got, err, existing)
		}
	})
	t.Run("missing env path", func(t *testing.T) {
		t.Setenv("PDF2MD_TEST_MODEL", existing+".nope")
		_, err := ResolveOrDownload("", "PDF2MD_TEST_MODEL", t.TempDir(), "m.onnx", "http://unused", "")
		if !errors.Is(err, ErrNotFound) {
			t.Errorf("error = %v, want ErrNotFound", err)
		}
	})
	t.Run("cache hit skips download", func(t *testing.T) {
		cache := filepath.Dir(existing)
		got, err := ResolveOrDownload("", "", cache, "m.onnx", "http://unused", "")
		if err != nil || got != existing {
			t.Errorf("got %q, %v; want %q", got, err, existing)
		}
	})
	t.Run("cache miss downloads", func(t *testing.T) {
		cache := t.TempDir()
		got, err := ResolveOrDownload("", "", cache, "m.onnx", newServer(t, http.StatusOK).URL, payloadHash())
		if want := filepath.Join(cache, "m.onnx"); err != nil || got != want {
			t.Errorf("got %q, %v; want %q", got, err, want)
		}
	})
}

func TestResolveRuntime(t *testing.T) {
	lib := RuntimeLibName()

	t.Run("config path", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), lib)
		if err := os.WriteFile(p, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		if got, err := ResolveRuntime(p); err != nil || got != p {
			t.Errorf("got %q, %v; want %q", got, err, p)
		}
		if _, err := ResolveRuntime(p + ".nope"); !errors.Is(err, ErrRuntimeNotFound) {
			t.Errorf("error = %v, want ErrRuntimeNotFound", err)
		}
	})
	t.Run("ORT_LIB_PATH directory", func(t *testing.T) {
		dir := t.TempDir()
		p := filepath.Join(dir, lib)
		if err := os.WriteFile(p, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("ORT_LIB_PATH", dir)
		if got, err := ResolveRuntime(""); err != nil || got != p {
			t.Errorf("got %q, %v; want %q", got, err, p)
		}
	})
	t.Run("ORT_LIB_PATH missing", func(t *testing.T) {
		t.Setenv("ORT_LIB_PATH", t.TempDir())
		if _, err := ResolveRuntime(""); !errors.Is(err, ErrRuntimeNotFound) {
			t.Errorf("error = %v, want ErrRuntimeNotFound", err)
		}
	})
	t.Run("home lib dir", func(t *testing.T) {
		home := t.TempDir()
		setHome(t, home)
		t.Setenv("ORT_LIB_PATH", "")
		p := filepath.Join(home, Root, "lib", lib)
		if err := os.MkdirAll(filepath.Dir(p), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		if got, err := ResolveRuntime(""); err != nil || got != p {
			t.Errorf("got %q, %v; want %q", got, err, p)
		}
	})
}

func TestStandardLibraryPaths(t *testing.T) {
	if len(StandardLibraryPaths()) == 0 {
		t.Error("StandardLibraryPaths() is empty")
	}
}

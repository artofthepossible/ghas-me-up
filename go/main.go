package main

// Go build for the GHAS Me Up demo.
// Renders go/index.html with html/template and lists the modules that were
// actually linked into this binary (runtime/debug build info).
//
// go.mod pins the language version and every module to old, vulnerable releases
// on purpose.

import (
	"crypto/md5"
	"crypto/tls"
	_ "embed"
	"encoding/hex"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sort"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

//go:embed template.html
var pageTemplate string

//go:embed config.yaml
var configYAML []byte

// ---------------------------------------------------------------------------
// BUILD-TIME FLAWS (intentional) — expect CodeQL findings below
// ---------------------------------------------------------------------------

const deployPassword = "correct-horse-battery-demo"
const jwtSigningKey = "demo-signing-key-do-not-use"

// buildStamp hashes a password-like value with MD5 (weak sensitive-data hashing).
func buildStamp() string {
	sum := md5.Sum([]byte(deployPassword))
	return hex.EncodeToString(sum[:])[:8]
}

// previewToken signs a JWT with a hard-coded key.
func previewToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "preview",
		"iat": time.Now().Unix(),
	})
	return token.SignedString([]byte(jwtSigningKey))
}

// statusCode checks a URL with TLS certificate verification disabled.
func statusCode(url string) (int, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

// ---------------------------------------------------------------------------

type config struct {
	Title  string `yaml:"title"`
	Accent string `yaml:"accent"`
	Locale string `yaml:"locale"`
}

type module struct {
	Path    string
	Version string
}

type page struct {
	Title        string
	Accent       string
	Locale       string
	GoVersion    string
	BuiltAt      string
	BuildStamp   string
	PreviewToken string
	BcryptOK     bool
	Modules      []module
}

func main() {
	out := flag.String("out", "../dist", "output directory")
	flag.Parse()

	var cfg config
	if err := yaml.Unmarshal(configYAML, &cfg); err != nil {
		log.Fatalf("config: %v", err)
	}

	token, err := previewToken()
	if err != nil {
		log.Fatalf("token: %v", err)
	}

	// Optional: only runs when a status URL is provided.
	if url := os.Getenv("GHAS_STATUS_URL"); url != "" {
		code, err := statusCode(url)
		log.Printf("status check %s -> %d (err=%v)", url, code, err)
	}

	// bcrypt is the right way to hash a password; shown next to the MD5 stamp.
	hash, err := bcrypt.GenerateFromPassword([]byte(deployPassword), bcrypt.MinCost)
	if err != nil {
		log.Fatalf("bcrypt: %v", err)
	}
	bcryptOK := bcrypt.CompareHashAndPassword(hash, []byte(deployPassword)) == nil

	var modules []module
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			modules = append(modules, module{Path: dep.Path, Version: dep.Version})
		}
	}
	sort.Slice(modules, func(i, j int) bool { return modules[i].Path < modules[j].Path })

	data := page{
		Title:        cfg.Title,
		Accent:       cfg.Accent,
		Locale:       language.Make(cfg.Locale).String(),
		GoVersion:    runtime.Version(),
		BuiltAt:      time.Now().UTC().Format(time.RFC3339),
		BuildStamp:   buildStamp(),
		PreviewToken: token,
		BcryptOK:     bcryptOK,
		Modules:      modules,
	}

	tmpl := template.Must(template.New("page").Parse(pageTemplate))

	target := filepath.Join(*out, "go")
	if err := os.MkdirAll(target, 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	file, err := os.Create(filepath.Join(target, "index.html"))
	if err != nil {
		log.Fatalf("create: %v", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		log.Fatalf("render: %v", err)
	}
	fmt.Printf("go build complete -> %s\n", filepath.Join(target, "index.html"))
}

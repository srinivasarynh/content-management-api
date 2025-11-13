package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Content struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Phone     string   `json:"phone"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
}

type ContentStore struct {
	mu       sync.RWMutex
	contents map[string]Content
	nextID   int
}

type Stats struct {
	TotalContents int            `json:"total_contents"`
	TagCounts     map[string]int `json:"tag_counts"`
	RecentIDs     [5]string
}

var store = &ContentStore{
	contents: make(map[string]Content),
	nextID:   1,
}

func (s *ContentStore) Create(c Content) Content {
	s.mu.Lock()
	defer s.mu.Unlock()

	c.ID = fmt.Sprintf("content_%d", s.nextID)
	c.CreatedAt = time.Now().Format(time.RFC3339)
	s.nextID++
	s.contents[c.ID] = c
	return c
}

func (s *ContentStore) GetAll() []Content {
	s.mu.RLock()
	defer s.mu.RUnlock()

	contents := make([]Content, 0, len(s.contents))
	for _, c := range s.contents {
		contents = append(contents, c)
	}
	return contents
}

func (s *ContentStore) GetByID(id string) (Content, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, exists := s.contents[id]
	return c, exists
}

func (s *ContentStore) GetStats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var recent [5]string
	tagCounts := make(map[string]int)

	i := 0
	for id, content := range s.contents {
		if i < 5 {
			recent[i] = id
			i++
		}

		for _, tag := range content.Tags {
			tagCounts[tag]++
		}
	}

	return Stats{
		TotalContents: len(s.contents),
		TagCounts:     tagCounts,
		RecentIDs:     recent,
	}
}

func createContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var content Content
	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if content.Name == "" || content.Email == "" {
		http.Error(w, "name and email required", http.StatusBadRequest)
		return
	}

	created := store.Create(content)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(created)
}

func listContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contents := store.GetAll()
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(contents)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	stats := store.GetStats()
	w.Header().Set("Content-type", "applicatioin/json")
	json.NewEncoder(w).Encode(stats)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	contents := store.GetAll()
	stats := store.GetStats()

	tmpl := template.Must(template.New("home").Funcs(template.FuncMap{
		"join": strings.Join,
	}).Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Content Manager</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 1200px; margin: 0 auto; padding: 20px; }
        .stats { background: #f0f0f0; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
        .form { background: #e8f4f8; padding: 20px; border-radius: 5px; margin-bottom: 20px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        th { background: #4CAF50; color: white; }
        input { padding: 8px; margin: 5px; }
        button { padding: 10px 20px; background: #4CAF50; color: white; border: none; cursor: pointer; }
        .tag { background: #2196F3; color: white; padding: 3px 8px; border-radius: 3px; margin: 2px; }
    </style>
</head>
<body>
    <h1>📇 Content Manager</h1>
    
    <div class="stats">
        <h3>📊 Statistics</h3>
        <p><strong>Total Contents:</strong> {{.Stats.TotalContents}}</p>
        <p><strong>Tag Distribution:</strong></p>
        <ul>
        {{range $tag, $count := .Stats.TagCounts}}
            <li>{{$tag}}: {{$count}}</li>
        {{end}}
        </ul>
    </div>

    <div class="form">
        <h3>➕ Add New Content</h3>
        <form id="contentForm">
            <input type="text" id="name" placeholder="Name" required>
            <input type="email" id="email" placeholder="Email" required>
            <input type="tel" id="phone" placeholder="Phone">
            <input type="text" id="tags" placeholder="Tags (comma-separated)">
            <button type="submit">Add Content</button>
        </form>
    </div>

    <h3>👥 All Contents ({{len .Contents}})</h3>
    <table>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Email</th>
            <th>Phone</th>
            <th>Tags</th>
            <th>Created</th>
        </tr>
        {{range .Contents}}
        <tr>
            <td>{{.ID}}</td>
            <td>{{.Name}}</td>
            <td>{{.Email}}</td>
            <td>{{.Phone}}</td>
            <td>
                {{range .Tags}}<span class="tag">{{.}}</span>{{end}}
            </td>
            <td>{{.CreatedAt}}</td>
        </tr>
        {{end}}
    </table>

    <script>
        document.getElementById('contentForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const tags = document.getElementById('tags').value.split(',').map(t => t.trim()).filter(t => t);
            const content = {
                name: document.getElementById('name').value,
                email: document.getElementById('email').value,
                phone: document.getElementById('phone').value,
                tags: tags
            };
            
            const res = await fetch('/api/contents', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(content)
            });
            
            if (res.ok) {
                location.reload();
            } else {
                alert('Error creating content');
            }
        });
    </script>
</body>
</html>
	`))

	data := struct {
		Contents []Content
		Stats    Stats
	}{
		Contents: contents,
		Stats:    stats,
	}

	tmpl.Execute(w, data)
}

func main() {
	store.Create(Content{
		Name:  "Alice Johnson",
		Email: "alice@example.com",
		Phone: "+1-555-0100",
		Tags:  []string{"client", "vip"},
	})
	store.Create(Content{
		Name:  "Bob Smith",
		Email: "bob@example.com",
		Phone: "+1-555-0101",
		Tags:  []string{"vendor", "tech"},
	})

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/contents", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listContentHandler(w, r)
		case http.MethodPost:
			createContentHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/stats", statsHandler)

	port := "8080"
	fmt.Printf("server starting on http://localhost:%s\n", port)
	fmt.Println("📖 Open browser to see the web UI")
	fmt.Println("🔌 API endpoints:")
	fmt.Println("   GET  /api/contacts - List all contacts")
	fmt.Println("   POST /api/contacts - Create contact")
	fmt.Println("   GET  /api/stats    - View statistics")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

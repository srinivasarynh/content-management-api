# 📇 Content Manager API

A simple production-grade content management system built with Go to demonstrate Chapter 4 concepts: Arrays, Slices, Maps, Structs, JSON, and Templates.

## ✨ Features

- 📝 Create and list contacts
- 🏷️ Tag-based organization
- 📊 Real-time statistics
- 🌐 RESTful JSON API
- 🎨 Web UI interface
- 🔒 Thread-safe operations

## 🚀 Quick Start

### Prerequisites
- Go 1.16 or higher

### Run the Application

```bash
go run main.go
```

The server starts on `http://localhost:8080`

## 📡 API Endpoints

### List All Contacts
```bash
GET /api/contents
```

### Create Contact
```bash
POST /api/contents
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1-555-0123",
  "tags": ["client", "vip"]
}
```

### View Statistics
```bash
GET /api/stats
```

## 🧪 Test with cURL

```bash
# Create a contact
curl -X POST http://localhost:8080/api/contents \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com","tags":["vendor"]}'

# List all contacts
curl http://localhost:8080/api/contents

# View stats
curl http://localhost:8080/api/stats
```

## 📚 Learning Concepts

This project demonstrates:

- **Arrays**: Fixed-size storage for recent contact IDs
- **Slices**: Dynamic lists for contacts and tags
- **Maps**: Key-value store for fast contact lookups
- **Structs**: Data models (Contact, Stats)
- **JSON**: API serialization/deserialization
- **Templates**: Dynamic HTML generation for web UI

## 📁 Project Structure

```
main.go          # Complete application
├── Structs      # Contact, ContactStore, Stats
├── Maps         # In-memory contact storage
├── Slices       # Dynamic contact lists
├── Arrays       # Fixed recent IDs
├── JSON         # API handlers
└── Templates    # Web UI rendering
```

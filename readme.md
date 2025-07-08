# 📚 Book Club API Spec

This project defines a simple Book Club API built with Go + Gin, using mock in-memory data for simplicity. It's designed to simulate a real-world backend system that supports a frontend web dashboard.


---

## 🎯 Objective

Build a RESTful API that allows users to:

- View book clubs
- Join clubs
- Propose books
- Vote on books
- Track the currently selected book

This API should support a mock frontend app (`book-club-ui`) and be testable with curl, Postman, Thunder Client, or REST Client extensions.

---

## 🧑 Assumed User

For now, the logged-in user is hardcoded:

```json
{
  "username": "alice"
}

All membership and voting logic should assume requests come from this user.

# 📘 Endpoint Reference

### `GET /clubs`

Returns all book clubs and their metadata.

**Response:**

```json
[
  {
    "id": "sf_club",
    "name": "Sci-Fi Enthusiasts",
    "members": ["alice", "bob"],
    "current_book": {
      "id": "neuromancer",
      "title": "Neuromancer",
      "author": "William Gibson",
      "votes": 5
    }
  }
]
```

---

### `POST /clubs`

Create a new book club.

**Request:**

```json
{
  "id": "fantasy_lovers",
  "name": "Fantasy Lovers"
}
```

**Response:**
Created club object.

---

### `GET /clubs/:id`

Get full detail of a club including members and proposals.

**Response:**
Club object matching the given ID.

---

### `POST /clubs/:id/members`

Add a user to the club.

**Request:**

```json
{
  "name": "alice"
}
```

**Response:**
Updated club object.

---

### `DELETE /clubs/:id/members/:name`

Remove a user from a club by name.

**Response:**
Updated club object.

---

### `GET /clubs/:id/books`

List all proposed books for a club.

**Response:**

```json
[
  {
    "id": "dune",
    "title": "Dune",
    "author": "Frank Herbert",
    "votes": 2
  }
]
```

---

### `POST /clubs/:id/books`

Propose a new book for the club.

**Request:**

```json
{
  "id": "dune",
  "title": "Dune",
  "author": "Frank Herbert"
}
```

**Response:**
Updated list of proposals.

---

### `POST /clubs/:id/books/:bookID/vote`

Vote for a book.

**Response:**
Updated vote count.

---

### `GET /clubs/:id/current`

Get the currently selected book (most-voted).

**Response:**

```json
{
  "id": "dune",
  "title": "Dune",
  "author": "Frank Herbert",
  "votes": 4
}
```

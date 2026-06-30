# Real-Time-Forum

**Product Requirements Document — Forum Platform v2**

## Executive Summary

This release upgrades the existing forum into a real-time, single-page application. The previous forum's codebase may be partially reused, but the upgrade is substantial enough that it should be treated as a new build rather than a patch: a new authentication model, a refreshed posts-and-comments experience, and an entirely new private messaging system delivered over WebSockets.

The product is delivered as a **single HTML file**. Every page transition — feed, post detail, chat, registration, login — is handled client-side in JavaScript. There is no server-rendered page navigation once the application has loaded.

## Objectives

Three product pillars define this release. All other decisions in this document support one of them:

1. **Registration and Login** — gate the forum behind authentication.
2. **Posts and Comments** — let users create, browse and discuss posts within categories.
3. **Private Messages** — let users message one another in real time, with a Discord-style online/offline presence list.

## Architecture Overview

The application is composed of five distinct layers. Each has a single, well-defined responsibility:

| Layer | Technology | Responsibility |
|---|---|---|
| Data | SQLite | Persistent storage for users, posts, comments and messages |
| Backend | Go (Golang) | Business logic, data handling, and WebSocket connections |
| Client logic | JavaScript | All frontend events and the client side of the WebSocket connection |
| Structure | HTML | A single page containing the markup for every view |
| Presentation | CSS | Styling of all elements |

**Architectural constraint:** because the product ships as one HTML file, view-switching (e.g. moving from the feed to a post, or from the feed to a chat) must be implemented entirely in JavaScript. This is a single-page application by design, not by convenience.

## Functional Requirements

### Epic 1 — Registration and Login

Until a user has registered and logged in, they are restricted to the registration or login screen — no other part of the forum is reachable.

**Registration**
- A registration form must capture, as a minimum:
  - Nickname
  - Age
  - Gender
  - First name
  - Last name
  - E-mail
  - Password

**Login**
- A user must be able to log in using **either their nickname or their e-mail**, combined with their password.

**Logout**
- A logged-in user must be able to log out from **any page** of the forum, not only from a dedicated account page.

### Epic 2 — Posts and Comments

This functionality carries over conceptually from the first forum, with the same category model.

- Users can create posts.
- Posts are organised by category, as in the first forum.
- Users can comment on a post.
- Posts are presented in a **feed display**.
- Comments are not shown in the feed — a user only sees a post's comments after clicking into that post.

### Epic 3 — Private Messages

This is the most significant new capability in this release, and the area with the most precise requirements.

**Presence list**
- A section of the forum shows which users are online and offline, and is always visible.
- This list is ordered by the **time of the last message sent**, in the same style as Discord.
- A user with no message history is placed in the list in **alphabetical order**.
- Users can send a private message to any user shown as online.

**Conversation view**
- Clicking on a user reloads the message history between the current user and that user.
- The conversation must show the previous messages exchanged between the two users.
- On open, the conversation loads the **last 10 messages**.
- Scrolling up loads **10 further messages** at a time.
- The scroll event that triggers loading must be **throttled or debounced** — it must not fire a fresh load on every pixel of scroll movement.

**Message format**
- Every message must display:
  - The **date** the message was sent.
  - The **user name** of the sender.

**Real-time delivery**
- Messages must be delivered in real time. If User A sends a message, User B must receive it — and any associated notification — **without refreshing the page**.
- Real-time delivery is implemented through WebSockets on both the Go backend and the JavaScript frontend.

## Technical Constraints

**Allowed packages**

- All packages from the Go standard library.
- [`gorilla/websocket`](https://pkg.go.dev/github.com/gorilla/websocket)
- [`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3)
- [`golang.org/x/crypto/bcrypt`](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [`gofrs/uuid`](https://github.com/gofrs/uuid) or [`google/uuid`](https://github.com/google/uuid)

**Frontend constraint**

- No frontend libraries or frameworks (React, Angular, Vue, or equivalent) are permitted. The frontend is plain HTML, CSS and JavaScript.

**Code reuse note**

- Reuse of code from the first forum is permitted, but the brief is explicit that **not all of it** can be carried forward — the registration/login model and the private messaging system are new requirements with no equivalent in the original build.

## Non-Functional Requirements

- **Sessions and cookies** govern authentication state across the single-page application.
- **Concurrency on the backend** is handled using Go routines and Go channels, particularly to support multiple simultaneous WebSocket connections.
- **Scroll performance** in the chat view is protected using throttling or debouncing, so that loading older messages does not spam the scroll event.

## Learning Outcomes

This project is also a vehicle for building foundational competence in the following areas:

- The basics of the web: HTML, HTTP, CSS, backend vs. frontend, and the DOM.
- Sessions and cookies.
- Go routines and Go channels.
- WebSockets, on both the Go backend and the JavaScript frontend.
- SQL and the manipulation of relational databases.

## Status

Requirements captured and ready for technical design. Implementation to proceed epic by epic: Registration & Login → Posts & Comments → Private Messages.